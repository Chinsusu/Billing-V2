# T251 - Provider test stability

Status: DONE
Owner: Codex
Branch: codex/t251-provider-test-stability
PR: https://github.com/Chinsusu/Billing-V2/pull/534
Risk: provider provisioning tests, CI reliability
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Stabilize the Cloudmini provider hardening test that flaked in GitHub Actions during the T250 marker PR.

## Scope

- Fix only the flaky test timing behavior.
- Preserve the existing provider behavior and assertions.
- Do not skip, weaken, or delete the test.

## Acceptance Criteria

- The focused provider test passes repeatedly.
- `go test ./internal/modules/provider -count=1` passes.
- `make test` passes.
- `go run ./cmd/taskguard` and `git diff --check` pass.

## Notes

- Do not change provider runtime behavior in this task.

## Agent Log

- 2026-05-18: Task created and claimed by Codex after PR #533 first CI run failed on `TestCloudminiV3AdapterProvisionRequiresUsableProxyStatus`; rerun passed, indicating timing flake rather than marker diff failure.
- 2026-05-18: Increased the not-usable-status test poll budget from 3ms to 75ms with a 1ms interval so operation polling has time to return on CI while still asserting partial-success for a non-usable `creating` proxy status.
- 2026-05-18: Validation passed: focused provider test with `-count=20`, `go test ./internal/modules/provider -count=1`, `make test`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-05-18: Opened PR #534 for review.
- 2026-05-18: PR #534 merged into `main`; marking task done.
