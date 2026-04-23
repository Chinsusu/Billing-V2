GO ?= go

.PHONY: fmt test build build-api build-migrate build-seed run-api migrate-validate seed-plan seed-dev

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

build: build-api build-migrate build-seed

build-api:
	$(GO) build -o bin/api ./cmd/api

build-migrate:
	$(GO) build -o bin/migrate ./cmd/migrate

build-seed:
	$(GO) build -o bin/seed ./cmd/seed

run-api:
	$(GO) run ./cmd/api

migrate-validate:
	$(GO) run ./cmd/migrate validate

seed-plan:
	$(GO) run ./cmd/seed plan

seed-dev:
	$(GO) run ./cmd/seed dev
