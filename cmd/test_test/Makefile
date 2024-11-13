#
# /cmd/test/Makefile
#

# ldflags variables to update --version
# short commit hash
COMMIT := $(shell git describe --always)
DATE := $(shell date -u +"%Y-%m-%d-%H:%M")
BINARY := test_test

all: clean build

test: TestComprehensive

clean:
	[ -f ${BINARY} ] && rm -rf ./${BINARY} || true

hyperfine:
	hyperfine \
		--ignore-failure \
		--runs 100 \
		'go test -v'

hyperfineDebug:
	hyperfine \
		--show-output \
		--runs 100 \
		'go test -v'

TestComprehensive:
	go test -run TestComprehensive -v

# end
