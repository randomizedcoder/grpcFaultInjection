package unaryServerFaultInjector

import (
	"testing"

	"google.golang.org/grpc/codes"
)

// This is a bit of a silly test.  Could probably remove this.
// Just trying to make sure every function has a test
// go test -run testRandomFaultCode -v
func TestRandomFaultCode(t *testing.T) {

	interations := 1000

	for i := 0; i < interations; i++ {
		code := randomFaultCode()
		switch code {
		case codes.OK:
			t.Errorf("TestRandomFaultCode found code:%s == codes.ok", code)
		default:
		}
	}
}

type RandomSuppliedFaultCodeTest struct {
	name       string
	iterations int
	cs         []codes.Code
}

// go test -run TestRandomSuppliedFaultCode -v
func TestRandomSuppliedFaultCode(t *testing.T) {
	tests := []RandomSuppliedFaultCodeTest{
		{
			name:       "single code",
			iterations: 100,
			cs:         []codes.Code{codes.Unimplemented},
		},
		{
			name:       "three codes",
			iterations: 100,
			cs: []codes.Code{
				codes.Aborted,
				codes.Unimplemented,
				codes.Unavailable,
			},
		},
		{
			name:       "eight codes",
			iterations: 100,
			cs: []codes.Code{
				codes.NotFound,
				codes.AlreadyExists,
				codes.PermissionDenied,
				codes.ResourceExhausted,
				codes.FailedPrecondition,
				codes.Aborted,
				codes.OutOfRange,
				codes.Unimplemented,
			},
		},
		{
			name:       "all codes",
			iterations: 100,
			cs: []codes.Code{
				codes.OK,
				codes.Canceled,
				codes.Unknown,
				codes.InvalidArgument,
				codes.DeadlineExceeded,
				codes.NotFound,
				codes.AlreadyExists,
				codes.PermissionDenied,
				codes.ResourceExhausted,
				codes.FailedPrecondition,
				codes.Aborted,
				codes.OutOfRange,
				codes.Unimplemented,
				codes.Internal,
				codes.Unavailable,
				codes.DataLoss,
				codes.Unauthenticated,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			myMap := ConvertSliceToMap(
				tt.cs,
				func(c codes.Code) codes.Code { return c },
				func(c codes.Code) codes.Code { return c },
			)

			for i := 0; i < tt.iterations; i++ {

				code := randomSuppliedFaultCode(&tt.cs)

				_, found := myMap[code]
				if !found {
					t.Error("TestRandomSuppliedFaultCode key not found", code)
				}
			}
		})
	}
}

func ConvertSliceToMap[T any, K comparable, V any](slice []T, keyMapper func(T) K, valueMapper func(T) V) map[K]V {
	result := make(map[K]V)
	for _, item := range slice {
		key := keyMapper(item)
		value := valueMapper(item)
		result[key] = value
	}
	return result
}
