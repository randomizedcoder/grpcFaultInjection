package unaryClientFaultInjector

import (
	"testing"
)

type logNoFaultRequestTest struct {
	name string
	s    uint64
	f    uint64
	log  string
}

// go test -run TestLogNoFaultRequest -v
func TestLogNoFaultRequest(t *testing.T) {
	tests := []logNoFaultRequestTest{
		{
			name: "s = 0, f = 1",
			s:    0,
			f:    1,
			log:  "no fault request success:0 fault:1",
		},
		{
			name: "s = 0, f = 10",
			s:    0,
			f:    10,
			log:  "no fault request success:0 fault:10",
		},
		{
			name: "s = 1, f = 10",
			s:    1,
			f:    10,
			log:  "no fault request success:1 fault:10 ~= 10",
		},
		{
			name: "s = 1, f = 1",
			s:    1,
			f:    1,
			log:  "no fault request success:1 fault:1 ~= 1",
		},
		{
			name: "s = 2, f = 1",
			s:    2,
			f:    1,
			log:  "no fault request success:2 fault:1 ~= 0.5",
		},
		{
			name: "s = 3, f = 1",
			s:    3,
			f:    1,
			log:  "no fault request success:3 fault:1 ~= 0.333",
		},
		{
			name: "s = 3, f = 2",
			s:    3,
			f:    2,
			log:  "no fault request success:3 fault:2 ~= 0.667",
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

type logFaultRequestTest struct {
	name string
	s    uint64
	f    uint64
	log  string
}

// go test -run TestLogFaultRequest -v
func TestLogFaultRequest(t *testing.T) {
	tests := []logFaultRequestTest{
		{
			name: "s = 0, f = 1",
			s:    0,
			f:    1,
			log:  "fault request success:0 fault:1",
		},
		{
			name: "s = 0, f = 10",
			s:    0,
			f:    10,
			log:  "fault request success:0 fault:10",
		},
		{
			name: "s = 1, f = 10",
			s:    1,
			f:    10,
			log:  "fault request success:1 fault:10 ~= 10",
		},
		{
			name: "s = 1, f = 1",
			s:    1,
			f:    1,
			log:  "fault request success:1 fault:1 ~= 1",
		},
		{
			name: "s = 2, f = 1",
			s:    2,
			f:    1,
			log:  "fault request success:2 fault:1 ~= 0.5",
		},
		{
			name: "s = 3, f = 1",
			s:    3,
			f:    1,
			log:  "fault request success:3 fault:1 ~= 0.333",
		},
		{
			name: "s = 3, f = 2",
			s:    3,
			f:    2,
			log:  "fault request success:3 fault:2 ~= 0.667",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logFaultRequest(tt.s, tt.f)
			if log != tt.log {
				t.Errorf("test: %s,log:%s != tt.log:%s", tt.name, log, tt.log)
			}
		})
	}
}
