GO ?= go

.PHONY: fmt test build build-api build-migrate run-api migrate-validate

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

build: build-api build-migrate

build-api:
	$(GO) build -o bin/api ./cmd/api

build-migrate:
	$(GO) build -o bin/migrate ./cmd/migrate

run-api:
	$(GO) run ./cmd/api

migrate-validate:
	$(GO) run ./cmd/migrate validate
