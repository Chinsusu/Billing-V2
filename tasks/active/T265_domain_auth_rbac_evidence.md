# T265 - Domain auth/RBAC evidence

Status: BLOCKED
Owner: Codex
Branch: codex/t265-domain-auth-rbac-evidence
PR: -
Risk: auth, RBAC, tenant isolation, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

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
