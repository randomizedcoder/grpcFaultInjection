package unaryServerFaultInjector

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

type readFaultModulusTest struct {
	name            string
	md              metadata.MD
	expectErr       bool
	found           bool
	validateModulus bool
	faultModulus    uint64
}

// go test -run TestReadFaultModulus -v
func TestReadFaultModulus(t *testing.T) {
	tests := []readFaultModulusTest{
		{
			name: "valid no fault modulus header",
			md: metadata.Pairs(
				"anotherHeader", "doesn_t_matter",
			),
			expectErr:       false,
			found:           false,
			validateModulus: false,
			faultModulus:    0,
		},
		{
			name: "valid, 1 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "1",
			),
			expectErr:       false,
			found:           true,
			validateModulus: true,
			faultModulus:    1,
		},
		{
			name: "valid, 50 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "50",
			),
			expectErr:       false,
			found:           true,
			validateModulus: true,
			faultModulus:    50,
		},
		{
			name: "valid, 100 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "100",
			),
			expectErr:       false,
			found:           true,
			validateModulus: true,
			faultModulus:    100,
		},
		{
			name: "valid, 10000 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "10000",
			),
			expectErr:       false,
			found:           true,
			validateModulus: true,
			faultModulus:    10000,
		},
		{
			name: "invalid, zero modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "0",
			),
			expectErr:       true,
			found:           true,
			validateModulus: true,
			faultModulus:    0,
		},
		{
			name: "invalid, 10001 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "10001",
			),
			expectErr:       true,
			found:           true,
			validateModulus: false,
			faultModulus:    10001,
		},
		{
			name: "invalid, 100001 modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "100001",
			),
			expectErr:       true,
			found:           true,
			validateModulus: false,
			faultModulus:    100001,
		},
		{
			name: "invalid, 50.5 modulus (non integer)",
			md: metadata.Pairs(
				faultmodulusHeader, "50.5",
			),
			expectErr:       true,
			found:           true,
			validateModulus: false,
			faultModulus:    0,
		},
		{
			name: "invalid, blah modulus",
			md: metadata.Pairs(
				faultmodulusHeader, "blah",
			),
			expectErr:       true,
			found:           true,
			validateModulus: false,
			faultModulus:    101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, faultModulus, err := readFaultModulus(&tt.md, 0)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
			if found != tt.found {
				t.Errorf("test: %s,found:%t != tt.found%t", tt.name, found, tt.found)
			}
			if tt.validateModulus {
				if faultModulus != tt.faultModulus {
					t.Errorf("test: %s,faultModulus:%d != tt.faultModulus:%d", tt.name, faultModulus, tt.faultModulus)
				}
			}
		})
	}

}
