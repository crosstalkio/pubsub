package pubsub

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/crosstalkio/log"
	"github.com/crosstalkio/pubsub/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	closed   bool
	conn     *grpc.ClientConn
	cancel   context.CancelFunc
	timeout  time.Duration
	metadata metadata.MD
	cancels  map[string]context.CancelFunc
}

func NewClient(logger log.Logger, timeout time.Duration) *Client {
	s := log.NewSugar(logger)
	return &Client{
		Sugar:    s,
		timeout:  timeout,
		metadata: make(metadata.MD),
		cancels:  make(map[string]context.CancelFunc),
	}
}

func (c *Client) Dial(addr string, opts ...grpc.DialOption) error {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(c.timeout)
	go func() {
		<-timer.C
		c.Errorf("Dial timed out: %v", c.timeout)
		cancel()
	}()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		return err
	}
	timer.Stop()
	c.cancel = cancel
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	for _, cancel := range c.cancels {
		cancel()
	}
	c.cancels = nil
	return nil
}

func (c *Client) SetMetadata(key string, vals ...string) {
	c.metadata.Set(key, vals...)
}

func (c *Client) PublishText(ch, msg string) error {
	c.Debugf("Publishing %d bytes text message to: %s", len(msg), ch)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, c.metadata)
	_, err := api.NewPubSubClient(c.conn).Publish(ctx, &api.PublishRequest{
		Channel: ch,
		Type:    &api.PublishRequest_Text{Text: msg},
	})
	return err
}

func (c *Client) PublishBinary(ch string, msg []byte) error {
	c.Debugf("Publishing %d bytes binary message to: %s", len(msg), ch)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, c.metadata)
	_, err := api.NewPubSubClient(c.conn).Publish(ctx, &api.PublishRequest{
		Channel: ch,
		Type:    &api.PublishRequest_Binary{Binary: msg},
	})
	return err
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

func (c *Client) subscribe(ch string, text func(*TextData), binary func(*BinaryData)) error {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = metadata.NewOutgoingContext(ctx, c.metadata)
	stream, err := api.NewPubSubClient(c.conn).Subscribe(ctx, &api.SubscribeRequest{Channel: ch})
	if err != nil {
		cancel()
		return err
	}
	c.cancels[ch] = cancel
	var addr net.Addr
	if peer, ok := peer.FromContext(stream.Context()); ok {
		addr = peer.Addr
	}
	errCh := make(chan error, 1)
	go func() {
		c.Infof("Receiving stream: %s", addr.String())
		defer func() {
			if text != nil {
				text(nil)
			}
			if text != nil {
				binary(nil)
			}
		}()
		errCh <- nil
		for {
			res, err := stream.Recv()
			if err != nil {
				c.Errorf("Failed to receive stream: %s", err.Error())
				return
			}
			switch v := res.Type.(type) {
			case *api.SubscribeResponse_Text:
				if text != nil {
					text(&TextData{
						Time:    time.Unix(0, res.NanoTime),
						From:    res.From,
						Channel: res.Channel,
						Data:    v.Text,
					})
				}
			case *api.SubscribeResponse_Binary:
				if binary != nil {
					binary(&BinaryData{
						Time:    time.Unix(0, res.NanoTime),
						From:    res.From,
						Channel: res.Channel,
						Data:    v.Binary,
					})
				}
			default:
				err = fmt.Errorf("Unknown type of subscribe response: %v", v)
				c.Errorf("%s", err.Error())
				return
			}
		}
	}()
	return <-errCh
}

func (c *Client) Unsubscribe(ch string) error {
	cancel, ok := c.cancels[ch]
	if !ok {
		err := fmt.Errorf("Not subscribed: %s", ch)
		c.Errorf("%s", err.Error())
		return err
	}
	delete(c.cancels, ch)
	cancel()
	return nil
}
