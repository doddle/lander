GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --tags 2>/dev/null || echo "undefined")

fmt:
	gofmt -s -w .
	goimports -w .
	golangci-lint run --timeout 3m --fix

backend-lint:
	gofmt -s -l .
	goimports -l .
	golangci-lint run --timeout 3m
	@echo run \'make fmt\' if lint returned changes to stdout

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
	cd frontend && npx eslint --ext .js,.vue --max-warnings 0 ./src

frontend-lint-fix:
	cd frontend && npx eslint --fix --ext .js,.vue --max-warnings 0 ./src

frontend-build: frontend-install
	cd frontend && npm run build

clean:
	-rm lander
	-rm -rf frontend/node_modules/

lint: backend-lint frontend-lint

test:
	go test ./...
