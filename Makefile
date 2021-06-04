GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --tags 2>/dev/null || echo "undefined")

lint-go:
	golangci-lint run

lint-vue:
	npm -g i eslint
	cd frontend && eslint --ext .js,.vue --max-warnings 0 ./src

#"-s -w" strips debug headers
build-go: lint-go
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux $(GO) build -v -a -installsuffix cgo -ldflags="-X main.VERSION=${VERSION} -s -w" -o lander .

build-node: lint-vue
	cd frontend && npm i --from-lockfile && npm run build

clean:
	-rm lander
