package main

import (
	"challenge/pkg/config"
	"challenge/pkg/grpc/challenge_server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg := config.MustLoadByPath("./configs/server.yaml")

	// Start gRPC challenge_server
	server := grpc.NewServer()
	challenge_server.Register(server)

	done := make(chan struct{})
	go mustRun(server, cfg.Port)

	<-done
}

func mustRun(server *grpc.Server, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
