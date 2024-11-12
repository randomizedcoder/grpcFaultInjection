package validate

import "testing"

func TestValidateModulus(t *testing.T) {
	tests := []struct {
		name      string
		modulus   int64
		expectErr bool
	}{
		{"Valid, low modulus", 1, false},
		{"Valid, mid modulus", 50, false},
		{"Valid, high modulus 100", 100, false},
		{"Valid, high modulus 1000", 1000, false},
		{"Valid, high modulus 10000", 10000, false},
		{"Invalid, negative modulus -1000", -1000, true},
		{"Invalid, negative modulus -10", -10, true},
		{"Invalid, low modulus", 0, true},
		{"Invalid, over 10000 modulus", 10001, true},
		{"Invalid, over 10000 modulus", 11000, true},
		{"Invalid, over 10000 modulus", 100000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateModulus(tt.modulus)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
		})
	}
}

func TestValidatePercent(t *testing.T) {
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
			_, err := ValidatePercent(tt.percent)
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
		{"Valid, code 8", 8, false},
		{"Valid, code 16", 16, false},
		{"Invalid, code -1000", -1000, true},
		{"Invalid, code -10", -10, true},
		{"Invalid, code 17", 17, true},
		{"Invalid, code 100", 100, true},
		{"Invalid, code 10000", 10000, true},
		{"Invalid, max uint32 code", 4294967295, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateCode(tt.code)
			if (err != nil) != tt.expectErr {
				t.Errorf("test:%s, expected error: %v, got: %v", tt.name, tt.expectErr,
					err != nil)
			}
		})
	}
}
