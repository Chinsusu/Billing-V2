# T267 - Target runtime systemd services

Status: REVIEW
Owner: Codex
Branch: codex/t267-target-runtime-services
PR: https://github.com/Chinsusu/Billing-V2/pull/565
Risk: deploy/runtime, auth, secrets, launch-readiness evidence
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Promote the selected target dev/staging-equivalent runtime from temporary T266 processes to protected systemd services, then record redacted evidence.

## Scope

- Deploy latest `origin/main` to the canonical `/opt/Billing` target path without printing secrets.
- Build the API binary and frontend production bundle.
- Run API and frontend through systemd units with protected environment files outside git.
- Verify domain root, `/backend/healthz`, and `dev-target-auth-rbac` through `https://billing.resvn.net/backend`.
- Update launch evidence docs with the long-lived runtime proof.
- Do not change product behavior, money flows, provider provisioning, or production customer data.

## Acceptance Criteria

- `billing-api` and `billing-frontend` are active systemd services.
- Service command lines do not contain raw DSNs, tokens, cookies, passwords, provider payloads, or credentials.
- Protected env file metadata is recorded without contents.
- Domain health and auth/RBAC smoke pass.
- Required repo checks and CI pass.

## Notes

- Runtime env values must stay in protected local files and must not be committed or printed.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-19: Updated canonical `/opt/Billing` to detached `origin/main` at merge commit `87a6584`, built `bin/api`, ran `npm --prefix /opt/Billing/frontend ci`, and built the frontend production bundle. `npm ci` reported one moderate advisory; `npm audit --audit-level=high` remains the required pass threshold for this repo.
- 2026-05-19: Created local service user `billing-svc`, promoted protected runtime env files to `/etc/billing/secrets/billing-api.env` and `/etc/billing/secrets/billing-frontend.env` with mode `600` owner `root:root`, and installed `billing-api.service` plus `billing-frontend.service`.
- 2026-05-19: Started API from `/opt/Billing/bin/api` and frontend from `/opt/Billing/frontend/.next/standalone/server.js` under systemd. Both services are active and enabled, run as `billing-svc`, and command lines contain no raw DSN, token, cookie, password, provider payload, or credential.
- 2026-05-19: Verified HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`.
- 2026-05-19: Reran `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac`; result PASS for cookie-only client session, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three RBAC denials. Smoke output states no provider or money mutation routes were called.
- 2026-05-19: Opened PR #565 for review.
