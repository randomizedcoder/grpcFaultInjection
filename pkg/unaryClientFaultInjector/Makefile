#
# /pkg/pkg/unaryClientFaultInjector/Makefile
#

test: TestCheckConfig TestValidateCodes TestLogNoFaultRequest TestLogFaultRequest

verbose:
	go test -v

TestCheckConfig:
	go test -run TestCheckConfig -v

TestValidateCodes:
	go test -run TestValidateCodes -v

TestLogNoFaultRequest:
	go test -run TestLogNoFaultRequest -v

TestLogFaultRequest:
	go test -run TestLogFaultRequest -v

FindTests:
	grep -R "func Test" ./

# end