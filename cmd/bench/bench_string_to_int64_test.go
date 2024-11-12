package main

import (
	"fmt"
	"strconv"
	"testing"
)

// Benchmark using strconv.ParseInt with base 10
func BenchmarkParseInt(b *testing.B) {
	str := "1234567890"
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseInt(str, 10, 64)
	}
}

// Benchmark using strconv.Atoi (which returns an int but can be converted to int64)
func BenchmarkAtoi(b *testing.B) {
	str := "1234567890"
	for i := 0; i < b.N; i++ {
		val, _ := strconv.Atoi(str)
		_ = int64(val)
	}
}

// Benchmark using custom loop for conversion (illustrative but less efficient)
func BenchmarkCustomConversion(b *testing.B) {
	str := "1234567890"
	for i := 0; i < b.N; i++ {
		var result int64
		for _, ch := range str {
			result = result*10 + int64(ch-'0')
		}
		_ = result
	}
}

func BenchmarkItoaSprintBig(b *testing.B) {
	j := "1234567890"
	for i := 0; i < b.N; i++ {
		val := fmt.Sprint(j)
		_ = val
	}
}

func BenchmarkItoaSprintfBig(b *testing.B) {
	j := 1234567890
	for i := 0; i < b.N; i++ {
		val := fmt.Sprintf("%d", j)
		_ = val
	}
}
