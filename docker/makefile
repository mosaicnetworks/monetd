MONETD_VERSION=0.3.1dev
GIVERNY_VERSION=0.3.1dev

all: solidity monetd-image giverny-image

solidity:
	wget -O solc https://github.com/ethereum/solidity/releases/download/v0.5.10/solc-static-linux
	chmod +x solc

monetd-image:
	go build \
		--ldflags '-extldflags "-static"' \
		-o monetd/monetd ../cmd/monetd/
	cp solc monetd/solc
	docker build --no-cache=true -t mosaicnetworks/monetd:$(MONETD_VERSION) -t mosaicnetworks/monetd:latest ./monetd/

giverny-image:
	go build \
		--ldflags '-extldflags "-static"' \
		-o giverny/giverny ../cmd/giverny/
	cp solc giverny/solc
	docker build --no-cache=true -t mosaicnetworks/giverny:$(GIVERNY_VERSION) -t mosaicnetworks/giverny:latest ./giverny/

.PHONY: solidity monetd-image giverny-image
