package challenge_server

import (
	"challenge/pkg/proto"
	"challenge/pkg/timer"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type UrlShortener interface {
	CreateShortLink(url string) (string, error)
}

const (
	metadataKey = "i-am-random-key"
)

type server struct {
	timer     *timer.Timer
	shortener UrlShortener
	proto.UnimplementedChallengeServiceServer
}

func Register(gRPC *grpc.Server, shortener UrlShortener, timer *timer.Timer) {
	proto.RegisterChallengeServiceServer(gRPC, &server{shortener: shortener, timer: timer})
}

func (s *server) MakeShortLink(ctx context.Context, in *proto.Link) (*proto.Link, error) {

	link, err := s.shortener.CreateShortLink(in.GetData())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get shortened link")
	}

	return &proto.Link{Data: link}, nil
}

func (s *server) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {

	ping, cancel, err := s.timer.StartOrSubscribe(timer.GetName(), int(timer.GetSeconds()), int(timer.GetFrequency()))
	if err != nil {
		log.Println("error when subscribing to timer: ", err)
		return status.Error(codes.Internal, "Couldn't start or subscribe to timer")
	}
	defer func() {
		cancel()
		fmt.Println("Ending streaming grpc method")
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case info, ok := <-ping:
			if !ok {
				log.Println("ping channel was closed")
				return nil
			}

			err = stream.Send(&proto.Timer{
				Name:      info.TimerName,
				Seconds:   int64(info.SecondsLeft),
				Frequency: timer.Frequency,
			})
			if err != nil {
				return status.Error(codes.Internal, "Failed to send streaming message")
			}
		}
	}
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
