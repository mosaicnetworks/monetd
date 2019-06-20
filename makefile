BUILD_TAGS?=evml


all: vendor install


# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install


update:
	glide update


# install compiles and places the binary in GOPATH/bin
install: installd installcli

installd:
	go install \
		--ldflags "-X github.com/mosaicnetworks/evm-lite/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/monetd

installcli:
	go install \
		--ldflags "-X github.com/mosaicnetworks/evm-lite/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/monetcli

test:
	glide novendor | xargs go test

.PHONY: all vendor install installd installcli test update
