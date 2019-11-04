BUILD_TAGS?=evml

all: vendor install

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	 (rm glide.lock || rm -rf vendor ) && glide install

# install compiles and places the binary in GOPATH/bin
install: installd installgiv

installd:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/monetd/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/monetd

installgiv:
	go install \
		--ldflags "-X github.com/mosaicnetworks/monetd/src/version.GitCommit=`git rev-parse HEAD` -X github.com/mosaicnetworks/monetd/src/version.GitBranch=`git symbolic-ref --short HEAD`" \
		./cmd/giverny

docker:
	$(MAKE) -C docker

test: testmonetd testevml testbabble

testmonetd:
	@echo "\nMonetd Tests\n\n" ; glide novendor | xargs go test | sed -e 's?github.com/mosaicnetworks/?.../?g'

testevml:
	@echo "\nEVM-Lite Tests\n\n" ; cd vendor/github.com/mosaicnetworks/evm-lite ; go test ./src/... -count=1 -tags=unit | sed -e 's?github.com/mosaicnetworks/monetd/vendor/github.com/mosaicnetworks/?.../vendor/.../?g'

testbabble:
	@echo "\nBabble Tests\n\n" ; cd vendor/github.com/mosaicnetworks/babble ;   go test ./src/... -count=1 -tags=unit | sed -e 's?github.com/mosaicnetworks/monetd/vendor/github.com/mosaicnetworks/?.../vendor/.../?g'

dist:
	xgo --targets=*/amd64 --dest=build/  ./cmd/monetd/ 
 
lint:
	glide novendor | xargs golint

.PHONY: all vendor install installd installcli installgiv test update docker testmonetd testevml testbabble lint
