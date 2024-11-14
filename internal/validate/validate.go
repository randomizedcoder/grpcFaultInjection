package validate

import "errors"

var (
	errInvalidModulus = errors.New("invalid modulus")
	errInvalidPercent = errors.New("invalid percent")
	errInvalidCode    = errors.New("invalid code")
)

// ValidateModulus ensure the modulus is between 1-10000 inclusive
func ValidateModulus(modulus int64) (modulusInt uint64, err error) {
	if modulus < 1 || modulus > 10000 {
		return modulusInt, errInvalidModulus
	}
	return uint64(modulus), nil
}

// ValidatePercent ensure the percentage is between 1-100 inclusive
func ValidatePercent(percent int64) (percentInt int, err error) {
	if percent < 1 || percent > 100 {
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
