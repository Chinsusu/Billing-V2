# T262 - Cloudmini cancel/delete rejected runtime evidence

Status: DONE
Owner: Codex
Branch: codex/t262-cloudmini-cancel-delete-runtime
Provider Branch: codex/v3-delete-rejected-fixture
Provider PR: https://github.com/Chinsusu/proxy-cloudmini/pull/11
PR: https://github.com/Chinsusu/Billing-V2/pull/556
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add and run safe Cloudmini V3 cancel/delete rejected evidence so Billing can close the remaining provider-controlled runtime error case without deleting a sellable resource.

## Scope

- Inspect the current Cloudmini V3 delete/action failure shape and Billing error mapping.
- Add provider-side fixture support only if no existing side-effect-free route can produce the required error safely.
- Add guarded Billing smoke support for exactly one cancel/delete rejected fixture request.
- Deploy or verify the provider fixture on the approved dev manager only long enough to collect evidence.
- Record redacted runtime evidence if the run passes.
- Do not create/delete proxies, do not break provider state, and do not print raw provider payloads, provider credentials, tokens, cookies, DSNs, proxy credentials, or raw provider IDs.

## Acceptance Criteria

- Evidence closes the cancel/delete rejected case with a stable provider error shape, normalized Billing provider code, retry safety, and `mutating_routes_called=false` unless the approved fixture contract explicitly requires otherwise.
- Provider fixture is disabled unless `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes` and requires a specific fixture header.
- Docs 66, 69, 70, and 77 reflect the new evidence or a concrete blocker.
- Required provider and Billing checks pass before PR.

## Notes

- Prefer a side-effect-free fixture. Do not create or delete a real proxy solely to force this error path.
- If provider semantics require owner decision between delete rejection and action rejection, document the tradeoff and choose the safer fixture path.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-18: Added Cloudmini provider `GET /api/v3/error-fixtures/delete-rejected` fixture in provider branch `codex/v3-delete-rejected-fixture`; targeted provider tests/builds passed. Full provider `go test -buildvcs=false ./...` remains blocked by pre-existing unrelated compile failures in root `test_list.go` and `cmd/test_ssh`.
- 2026-05-18: Opened, self-reviewed, and merged provider PR #11.
- 2026-05-18: Deployed provider PR #11 manager binary to the approved dev manager, temporarily enabled `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes`, ran one guarded cancel/delete rejected smoke, and removed the fixture env afterward. Post-run public/local capabilities checks returned HTTP `401` JSON and `vpm-manager` remained active.
- 2026-05-18: Billing smoke returned PASS for `cancel_delete_rejected_fixture`: HTTP `200`, provider failed operation code `DELETE_FAILED`, normalized `PROVIDER_PARTIAL_SUCCESS`, retry `manual_review_required`, `mutating_routes_called=false`, one fixture request, and no raw secrets/provider payloads printed.
- 2026-05-18: Billing validation passed: `GOFLAGS=-buildvcs=false make fmt`, `GOFLAGS=-buildvcs=false go test ./cmd/smoke ./internal/modules/provider -run 'CloudminiErrorEvidence|CloudminiV3AdapterTerminateFailedOperation'`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false make build`, `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff secret scan, and touched-file line count check.
- 2026-05-18: Opened Billing PR #556 and moved task to REVIEW.
- 2026-05-18: Billing PR #556 merged into `main`; task marked DONE.
