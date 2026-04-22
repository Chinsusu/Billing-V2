GO ?= go

.PHONY: fmt test build build-api run-api

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

build: build-api

build-api:
	$(GO) build -o bin/api ./cmd/api

run-api:
	$(GO) run ./cmd/api
