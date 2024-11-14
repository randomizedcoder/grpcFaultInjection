package unaryServerFaultInjector

import (
	"strconv"

	_ "unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/randomizedcoder/grpcFaultInjection/internal/validate"
)

const (
	faultmodulusHeader = "faultmodulus"
)

// readFaultModulus reads the "faultmodulus", including validation
// percentage needs to be a integer between 1-10000
// e.g. faultmodulus = 1 ( 100% )
// e.g. faultmodulus = 2 ( 50% )
// e.g. faultmodulus = 10 ( 10% ) Every 10th is a fault
// e.g. faultmodulus = 100 Every 100th is a fault
// e.g. faultmodulus = 1000 Every 1000th is a fault
func readFaultModulus(md *metadata.MD, debugLevel int) (found bool, faultModulus uint64, err error) {

	// metadata keys are always lower case
	// https://github.com/grpc/grpc-go/blob/v1.68.0/metadata/metadata.go#L207
	var faultModulusValue []string

	if faultModulusValue, found = (*md)[faultmodulusHeader]; found {

		fm, err := strconv.ParseInt(faultModulusValue[0], 0, 64)
		if err != nil {
			return found, 0, status.Error(codes.InvalidArgument,
				"readFaultModulus ParseInt error")
		}

		var errV error
		faultModulus, errV = validate.ValidateModulus(fm)
		if errV != nil {
			return found, 0, status.Error(codes.InvalidArgument,
				"readfaultpercent validateFaultPercent error")
		}

		if debugLevel > 10 {
			logger.Printf("readFaultModulus faultModulus:%d", faultModulus)
		}

		return found, faultModulus, nil
	}

	// faultmodulusHeader does not exist
	return found, 0, nil
}
