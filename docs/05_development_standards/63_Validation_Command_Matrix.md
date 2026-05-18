# Validation Command Matrix

**Scope:** One place for local and CI validation commands by change type.

Use this before opening a PR. Add the commands you actually ran to the task log and PR body.

## Ordering Rules

- Run commands from the repo root unless a row says `cd frontend`.
- Run focused tests first when they exist, then broader repo checks.
- Do not run `npm run build`, `npm run smoke:admin`, or `npm run smoke:admin:ci` in parallel. They read or write `.next`.
- For frontend CI smoke, run `npm run build` before `npm run smoke:admin:ci`.
- If `make` is not available on Windows, use `go run ./cmd/gopackages` to get the repo-scoped Go package list and pass it to `go fmt` or `go test`.
- Do not use production DSNs, provider accounts, or real customer data for local smoke.

## Command Reference

| Purpose | Preferred command | Windows/no-make equivalent | When to run |
| --- | --- | --- | --- |
| Repo Go package list | `make go-packages` | `go run ./cmd/gopackages` | Before manual Go fmt/test commands when `frontend/node_modules` exists. |
| Go format | `make fmt` | `$pkgs = go run ./cmd/gopackages; go fmt $pkgs` | Go code changed. |
| Go tests | `make test` | `$pkgs = go run ./cmd/gopackages; go test $pkgs` | Backend, DB, provider, shared code, or full-stack changed. |
| Go build | `make build` | `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker` | Backend entrypoints or shared Go packages changed. |
| Migration syntax | `make migrate-validate` | `go run ./cmd/migrate validate` | Migration files or migrator changed. |
| API contract guard | `make contract-guard` | `go run ./cmd/contractguard` | Backend route, permission, query, response shape, or API docs changed. |
| Error code guard | `make error-code-guard` | `go run ./cmd/errorcodeguard` | Error envelope, public error code, handler error behavior, or error docs changed. |
| Task board guard | `make task-guard` | `go run ./cmd/taskguard` | `TASKS.md`, task files, task batches, or board cleanup changed. |
| Frontend install | `cd frontend && npm ci` | same | Fresh checkout, dependency change, or CI parity check. |
| Frontend dependency audit | `cd frontend && npm audit --omit=dev` | same | Dependency or lockfile changed; useful before frontend PRs. |
| Frontend sensitive text guard | `cd frontend && npm run check:sensitive-text` | same | Frontend source, mocks, API mapping, or smoke scripts changed. |
| Frontend lint | `cd frontend && npm run lint` | same | Frontend source changed. |
| Frontend build | `cd frontend && npm run build` | same | Frontend source, config, dependency, or app shell changed. |
| Frontend admin browser smoke | `cd frontend && npm run smoke:admin` | same | Admin navigation, screens, API adapter, mock data, or response mapping changed. |
| Frontend CI browser smoke | `cd frontend && npm run smoke:admin:ci` | same, after `npm run build` | CI parity or standalone output changed. |
| Install Playwright Chromium | `cd frontend && npx playwright install chromium` | same | Browser smoke fails because Chromium is missing locally. |
| DB smoke | `make smoke-dev-db` | `go run ./cmd/smoke dev-db` | Migration, seed, DB lifecycle, or billing seed data changed. |
| API smoke | `make smoke-dev-api` | `go run ./cmd/smoke dev-api` | API read paths or seeded billing API behavior changed and local API is running. |
| Billing mutation smoke | `make smoke-dev-billing` | `go run ./cmd/smoke dev-billing` | Checkout, wallet, payment, order finalization, job creation, provisioning, or service activation changed. |
| Top-up review smoke | `make smoke-dev-topup-review` | `go run ./cmd/smoke dev-topup-review` | Top-up request create, approve, reject, ledger credit, or top-up audit path changed. |
| Target auth/RBAC smoke | `make smoke-dev-target-auth-rbac` | `go run ./cmd/smoke dev-target-auth-rbac` | Target auth session, 2FA gate, RBAC denial, or cross-tenant denial evidence changed. |
| Target credential reveal smoke | `make smoke-dev-target-credential-reveal` | `go run ./cmd/smoke dev-target-credential-reveal` | Target credential reveal, no-store response, reveal audit, redaction, or reveal rate-limit evidence changed. |
| Target finance reconciliation smoke | `make smoke-dev-target-finance-reconciliation` | `go run ./cmd/smoke dev-target-finance-reconciliation` | Target payment reconciliation, daily reconciliation, wallet/ledger balance report, or finance launch evidence changed. |
| Cloudmini idempotency evidence smoke | n/a | `go run ./cmd/smoke cloudmini-idempotency-evidence` | Owner-approved non-production Cloudmini duplicate-create or timeout-after-send evidence collection. |
| Cloudmini error evidence smoke | n/a | `go run ./cmd/smoke cloudmini-error-evidence` | Owner-approved non-production Cloudmini redacted provider error evidence collection. |
| Dev wallet projection repair | n/a | `APP_ENV=dev BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED=yes scripts/dev_wallet_projection_repair.sh` | Approved non-production wallet projection drift must be repaired from posted ledger source-of-truth before finance smoke evidence. |
| Full E2E launch gate | `make full-e2e-quality-gate` | `bash scripts/full_e2e_quality_gate.sh` | T204/T205/T243 launch-readiness validation on an approved local/dev database, including T206 renewal coverage. |
| Provider sandbox contract | n/a | `go test ./internal/modules/provider -run SandboxContract` | Provider adapter behavior or provider sandbox readiness changed. |
| Whitespace check | n/a | `git diff --check` | Every PR before commit or review. |

## Matrix By Change Type

| Change type | Required local validation | Add when touched |
| --- | --- | --- |
| Docs only | `git diff --check` | `go run ./cmd/taskguard` for task docs; `go run ./cmd/contractguard` or `go run ./cmd/errorcodeguard` if API/error docs changed. |
| Task board or task batch | `go run ./cmd/taskguard`, `git diff --check` | None unless code changed too. |
| Backend API or service | `make fmt`, `make test`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `git diff --check` | `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, smoke command matching the flow. |
| Shared Go platform | `make fmt`, `make test`, Go build command above, `git diff --check` | API and error guards if HTTP response, middleware, or routing changed. |
| DB migration or seed | `go run ./cmd/migrate validate`, `go test ./cmd/migrate ./internal/platform/db ./internal/seed`, `git diff --check` | `go run ./cmd/smoke dev-db` on local DB; `go run ./cmd/smoke dev-billing` for billing seed flows. |
| Provider adapter | `go test ./internal/modules/provider`, `go test ./internal/modules/provider -run SandboxContract`, `make test`, `git diff --check` | Billing smoke if provisioning outcome changes; update provider sandbox docs. |
| Frontend UI or app shell | `cd frontend && npm run check:sensitive-text`, `npm run lint`, `npm run build`, `git diff --check` | `npm run smoke:admin` for admin screens/nav/API mapping; `npm audit --omit=dev` for dependency changes. |
| Frontend CI or standalone smoke | `cd frontend && npm run check:sensitive-text`, `npm run lint`, `npm run build`, `npm run smoke:admin:ci`, `git diff --check` | `npx playwright install chromium` when local browser runtime is missing. |
| Full-stack billing flow | Backend API/service commands, frontend commands, `go run ./cmd/smoke dev-billing`, `git diff --check` | `go run ./cmd/smoke dev-db` when DB setup changed. |
| Top-up review flow | `make fmt`, `make test`, `make build`, `go run ./cmd/smoke dev-topup-review`, `git diff --check` | `go run ./cmd/contractguard` and `go run ./cmd/errorcodeguard` when API route/error behavior changes. |
| Target auth/RBAC flow | `make fmt`, `make test`, `make build`, `go run ./cmd/smoke dev-target-auth-rbac`, `git diff --check` | `go run ./cmd/contractguard` and `go run ./cmd/errorcodeguard` when API route/error behavior changes. |
| CI or workflow | Local equivalent of the changed job plus `git diff --check` | Run the exact frontend/backend commands the workflow invokes. |

## Smoke Prerequisites

`dev-db` requires a local or approved sandbox PostgreSQL database and `DB_DSN`.

`dev-api` requires the API to be running against the same seeded database. Use `-base-url` when it is not on `http://localhost:8080`.

`dev-billing` requires the same database for the API, smoke command, and local fake-provider worker. It must never point at production.

`dev-topup-review` requires the API and database to point at the same approved dev/test environment. It creates two top-up requests, approves one, rejects one, verifies wallet ledger/audit behavior, and must never point at production.

`dev-target-auth-rbac` requires the API and database to point at the same approved dev/test environment. It creates dev/test auth sessions, verifies cookie-only client access, 2FA admin blocking, invalid session denial, missing actor denial, cross-tenant denial, and RBAC permission denial. It must never point at production.

`dev-target-credential-reveal` requires the API and database to point at the same approved dev/test environment with `ENCRYPTION_KEY` available to the smoke runner. It creates or refreshes one encrypted dev/test credential fixture for the seeded demo service, logs in as the seeded client, reveals the fixture through the client API, verifies no-store headers, reveal metadata, rate-limit state, and redacted audit evidence. It must never point at production or real customer data.

`dev-target-finance-reconciliation` requires the API and database to point at the same approved dev/test environment. It is read-only: it selects an existing posted wallet payment, verifies payment reconciliation list/detail and daily reconciliation API evidence, and verifies database counters and wallet balance projection did not change. It must never point at production or real customer data.

`cloudmini-idempotency-evidence` requires owner-approved non-production Cloudmini credentials and explicit guardrail env. It calls Cloudmini mutating routes for exactly one scenario: `duplicate-create` uses two create attempts with the same idempotency key and must clean up one redacted resource reference, while `timeout-after-send` uses one create attempt, expects `PROVIDER_TIMEOUT_REQUEST_KNOWN`, and must clean up the created resource. It writes raw provider cleanup references only to `CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH`, which must be an absolute path outside the repo, and stdout intentionally excludes raw DSNs, tokens, provider IDs, provider payloads, and proxy credentials. It must never point at production or real customer data.

`cloudmini-error-evidence` requires owner-approved non-production Cloudmini credentials and explicit approval env. It captures redacted error metadata only: auth missing, invalid auth, not-found, optional malformed-create validation, optional permission-denied evidence, optional out-of-capacity evidence, optional rate-limit fixture evidence, and optional provider 5xx fixture evidence. The malformed-create validation uses an invalid JSON body and must be enabled with `CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE=yes`, `CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED=yes`, and `CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS=1`. Permission-denied evidence creates one temporary low-scope provider API key, calls a `proxy_crud` read route, revokes the temporary key, and must be enabled with `CLOUDMINI_ERROR_EVIDENCE_ALLOW_PERMISSION_DENIED=yes`, `CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED=yes`, and `CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE=1`. Out-of-capacity evidence uses one exhausted-group reservation probe and must be enabled with `CLOUDMINI_ERROR_EVIDENCE_ALLOW_OUT_OF_CAPACITY=yes`, `CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_APPROVED=yes`, `CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_MAX_RESERVATIONS=1`, `CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_KIND=ipv4_dc|residential`, and `CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS` between `1` and `60`. Rate-limit fixture evidence must use a provider-owned side-effect-free GET fixture path under `/api/v3/` that includes `fixture` and `rate`; it must be enabled with `CLOUDMINI_ERROR_EVIDENCE_ALLOW_RATE_LIMITED=yes`, `CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_APPROVED=yes`, `CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_MAX_REQUESTS=1`, and `CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH=/api/v3/...`. Provider 5xx fixture evidence must use a provider-owned side-effect-free GET fixture path under `/api/v3/` that includes `fixture`, `internal`, and `error`; it must be enabled with `CLOUDMINI_ERROR_EVIDENCE_ALLOW_PROVIDER_5XX=yes`, `CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_APPROVED=yes`, `CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_MAX_REQUESTS=1`, and `CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_FIXTURE_PATH=/api/v3/...`. It must not print raw response bodies, tokens, provider IDs, provider payloads, proxy credentials, cookies, or file contents. It must never point at production or real customer data.

`scripts/dev_wallet_projection_repair.sh` is a non-production repair tool for approved dev/test projection drift only. It fails closed without `APP_ENV`, `DB_DSN`, and `BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED=yes`, updates wallet projection from posted ledger source-of-truth, writes an audit row, and must never point at production or real customer data.

`full_e2e_quality_gate.sh` requires a Git worktree for `git diff --check` by default. On a deploy-copy target with no `.git`, set `BILLING_E2E_SKIP_GIT_DIFF_CHECK=1` only if `git diff --check` is run separately on the task branch and recorded in the same evidence packet.

Frontend browser smoke uses mock/intercepted data and does not need a backend or provider account.

## PR Validation Note Template

Use a short list, not a paragraph:

```text
Validation:
- make test
- go run ./cmd/contractguard
- cd frontend && npm run build
- cd frontend && npm run smoke:admin:ci
- git diff --check
```

If a command was not run, write the reason and the risk.

## Related Docs

- Git workflow: `docs/05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
- Testing strategy: `docs/05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`
- Database workflow: `docs/05_development_standards/52_Database_Migration_Seed_Data_Workflow.md`
- Frontend standard: `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`
- Local runbook: `docs/05_development_standards/55_Local_Development_Runbook.md`
- API contract guard: `docs/05_development_standards/59_API_Contract_Drift_Guard.md`
- Provider sandbox: `docs/05_development_standards/60_Provider_Sandbox_Contract_Checklist.md`
- Task board guard: `docs/05_development_standards/61_Task_Board_Consistency_Guard.md`
- API error guard: `docs/05_development_standards/62_API_Error_Code_Drift_Guard.md`
- Public ID policy: `docs/05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md`
