package unaryServerFaultInjector

import (
	"testing"

	"google.golang.org/grpc/codes"
)

type LogNoFaultRequestTest struct {
	name string
	s    uint64
	f    uint64
	log  string
}

// go test -run TestLogNoFaultRequest -v
func TestLogNoFaultRequest(t *testing.T) {
	tests := []LogNoFaultRequestTest{
		{
			name: "s = 0, f = 1",
			s:    0,
			f:    1,
			log:  "request success:0 fault:1",
		},
		{
			name: "s = 0, f = 10",
			s:    0,
			f:    10,
			log:  "request success:0 fault:10",
		},
		{
			name: "s = 1, f = 10",
			s:    1,
			f:    10,
			log:  "request success:1 fault:10 ~= 10",
		},
		{
			name: "s = 1, f = 1",
			s:    1,
			f:    1,
			log:  "request success:1 fault:1 ~= 1",
		},
		{
			name: "s = 2, f = 1",
			s:    2,
			f:    1,
			log:  "request success:2 fault:1 ~= 0.5",
		},
		{
			name: "s = 3, f = 1",
			s:    3,
			f:    1,
			log:  "request success:3 fault:1 ~= 0.333",
		},
		{
			name: "s = 3, f = 2",
			s:    3,
			f:    2,
			log:  "request success:3 fault:2 ~= 0.667",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logNoFaultRequest(tt.s, tt.f)
			if log != tt.log {
				t.Errorf("test: %s,log:%s != tt.log:%s", tt.name, log, tt.log)
			}
		})
	}
}

type testLogFaultRequest struct {
	name string
	s    uint64
	f    uint64
	code codes.Code
	log  string
}

// go test -run TestLogFaultRequest -v
func TestLogFaultRequest(t *testing.T) {
	tests := []testLogFaultRequest{
		{
			name: "s = 0, f = 1",
			s:    0,
			f:    1,
			code: 1,
			log:  "request code:Canceled success:0 fault:1",
		},
		{
			name: "s = 0, f = 10",
			s:    0,
			f:    10,
			code: 1,
			log:  "request code:Canceled success:0 fault:10",
		},
		{
			name: "s = 1, f = 10",
			s:    1,
			f:    10,
			code: 1,
			log:  "request code:Canceled success:1 fault:10 ~= 10",
		},
		{
			name: "s = 1, f = 1",
			s:    1,
			f:    1,
			code: 1,
			log:  "request code:Canceled success:1 fault:1 ~= 1",
		},
		{
			name: "s = 2, f = 1",
			s:    2,
			f:    1,
			code: 1,
			log:  "request code:Canceled success:2 fault:1 ~= 0.5",
		},
		{
			name: "s = 3, f = 1",
			s:    3,
			f:    1,
			code: 1,
			log:  "request code:Canceled success:3 fault:1 ~= 0.333",
		},
		{
			name: "s = 3, f = 2",
			s:    3,
			f:    2,
			code: 1,
			log:  "request code:Canceled success:3 fault:2 ~= 0.667",
		},
		{
			name: "s = 3, f = 2, code = 14",
			s:    3,
			f:    2,
			code: 14,
			log:  "request code:Unavailable success:3 fault:2 ~= 0.667",
		},
		{
			name: "s = 3, f = 2, code = 14",
			s:    3,
			f:    2,
			code: 16,
			log:  "request code:Unauthenticated success:3 fault:2 ~= 0.667",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logFaultRequest(tt.s, tt.f, tt.code)
			if log != tt.log {
				t.Errorf("test: %s,log:%s != tt.log:%s", tt.name, log, tt.log)
			}
		})
	}
}
