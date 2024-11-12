package rand

// This .go file holds the functions performing random functions
// allowing "unsafe" to only be used in this file

import (
	_ "unsafe"

	"google.golang.org/grpc/codes"
)

// unsafe for FastRandN

const (

	// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	// We want a range between 1 and 16, so our maxCode is 15, cos we will +1
	maxCode = 15
)

// https://cs.opensource.google/go/go/+/master:src/runtime/stubs.go;l=151?q=FastRandN&ss=go%2Fgo
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/

//go:linkname FastRandN runtime.fastrandn
func FastRandN(n uint32) uint32

func FastRandNInt() int {
	return int(FastRandN(100))
}

// randomFaultCode returns ANY random fault code ( 1-16 )
// does NOT return code 0
func RandomFaultCode() (code codes.Code) {
	return codes.Code(FastRandN(maxCode) + 1)
}

// randomSuppliedFaultCode randomly selects one of the "faultcodes"
func RandomSuppliedFaultCode(cs *[]codes.Code) (code codes.Code) {
	return (*cs)[int(FastRandN(uint32(len(*cs))))]
}
