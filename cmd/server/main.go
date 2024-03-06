package main

import (
	"challenge/pkg/grpc/challenge_server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// Start gRPC challenge_server
	server := grpc.NewServer()
	challenge_server.Register(server)

	done := make(chan struct{})
	go mustRun(server)

	<-done
}

func mustRun(server *grpc.Server) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 5000))
	if err != nil {
		panic(err)
	}

	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
