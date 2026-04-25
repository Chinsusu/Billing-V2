GO ?= go

.PHONY: fmt test build build-api build-migrate build-seed build-smoke build-worker run-api run-worker migrate-validate seed-plan seed-dev smoke-dev-db smoke-dev-api smoke-dev-billing contract-guard error-code-guard task-guard

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

build: build-api build-migrate build-seed build-smoke build-worker

build-api:
	$(GO) build -o bin/api ./cmd/api

build-migrate:
	$(GO) build -o bin/migrate ./cmd/migrate

build-seed:
	$(GO) build -o bin/seed ./cmd/seed

build-smoke:
	$(GO) build -o bin/smoke ./cmd/smoke

build-worker:
	$(GO) build -o bin/worker ./cmd/worker

run-api:
	$(GO) run ./cmd/api

run-worker:
	$(GO) run ./cmd/worker provision-once

migrate-validate:
	$(GO) run ./cmd/migrate validate

contract-guard:
	$(GO) run ./cmd/contractguard

error-code-guard:
	$(GO) run ./cmd/errorcodeguard

task-guard:
	$(GO) run ./cmd/taskguard

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
