package pubsub

import (
	"io"
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
