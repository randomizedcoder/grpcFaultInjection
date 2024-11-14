package unaryClientFaultInjector

import (
	"fmt"
	"strings"
)

func logNoFaultRequest(s uint64, f uint64) string {
	if s == 0 {
		return fmt.Sprintf("no fault request success:%d fault:%d", s, f)
	}
	return strings.TrimRight(
		strings.TrimRight(
			fmt.Sprintf("no fault request success:%d fault:%d ~= %.3f", s, f, float64(f)/float64(s)),
			"0"),
		".")
}

func logFaultRequest(s uint64, f uint64) string {
	if s == 0 {
		return fmt.Sprintf("fault request success:%d fault:%d", s, f)
	}
	return strings.TrimRight(
		strings.TrimRight(fmt.Sprintf("fault request success:%d fault:%d ~= %.3f", s, f, float64(f)/float64(s)),
			"0"),
		".")
}
