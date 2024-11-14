package unaryServerFaultInjector

import (
	"context"
	"log"
	"os"
	"sync/atomic"

	_ "unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/randomizedcoder/grpcFaultInjection/internal/rand"
)

var (
	count   atomic.Uint64
	fault   atomic.Uint64
	success atomic.Uint64

	errMetadata = status.Errorf(codes.InvalidArgument, "error metadata")

	logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
)

// https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryServerInterceptor
func UnaryServerFaultInjector(debugLevel int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		counter := count.Add(1)

		// https://grpc.io/docs/guides/metadata/
		// https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/server/main.go
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMetadata
		}

		var (
			foundModulus bool
			faultModulus uint64
			errM         error
		)
		foundModulus, faultModulus, errM = readFaultModulus(&md, debugLevel)
		if errM != nil {
			return nil, errM
		}

		if foundModulus {
			if counter%faultModulus == 0 {
				return faultInject(counter, &md, debugLevel)
			}
			return noFaultInject(ctx, req, handler, debugLevel)
		}

		return faultPercentInject(ctx, req, handler, counter, &md, debugLevel)
	}
}

func noFaultInject(
	ctx context.Context, req any, handler grpc.UnaryHandler, debugLevel int) (any, error) {

	s := success.Add(1)
	f := fault.Load()

	if debugLevel > 11 {
		logger.Print(logNoFaultRequest(s, f))
	}

	return handler(ctx, req)
}

func faultPercentInject(
	ctx context.Context,
	req any,
	handler grpc.UnaryHandler,
	counter uint64,
	md *metadata.MD,
	debugLevel int) (any, error) {

	var (
		foundPercent bool
		faultPercent int
		errP         error
	)

	foundPercent, faultPercent, errP = readFaultPercent(md, debugLevel)
	if errP != nil {
		return nil, errP
	}

	if !foundPercent {
		return noFaultInject(ctx, req, handler, debugLevel)
	}

	if faultPercent == 100 {
		return faultInject(counter, md, debugLevel)
	}

	if rand.FastRandNInt() > faultPercent {
		return noFaultInject(ctx, req, handler, debugLevel)
	}

	return faultInject(counter, md, debugLevel)
}

func faultInject(
	counter uint64, md *metadata.MD, debugLevel int) (any, error) {

	f := fault.Add(1)
	s := success.Load()

	faultCodes, errC := readFaultCodes(md)
	if errC != nil {
		return nil, errC
	}

	var code codes.Code
	switch len(faultCodes) {
	case 0:
		code = rand.RandomFaultCode()
	case 1:
		code = faultCodes[0]
	default:
		code = rand.RandomSuppliedFaultCode(&faultCodes)
	}

	if debugLevel > 10 {
		logger.Print(logFaultRequest(s, f, code))
	}

	return nil, status.Errorf(
		code,
		"intercept fault code:%d counter:%d success:%d fault:%d",
		uint32(code), counter, s, f)
}
