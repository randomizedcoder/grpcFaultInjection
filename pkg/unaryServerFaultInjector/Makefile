#
# /pkg/pkg/unaryServerFaultInjector/Makefile
#

test: TestLogNoFaultRequest TestLogFaultRequest TestReadFaultCodes TestReadFaultPercent

verbose:
	go test -v

TestLogNoFaultRequest:
	go test -run TestLogNoFaultRequest -v

TestLogFaultRequest:
	go test -run TestLogFaultRequest -v

TestReadFaultCodes:
	go test -run TestReadFaultCodes -v

TestReadFaultPercent:
	go test -run TestReadFaultPercent -v

FindTests:
	grep -R "func Test" ./

# end