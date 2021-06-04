GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --tags 2>/dev/null || echo "undefined")

backend-lint:
	golangci-lint run --timeout 3m

#"-s -w" strips debug headers
backend-build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux $(GO) build \
		-v \
		-a \
		-installsuffix cgo \
		-ldflags="-X main.VERSION=${VERSION} -s -w" \
		-o lander .


frontend-install:
	cd frontend && npm i --from-lockfile

frontend-lint: frontend-install
	cd frontend && npm i eslint -g && eslint --ext .js,.vue --max-warnings 0 ./src

frontend-lint-fix:
	cd frontend && npm i eslint -g && eslint --fix --ext .js,.vue --max-warnings 0 ./src

frontend-build: frontend-install
	cd frontend && npm run build

clean:
	-rm lander
	-rm -rf frontend/node_modules/
