package unaryServerFaultInjector

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	_ "unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// unsafe for FastRandN

// https://cs.opensource.google/go/go/+/master:src/runtime/stubs.go;l=151?q=FastRandN&ss=go%2Fgo
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/

//go:linkname FastRandN runtime.fastrandn
func FastRandN(n uint32) uint32

const (
	faultpercentHeader = "faultpercent"
	faultcodesHeader   = "faultcodes"

	// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	// We want a range between 1 and 16, so our maxCode is 15, cos we will +1
	maxCode = 15
)

var (
	fault   atomic.Uint64
	success atomic.Uint64

	errMetadata = status.Errorf(codes.InvalidArgument, "error metadata")

	logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
)

// https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#UnaryServerInterceptor
func UnaryServerFaultInjector(debugLevel int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var (
			code codes.Code
		)

		// https://grpc.io/docs/guides/metadata/
		// https://github.com/grpc/grpc-go/blob/master/examples/features/metadata/server/main.go
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMetadata
		}

		fp, errG := readFaultPercent(&md, debugLevel)
		if errG != nil {
			return nil, errG
		}

		if fp <= 0 {
			s := success.Add(1)
			f := fault.Load()
			if debugLevel > 11 {
				logRequest(s, f, code)
			}

			return handler(ctx, req)
		}

		r := FastRandN(100)

		if r > uint32(fp) {
			s := success.Add(1)
			f := fault.Load()
			if debugLevel > 11 {
				logRequest(s, f, code)
			}

			return handler(ctx, req)
		}

		f := fault.Add(1)
		s := success.Load()
		fcs, errC := readFaultCodes(&md)
		if errC != nil {
			return nil, errC
		}

		switch len(fcs) {
		case 0:
			code = anyRandomCode()
		case 1:
			code = fcs[0]
		default:
			code = returnAnySuppliedCode(&fcs)
		}

		if debugLevel > 10 {
			logRequest(s, f, code)
		}

		return nil, status.Errorf(
			code,
			"intercept fault code:%d r:%d success:%d fault:%d",
			uint32(code),
			r,
			s,
			f)

	}
}

func logRequest(s uint64, f uint64, code codes.Code) {
	if s > 0 {
		logger.Printf("request code:%s success:%d fault:%d ~= %.3f", code.String(), s, f, float64(f)/float64(s))
	} else {
		logger.Printf("request code:%s success:%d fault:%d", code.String(), s, f)
	}
}

// readFaultPercent reads the "Faultpercent" percentage, including validation
// percentage needs to be a integer between 0-100
// e.g. faultpercent = 10 ( 10% )
// e.g. faultpercent = 90 ( 90% )
func readFaultPercent(md *metadata.MD, debugLevel int) (int, error) {

	// metadata keys are always lower case
	// https://github.com/grpc/grpc-go/blob/v1.68.0/metadata/metadata.go#L207
	if t, ok := (*md)[faultpercentHeader]; ok {
		i, err := strconv.ParseInt(t[0], 0, 64)
		if err != nil {
			return 0, status.Error(codes.InvalidArgument, "readfaultpercent ParseInt error")
		}
		errV := validateFaultPercent(int(i))
		if errV != nil {
			return 0, status.Error(codes.InvalidArgument, "readfaultpercent validateFaultPercent error")
		}
		if debugLevel > 10 {
			logger.Printf("readFaultPercent:%d from metadata:\n", i)
		}

		return int(i), nil
	}
	return 0, nil
}

// validateFaultPercent ensure the percentage is between 0-100 inclusive
func validateFaultPercent(FaultureRate int) error {
	if FaultureRate < 0 || FaultureRate > 100 {
		return errors.New("invalid faultpercent")
	}
	return nil
}

// readFaultCodes returns a slice of codes.Code from the metadata "faultcodes"
// calls can include a single code, or a comma seperated list of codes
// e.g. Faultcodes = 14 (unavailable)
// e.g. Faultcodes = 10,12,14
// valid codes: https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
func readFaultCodes(md *metadata.MD) (cs []codes.Code, errR error) {

	if fc, ok := (*md)[faultcodesHeader]; ok {

		if !strings.Contains(fc[0], ",") {
			c, err := strconv.ParseInt(fc[0], 0, 64)
			if err != nil {
				return cs, status.Error(codes.InvalidArgument, "faultcodes ParseInt error")
			}
			if validateCodeUint32(uint32(c)) != nil {
				return cs, status.Error(codes.InvalidArgument, "faultcodes validate error")
			}
			cs = append(cs, codes.Code(uint32(c)))
			return cs, nil
		}

		parts := strings.Split(fc[0], ",")
		for i := 0; i < len(parts); i++ {
			i, err := strconv.ParseInt(parts[i], 0, 64)
			if err != nil {
				return cs, status.Error(codes.InvalidArgument, "faultcodes ParseInt error")
			}
			if validateCodeUint32(uint32(i)) != nil {
				return cs, status.Error(codes.InvalidArgument, "faultcodes validate error")
			}
			cs = append(cs, codes.Code(uint32(i)))
		}
		return cs, nil
	}
	return cs, nil
}

// validateCodeUint32 ensure the code is between 0-16 inclusive
// code can't be < 0 because it's a uint32
func validateCodeUint32(code uint32) error {
	if code > 16 {
		return errors.New("invalid code")
	}
	return nil
}

// anyRandomCode returns ANY random valid code ( 0-16 )
// if the request metadata does not contain "Faultcodes"
func anyRandomCode() (code codes.Code) {
	return codes.Code(FastRandN(maxCode) + 1)
}

// returnAnySuppliedCode randomly selects from the
// metadata supplied list of codes in "Faultcodes"
func returnAnySuppliedCode(cs *[]codes.Code) (code codes.Code) {
	return (*cs)[int(FastRandN(uint32(len(*cs))))]
}
