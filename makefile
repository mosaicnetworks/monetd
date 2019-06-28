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
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/monetd



installcfg:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/monetcfgsrv


installcli:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/monetcli


docker:
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./docker/monetd ./cmd/monetd/
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./docker/monetcli ./cmd/monetcli/



test:
	glide novendor | xargs go test

.PHONY: all vendor install installd installcli test update docker
