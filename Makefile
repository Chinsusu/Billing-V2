GO ?= go

.PHONY: fmt test build build-api build-migrate build-seed build-smoke run-api migrate-validate seed-plan seed-dev smoke-dev-db smoke-dev-api smoke-dev-billing

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

build: build-api build-migrate build-seed build-smoke

build-api:
	$(GO) build -o bin/api ./cmd/api

build-migrate:
	$(GO) build -o bin/migrate ./cmd/migrate

build-seed:
	$(GO) build -o bin/seed ./cmd/seed

build-smoke:
	$(GO) build -o bin/smoke ./cmd/smoke

run-api:
	$(GO) run ./cmd/api

migrate-validate:
	$(GO) run ./cmd/migrate validate

seed-plan:
	$(GO) run ./cmd/seed plan

seed-dev:
	$(GO) run ./cmd/seed dev

smoke-dev-db:
	$(GO) run ./cmd/smoke dev-db

smoke-dev-api:
	$(GO) run ./cmd/smoke dev-api

smoke-dev-billing:
	$(GO) run ./cmd/smoke dev-billing
