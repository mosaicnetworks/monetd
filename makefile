BUILD_TAGS?=evml

all: vendor install

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

update:
	glide update

# install compiles and places the binary in GOPATH/bin
install: installd installcfg installgiv

installd:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/monetd/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/monetd

installcfg:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/monetd/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/monetcfgsrv

installcli:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/monetd/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/monetcli

installgiv:
	go install \
		--ldflags "-X github.com/mosaicnetworks/giverny/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/giverny/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/giverny

docker:
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./docker/monetd ./cmd/monetd/
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./docker/monetcli ./cmd/monetcli/

test: testmonetd testevml testbabble

testmonetd:
	@echo "\nMonet Tests\n\n" ; glide novendor | xargs go test | sed -e 's?github.com/mosaicnetworks/?.../?g'


testevml:
	@echo "\nEVM-Lite Tests\n\n" ; cd vendor/github.com/mosaicnetworks/evm-lite ; go test ./src/...| sed -e 's?github.com/mosaicnetworks/monetd/vendor/github.com/mosaicnetworks/?.../vendor/.../?g'


testbabble:
	@echo "\nBabble Tests\n\n" ; cd vendor/github.com/mosaicnetworks/babble ;   go test ./src/... -count=1 -tags=unit  | sed -e 's?github.com/mosaicnetworks/monetd/vendor/github.com/mosaicnetworks/?.../vendor/.../?g'



.PHONY: all vendor install installd installcli installgiv test update docker testmonetd testevml testbabble
