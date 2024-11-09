#
# /pkg/pkg/validate/Makefile
#

test: TestValidatePercent TestValidateCode

verbose:
	go test -v

TestValidatePercent:
	go test -run TestValidatePercent -v

TestValidateCode:
	go test -run TestValidateCode -v

FindTests:
	grep -R "func Test" ./

# end