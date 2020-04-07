package pubsub

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/crosstalkio/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	RequestTimeout = 10 * time.Second
)

type Client struct {
	log.Sugar
	u         *url.URL
	closed    bool
	localAddr net.Addr
	conn      *websocket.Conn
	subs      map[string]func(*Data)
	reqs      map[string]chan error
	writer    *clientWriter
}

func NewClient(logger log.Logger, u *url.URL) *Client {
	return &Client{
		Sugar: log.NewSugar(logger),
		u:     u,
		subs:  make(map[string]func(*Data)),
		reqs:  make(map[string]chan error),
	}
}

func (c *Client) Dial(header http.Header) error {
	c.Debugf("Dialing: %s", c.u.String())
	conn, res, err := websocket.DefaultDialer.Dial(c.u.String(), header)
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
	return c.request(&Request{
		Payload: &Request_Publish{
			Publish: &Publish{
				Channel: ch,
				Payload: &Publish_Text{Text: msg},
			},
		},
	})
}

func (c *Client) PublishBinary(ch string, msg []byte) error {
	c.Debugf("Publishing %d bytes binary message to: %s", len(msg), ch)
	return c.request(&Request{
		Payload: &Request_Publish{
			Publish: &Publish{
				Channel: ch,
				Payload: &Publish_Binary{Binary: msg},
			},
		},
	})
}

func (c *Client) Subscribe(ch string, cb func(*Data)) error {
	c.Debugf("Subscribing: %s", ch)
	c.subs[ch] = cb
	return c.request(&Request{
		Payload: &Request_Subscribe{
			Subscribe: &Subscribe{
				Channel: ch,
			},
		},
	})
}

func (c *Client) Unsubscribe(ch string) error {
	c.Debugf("Unsubscribing: %s", ch)
	return c.request(&Request{
		Payload: &Request_Unsubscribe{
			Unsubscribe: &Unsubscribe{
				Channel: ch,
			},
		},
	})
}

func (c *Client) request(req *Request) error {
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
	err = c.writer.write(&Message{
		Payload: &Message_Control{
			Control: &Control{
				Payload: &Control_Request{Request: req},
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
			msg := &Message{}
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

func (c *Client) handle(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *Message_Control:
		ctl := v.Control
		switch v := ctl.Payload.(type) {
		case *Control_Request:
			err := fmt.Errorf("Unexpected payload of Request")
			c.Errorf(err.Error())
			return err
		case *Control_Response:
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
			switch v := res.Payload.(type) {
			case *Response_Success:
				c.Debugf("Request success: %s", id)
				req <- nil
			case *Response_Error:
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
	case *Message_Data:
		data := v.Data
		ch := data.GetChannel()
		sub := c.subs[ch]
		if sub == nil {
			err := fmt.Errorf("Not subscribed: %s", ch)
			c.Errorf("%s", err.Error())
			return err
		}
		switch v := data.Payload.(type) {
		case *Data_Text:
			txt := v.Text
			c.Debugf("Dispatching %d bytes text data from: %s", len(txt), ch)
		case *Data_Binary:
			bin := v.Binary
			c.Debugf("Dispatching %d bytes binary data from: %s", len(bin), ch)
		default:
			err := fmt.Errorf("Unexpected type of Data: %v", v)
			c.Errorf(err.Error())
			return err
		}
		go sub(data)
	default:
		err := fmt.Errorf("Unexpected type of Message: %v", v)
		c.Errorf(err.Error())
		return err
	}
	return nil
}
