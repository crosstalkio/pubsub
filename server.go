package pubsub

import (
	"fmt"
	"net"
	"time"

	"github.com/crosstalkio/log"
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
			msg := &Message{}
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

func (s *Server) handle(msg *Message) (*Message, error) {
	switch v := msg.Payload.(type) {
	case *Message_Control:
		ctl := v.Control
		switch v := ctl.Payload.(type) {
		case *Control_Request:
			req := v.Request
			rid := req.GetId()
			if rid == "" {
				err := fmt.Errorf("Missing request ID")
				s.Errorf(err.Error())
				return nil, err
			}
			s.Debugf("Handling request %s from %s", rid, s.remoteAddr.String())
			switch v := req.Payload.(type) {
			case *Request_Publish:
				pub := v.Publish
				return s.publish(rid, pub.GetChannel(), pub.GetPayload())
			case *Request_Subscribe:
				sub := v.Subscribe
				return s.subscribe(rid, sub.GetChannel())
			case *Request_Unsubscribe:
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

func (s *Server) subscribe(id, ch string) (*Message, error) {
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
			msg := &Data{}
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
			pmsg := &Message{
				Payload: &Message_Data{
					Data: &Data{
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

func (s *Server) unsubscribe(id, ch string) (*Message, error) {
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

func (s *Server) publish(id, ch string, payload isPublish_Payload) (*Message, error) {
	if id == "" {
		return s.error(id, 400, "Missing request ID"), nil
	}
	if ch == "" {
		return s.error(id, 400, "Missing channel to publish"), nil
	}
	if payload == nil {
		return s.error(id, 400, "Missing payload to publish"), nil
	}
	msg := &Data{
		NanoTime: time.Now().UnixNano(),
		From:     s.Identity,
	}
	switch v := payload.(type) {
	case *Publish_Text:
		s.Debugf("Publishing %d bytes of text data from '%s': %s", len(v.Text), msg.From, s.remoteAddr.String())
		msg.Payload = &Data_Text{Text: v.Text}
	case *Publish_Binary:
		s.Debugf("Publishing %d bytes of binary data from '%s': %s", len(v.Binary), msg.From, s.remoteAddr.String())
		msg.Payload = &Data_Binary{Binary: v.Binary}
	default:
		err := fmt.Errorf("Unexpected type of Publish: %v", v)
		s.Errorf(err.Error())
		return nil, err
	}
	if s.perm != nil {
		authz := false
		s.Debugf("Checking publish permission: %s vs %v", ch, s.perm.Write)
		for _, pattern := range s.perm.Write {
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
			msg := fmt.Sprintf("Unauthorized publish: %s", ch)
			s.Warningf(msg)
			return s.error(id, 401, msg), nil
		}
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed to marshal payload for backbone: %s", err.Error)
		return nil, err
	}
	err = s.backbone.Publish(ch, data)
	if err != nil {
		s.Errorf("Failed to publish to backbone: %s", err.Error())
		return s.error(id, 500, "%s", err.Error()), nil
	}
	return s.success(id), nil
}

func (s *Server) send(msg *Message) error {
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

func (s *Server) control(ctl *Control) *Message {
	return &Message{
		Payload: &Message_Control{
			Control: ctl,
		},
	}
}

func (s *Server) success(id string) *Message {
	return s.control(&Control{
		Payload: &Control_Response{
			Response: &Response{
				Id: id,
				Payload: &Response_Success{
					Success: &Success{},
				},
			},
		},
	})
}

func (s *Server) error(id string, code int32, msg string, args ...interface{}) *Message {
	return s.control(&Control{
		Payload: &Control_Response{
			Response: &Response{
				Id: id,
				Payload: &Response_Error{
					Error: &Error{
						Code:   code,
						Reason: fmt.Sprintf(msg, args...),
					},
				},
			},
		},
	})
}
