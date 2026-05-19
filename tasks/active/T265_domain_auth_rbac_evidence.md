# T265 - Domain auth/RBAC evidence

Status: DONE
Owner: Codex
Branch: codex/t265-domain-auth-rbac-evidence
PR: #562; unblocked by T266
Risk: auth, RBAC, tenant isolation, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-19

## Summary

Run and record target-domain auth/session/RBAC evidence so launch docs no longer rely only on local target API auth checks for the selected bounded pilot.

## Scope

- Run the existing `dev-target-auth-rbac` smoke against approved non-production domain inputs.
- Record only redacted evidence: status categories, pass/fail checks, and no raw cookies/session tokens/passwords/DSNs.
- Update docs 69 and 70 with the domain evidence result or a concrete blocker.
- Do not use production customer data.
- Do not change auth/RBAC behavior unless the smoke exposes a code defect.

## Acceptance Criteria

- Domain auth/RBAC smoke evidence proves cookie-only client access, admin 2FA gate, invalid session denial, missing actor denial, cross-tenant denial, and low-permission RBAC denial on approved domain input.
- Evidence states no provider or money mutation routes were called.
- Task board stays consistent.
- Required target-auth/docs checks pass or blockers are recorded.

## Notes

- Use protected target server secrets without printing file contents or raw environment values.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-18: Domain route preflight returned HTTP `500` for `https://billing.resvn.net/`, `/healthz`, and `/backend/healthz`; local `http://localhost:3000/` also returned HTTP `500`. Current host has `cloudflared` active but no `billing-api` or `billing-worker` systemd unit, and no protected target env file was present under `/opt/Billing`.
- 2026-05-18: `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 20s dev-target-auth-rbac` failed before auth assertions with `target auth login expected HTTP 200, got 500`. No raw cookies, session tokens, passwords, DSNs, provider payloads, or credentials were printed.
- 2026-05-18: Updated docs 69 and 70 to record the domain auth/RBAC blocker. Task is BLOCKED until the approved domain/backend route returns a usable API response and the domain smoke can be rerun.
- 2026-05-19: T266 recovered the target dev/staging-equivalent route by migrating and seeding the empty `billing_smoke` database, starting the Billing API on `127.0.0.1:8080`, and serving the frontend on `3000` with `/backend` proxying to that API. Runtime env was loaded from root-only `/run` files and no raw DSN, token, cookie, password, provider payload, or credential was printed.
- 2026-05-19: `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac` passed: client session cookie-only access, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three low-permission RBAC denials. The smoke reported no provider or money mutation routes called.
