# T288 - UAT consolidation and auth smoke credential override

Status: REVIEW
Owner: Codex
Branch: codex/t288-uat-consolidation-auth-smoke
PR: -
Risk: UAT, auth/RBAC, credential handling, launch evidence
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Consolidate client/reseller/admin UAT evidence and make the target auth/RBAC smoke command usable when selected target credentials differ from the original dev seed defaults.

## Scope

- In scope: create a redacted consolidated UAT evidence packet for T285/T286/T287.
- In scope: parameterize target auth/RBAC smoke client/admin credentials through protected environment variables with safe dev defaults.
- In scope: update relevant validation docs and tests.
- Out of scope: production approval, production customer data, real Cloudmini provisioning, human UAT sign-off, or printing/committing credential values.

## Acceptance Criteria

- Client, reseller, and admin UAT evidence is summarized in one packet with explicit PASS boundary and residual risk.
- `dev-target-auth-rbac` can use operator-provided client/admin email/password values from environment variables without exposing them in output.
- Existing dev seed defaults continue to work when no overrides are configured.
- Unit tests cover credential override request behavior without real secret values.
- Required local validation passes before PR: focused Go smoke tests, repo Go tests/build, `taskguard`, `git diff --check`, touched-file line-count check, and added-line secret scan.

## Notes

- This task does not run or approve production traffic.
- Do not print raw passwords, cookies, session tokens, provider payloads, DSNs, TOTP values, or credential payloads.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t288-uat-consolidation-auth-smoke`.
- 2026-05-21: Added protected env override support for target auth smoke credentials and updated credential reveal smoke to share the same client credential source.
- 2026-05-21: Added consolidated UAT evidence packet for client, reseller, and admin selected-environment evidence.
- 2026-05-21: Validation passed: `go test ./cmd/smoke`, `GOFLAGS=-buildvcs=false make fmt`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false make build`, `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret scan, and added-line UUID scan for docs/task changes.
- 2026-05-21: Did not run live `dev-target-auth-rbac`; protected target credential override variables are not configured in `/etc/billing/secrets/billing-api.env`, and fallback seed credentials are already known stale on the selected target.
- 2026-05-21: Moved task to REVIEW pending PR.
