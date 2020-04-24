package google

import (
	"encoding/base64"
	"fmt"

	"github.com/crosstalkio/log"
	"github.com/crosstalkio/pubsub"
	"github.com/go-redis/redis/v7"
)

type backbone struct {
	log.Sugar
	client  *redis.Client
	pubsub  *redis.PubSub
	prefix  string
	subs    map[string][]*subscription
	msgCh   chan *redis.Message
	subCh   chan *subscription
	unsCh   chan *subscription
	closeCh chan bool
}

type Options struct {
	redis.Options
	Prefix string
}

func NewBackbone(logger log.Logger, options *Options) (pubsub.Backbone, error) {
	s := log.NewSugar(logger)
	client := redis.NewClient(&options.Options)
	err := client.Ping().Err()
	if err != nil {
		s.Errorf("Failed to ping redis: %s", err.Error())
		return nil, err
	}
	b := &backbone{
		Sugar:   s,
		client:  client,
		pubsub:  client.Subscribe(),
		prefix:  options.Prefix,
		subs:    make(map[string][]*subscription),
		msgCh:   make(chan *redis.Message),
		subCh:   make(chan *subscription),
		unsCh:   make(chan *subscription),
		closeCh: make(chan bool, 1),
	}
	go b.receive()
	go b.loop()
	return b, nil
}

func (b *backbone) Close() error {
	b.closeCh <- true
	return b.client.Close()
}

func (b *backbone) Publish(ch string, payload []byte) error {
	ch = b.prefix + ch
	b.Debugf("Publishing to redis: %s (%d bytes)", ch, len(payload))
	err := b.client.Publish(ch, base64.StdEncoding.EncodeToString(payload)).Err()
	if err != nil {
		b.Errorf("Failed to publish to '%s': %s", ch, err.Error())
		return err
	}
	return nil
}

func (b *backbone) Subscribe(ch string) (pubsub.Subscription, error) {
	b.Debugf("Subscribing from redis: %s", b.prefix+ch)
	pubsub := &subscription{
		Sugar:    b,
		backbone: b,
		channel:  b.prefix + ch,
		msgCh:    make(chan string),
	}
	b.subCh <- pubsub
	return pubsub, nil
}

func (b *backbone) loop() {
	defer func() {
		for _, subs := range b.subs {
			for _, sub := range subs {
				sub.msgCh <- ""
			}
		}
		b.pubsub.Close()
		b.Debugf("Exiting redis backbone loop")
	}()
	for {
		select {
		case <-b.closeCh:
			return
		case msg := <-b.msgCh:
			subs := b.subs[msg.Channel]
			if subs != nil {
				b.Debugf("Dispatching message: %s => %d subscribers", msg.Channel, len(subs))
				for _, sub := range subs {
					sub.msgCh <- msg.Payload
				}
			} else {
				b.Warningf("Not subscribed: %s", msg.Channel)
			}
		case sub := <-b.subCh:
			b.Debugf("Adding subscriber: %s", sub.channel)
			subs := b.subs[sub.channel]
			if subs == nil {
				subs = []*subscription{sub}
				err := b.pubsub.Subscribe(sub.channel)
				if err != nil {
					b.Errorf("Failed to subscribe '%s': %s", sub.channel, err.Error())
				}
			} else {
				subs = append(subs, sub)
			}
			b.subs[sub.channel] = subs
		case sub := <-b.unsCh:
			b.Debugf("Removing subscriber: %s", sub.channel)
			subs := b.subs[sub.channel]
			if subs == nil {
				b.Warningf("No subscriber: %s", sub.channel)
			} else {
				at := -1
				for i, c := range subs {
					if c == sub {
						at = i
						break
					}
				}
				if at == -1 {
					b.Warningf("Not subscribed: %s", sub.channel)
				} else {
					subs = append(subs[:at], subs[at+1:]...)
					if len(subs) <= 0 {
						err := b.pubsub.Unsubscribe(sub.channel)
						if err != nil {
							b.Errorf("Failed to unsubscribe '%s': %s", sub.channel, err.Error())
						}
						delete(b.subs, sub.channel)
					} else {
						b.subs[sub.channel] = subs
					}
				}
			}
		}
	}
}

func (b *backbone) receive() {
	for {
		msg, err := b.pubsub.ReceiveMessage()
		if err != nil {
			b.Errorf("Failed to receive redis message: %s", err.Error())
			break
		} else {
			b.msgCh <- msg
		}
	}
}

type subscription struct {
	log.Sugar
	backbone *backbone
	channel  string
	msgCh    chan string
	closed   bool
}

func (p *subscription) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	p.Debugf("Unsubscribing redis: %s", p.channel)
	p.backbone.unsCh <- p
	p.msgCh <- ""
	return nil
}

func (p *subscription) Receive() ([]byte, error) {
	msg := <-p.msgCh
	if msg == "" {
		return nil, fmt.Errorf("Subscription closed")
	}
	data, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		p.Errorf("Failed to decode redis payload: %s", err.Error())
		return nil, err
	}
	p.Debugf("Received message: %s => %d bytes", p.channel, len(data))
	return data, nil
}
