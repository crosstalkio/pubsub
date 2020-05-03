package pubsub

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/crosstalkio/log"
	"github.com/crosstalkio/pubsub/api"
	"github.com/mb0/glob"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	log.Sugar
	Identity string
	backbone Backbone
	perm     *Permission
}

func NewServer(logger log.Logger, backbone Backbone) *Server {
	server := &Server{
		Sugar:    log.NewSugar(logger),
		backbone: backbone,
	}
	return server
}

func (s *Server) Authorize(perm *Permission) {
	s.perm = perm
}

func (s *Server) Publish(ctx context.Context, req *api.PublishRequest) (*api.PublishResponse, error) {
	var addr net.Addr
	if peer, ok := peer.FromContext(ctx); ok {
		addr = peer.Addr
	}
	if req.Channel == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Missing 'channel' to publish")
	}
	if req.Type == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Missing 'data' to publish")
	}
	if s.perm != nil {
		authz := false
		s.Debugf("Checking publish permission: %s vs %v", req.Channel, s.perm.Write)
		for _, pattern := range s.perm.Write {
			match, err := glob.Match(pattern, req.Channel)
			if err != nil {
				s.Errorf("Failed to match '%s' vs '%s': %s", pattern, req.Channel, err.Error)
				return nil, err
			}
			if match {
				authz = true
				break
			}
		}
		if !authz {
			msg := fmt.Sprintf("Unauthorized publish: %s", req.Channel)
			s.Warningf(msg)
			return nil, status.Errorf(codes.PermissionDenied, msg)
		}
	}
	msg := &api.Data{
		NanoTime: time.Now().UnixNano(),
		From:     s.Identity,
	}
	switch v := req.Type.(type) {
	case *api.PublishRequest_Text:
		s.Debugf("Publishing %d bytes of text data from '%s': %s", len(v.Text), msg.From, addr.String())
		msg.Type = &api.Data_Text{Text: v.Text}
	case *api.PublishRequest_Binary:
		s.Debugf("Publishing %d bytes of binary data from '%s': %s", len(v.Binary), msg.From, addr.String())
		msg.Type = &api.Data_Binary{Binary: v.Binary}
	default:
		err := fmt.Errorf("Unexpected type of publish: %v", v)
		s.Errorf(err.Error())
		return nil, err
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		s.Errorf("Failed to marshal payload for backbone: %s", err.Error)
		return nil, err
	}
	err = s.backbone.Publish(req.Channel, data)
	if err != nil {
		s.Errorf("Failed to publish to backbone: %s", err.Error())
		return nil, err
	}
	return &api.PublishResponse{}, nil
}

func (s *Server) Subscribe(req *api.SubscribeRequest, stream api.PubSub_SubscribeServer) error {
	var addr net.Addr
	if peer, ok := peer.FromContext(stream.Context()); ok {
		addr = peer.Addr
	}
	if req.Channel == "" {
		return status.Errorf(codes.InvalidArgument, "Missing 'channel' to subscribe")
	}
	if s.perm != nil {
		authz := false
		s.Debugf("Checking subscribe permission: %s vs %v", req.Channel, s.perm.Read)
		for _, pattern := range s.perm.Read {
			match, err := glob.Match(pattern, req.Channel)
			if err != nil {
				s.Errorf("Failed to match '%s' vs '%s': %s", pattern, req.Channel, err.Error)
				return err
			}
			if match {
				authz = true
				break
			}
		}
		if !authz {
			msg := fmt.Sprintf("Unauthorized subscribe: %s", req.Channel)
			s.Warningf(msg)
			return status.Errorf(codes.PermissionDenied, msg)
		}
	}
	sub, err := s.backbone.Subscribe(req.Channel)
	if err != nil {
		s.Errorf("Failed to subscribe backbone: %s", err.Error())
		return err
	}
	defer func() {
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
			return err
		}
		from := msg.GetFrom()
		if s.Identity != "" && from == s.Identity {
			s.Debugf("Skipping self-published message: %s", from)
			continue
		}
		res := &api.SubscribeResponse{
			NanoTime: msg.GetNanoTime(),
			Channel:  req.Channel,
			From:     from,
		}
		switch v := msg.Type.(type) {
		case *api.Data_Text:
			res.Type = &api.SubscribeResponse_Text{Text: v.Text}
		case *api.Data_Binary:
			res.Type = &api.SubscribeResponse_Binary{Binary: v.Binary}
		default:
			err = fmt.Errorf("Unknown type of data: %v", v)
			s.Errorf("%s", err.Error())
			return err
		}
		s.Debugf("Dispatching received message to %s => %s", addr.String(), req.Channel)
		err = stream.Send(res)
		if err != nil {
			s.Errorf("Failed to dispatch message: %s", err.Error())
			return err
		}
	}
	return nil
}
