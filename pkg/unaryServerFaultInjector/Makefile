#
# /pkg/pkg/unaryServerFaultInjector/Makefile
#

test:
	go test

verbose:
	go test -v

mytests: TestValidateFaultPercent TestValidateCode TestReadFaultCodes TestReadFaultPercent

TestValidateFaultPercent:
	go test -run TestValidateFaultPercent -v

TestValidateCode:
	go test -run TestValidateCode -v

TestReadFaultCodes:
	go test -run TestReadFaultCodes -v

TestReadFaultPercent:
	go test -run TestReadFaultPercent -v

# end