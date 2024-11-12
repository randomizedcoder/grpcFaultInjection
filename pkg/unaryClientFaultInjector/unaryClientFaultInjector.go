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

	"randomizedcoder/grpcFaultInjection/pkg/rand"
)

const (
	faultmodulusHeader = "faultmodulus"
	faultpercentHeader = "faultpercent"
	faultcodesHeader   = "faultcodes"
)

type UnaryClientInterceptorConfig struct {
	ClientFaultModulus int
	ClientFaultPercent int
	ServerFaultModulus int
	ServerFaultPercent int
	ServerFaultCodes   string
}

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

		if counter%uint64(config.ClientFaultModulus) == 0 {
			return faultInject(ctx, config, debugLevel, method, req, reply, cc, invoker, opts...)
		}

		if config.ClientFaultPercent == 100 {
			return faultInject(ctx, config, debugLevel, method, req, reply, cc, invoker, opts...)
		}

		if rand.FastRandNInt() > config.ClientFaultPercent {
			return noFaultInject(ctx, debugLevel, method, req, reply, cc, invoker, opts...)
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
	if config.ServerFaultModulus > 0 {
		md = metadata.Pairs(
			faultmodulusHeader, strconv.FormatInt(int64(config.ServerFaultModulus), 10),
		)
	} else {
		md = metadata.Pairs(
			faultpercentHeader, strconv.FormatInt(int64(config.ServerFaultPercent), 10),
		)
	}

	if config.ServerFaultCodes != "" {
		md.Append(faultcodesHeader, config.ServerFaultCodes)
	}

	ctxMD := metadata.NewOutgoingContext(ctx, md)

	return invoker(ctxMD, method, req, reply, cc, opts...)
}
