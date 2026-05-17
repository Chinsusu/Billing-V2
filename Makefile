GO ?= go

.PHONY: go-packages fmt test build build-api build-migrate build-seed build-smoke build-worker run-api run-worker migrate-validate seed-plan seed-dev smoke-dev-db smoke-dev-api smoke-dev-billing smoke-dev-topup-review smoke-dev-target-auth-rbac smoke-dev-target-credential-reveal backup-restore-drill-plan backup-restore-drill full-e2e-quality-gate contract-guard error-code-guard task-guard

go-packages:
	$(GO) run ./cmd/gopackages

fmt:
	$(GO) fmt $$($(GO) run ./cmd/gopackages)

test:
	$(GO) test $$($(GO) run ./cmd/gopackages)

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

smoke-dev-topup-review:
	$(GO) run ./cmd/smoke dev-topup-review

smoke-dev-target-auth-rbac:
	$(GO) run ./cmd/smoke dev-target-auth-rbac

smoke-dev-target-credential-reveal:
	$(GO) run ./cmd/smoke dev-target-credential-reveal

backup-restore-drill-plan:
	bash scripts/backup_restore_drill.sh --plan

backup-restore-drill:
	bash scripts/backup_restore_drill.sh --run

full-e2e-quality-gate:
	bash scripts/full_e2e_quality_gate.sh
