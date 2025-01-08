package grpc

import (
	"local/othello/gateway/grpc/gen"

	"google.golang.org/grpc"
)

type Service struct {
	gen.UnimplementedOthelloServer
	stream chan grpc.BidiStreamingServer[gen.Action, gen.Action]
}

var _ gen.OthelloServer = &Service{}

func NewService() *Service {
	return &Service{
		stream: make(chan grpc.BidiStreamingServer[gen.Action, gen.Action]),
	}
}

func (s *Service) Sync(stream grpc.BidiStreamingServer[gen.Action, gen.Action]) error {
	s.stream <- stream

	<-s.stream

	return nil
}
