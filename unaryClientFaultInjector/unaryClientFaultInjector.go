package unaryClientFaultInjector

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/randomizedcoder/grpcFaultInjection/internal/rand"
)

const (
	faultmodulusHeader = "faultmodulus"
	faultpercentHeader = "faultpercent"
	faultcodesHeader   = "faultcodes"
)

var (
	count   atomic.Uint64
	fault   atomic.Uint64
	success atomic.Uint64

	once        sync.Once
	configError atomic.Uint64

	logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
)

// unaryClientInterceptor allows a GRPC client to randomly inject metadata(headers) into
// the GRPC request.  The metadata headers themselves make a request to a similar intercetpor
// on the GRPC server side, which will randomly inject failures into the GRPC responses
// this is designed for testing, to allow the client to request failures from the GRPC server
// ultimately to test the client side error handling behavior
// https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryClientInterceptor
func UnaryClientFaultInjector(config UnaryClientInterceptorConfig, debugLevel int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		counter := count.Add(1)

		once.Do(func() {
			err := CheckConfig(config)
			if err != nil {
				configError.Add(1)
				log.Print("checkConfig(config) fails")
			}
		})

		if configError.Load() > 0 {
			c := configError.Add(1)
			return fmt.Errorf("config error:%d", c)
		}

		switch config.Client.Mode {
		case Modulus:
			if counter%uint64(config.Client.Value) == 0 {

				if debugLevel > 10 {
					logger.Printf("UnaryClientFaultInjector counter:%d", counter)
				}

				return faultInject(ctx, config, debugLevel, method, req, reply, cc, invoker, opts...)
			}
			return noFaultInject(ctx, debugLevel, method, req, reply, cc, invoker, opts...)

		case Percent:
			if config.Client.Value == 100 {
				return faultInject(ctx, config, debugLevel, method, req, reply, cc, invoker, opts...)
			}

			if rand.FastRandNInt() > config.Client.Value {
				return noFaultInject(ctx, debugLevel, method, req, reply, cc, invoker, opts...)
			}
		default:
			return fmt.Errorf("config error: must have modulus or percent")
		}

		return faultInject(ctx, config, debugLevel, method, req, reply, cc, invoker, opts...)
	}
}

func noFaultInject(ctx context.Context, debugLevel int, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	s := success.Add(1)
	f := fault.Load()

	if debugLevel > 10 {
		logger.Print(logNoFaultRequest(s, f))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func faultInject(ctx context.Context, config UnaryClientInterceptorConfig, debugLevel int,
	method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	f := fault.Add(1)
	s := success.Load()

	if debugLevel > 10 {
		logger.Print(logFaultRequest(s, f))
	}

	// https://grpc.io/docs/guides/metadata/
	// https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/client/main.go
	var md metadata.MD

	switch config.Server.Mode {
	case Modulus:
		md = metadata.Pairs(
			faultmodulusHeader, strconv.FormatInt(int64(config.Server.Value), 10),
		)
	case Percent:
		md = metadata.Pairs(
			faultpercentHeader, strconv.FormatInt(int64(config.Server.Value), 10),
		)
	}

	if len(config.Codes) > 0 {
		md.Append(faultcodesHeader, config.Codes)
	}

	if debugLevel > 10 {
		logger.Print("md:", md)
	}

	ctxMD := metadata.NewOutgoingContext(ctx, md)

	return invoker(ctxMD, method, req, reply, cc, opts...)
}
