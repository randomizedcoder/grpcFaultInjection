package unaryClientFaultInjector

import "testing"

type CheckConfigTest struct {
	name      string
	conf      UnaryClientInterceptorConfig
	expectErr bool
}

// go test -run TestCheckConfig -v
func TestCheckConfig(t *testing.T) {
	tests := []CheckConfigTest{
		{
			name: "valid, clean 1",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 1,
				ServerFaultPercent: 1,
				ServerFaultCodes:   "10",
			},
			expectErr: false,
		},
		{
			name: "valid, clean 10",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 10,
				ServerFaultPercent: 10,
				ServerFaultCodes:   "10",
			},
			expectErr: false,
		},
		{
			name: "valid, clean 100",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: 100,
				ServerFaultCodes:   "10,12,14",
			},
			expectErr: false,
		},
		{
			name: "invalid, clean 0",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 0,
				ServerFaultPercent: 0,
				ServerFaultCodes:   "10,12,14",
			},
			expectErr: true,
		},
		{
			name: "invalid, clean 10000",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 10000,
				ServerFaultPercent: 10000,
				ServerFaultCodes:   "10,12,14",
			},
			expectErr: true,
		},
		{
			name: "invalid, clean client -100",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: -100,
				ServerFaultPercent: 100,
				ServerFaultCodes:   "10,12,14",
			},
			expectErr: true,
		},
		{
			name: "invalid, clean server -100",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: -100,
				ServerFaultCodes:   "10,12,14",
			},
			expectErr: true,
		},
		{
			name: "invalid, clean 100, invalid codes",
			conf: UnaryClientInterceptorConfig{
				ClientFaultPercent: 100,
				ServerFaultPercent: 100,
				ServerFaultCodes:   "not_a_code",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckConfig(tt.conf)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
		})
	}
}

type validateCodesTest struct {
	codes     string
	expectErr bool
}

// go test -run TestValidateCodes -v
func TestValidateCodes(t *testing.T) {
	tests := []validateCodesTest{
		{
			codes:     "14",
			expectErr: false,
		},
		{
			codes:     "10,12,14",
			expectErr: false,
		},
		{
			codes:     "-10",
			expectErr: true,
		},
		{
			codes:     "17",
			expectErr: true,
		},
		{
			codes:     "1700",
			expectErr: true,
		},
		{
			codes:     "invalid blah",
			expectErr: true,
		},
		{
			codes:     "",
			expectErr: true,
		},
		{
			codes:     ",,,,",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.codes, func(t *testing.T) {
			err := validateCodes(tt.codes)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.codes, tt.expectErr, err != nil)
			}
		})
	}
}
