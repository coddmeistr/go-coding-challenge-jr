package main

import (
	"challenge/pkg/api/bilty"
	"challenge/pkg/api/timercheck"
	"challenge/pkg/config"
	"challenge/pkg/grpc/challenge_server"
	"challenge/pkg/timer"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const defaultConfigPath = "./configs/server.yaml"

func main() {
	// Init config
	cfg := config.MustLoadByPath(defaultConfigPath)

	// Init and inject all dependencies
	bil := bilty.NewBilty(cfg.BitlyOAuthToken, http.DefaultClient)
	timerChecker := timercheck.NewTimerCheck(http.DefaultClient)
	t := timer.NewTimer(*timerChecker)

	// Create gRPC server
	server := grpc.NewServer()
	challenge_server.Register(server, bil, t)

	// Start gRPC server
	go mustRun(server, cfg.Port)

	// Gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	log.Printf("starting gracefull shutdown. Signal: %v\n", sig)
	server.GracefulStop()

	log.Println("gracefully stopped")
}

func mustRun(server *grpc.Server, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	log.Printf("starting gRPC server on port :%d\n", port)
	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
