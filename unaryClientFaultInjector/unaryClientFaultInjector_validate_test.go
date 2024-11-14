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
			name: "valid, modulus 1",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Modulus,
					Value: 1,
				},
				Server: ModeValue{
					Mode:  Modulus,
					Value: 1,
				},
				Codes: "10",
			},
			expectErr: false,
		},
		{
			name: "valid, percent 1",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: 1,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: 1,
				},
				Codes: "10",
			},
			expectErr: false,
		},
		{
			name: "valid, modulus 1",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Modulus,
					Value: 10,
				},
				Server: ModeValue{
					Mode:  Modulus,
					Value: 10,
				},
				Codes: "10",
			},
			expectErr: false,
		},
		{
			name: "valid, percent 100",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: 100,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: 100,
				},
				Codes: "10",
			},
			expectErr: false,
		},
		{
			name: "invalid, percent 0",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: 0,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: 0,
				},
				Codes: "10",
			},
			expectErr: true,
		},
		{
			name: "invalid,  percent 0",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: 0,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: 0,
				},
				Codes: "10",
			},
			expectErr: true,
		},
		{
			name: "invalid, percent 101",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: 101,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: 101,
				},
				Codes: "10",
			},
			expectErr: true,
		},
		{
			name: "invalid, percent -1010",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Percent,
					Value: -1010,
				},
				Server: ModeValue{
					Mode:  Percent,
					Value: -1010,
				},
				Codes: "10",
			},
			expectErr: true,
		},
		{
			name: "invalid, modulus -1",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Modulus,
					Value: -1,
				},
				Server: ModeValue{
					Mode:  Modulus,
					Value: -1,
				},
				Codes: "10",
			},
			expectErr: true,
		},
		{
			name: "invalid, modulus 10001",
			conf: UnaryClientInterceptorConfig{
				Client: ModeValue{
					Mode:  Modulus,
					Value: 10001,
				},
				Server: ModeValue{
					Mode:  Modulus,
					Value: 10001,
				},
				Codes: "10",
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
