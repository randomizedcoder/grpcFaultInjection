package unaryServerFaultInjector

import "testing"

func TestValidateFaultPercent(t *testing.T) {
	tests := []struct {
		name      string
		percent   int64
		expectErr bool
	}{
		{"Valid, low percent", 1, false},
		{"Valid, mid percent", 50, false},
		{"Valid, high percent", 100, false},
		{"Invalid, negative percent", -10, true},
		{"Invalid, low percent", 0, true},
		{"Invalid, over 100 percent", 101, true},
		{"Invalid, over 100 percent", 110, true},
		{"Invalid, over 100 percent", 11000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateFaultPercent(tt.percent)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
		})
	}
}

func TestValidateCode(t *testing.T) {
	tests := []struct {
		name      string
		code      int64
		expectErr bool
	}{
		{"Valid, code 0", 0, false},
		{"Valid, code 16", 16, false},
		{"Invalid, code -10", -10, true},
		{"Invalid, code 17", 17, true},
		{"Invalid, code 17", 100, true},
		{"Invalid, max uint32 code", 4294967295, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateCode(tt.code)
			if (err != nil) != tt.expectErr {
				t.Errorf("test:%s, expected error: %v, got: %v", tt.name, tt.expectErr,
					err != nil)
			}
		})
	}
}
