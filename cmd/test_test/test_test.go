package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"

	"randomizedcoder/grpcFaultInjection/pkg/unaryClientFaultInjector"
	"randomizedcoder/grpcFaultInjection/pkg/unaryServerFaultInjector"

	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/features/proto/echo"
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func newEchoServer() (s echoServer) {
	s = *new(echoServer)
	return s
}

func (s echoServer) UnaryEcho(_ context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {

	success.Add(1)
	//s := success.Add(1)
	//log.Println("request succeeded count:", s)

	return &pb.EchoResponse{Message: req.Message}, nil
}

var (
	success atomic.Uint64
)

func TestTest(t *testing.T) {

	address := "localhost:50052"
	policy := "grpc_client_policy.yaml"
	port := 50052
	debugLevel := 0
	debugLevelGRPCServer := 0
	debugLevelGRPCClient := 0

	ctx := context.Background()

	//------------------------------------------------
	// Server setup

	serverAddress := fmt.Sprintf(":%v", port)

	lis, err := net.Listen("tcp", serverAddress)

	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	t.Log("listen on address", serverAddress)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			unaryServerFaultInjector.UnaryServerFaultInjector(debugLevelGRPCServer),
		),
	)

	srv := newEchoServer()

	pb.RegisterEchoServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	//------------------------------------------------
	// Client setup

	// https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
	// https://github.com/grpc/grpc-go/blob/11feb0a9afd8/examples/features/retry/client/main.go#L36
	// https://grpc.github.io/grpc/core/md_doc_statuscodes.html
	servicePolicyBytes, err := os.ReadFile(policy)
	if err != nil {
		t.Fatal(err)
	}

	//------------------------------------------------
	// Tests setup

	type myTest struct {
		name            string
		config          unaryClientFaultInjector.UnaryClientInterceptorConfig
		loops           int
		checkMinSuccess bool
		minSuccess      int
		checkMaxSuccess bool
		maxSuccess      int
		checkMinFault   bool
		minFault        int
		checkMaxFault   bool
		maxFault        int
	}

	// These table tests loosely follow the "config matrix" section of the main readme.md
	//
	// Please note
	// We don't want these tests to fail, unless there's really a problem
	// This is tricky when we're talking about probabilities
	// Therefore, the general strategy is to just make sure we aren't too far off
	// Considering increase the 'loops' to increase the change of converging close to the target
	// We're prefer a little slop in the numbers, verse failing/flakey tests
	tests := []myTest{
		{
			name: "100 client, 0 server fault, loops 100",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: 0,
			},
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      100,
			checkMaxSuccess: true,
			maxSuccess:      100,
			checkMinFault:   true,
			minFault:        0,
			checkMaxFault:   true,
			maxFault:        0,
		},
		{
			name: "0 client, 100 server fault, loops 100",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 0,
				ServerFaultPercent: 100,
			},
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      100,
			checkMaxSuccess: true,
			maxSuccess:      100,
			checkMinFault:   true,
			minFault:        0,
			checkMaxFault:   true,
			maxFault:        0,
		},
		{
			// the odds a fault for this test are pretty low (10% * 10% = ~1%)
			// increase the interations, and don't make too many promises, cos
			// otherwise this test will be flake/annoying
			name: "10 client, 10 server fault, loops 500",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 10,
				ServerFaultPercent: 10,
			},
			loops:           500,
			checkMinSuccess: true,
			minSuccess:      1,
			checkMaxSuccess: false,
			maxSuccess:      99,
			checkMinFault:   true,
			minFault:        1,
			checkMaxFault:   true,
			maxFault:        int(500 * 0.2), // target is ~1%
		},
		{
			// the odds a fault for this test are pretty low (50% * 50% = ~25%)
			// increase the interations, and don't make too many promises, cos
			// otherwise this test will be flake/annoying
			name: "50 client, 50 server fault, loops 100",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 50,
				ServerFaultPercent: 50,
			},
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      15,
			checkMaxSuccess: false,
			maxSuccess:      int(100 * 0.85), // target is ~75%
			checkMinFault:   true,
			minFault:        15,
			checkMaxFault:   true,
			maxFault:        int(100 * 0.4), // target is ~25%
		},
		{
			name: "100 client, 50 server fault, loops 100",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: 50,
			},
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.3), // target is ~50%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.6), // target is ~50%
			checkMinFault:   true,
			minFault:        int(100 * 0.3), // target is ~50%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.6), // target is ~50%
		},
		{
			name: "100 client, 10 server fault, loops 100",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: 50,
			},
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.05), // target is ~10%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.95), // target is ~10%
			checkMinFault:   true,
			minFault:        int(100 * 0.05), // target is ~10%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.15), // target is ~10%
		},
		// {"Valid, mid percent", 50, false},
		// {"Valid, high percent", 100, false},
		// {"Invalid negative percent", -10, true},
		// {"Invalid over 100 percent", 110, true},
	}

	//------------------------------------------------
	// Run tests

	t.Log("run tests")

	for _, tt := range tests {

		t.Logf("tt.Name:%s", tt.name)

		t.Run(tt.name, func(t *testing.T) {

			var (
				success int
				fault   int
			)

			conn, err := grpc.NewClient(
				address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultServiceConfig(string(servicePolicyBytes)),
				grpc.WithUnaryInterceptor(
					unaryClientFaultInjector.UnaryClientFaultInjector(
						tt.config,
						debugLevelGRPCClient,
					),
				),
			)
			if err != nil {
				t.Fatalf("did not connect: %v", err)
			}
			defer func() {
				if e := conn.Close(); e != nil {
					t.Logf("failed to close connection: %s", e)
				}
			}()

			c := pb.NewEchoClient(conn)
			for i := 0; i < tt.loops; i++ {

				if debugLevel > 10 {
					t.Logf("tt.Name:%s i:%d", tt.name, i)
				}

				ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()

				reply, err := c.UnaryEcho(ctx,
					&pb.EchoRequest{
						Message: "Try and Success",
					},
				)
				if err != nil {
					fault++
					if debugLevel > 10 {
						t.Logf("i:%d success:%d fault:%d UnaryEcho error: %v", success, fault, i, err)
					}
					continue
				}
				success++
				if debugLevel > 10 {
					t.Logf("i:%d success:%d fault:%d UnaryEcho reply: %v", success, fault, i, reply)
				}
			}
			if tt.checkMinSuccess {
				if success < tt.minSuccess {
					t.Errorf("tt.Name:%s success:%d < tt.minSuccess:%d = ERROR", tt.name, success, tt.minSuccess)
				} else {
					t.Logf("tt.Name:%s success:%d minSuccess:%d = good", tt.name, success, tt.minSuccess)
				}
			}

			if tt.checkMaxSuccess {
				if success > tt.maxSuccess {
					t.Errorf("tt.Name:%s success:%d > tt.maxSuccess:%d = ERROR", tt.name, success, tt.maxSuccess)
				} else {
					t.Logf("tt.Name:%s success:%d maxSuccess:%d = good", tt.name, success, tt.maxSuccess)
				}
			}

			if tt.checkMinFault {
				if fault < tt.minFault {
					t.Errorf("tt.Name:%s fault:%d < tt.minFault:%d = ERROR", tt.name, fault, tt.minFault)
				} else {
					t.Logf("tt.Name:%s fault:%d minFault:%d = good", tt.name, fault, tt.minFault)
				}
			}

			if tt.checkMaxFault {
				if fault > tt.maxFault {
					t.Errorf("tt.Name:%s fault:%d > tt.maxFault:%d = ERROR", tt.name, fault, tt.maxFault)
				} else {
					t.Logf("tt.Name:%s fault:%d minFault:%d = good", tt.name, fault, tt.minFault)
				}
			}
		})
	}
}
