GO ?= go
GOLANGCILINT ?= golangci-lint

BINARY := lander
VERSION ?= $(shell git describe --always --tags 2>/dev/null || echo "undefined")

## Print this help
#  eg: 'make' or 'make help'
help:
	@awk -v skip=1 \
		'/^##/ { sub(/^[#[:blank:]]*/, "", $$0); doc_h=$$0; doc=""; skip=0; next } \
		skip  { next } \
		/^#/  { doc=doc "\n" substr($$0, 2); next } \
		/:/   { sub(/:.*/, "", $$0); printf "\033[1m%-30s\033[0m\033[1m%s\033[0m %s\n\n", $$0, doc_h, doc; skip=1 }' \
		$(MAKEFILE_LIST)

## lint and fix the golang (backend)
backend-lint-fix:
	gofmt -s -w .
	goimports -w .
	golangci-lint run --timeout 3m --fix

## lint without fixing the backend
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

## lint without fixing the frontend
frontend-lint: frontend-install
	cd frontend && npx eslint --ext .js,.vue --max-warnings 0 ./src

## lint and fix the vuejs (frontend)
frontend-lint-fix:
	cd frontend && npx eslint --fix --ext .js,.vue --max-warnings 0 ./src

frontend-build: frontend-install
	cd frontend && npm run build

clean:
	-rm lander
	-rm -rf frontend/node_modules/

## lint everything
lint: backend-lint frontend-lint

## lint and fix everything
lint-fix: backend-lint-fix frontend-lint-fix

test:
	go test ./...

## run the frontend
run-frontend:
	cd frontend && npm run serve

## run the backend
run-backend:
	go run . \
		-labels kubernetes.io/role,node.kubernetes.io/instance-type,node.kubernetes.io/instancegroup,topology.kubernetes.io/zone,node.acmecorp.org/app \
		-clusters main.cluster-a.acmecorp.org,main.cluster-b.acmecorp.org,main.cluster-c.acmecorp.org \
		-debug
