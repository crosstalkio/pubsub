package pubsub

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/crosstalkio/log"
	"github.com/crosstalkio/pubsub/api"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	RequestTimeout = 10 * time.Second
)

type TextData struct {
	Time    time.Time
	From    string
	Channel string
	Data    string
}

type BinaryData struct {
	Time    time.Time
	From    string
	Channel string
	Data    []byte
}

type Client struct {
	log.Sugar
	closed    bool
	localAddr net.Addr
	conn      *websocket.Conn
	subs      map[string]*subscriber
	reqs      map[string]chan error
	writer    *clientWriter
}

type subscriber struct {
	text   func(*TextData)
	binary func(*BinaryData)
}

func NewClient(logger log.Logger, timeout time.Duration) *Client {
	return &Client{
		Sugar: log.NewSugar(logger),
		subs:  make(map[string]*subscriber),
		reqs:  make(map[string]chan error),
	}
}

func (c *Client) Dial(u *url.URL, header http.Header) error {
	c.Debugf("Dialing: %s", u.String())
	conn, res, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		if res != nil {
			c.Errorf("Failed to dial: %s", res.Status)
			return fmt.Errorf("%s", res.Status)
		} else {
			c.Errorf("Failed to dial: %s", err.Error())
			return err
		}
	}
	c.localAddr = conn.LocalAddr()
	c.conn = conn
	c.writer = newClientWriter(c, conn)
	go c.writer.loop()
	go c.loop(conn)
	return nil
}

func (c *Client) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	err := fmt.Errorf("Closed")
	for _, req := range c.reqs {
		req <- err
	}
	c.reqs = make(map[string]chan error)
	if c.writer != nil {
		c.writer.exit()
		c.writer = nil
	}
	if c.conn != nil {
		err = c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(3*time.Second))
		if err != nil {
			c.Warningf("Failed to write close: %s", err.Error())
		}
		c.conn.Close()
		c.conn = nil
	}
	return nil
}

func (c *Client) PublishText(ch, msg string) error {
	c.Debugf("Publishing %d bytes text message to: %s", len(msg), ch)
	return c.request(&api.Request{
		Type: &api.Request_Publish{
			Publish: &api.Publish{
				Channel: ch,
				Type:    &api.Publish_Text{Text: msg},
			},
		},
	})
}

func (c *Client) PublishBinary(ch string, msg []byte) error {
	c.Debugf("Publishing %d bytes binary message to: %s", len(msg), ch)
	return c.request(&api.Request{
		Type: &api.Request_Publish{
			Publish: &api.Publish{
				Channel: ch,
				Type:    &api.Publish_Binary{Binary: msg},
			},
		},
	})
}

func (c *Client) SubscribeText(ch string) (chan *TextData, error) {
	c.Debugf("Subscribing text: %s", ch)
	msgCh := make(chan *TextData)
	err := c.subscribe(ch, func(msg *TextData) {
		msgCh <- msg
	}, nil)
	if err != nil {
		return nil, err
	}
	return msgCh, nil
}

func (c *Client) SubscribeBinary(ch string) (chan *BinaryData, error) {
	c.Debugf("Subscribing binary: %s", ch)
	msgCh := make(chan *BinaryData)
	err := c.subscribe(ch, nil, func(msg *BinaryData) {
		msgCh <- msg
	})
	if err != nil {
		return nil, err
	}
	return msgCh, nil
}

func (c *Client) Unsubscribe(ch string) error {
	c.Debugf("Unsubscribing: %s", ch)
	return c.request(&api.Request{
		Type: &api.Request_Unsubscribe{
			Unsubscribe: &api.Unsubscribe{
				Channel: ch,
			},
		},
	})
}

func (c *Client) subscribe(ch string, text func(*TextData), binary func(*BinaryData)) error {
	c.Debugf("Subscribing: %s", ch)
	c.subs[ch] = &subscriber{text: text, binary: binary}
	return c.request(&api.Request{
		Type: &api.Request_Subscribe{
			Subscribe: &api.Subscribe{
				Channel: ch,
			},
		},
	})
}

func (c *Client) request(req *api.Request) error {
	u, err := uuid.NewRandom()
	if err != nil {
		c.Errorf("Failed to create UUID: %s", err.Error())
		return err
	}
	id := hex.EncodeToString(u[:])
	req.Id = id
	c.Debugf("Making request: %s", id)
	resCh := make(chan error, 1)
	c.reqs[id] = resCh
	err = c.writer.write(&api.Message{
		Type: &api.Message_Control{
			Control: &api.Control{
				Type: &api.Control_Request{Request: req},
			},
		},
	})
	if err != nil {
		return err
	}
	t := time.NewTimer(RequestTimeout)
	select {
	case <-t.C:
		err = fmt.Errorf("Request timed out: %s", id)
		c.Error(err.Error())
	case err = <-resCh:
		t.Stop()
		c.Debugf("Received response: %s", id)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) loop(conn *websocket.Conn) {
	defer c.Debugf("Exiting reader: %s", c.localAddr.String())
	defer c.Close()
	for {
		c.Debugf("Receiving websocket message: %s", c.localAddr.String())
		mt, p, err := conn.ReadMessage()
		if err != nil {
			if !c.closed {
				c.Errorf("Failed to read: %s", err.Error())
			}
			return
		}
		switch mt {
		case websocket.BinaryMessage:
			c.Debugf("Received %d bytes websocket binary message: %s", len(p), c.localAddr.String())
			msg := &api.Message{}
			err := proto.Unmarshal(p, msg)
			if err != nil {
				c.Errorf("Failed to parse proto message: %s", err.Error())
				return
			}
			err = c.handle(msg)
			if err != nil {
				return
			}
		default:
			c.Errorf("Unexpected websocket message type %d from %s", mt, c.localAddr.String())
			return
		}
	}
}

func (c *Client) handle(msg *api.Message) error {
	switch v := msg.Type.(type) {
	case *api.Message_Control:
		ctl := v.Control
		switch v := ctl.Type.(type) {
		case *api.Control_Request:
			err := fmt.Errorf("Unexpected payload of Request")
			c.Errorf(err.Error())
			return err
		case *api.Control_Response:
			res := v.Response
			id := res.GetId()
			if id == "" {
				err := fmt.Errorf("Missing request ID")
				c.Errorf("%s", err.Error())
				return err
			}
			req := c.reqs[id]
			if req == nil {
				err := fmt.Errorf("No such request: %s", id)
				c.Errorf("%s", err.Error())
				return err
			}
			switch v := res.Type.(type) {
			case *api.Response_Success:
				c.Debugf("Request success: %s", id)
				req <- nil
			case *api.Response_Error:
				e := v.Error
				err := fmt.Errorf("%d %s", e.GetCode(), e.GetReason())
				c.Debugf("Request failed: %s", err.Error())
				req <- err
			default:
				err := fmt.Errorf("Unexpected type of Response: %v", v)
				c.Errorf(err.Error())
				return err
			}
		default:
			err := fmt.Errorf("Unexpected type of Control: %v", v)
			c.Errorf(err.Error())
			return err
		}
	case *api.Message_Data:
		data := v.Data
		ch := data.GetChannel()
		sub := c.subs[ch]
		if sub == nil {
			err := fmt.Errorf("Not subscribed: %s", ch)
			c.Errorf("%s", err.Error())
			return err
		}
		switch v := data.Type.(type) {
		case *api.Data_Text:
			txt := v.Text
			if sub.text != nil {
				c.Debugf("Dispatching %d bytes text data from: %s", len(txt), ch)
				go sub.text(&TextData{
					Time:    time.Unix(0, data.NanoTime),
					From:    data.From,
					Channel: data.Channel,
					Data:    v.Text,
				})
			} else {
				c.Debugf("Skipping %d bytes text data from: %s", len(txt), ch)
			}
		case *api.Data_Binary:
			bin := v.Binary
			if sub.binary != nil {
				c.Debugf("Dispatching %d bytes binary data from: %s", len(bin), ch)
				go sub.binary(&BinaryData{
					Time:    time.Unix(0, data.NanoTime),
					From:    data.From,
					Channel: data.Channel,
					Data:    v.Binary,
				})
			} else {
				c.Debugf("Skipping %d bytes binary data from: %s", len(bin), ch)
			}
		default:
			err := fmt.Errorf("Unexpected type of Data: %v", v)
			c.Errorf(err.Error())
			return err
		}
	default:
		err := fmt.Errorf("Unexpected type of Message: %v", v)
		c.Errorf(err.Error())
		return err
	}
	return nil
}
