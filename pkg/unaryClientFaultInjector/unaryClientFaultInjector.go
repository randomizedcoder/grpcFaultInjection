package unaryClientFaultInjector

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync/atomic"

	_ "unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// unsafe for FastRandN

// https://cs.opensource.google/go/go/+/master:src/runtime/stubs.go;l=151?q=FastRandN&ss=go%2Fgo
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/

//go:linkname FastRandN runtime.fastrandn
func FastRandN(n uint32) uint32

const (
	faultpercentHeader = "faultpercent"
	faultcodesHeader   = "faultcodes"
)

type UnaryClientInterceptorConfig struct {
	ClientFaultPercent int
	ServerFaultPercent int
	ServerFaultCodes   string
}

var (
	fault   atomic.Uint64
	success atomic.Uint64

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

		r := FastRandN(100)

		if r > uint32(config.ClientFaultPercent) {
			success.Add(1)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		f := fault.Add(1)
		s := success.Load()

		if debugLevel > 10 {
			logRequest(s, f)
		}

		// https://grpc.io/docs/guides/metadata/
		// https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/client/main.go
		md := metadata.Pairs(
			faultpercentHeader, strconv.FormatInt(int64(config.ServerFaultPercent), 10),
			faultcodesHeader, config.ServerFaultCodes,
		)
		ctxMD := metadata.NewOutgoingContext(ctx, md)

		return invoker(ctxMD, method, req, reply, cc, opts...)

	}
}

func logRequest(s uint64, f uint64) {
	if s > 0 {
		logger.Printf("request success:%d fault:%d ~= %.3f", s, f, float64(f)/float64(s))
	} else {
		logger.Printf("request success:%d fault:%d", s, f)
	}
}
