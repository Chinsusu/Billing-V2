# T266 - Domain backend route recovery for auth/RBAC evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t266-domain-backend-route-fix
PR: -
Risk: auth, RBAC, tenant isolation, deploy/runtime, launch-readiness evidence
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Recover the approved target domain/backend route that blocked T265, then rerun the domain auth/RBAC smoke and update the launch evidence.

## Scope

- Diagnose why `https://billing.resvn.net/backend` returns HTTP `500` before auth assertions.
- Recover the target runtime or route without printing secrets or file contents.
- Rerun `dev-target-auth-rbac` against approved domain input if the route becomes usable.
- Update T265 and launch docs with either PASS evidence or a concrete remaining blocker.
- Do not use production customer data.

## Acceptance Criteria

- Domain health/auth route is usable, or the remaining blocker is documented with exact redacted evidence.
- If usable, domain auth/RBAC smoke proves cookie-only client access, admin 2FA gate, invalid session denial, missing actor denial, cross-tenant denial, and low-permission RBAC denial.
- Evidence states no provider or money mutation routes were called.
- Required checks pass.

## Notes

- Do not print `.env`, secret-store files, cookies, tokens, passwords, DSNs, provider payloads, or credential contents.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-19: Found `cloudflared` active and forwarding the domain to local frontend, but no Billing API was listening on `8080`; the frontend process on `3000` had no usable backend route and returned HTTP `500`.
- 2026-05-19: Recovered the target dev/staging-equivalent database path by creating required `pgcrypto` extension with the local Postgres superuser, applying 25 migrations to the empty `billing_smoke` database, running dev seed twice, and passing `dev-db` smoke with 20 checks. No DB DSN or password was printed.
- 2026-05-19: Started the Billing API on `127.0.0.1:8080` and frontend on `0.0.0.0:3000` using root-only `/run/billing-t266-*.env` files so secrets stayed out of process arguments and logs. Verified `http://127.0.0.1:8080/healthz`, `http://127.0.0.1:3000/`, `http://127.0.0.1:3000/backend/healthz`, `https://billing.resvn.net/`, and `https://billing.resvn.net/backend/healthz` all returned HTTP `200`.
- 2026-05-19: Reran `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac`; result PASS for cookie-only client session, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three RBAC denials. Smoke output states no provider or money mutation routes were called.
