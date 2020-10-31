GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")

#"-s -w" strips debug headers
build-go:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux $(GO) build -v -a -installsuffix cgo -ldflags="-X main.VERSION=${VERSION} -s -w" -o lander .

build-node:
	cd frontend && npm i --from-lockfile && npm run build

clean:
	-rm lander
