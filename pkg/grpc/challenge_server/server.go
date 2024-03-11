package challenge_server

import (
	"challenge/pkg/proto"
	"challenge/pkg/timer"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"sync"
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
	mu *sync.Mutex
}

func Register(gRPC *grpc.Server, shortener UrlShortener, timer *timer.Timer) {
	proto.RegisterChallengeServiceServer(gRPC, &server{shortener: shortener, timer: timer, mu: &sync.Mutex{}})
}

func (s *server) MakeShortLink(_ context.Context, in *proto.Link) (*proto.Link, error) {

	link, err := s.shortener.CreateShortLink(in.GetData())
	if err != nil {
		log.Printf("failed to get shortened link. err: %v\n", err)
		return nil, status.Error(codes.Internal, "Failed to get shortened link")
	}

	return &proto.Link{Data: link}, nil
}

func (s *server) StartTimer(timer *proto.Timer, stream proto.ChallengeService_StartTimerServer) error {

	// Preventing parallel calls to api. May lead to errors with simultaneous calls
	s.mu.Lock()
	ping, err := s.timer.Subscribe(timer.GetName(), int(timer.GetSeconds()), int(timer.GetFrequency()))
	if err != nil {
		log.Println("error when subscribing to timer: ", err)
		return status.Error(codes.Internal, "Couldn't start or subscribe to timer")
	}
	s.mu.Unlock()

	defer func() {
		s.timer.Unsubscribe(timer.GetName(), ping)
		log.Println("ending streaming grpc method")
	}()

	for {
		select {
		case <-stream.Context().Done():
			log.Println("connection was closed from client side")
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
				log.Printf("failed to send message to stream. err: %v\n", err)
				return status.Error(codes.Internal, "Failed to send streaming message")
			}
		}
	}
}

func (s *server) ReadMetadata(ctx context.Context, _ *proto.Placeholder) (*proto.Placeholder, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("failed to get metadata from context")
		return nil, status.Errorf(codes.DataLoss, "Failed to get metadata")
	}
	mds, ok := md[metadataKey]
	if !ok {
		log.Println("metadata not exists")
		return nil, status.Errorf(codes.NotFound, "String metadata not found")
	}
	if len(mds) == 0 {
		log.Println("metadata len is 0")
		return nil, status.Errorf(codes.NotFound, "String metadata found but empty")
	}

	// We consider first existing value of metadata our needed value
	return &proto.Placeholder{Data: mds[0]}, nil
}
