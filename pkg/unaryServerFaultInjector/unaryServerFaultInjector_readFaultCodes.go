package unaryServerFaultInjector

import (
	"strconv"
	"strings"

	_ "unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"randomizedcoder/grpcFaultInjection/pkg/validate"
)

const (
	faultcodesHeader = "faultcodes"
)

// readFaultCodes returns a slice of codes.Code from the metadata "faultcodes"
// "faultcodes" can be a single code, or a comma seperated list of codes
// e.g. "faultcodes" = 14 (unavailable)
// e.g. "faultcodes" = 10,12,14
// valid codes: https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
func readFaultCodes(md *metadata.MD) (cs []codes.Code, err error) {

	fc, found := (*md)[faultcodesHeader]
	if !found {
		return cs, nil
	}

	parts := strings.Split(fc[0], ",")
	for i := 0; i < len(parts); i++ {
		c, err := strconv.ParseInt(parts[i], 0, 64)
		if err != nil {
			return cs, status.Error(codes.InvalidArgument, "faultcodes ParseInt error")
		}
		code, errV := validate.ValidateCode(c)
		if errV != nil {
			return cs, status.Error(codes.InvalidArgument, "faultcodes validate error")
		}
		cs = append(cs, codes.Code(code))
	}

	return cs, nil
}
