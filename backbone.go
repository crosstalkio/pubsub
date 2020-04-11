package pubsub

import (
	"fmt"
	"io"
	"time"

	api "github.com/crosstalkio/pubsub/api/pubsub"
	"google.golang.org/protobuf/proto"
)

type Subscription interface {
	io.Closer
	Receive() ([]byte, error)
}

type Backbone interface {
	io.Closer
	Publish(ch string, payload []byte) error
	Subscribe(ch string) (Subscription, error)
}

func PublishText(b Backbone, as, ch, text string) error {
	if ch == "" {
		err := fmt.Errorf("Missing channel to publish")
		return err
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     as,
		Payload:  &api.Data_Text{Text: text},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = b.Publish(ch, data)
	if err != nil {
		return err
	}
	return nil
}

func PublishBinary(b Backbone, as, ch string, bin []byte) error {
	if ch == "" {
		err := fmt.Errorf("Missing channel to publish")
		return err
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     as,
		Payload:  &api.Data_Binary{Binary: bin},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = b.Publish(ch, data)
	if err != nil {
		return err
	}
	return nil
}
