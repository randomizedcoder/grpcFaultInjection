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
	"google.golang.org/grpc/examples/features/proto/echo"

	"github.com/randomizedcoder/grpcFaultInjection/unaryClientFaultInjector"
)

var (
	loops = flag.Int("loops", 10, "loops")

	clientmode  = flag.String("clientmode", "Modulus", "clientmode 'modulus/mod/m' or 'percent/per/p'")
	clientvalue = flag.Int("clientvalue", 2, "clientvalue integers only, modulus 1-10000, percent 1-100")
	servermode  = flag.String("servermode", "Modulus", "servermode 'modulus/mod/m' or 'percent/per/p'")
	servervalue = flag.Int("servervalue", 2, "servervalue integers only, modulus 1-10000, percent 1-100")

	codes = flag.String("codes", "10,12,14", "GRPC status codes to return. comma seperated")

	addr   = flag.String("addr", "localhost:50052", "the address to connect to")
	policy = flag.String("policy", "grpc_client_policy.yaml", "filename of the grpc client policy.yaml")

	debugLevel = flag.Int("debugLevel", 11, "debugLevel. > 10 for output")
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
		Client: unaryClientFaultInjector.ModeValue{
			Mode:  unaryClientFaultInjector.StringToMode(*clientmode),
			Value: *clientvalue,
		},
		Server: unaryClientFaultInjector.ModeValue{
			Mode:  unaryClientFaultInjector.StringToMode(*servermode),
			Value: *servervalue,
		},
		Codes: *codes,
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

	var (
		success int
		fault   int
	)

	c := echo.NewEchoClient(conn)
	for i := 0; i < *loops; i++ {

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		reply, err := c.UnaryEcho(ctx,
			&echo.EchoRequest{
				Message: "Try and Success",
			},
		)
		if err != nil {
			log.Printf("i:%d UnaryEcho error: %v", i, err)
			fault++
			continue
		}
		log.Printf("i:%d UnaryEcho reply: %v", i, reply)
		success++
	}

	log.Printf("Complete.  success:%d fault:%d", success, fault)
}
