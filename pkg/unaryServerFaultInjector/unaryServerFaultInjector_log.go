package unaryServerFaultInjector

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
)

func logNoFaultRequest(s uint64, f uint64) string {
	if s == 0 {
		return fmt.Sprintf("request success:%d fault:%d", s, f)
	}
	return strings.TrimRight(
		strings.TrimRight(
			fmt.Sprintf("request success:%d fault:%d ~= %.3f", s, f, float64(f)/float64(s)),
			"0"),
		".")
}

func logFaultRequest(s uint64, f uint64, code codes.Code) string {
	if s == 0 {
		return fmt.Sprintf("request code:%s success:%d fault:%d", code.String(), s, f)
	}
	return strings.TrimRight(
		strings.TrimRight(
			fmt.Sprintf("request code:%s success:%d fault:%d ~= %.3f", code.String(), s, f, float64(f)/float64(s)),
			"0"),
		".")

}
