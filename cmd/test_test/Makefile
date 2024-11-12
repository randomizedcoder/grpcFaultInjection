#
# /cmd/test/Makefile
#

# ldflags variables to update --version
# short commit hash
COMMIT := $(shell git describe --always)
DATE := $(shell date -u +"%Y-%m-%d-%H:%M")
BINARY := test_test

all: clean build

test:
	go test

clean:
	[ -f ${BINARY} ] && rm -rf ./${BINARY} || true

build:
	go build -ldflags \
		"-X main.commit=${COMMIT} -X main.date=${DATE} -X main.version=${VERSION}" \
		-o ./${BINARY} \
		./${BINARY}.go

loop:
	./loop.bash

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

# end
