#
# /pkg/pkg/validate/Makefile
#

test: TestValidateModulus TestValidatePercent TestValidateCode

simpleTest:
	go test .

verbose:
	go test -v

TestValidateModulus:
	go test -run TestValidateModulus -v

TestValidatePercent:
	go test -run TestValidatePercent -v

TestValidateCode:
	go test -run TestValidateCode -v

FindTests:
	grep -R "func Test" ./

# end