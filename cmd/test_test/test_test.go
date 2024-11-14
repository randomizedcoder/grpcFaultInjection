package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"

	"randomizedcoder/grpcFaultInjection/unaryClientFaultInjector"
	"randomizedcoder/grpcFaultInjection/unaryServerFaultInjector"

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

// go test -run TestComprehensive -v
func TestComprehensive(t *testing.T) {

	address := "localhost:50053"
	policy := "grpc_client_policy.yaml"

	debugLevel := 0
	debugLevelGRPCServer := 0
	debugLevelGRPCClient := 0
	//debugLevelGRPCClient := 11

	ctx := context.Background()

	//------------------------------------------------
	// Server setup

	lis, err := net.Listen("tcp", address)

	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	t.Log("listen on address", address)

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
		expectErr       bool
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
			name: "1/1 client, 1/1 server fault, loops 100, = 100%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      0,
			checkMaxSuccess: true,
			maxSuccess:      0,
			checkMinFault:   true,
			minFault:        100,
			checkMaxFault:   true,
			maxFault:        100,
		},
		{
			name: "1/1 client, 1/2 server fault, loops 100, = 50%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 2,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.5), // target is ~50%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.5), // target is ~50%
			checkMinFault:   true,
			minFault:        int(100 * 0.5), // target is ~50%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.5), // target is ~50%
		},
		{
			name: "1/2 client, 1/1 server fault, loops 100, = 50%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 2,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.5), // target is ~50%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.5), // target is ~50%
			checkMinFault:   true,
			minFault:        int(100 * 0.5), // target is ~50%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.5), // target is ~50%
		},
		{
			name: "1/2 client, 1/2 server fault, loops 100, = 25%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 2,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 2,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      50,
			checkMaxSuccess: false,
			maxSuccess:      int(100 * 0.5), // target is ~50%
			checkMinFault:   true,
			minFault:        50,
			checkMaxFault:   true,
			maxFault:        int(100 * 0.5), // target is ~50%
		},
		{
			// the odds a fault for this test are pretty low (1/10 * 1/10 = ~1%)
			// increase the interations, and don't make too many promises, cos
			// otherwise this test will be flake/annoying
			name: "1/10 client, 1/10 server fault, loops 100, = 10%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 10,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 10,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.1),
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.9),
			checkMinFault:   true,
			minFault:        int(100 * 0.1),
			checkMaxFault:   true,
			maxFault:        int(100 * 0.1), // target is ~1%
		},
		{
			name: "1/1 client, 1/10 server fault, loops 100, 10%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 10,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.9), // target is ~10%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.9), // target is ~10%
			checkMinFault:   true,
			minFault:        int(100 * 0.1), // target is ~10%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.1), // target is ~10%
		},
		{
			name: "1/1 client, 1/3 server fault, loops 100, 33.333...%",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 3,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.66), // target is ~66.66%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.67), // target is ~66.66%
			checkMinFault:   true,
			minFault:        int(100 * 0.33), // target is ~33.333%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.34), // target is ~33.333%
		},
		{
			name: "1/1 client, 1/6 server fault, loops 100, 16.666...",
			config: unaryClientFaultInjector.UnaryClientInterceptorConfig{
				Client: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 1,
				},
				Server: unaryClientFaultInjector.ModeValue{
					Mode:  unaryClientFaultInjector.Modulus,
					Value: 6,
				},
				Codes: "10",
			},
			expectErr:       false,
			loops:           100,
			checkMinSuccess: true,
			minSuccess:      int(100 * 0.83), // target is ~83.334%
			checkMaxSuccess: true,
			maxSuccess:      int(100 * 0.84), // target is ~83.334%
			checkMinFault:   true,
			minFault:        int(100 * 0.16), // target is ~16.666%
			checkMaxFault:   true,
			maxFault:        int(100 * 0.17), // target is ~16.666%
		},
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
			if debugLevel > 11 {
				t.Logf("tt.Name:%s success:%d fault:%d", tt.name, success, fault)
			}

			err := unaryClientFaultInjector.CheckConfig(tt.config)
			if err != nil {

				if !tt.expectErr {
					t.Fatalf("checkConfig(config) fails: %v", err)
				}
			}

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
					t.Logf("tt.Name:%s i:%d, success:%d fault:%d", tt.name, i, success, fault)
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
