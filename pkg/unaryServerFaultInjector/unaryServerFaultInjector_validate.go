package unaryServerFaultInjector

import "errors"

var (
	errInvalidFaultPercent = errors.New("invalid faultpercent")
	errInvalidCode         = errors.New("invalid code")
)

// validateFaultPercent ensure the percentage is between 1-100 inclusive
func validateFaultPercent(FaultureRate int64) (faultPercent int, err error) {
	if FaultureRate <= 0 || FaultureRate > 100 {
		return faultPercent, errInvalidFaultPercent
	}
	faultPercent = int(FaultureRate)
	return faultPercent, nil
}

// validateCode ensures the code is between 0-16 inclusive
func validateCode(c int64) (code uint32, err error) {
	if c < 0 || c > 16 {
		return code, errInvalidCode
	}
	return uint32(c), nil
}
