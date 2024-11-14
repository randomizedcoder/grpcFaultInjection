package unaryClientFaultInjector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/randomizedcoder/grpcFaultInjection/internal/validate"
)

// checkConfig is a simple configuration validator
// it is recommended to call this BEFORE instanciating the interceptor
// in the GRPC client
func CheckConfig(config UnaryClientInterceptorConfig) error {

	switch config.Client.Mode {
	case Modulus:
		if _, err := validate.ValidateModulus(int64(config.Client.Value)); err != nil {
			return fmt.Errorf("ValidateModulus config.Client.Value error: %w", err)
		}
	case Percent:
		if _, err := validate.ValidatePercent(int64(config.Client.Value)); err != nil {
			return fmt.Errorf("ValidatePercent config.Client.Value error: %w", err)
		}
	}

	switch config.Server.Mode {
	case Modulus:
		if _, err := validate.ValidateModulus(int64(config.Server.Value)); err != nil {
			return fmt.Errorf("ValidateModulus config.Server.Value error: %w", err)
		}
	case Percent:
		if _, err := validate.ValidatePercent(int64(config.Server.Value)); err != nil {
			return fmt.Errorf("ValidatePercent config.Server.Value error: %w", err)
		}
	}

	if len(config.Codes) > 0 {
		if err := validateCodes(config.Codes); err != nil {
			return fmt.Errorf("config.Codes error: %w", err)
		}
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
