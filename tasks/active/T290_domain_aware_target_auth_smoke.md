# T290 - Domain-aware target auth smoke

Status: REVIEW
Owner: Codex
Branch: codex/t290-domain-aware-auth-smoke
PR: https://github.com/Chinsusu/Billing-V2/pull/612
Risk: auth/RBAC, tenant domain resolution, credential handling, launch evidence
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Make `dev-target-auth-rbac` support public-domain target checks where client and platform admin login must use different tenant domains.

## Scope

- In scope: add separate target auth client/admin base URL flags and env overrides.
- In scope: preserve existing `-base-url` behavior for local API smoke.
- In scope: update tests and validation docs.
- Out of scope: API behavior changes, production approval, production customer data, real provider provisioning, or committing credential values.

## Acceptance Criteria

- `dev-target-auth-rbac` can run with one local `-base-url` as before.
- `dev-target-auth-rbac` can run with separate client/admin public base URLs.
- The command continues to exclude passwords, cookies, session tokens, DSNs, provider payloads, and credentials from output.
- Unit tests cover separate client/admin base URL routing.
- Required local validation passes before PR: focused smoke tests, repo Go tests/build, `taskguard`, `git diff --check`, touched-file line-count check, added-line secret scan, and added-line UUID scan for docs/task changes.

## Notes

- Public-domain auth is domain-first: client login belongs on the client/reseller domain, platform admin login belongs on the admin/platform domain.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t290-domain-aware-auth-smoke`.
- 2026-05-21: Added split client/admin base URL handling for `dev-target-auth-rbac`, unit coverage for separate base URLs, and docs/env references for protected overrides.
- 2026-05-21: Validation passed: `APP_ENV=dev GOFLAGS=-buildvcs=false go test ./cmd/smoke`, `GOFLAGS=-buildvcs=false make fmt`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false make build`, `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret-pattern scan, added-line docs/task UUID scan, and public-domain `dev-target-auth-rbac` with separate client/admin base URLs.
- 2026-05-21: Opened PR #612 and moved task to `REVIEW`.
