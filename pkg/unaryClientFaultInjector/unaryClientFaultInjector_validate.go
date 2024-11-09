package unaryClientFaultInjector

import (
	"fmt"
	"strconv"
	"strings"

	"randomizedcoder/grpcFaultInjection/pkg/validate"
)

// checkConfig is a simple configuration validator
// it is recommended to call this BEFORE instanciating the interceptor
// in the GRPC client
func CheckConfig(config UnaryClientInterceptorConfig) error {

	if _, err := validate.ValidatePercent(int64(config.ClientFaultPercent)); err != nil {
		return fmt.Errorf("config.ClientFaultPercent error: %w", err)
	}
	if _, err := validate.ValidatePercent(int64(config.ServerFaultPercent)); err != nil {
		return fmt.Errorf("config.ServerFaultPercent error: %w", err)
	}
	if err := validateCodes(config.ServerFaultCodes); err != nil {
		return fmt.Errorf("config.ServerFaultCodes error: %w", err)
	}

	return nil
}

// validateCodes converts a code to an int64 and then validates the code
// when multiple codes, comma seperated, are inputted ths function
// iterates over them, and validates each
func validateCodes(codes string) error {

	parts := strings.Split(codes, ",")
	for i := 0; i < len(parts); i++ {
		c, err := strconv.ParseInt(parts[i], 0, 64)
		if err != nil {
			return err
		}
		if _, err := validate.ValidateCode(c); err != nil {
			return err
		}
	}
	return nil
}
