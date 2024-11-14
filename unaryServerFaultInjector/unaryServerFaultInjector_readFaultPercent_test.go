package unaryServerFaultInjector

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

type testReadFaultPercent struct {
	name            string
	md              metadata.MD
	expectErr       bool
	found           bool
	validatePercent bool
	faultPercent    int
}

// go test -run TestReadFaultPercent -v
func TestReadFaultPercent(t *testing.T) {
	tests := []testReadFaultPercent{
		{
			name: "valid no fault percent header",
			md: metadata.Pairs(
				"anotherHeader", "doesn_t_matter",
			),
			expectErr:       false,
			found:           false,
			validatePercent: false,
			faultPercent:    0,
		},
		{
			name: "valid, 1 percent",
			md: metadata.Pairs(
				faultpercentHeader, "1",
			),
			expectErr:       false,
			found:           true,
			validatePercent: true,
			faultPercent:    1,
		},
		{
			name: "valid, 50 percent",
			md: metadata.Pairs(
				faultpercentHeader, "50",
			),
			expectErr:       false,
			found:           true,
			validatePercent: true,
			faultPercent:    50,
		},
		{
			name: "valid, 100 percent",
			md: metadata.Pairs(
				faultpercentHeader, "100",
			),
			expectErr:       false,
			found:           true,
			validatePercent: true,
			faultPercent:    100,
		},
		{
			name: "invalid, zero percent",
			md: metadata.Pairs(
				faultpercentHeader, "0",
			),
			expectErr:       true,
			found:           true,
			validatePercent: true,
			faultPercent:    0,
		},
		{
			name: "invalid, negative percent",
			md: metadata.Pairs(
				faultpercentHeader, "-10",
			),
			expectErr:       true,
			found:           true,
			validatePercent: false,
			faultPercent:    -10,
		},
		{
			name: "invalid, 101 percent",
			md: metadata.Pairs(
				faultpercentHeader, "101",
			),
			expectErr:       true,
			found:           true,
			validatePercent: false,
			faultPercent:    101,
		},
		{
			name: "invalid, 10001 percent",
			md: metadata.Pairs(
				faultpercentHeader, "10001",
			),
			expectErr:       true,
			found:           true,
			validatePercent: false,
			faultPercent:    10001,
		},
		{
			name: "invalid, 50.5 percent (non integer)",
			md: metadata.Pairs(
				faultpercentHeader, "50.5",
			),
			expectErr:       true,
			found:           true,
			validatePercent: false,
			faultPercent:    0,
		},
		{
			name: "invalid, blah percent",
			md: metadata.Pairs(
				faultpercentHeader, "blah",
			),
			expectErr:       true,
			found:           true,
			validatePercent: false,
			faultPercent:    101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, faultPercent, err := readFaultPercent(&tt.md, 0)
			if (err != nil) != tt.expectErr {
				t.Errorf("test: %s, expected error: %v, got: %v", tt.name, tt.expectErr, err != nil)
			}
			if found != tt.found {
				t.Errorf("test: %s,found:%t != tt.found%t", tt.name, found, tt.found)
			}
			if tt.validatePercent {
				if faultPercent != tt.faultPercent {
					t.Errorf("test: %s,faultPercent:%d != tt.faultPercent:%d", tt.name, faultPercent, tt.faultPercent)
				}
			}
		})
	}

}
