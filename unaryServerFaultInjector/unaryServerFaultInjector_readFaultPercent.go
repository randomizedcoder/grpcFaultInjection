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
	faultpercentHeader = "faultpercent"
)

// readFaultPercent reads the "Faultpercent" percentage, including validation
// percentage needs to be a integer between 0-100
// e.g. faultpercent = 10 ( 10% )
// e.g. faultpercent = 90 ( 90% )
func readFaultPercent(md *metadata.MD, debugLevel int) (found bool, faultPercent int, err error) {

	// metadata keys are always lower case
	// https://github.com/grpc/grpc-go/blob/v1.68.0/metadata/metadata.go#L207
	var faultPercentValue []string

	if faultPercentValue, found = (*md)[faultpercentHeader]; found {

		fp, err := strconv.ParseInt(faultPercentValue[0], 0, 64)
		if err != nil {
			return found, 0, status.Error(codes.InvalidArgument,
				"readfaultpercent ParseInt error")
		}

		var errV error
		faultPercent, errV = validate.ValidatePercent(fp)
		if errV != nil {
			return found, 0, status.Error(codes.InvalidArgument,
				"readfaultpercent validateFaultPercent error")
		}

		if debugLevel > 10 {
			logger.Printf("readFaultPercent faultPercent:%d", faultPercent)
		}

		return found, faultPercent, nil
	}

	// faultpercentHeader does not exist
	return found, 0, nil
}
