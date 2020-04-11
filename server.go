package pubsub

import (
	"fmt"
	"net"
	"time"

	"github.com/crosstalkio/log"
	api "github.com/crosstalkio/pubsub/api/pubsub"
	"github.com/gorilla/websocket"
	"github.com/mb0/glob"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	log.Sugar
	Identity   string
	backbone   Backbone
	closed     bool
	remoteAddr net.Addr
	conn       *websocket.Conn
	perm       *Permission
	writer     *serverWriter
	subs       map[string]Subscription
}

func NewServer(logger log.Logger, backbone Backbone, conn *websocket.Conn) *Server {
	return &Server{
		Sugar:      log.NewSugar(logger),
		backbone:   backbone,
		remoteAddr: conn.RemoteAddr(),
		conn:       conn,
		subs:       make(map[string]Subscription),
		writer:     newServerWriter(logger, conn),
	}
}

func (s *Server) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	s.Infof("Closing connection: %s", s.remoteAddr.String())
	err := s.conn.Close()
	if err != nil {
		s.Errorf("Failed to close websocket: %s", err.Error())
	}
	s.conn = nil
	for _, sub := range s.subs {
		suberr := sub.Close()
		if suberr != nil {
			s.Errorf("Failed to close subscription: %s", suberr.Error())
			err = suberr
		}
	}
	// we must close writer in the end of close() or there might leak go routine(s)
	s.writer.exit()
	return err
}

func (s *Server) Authorize(perm *Permission) {
	s.perm = perm
}

func (s *Server) Loop() {
	defer s.Debugf("Exiting reader: %s", s.remoteAddr.String())
	go s.writer.loop()
	for {
		mt, p, err := s.conn.ReadMessage()
		if err != nil {
			if !s.closed {
				s.Errorf("Failed to read message: %s", err.Error())
			}
			return
		}
		switch mt {
		case websocket.BinaryMessage:
			s.Debugf("Read %d bytes binary message: %s", len(p), s.remoteAddr.String())
			msg := &api.Message{}
			err := proto.Unmarshal(p, msg)
			if err != nil {
				s.Errorf("Failed to parse proto message: %s", err.Error())
				return
			}
			msg, err = s.handle(msg)
			if err != nil {
				return
			}
			if msg != nil {
				err = s.send(msg)
				if err != nil {
					return
				}
			}
		default:
			s.Errorf("Read %d bytes unexpected type (%d) message: %s", len(p), mt, s.remoteAddr.String())
			return
		}
	}
}

func (s *Server) PublishText(as, ch, text string) error {
	if ch == "" {
		err := fmt.Errorf("Missing channel to publish")
		s.Errorf("%s", err.Error())
		return err
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     as,
		Payload:  &api.Data_Text{Text: text},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed to marshal payload for backbone: %s", err.Error)
		return err
	}
	err = s.backbone.Publish(ch, data)
	if err != nil {
		s.Errorf("Failed to publish to backbone: %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) PublishBinary(as, ch string, bin []byte) error {
	if ch == "" {
		err := fmt.Errorf("Missing channel to publish")
		s.Errorf("%s", err.Error())
		return err
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     as,
		Payload:  &api.Data_Binary{Binary: bin},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed to marshal payload for backbone: %s", err.Error)
		return err
	}
	err = s.backbone.Publish(ch, data)
	if err != nil {
		s.Errorf("Failed to publish to backbone: %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) handle(msg *api.Message) (*api.Message, error) {
	switch v := msg.Payload.(type) {
	case *api.Message_Control:
		ctl := v.Control
		switch v := ctl.Payload.(type) {
		case *api.Control_Request:
			req := v.Request
			rid := req.GetId()
			if rid == "" {
				err := fmt.Errorf("Missing request ID")
				s.Errorf(err.Error())
				return nil, err
			}
			s.Debugf("Handling request %s from %s", rid, s.remoteAddr.String())
			switch v := req.Payload.(type) {
			case *api.Request_Publish:
				pub := v.Publish
				return s.publish(rid, pub)
			case *api.Request_Subscribe:
				sub := v.Subscribe
				return s.subscribe(rid, sub.GetChannel())
			case *api.Request_Unsubscribe:
				uns := v.Unsubscribe
				return s.unsubscribe(rid, uns.GetChannel())
			default:
				err := fmt.Errorf("Unexpected type of Request: %v", ctl)
				s.Errorf(err.Error())
				return nil, err
			}
		default:
			err := fmt.Errorf("Unexpected type of Control: %v", ctl)
			s.Errorf(err.Error())
			return nil, err
		}
	default:
		err := fmt.Errorf("Unexpected type of Message: %v", msg)
		s.Errorf(err.Error())
		return nil, err
	}
}

func (s *Server) subscribe(id, ch string) (*api.Message, error) {
	if id == "" {
		return s.error(id, 400, "Missing request ID"), nil
	}
	if ch == "" {
		return s.error(id, 400, "Missing channel to subscribe"), nil
	}
	if s.perm != nil {
		authz := false
		s.Debugf("Checking subscribe permission: %s vs %v", ch, s.perm.Read)
		for _, pattern := range s.perm.Read {
			match, err := glob.Match(pattern, ch)
			if err != nil {
				s.Errorf("Failed to match '%s' vs '%s': %s", pattern, ch, err.Error)
				return nil, err
			}
			if match {
				authz = true
				break
			}
		}
		if !authz {
			msg := fmt.Sprintf("Unauthorized subscribe: %s", ch)
			s.Warningf(msg)
			return s.error(id, 401, msg), nil
		}
	}
	if s.subs[ch] != nil {
		return s.error(id, 409, "Already subscribed: %s", ch), nil
	}
	sub, err := s.backbone.Subscribe(ch)
	if err != nil {
		s.Errorf("Failed to subscribe backbone: %s", err.Error())
		return s.error(id, 500, "%s", err.Error()), nil
	}
	s.subs[ch] = sub
	go func() {
		defer func() {
			delete(s.subs, ch)
			sub.Close()
		}()
		for {
			data, err := sub.Receive()
			if err != nil {
				break
			}
			msg := &api.Data{}
			err = proto.Unmarshal(data, msg)
			if err != nil {
				s.Errorf("Failed to unmarshal payload from backbone: %s", err.Error())
				continue
			}
			from := msg.GetFrom()
			if s.Identity != "" && from == s.Identity {
				s.Debugf("Skipping self-published message: %s", from)
				continue
			}
			pmsg := &api.Message{
				Payload: &api.Message_Data{
					Data: &api.Data{
						NanoTime: msg.GetNanoTime(),
						Channel:  ch,
						From:     from,
						Payload:  msg.GetPayload(),
					},
				},
			}
			s.Debugf("Dispatching received message to %s => %s", s.remoteAddr.String(), ch)
			err = s.send(pmsg)
			if err != nil {
				s.Errorf("Failed to dispatch message: %s", err.Error())
				continue
			}
		}
	}()
	return s.success(id), nil
}

func (s *Server) unsubscribe(id, ch string) (*api.Message, error) {
	if id == "" {
		return s.error(id, 400, "Missing request ID"), nil
	}
	if ch == "" {
		return s.error(id, 400, "Missing channel to unsubscribe"), nil
	}
	ps := s.subs[ch]
	if ps == nil {
		return s.error(id, 404, "Not subscribed: %s", ch), nil
	} else {
		delete(s.subs, ch)
		err := ps.Close()
		if err != nil {
			s.Errorf("Failed to close subscription: %s", err.Error())
			return s.error(id, 500, "%s", err.Error()), nil
		}
		return s.success(id), nil
	}
}

func (s *Server) publish(id string, pub *api.Publish) (*api.Message, error) {
	if id == "" {
		return s.error(id, 400, "Missing request ID"), nil
	}
	if pub.Channel == "" {
		return s.error(id, 400, "Missing channel to publish"), nil
	}
	if pub.Payload == nil {
		return s.error(id, 400, "Missing payload to publish"), nil
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     s.Identity,
	}
	switch v := pub.GetPayload().(type) {
	case *api.Publish_Text:
		s.Debugf("Publishing %d bytes of text data from '%s': %s", len(v.Text), msg.From, s.remoteAddr.String())
		msg.Payload = &api.Data_Text{Text: v.Text}
	case *api.Publish_Binary:
		s.Debugf("Publishing %d bytes of binary data from '%s': %s", len(v.Binary), msg.From, s.remoteAddr.String())
		msg.Payload = &api.Data_Binary{Binary: v.Binary}
	default:
		err := fmt.Errorf("Unexpected type of Payload: %v", v)
		s.Errorf(err.Error())
		return nil, err
	}
	if s.perm != nil {
		authz := false
		s.Debugf("Checking publish permission: %s vs %v", pub.Channel, s.perm.Write)
		for _, pattern := range s.perm.Write {
			match, err := glob.Match(pattern, pub.Channel)
			if err != nil {
				s.Errorf("Failed to match '%s' vs '%s': %s", pattern, pub.Channel, err.Error)
				return nil, err
			}
			if match {
				authz = true
				break
			}
		}
		if !authz {
			msg := fmt.Sprintf("Unauthorized publish: %s", pub.Channel)
			s.Warningf(msg)
			return s.error(id, 401, msg), nil
		}
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed to marshal payload for backbone: %s", err.Error)
		return nil, err
	}
	err = s.backbone.Publish(pub.Channel, data)
	if err != nil {
		s.Errorf("Failed to publish to backbone: %s", err.Error())
		return s.error(id, 500, "%s", err.Error()), nil
	}
	return s.success(id), nil
}

func (s *Server) send(msg *api.Message) error {
	ctl := msg.GetControl()
	if ctl != nil && ctl.GetResponse() != nil {
		s.Debugf("Replying request %s from %s", ctl.GetResponse().GetId(), s.remoteAddr.String())
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed marshal message: %s", err.Error())
		return err
	}
	s.writer.write(websocket.BinaryMessage, data)
	return nil
}

func (s *Server) control(ctl *api.Control) *api.Message {
	return &api.Message{
		Payload: &api.Message_Control{
			Control: ctl,
		},
	}
}

func (s *Server) success(id string) *api.Message {
	return s.control(&api.Control{
		Payload: &api.Control_Response{
			Response: &api.Response{
				Id: id,
				Payload: &api.Response_Success{
					Success: &api.Success{},
				},
			},
		},
	})
}

func (s *Server) error(id string, code int32, msg string, args ...interface{}) *api.Message {
	return s.control(&api.Control{
		Payload: &api.Control_Response{
			Response: &api.Response{
				Id: id,
				Payload: &api.Response_Error{
					Error: &api.Error{
						Code:   code,
						Reason: fmt.Sprintf(msg, args...),
					},
				},
			},
		},
	})
}
