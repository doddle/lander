GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")

run:
	CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go run .

clean:
	-rm lander

lander: clean
	GO111MODULE=on CGO_ENABLED=0 $(GO) build -v -a -installsuffix cgo -ldflags="-X main.VERSION=${VERSION}" -o $@ github.com/digtux/lander
