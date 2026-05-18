# T261 - Cloudmini provider 5xx runtime evidence

Status: REVIEW
Owner: Codex
Branch: codex/t261-cloudmini-internal-error-runtime
Provider Branch: codex/v3-internal-error-fixture
Provider PR: https://github.com/Chinsusu/proxy-cloudmini/pull/10
PR: https://github.com/Chinsusu/Billing-V2/pull/554
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add and run a safe Cloudmini V3 provider `INTERNAL_ERROR` fixture so Billing can record redacted provider 5xx runtime evidence without breaking the shared provider service.

## Scope

- Add a side-effect-free provider V3 fixture that returns HTTP `500` with provider code `INTERNAL_ERROR`.
- Add guarded Billing smoke support for exactly one provider 5xx fixture request.
- Deploy or verify the provider fixture on the approved dev manager only long enough to collect evidence.
- Record redacted runtime evidence if the run passes.
- Do not create/delete proxies, do not break the provider database/service, and do not print raw provider payloads, provider credentials, tokens, cookies, DSNs, proxy credentials, or raw provider IDs.

## Acceptance Criteria

- Provider fixture is disabled unless `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes` and requires a specific fixture header.
- Runtime evidence closes the provider 5xx case with HTTP `500`, provider code `INTERNAL_ERROR`, normalized `PROVIDER_TEMPORARY_ERROR`, retry safety `safe_retry`, and `mutating_routes_called=false`.
- Docs 66, 69, 70, and 77 reflect the new evidence or a concrete blocker.
- Required provider and Billing checks pass before PR.

## Notes

- Provider 5xx evidence must use a fixture. Do not force a real storage/service failure.
- Cancel/delete rejected remains separate follow-up work unless a safe fixture exists without expanding this task.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-18: Added Cloudmini provider `GET /api/v3/error-fixtures/internal-error` fixture in provider branch `codex/v3-internal-error-fixture`; targeted provider tests/builds passed. Full provider `go test -buildvcs=false ./...` remains blocked by pre-existing unrelated compile failures in root `test_list.go` and `cmd/test_ssh`.
- 2026-05-18: Opened, self-reviewed, and merged provider PR #10.
- 2026-05-18: Deployed the provider PR #10 manager binary to the approved dev manager, temporarily enabled `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes`, ran one guarded provider 5xx smoke, and removed the fixture env afterward. Post-run public/local capabilities checks returned HTTP `401` JSON and `vpm-manager` remained active.
- 2026-05-18: Billing smoke returned PASS for `provider_5xx_fixture`: HTTP `500`, provider code `INTERNAL_ERROR`, normalized `PROVIDER_TEMPORARY_ERROR`, retry `safe_retry`, `mutating_routes_called=false`, one fixture request, and no raw secrets/provider payloads printed.
- 2026-05-18: Billing validation passed: `GOFLAGS=-buildvcs=false make fmt`, `GOFLAGS=-buildvcs=false go test ./cmd/smoke -run 'CloudminiErrorEvidence'`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false make build`, `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff secret scan, and touched-file line count check. Plain `make fmt` was attempted first but failed in the temporary worktree due Go VCS stamping, so it was rerun with `GOFLAGS=-buildvcs=false`.
- 2026-05-18: Opened Billing PR #554 for review.
