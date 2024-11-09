package validate

import "errors"

var (
	errInvalidPercent = errors.New("invalid percent")
	errInvalidCode    = errors.New("invalid code")
)

// ValidatePercent ensure the percentage is between 1-100 inclusive
func ValidatePercent(percent int64) (percentInt int, err error) {
	if percent <= 0 || percent > 100 {
		return percentInt, errInvalidPercent
	}
	return int(percent), nil
}

// ValidatePercent ensures the code is between 0-16 inclusive
func ValidateCode(c int64) (code uint32, err error) {
	if c < 0 || c > 16 {
		return code, errInvalidCode
	}
	return uint32(c), nil
}
