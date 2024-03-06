package challenge_server

import (
	"challenge/pkg/proto"
	"context"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedChallengeServiceServer
}

func Register(gRPC *grpc.Server) {
	proto.RegisterChallengeServiceServer(gRPC, &server{})
}

func (s *server) MakeShortLink(ctx context.Context, in *proto.Link) (*proto.Link, error) {

	return &proto.Link{}, nil
}

func (s *server) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {

	return nil
}

func (s *server) ReadMetadata(ctx context.Context, in *proto.Placeholder) (*proto.Placeholder, error) {

	return &proto.Placeholder{}, nil
}
