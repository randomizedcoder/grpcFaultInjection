#
# /pkg/pkg/rand/Makefile
#

test: TestRandomFaultCode TestRandomSuppliedFaultCode

verbose:
	go test -v

TestRandomFaultCode:
	go test -run TestRandomFaultCode -v

TestRandomSuppliedFaultCode:
	go test -run TestRandomSuppliedFaultCode -v

FindTests:
	grep -R "func Test" ./

# end