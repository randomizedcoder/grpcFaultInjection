package main

// This is a basic demontration of using the unaryServerFaultInjector

// Originally adpated from:
// https://github.com/grpc/grpc-go/blob/master/examples/features/retry/server/main.go

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"

	"github.com/randomizedcoder/grpcFaultInjection/unaryServerFaultInjector"

	"google.golang.org/grpc/examples/features/proto/echo"
)

var (
	success atomic.Uint64
)

type echoServer struct {
	echo.UnimplementedEchoServer
}

func newEchoServer() (s echoServer) {
	s = *new(echoServer)
	return s
}

func (s echoServer) UnaryEcho(_ context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {

	log.Println("request succeeded count:", success.Add(1))

	return &echo.EchoResponse{Message: req.Message}, nil
}

func main() {

	port := flag.Int("port", 50052, "port number")
	debugLevel := flag.Int("debugLevel", 11, "debugLevel.  > 10 for output")

	flag.Parse()

	address := fmt.Sprintf(":%v", *port)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("listen on address", address)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			unaryServerFaultInjector.UnaryServerFaultInjector(*debugLevel),
		),
	)

	srv := newEchoServer()

	echo.RegisterEchoServer(s, srv)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
