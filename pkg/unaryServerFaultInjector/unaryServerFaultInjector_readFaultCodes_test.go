package unaryServerFaultInjector

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type testReadFaultCode struct {
	name          string
	md            metadata.MD
	expectErr     bool
	validateCodes bool
	cs            []codes.Code
}

// https://grpc.io/docs/guides/status-codes/
func TestReadFaultCodes(t *testing.T) {
	tests := []testReadFaultCode{
		{
			name: "valid no fault codes",
			md: metadata.Pairs(
				"anotherHeader", "doesn_t_matter",
			),
			expectErr:     false,
			validateCodes: true,
			cs:            []codes.Code{},
		},
		{
			name: "valid 14",
			md: metadata.Pairs(
				faultcodesHeader, "14",
			),
			expectErr:     false,
			validateCodes: true,
			cs:            []codes.Code{codes.Unavailable},
		},
		{
			name: "valid 10,12,14",
			md: metadata.Pairs(
				faultcodesHeader, "10,12,14",
			),
			expectErr:     false,
			validateCodes: true,
			cs:            []codes.Code{codes.Aborted, codes.Unimplemented, codes.Unavailable},
		},
		{
			name: "invalid -10",
			md: metadata.Pairs(
				faultcodesHeader, "-10",
			),
			expectErr:     true,
			validateCodes: false,
		},
		{
			name: "invalid 17",
			md: metadata.Pairs(
				faultcodesHeader, "17",
			),
			expectErr:     true,
			validateCodes: false,
		},
		{
			name: "invalid 1700",
			md: metadata.Pairs(
				faultcodesHeader, "1700",
			),
			expectErr:     true,
			validateCodes: false,
		},
		{
			name: "invalid blah",
			md: metadata.Pairs(
				faultcodesHeader, "blah",
			),
			expectErr:     true,
			validateCodes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := readFaultCodes(&tt.md)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
			if tt.validateCodes {
				if len(cs) != len(tt.cs) {
					t.Errorf("test: %s,len(cs:%d) != len(tt.cs:%d)", tt.name, len(cs), len(tt.cs))
				}
				// t.Logf("cs:%v", cs)
				// t.Logf("tt.cs:%v", tt.cs)
				if len(tt.cs) > 0 {
					if !reflect.DeepEqual(cs, tt.cs) {
						t.Errorf("test: %s,!reflect.DeepEqual(cs:%v, tt.cs:%v)", tt.name, cs, tt.cs)
					}
				}
			}
		})
	}

}
