package challenge_server

import (
	"challenge/pkg/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UrlShortener interface {
	CreateShortLink(url string) (string, error)
}

const (
	metadataKey = "i-am-random-key"
)

type server struct {
	shortener UrlShortener
	proto.UnimplementedChallengeServiceServer
}

func Register(gRPC *grpc.Server, shortener UrlShortener) {
	proto.RegisterChallengeServiceServer(gRPC, &server{shortener: shortener})
}

func (s *server) MakeShortLink(ctx context.Context, in *proto.Link) (*proto.Link, error) {

	link, err := s.shortener.CreateShortLink(in.GetData())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get shortened link")
	}

	return &proto.Link{Data: link}, nil
}

func (s *server) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {

	return nil
}

func (s *server) ReadMetadata(ctx context.Context, in *proto.Placeholder) (*proto.Placeholder, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Failed to get metadata")
	}
	mds, ok := md[metadataKey]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "String metadata not found")
	}
	if len(mds) == 0 {
		return nil, status.Errorf(codes.NotFound, "String metadata found but empty")
	}

	// We consider first existing value of metadata our needed value
	return &proto.Placeholder{Data: mds[0]}, nil
}
