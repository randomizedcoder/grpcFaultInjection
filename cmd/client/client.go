package main

// This is a basic demontration of using the unaryClientFaultInjector

// Originally adpated from:
// https://github.com/grpc/grpc-go/blob/master/examples/features/retry/client/main.go

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/features/proto/echo"

	"randomizedcoder/grpcFaultInjection/pkg/unaryClientFaultInjector"
)

var (
	loops              = flag.Int("loops", 10, "loops")
	addr               = flag.String("addr", "localhost:50052", "the address to connect to")
	policy             = flag.String("policy", "grpc_client_policy.yaml", "filename of the grpc client policy.yaml")
	clientfaultmodulus = flag.Int("clientfaultmodulus", 2, "clientfaultmodulus integers only between 1-10,000")
	clientfaultpercent = flag.Int("clientfaultpercent", 50, "clientfaultpercent integers only between 0-100")
	faultmodulus       = flag.Int("faultmodulus", 2, "faultmodulus integers only between 1-10,000")
	faultpercent       = flag.Int("faultpercent", 50, "faultpercent integers only between 0-100")
	faultcodes         = flag.String("faultcodes", "4,8,14", "faultcodes header to insert. single code, or comma seperated")
	debugLevel         = flag.Int("debugLevel", 11, "debugLevel. > 10 for output")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// https://github.com/grpc/grpc/blob/master/doc/service_config.md to know more about service config
	// https://github.com/grpc/grpc-go/blob/11feb0a9afd8/examples/features/retry/client/main.go#L36
	// https://grpc.github.io/grpc/core/md_doc_statuscodes.html
	servicePolicyBytes, err := os.ReadFile(*policy)
	if err != nil {
		log.Fatal(err)
	}

	conf := unaryClientFaultInjector.UnaryClientInterceptorConfig{
		ClientFaultModulus: *clientfaultmodulus,
		ClientFaultPercent: *clientfaultpercent,
		ServerFaultModulus: *faultmodulus,
		ServerFaultPercent: *faultpercent,
		ServerFaultCodes:   *faultcodes,
	}

	if err := unaryClientFaultInjector.CheckConfig(conf); err != nil {
		log.Fatal(err)
	}

	// Set up a connection to the server with service config and create the channel.
	// However, the recommended approach is to fetch the retry configuration
	// (which is part of the service config) from the name resolver rather than
	// defining it on the client side.
	conn, err := grpc.NewClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(string(servicePolicyBytes)),
		grpc.WithUnaryInterceptor(
			unaryClientFaultInjector.UnaryClientFaultInjector(conf, *debugLevel),
		),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := pb.NewEchoClient(conn)
	for i := 0; i < *loops; i++ {

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		reply, err := c.UnaryEcho(ctx,
			&pb.EchoRequest{
				Message: "Try and Success",
			},
		)
		if err != nil {
			log.Printf("i:%d UnaryEcho error: %v", i, err)
			continue
		}
		log.Printf("i:%d UnaryEcho reply: %v", i, reply)
	}
}
