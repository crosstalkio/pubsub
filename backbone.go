package pubsub

import (
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

func PublishText(backbone Backbone, channel, from string, time time.Time, data string) error {
	return publish(backbone, channel, &api.Data{
		NanoTime: time.UnixNano(),
		From:     from,
		Type:     &api.Data_Text{Text: data},
	})
}

func PublishBinary(backbone Backbone, channel, from string, time time.Time, data []byte) error {
	return publish(backbone, channel, &api.Data{
		NanoTime: time.UnixNano(),
		From:     from,
		Type:     &api.Data_Binary{Binary: data},
	})
}

func publish(backbone Backbone, channel string, data *api.Data) error {
	bytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	err = backbone.Publish(channel, bytes)
	if err != nil {
		return err
	}
	return nil
}
