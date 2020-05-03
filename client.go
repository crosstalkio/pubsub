package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/crosstalkio/log"
	"github.com/crosstalkio/pubsub/api"
	"google.golang.org/grpc"
)

const (
	RequestTimeout = 10 * time.Second
)

type Client struct {
	log.Sugar
	conn *grpc.ClientConn
}

type Data struct {
	*api.Data
}

func NewClient(logger log.Logger, conn *grpc.ClientConn) *Client {
	return &Client{
		Sugar: log.NewSugar(logger),
		conn:  conn,
	}
}

func (c *Client) PublishText(ch, msg string) error {
	c.Debugf("Publishing %d bytes text message to: %s", len(msg), ch)
	_, err := api.NewPubSubClient(c.conn).Publish(context.Background(), &api.PublishRequest{
		Channel: ch,
		Type:    &api.PublishRequest_Text{Text: msg},
	})
	return err
}

func (c *Client) PublishBinary(ch string, msg []byte) error {
	c.Debugf("Publishing %d bytes binary message to: %s", len(msg), ch)
	_, err := api.NewPubSubClient(c.conn).Publish(context.Background(), &api.PublishRequest{
		Channel: ch,
		Type:    &api.PublishRequest_Binary{Binary: msg},
	})
	return err
}

func (c *Client) Subscribe(ch string, cb func(*Data)) error {
	c.Debugf("Subscribing: %s", ch)
	stream, err := api.NewPubSubClient(c.conn).Subscribe(context.Background(), &api.SubscribeRequest{Channel: ch})
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		data := &api.Data{
			Channel: ch,
		}
		switch v := res.Type.(type) {
		case *api.SubscribeResponse_Text:
			data.Type = &api.Data_Text{Text: v.Text}
		case *api.SubscribeResponse_Binary:
			data.Type = &api.Data_Binary{Binary: v.Binary}
		default:
			err = fmt.Errorf("Unknown type of subscribe response: %v", v)
			c.Errorf("%s", err.Error())
			return err
		}
	}
}
