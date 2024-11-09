package unaryServerFaultInjector

import (
	"testing"
)

func TestValidateFaultPercent(t *testing.T) {
	tests := []struct {
		name      string
		percent   int
		expectErr bool
	}{
		{"Valid, low percent", 0, false},
		{"Valid, mid percent", 50, false},
		{"Valid, high percent", 100, false},
		{"Invalid negative percent", -10, true},
		{"Invalid over 100 percent", 110, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFaultPercent(tt.percent)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
		})
	}
}

// func TestReadFailPercent(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		envValue     string
// 		expectResult int
// 		expectErr    bool
// 	}{
// 		{"Valid percent 0", "0", 0, false},
// 		{"Valid percent 50", "50", 50, false},
// 		{"Valid percent 100", "100", 100, false},
// 		{"Invalid percent non-integer", "abc", 0, true},
// 		{"Invalid percent over 100", "110", 0, true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Setenv("FAIL_PERCENT", tt.envValue) // Set the environment variable
// 			result, err := readFailPercent(tt.envValue)
// 			if (err != nil) != tt.expectErr || result != tt.expectResult {
// 				t.Errorf("expected result: %v, got: %v, expected error: %v, got: %v",
// 					tt.expectResult, result, tt.expectErr, err != nil)
// 			}
// 		})
// 	}
// }

func TestValidateCodeUint32(t *testing.T) {
	tests := []struct {
		name      string
		code      uint32
		expectErr bool
	}{
		{"Valid, code 0", 0, false},
		{"Valid, code 16", 16, false},
		{"Invalid, code 17", 17, true},
		{"Invalid, code 17", 100, true},
		{"Invalid, max uint32 code", 4294967295, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCodeUint32(tt.code)
			if (err != nil) != tt.expectErr {
				t.Errorf("test:%s, expected error: %v, got: %v", tt.name, tt.expectErr,
					err != nil)
			}
		})
	}
}

// func TestReadFailCodes(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		envValue     string
// 		expectResult []uint32
// 		expectErr    bool
// 	}{
// 		{"Valid single code", "404", []uint32{404}, false},
// 		{"Valid multiple codes", "404,500,503", []uint32{404, 500, 503}, false},
// 		{"Invalid non-integer code", "404,abc", nil, true},
// 		{"Empty input", "", nil, false}, // Assuming empty input is valid
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Setenv("FAIL_CODES", tt.envValue) // Set the environment variable
// 			result, err := readFailCodes()
// 			if (err != nil) != tt.expectErr || !equalUint32Slices(result, tt.expectResult) {
// 				t.Errorf("expected result: %v, got: %v, expected error: %v, got: %v", tt.expectResult, result, tt.expectErr, err != nil)
// 			}
// 		})
// 	}
// }

// Helper function to compare slices of uint32
func equalUint32Slices(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
