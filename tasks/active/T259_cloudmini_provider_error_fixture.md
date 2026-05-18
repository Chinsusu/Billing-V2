# T259 - Cloudmini provider error fixture path

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t259-cloudmini-provider-error-fixture
Provider Branch: codex/v3-error-fixtures
Provider PR: https://github.com/Chinsusu/proxy-cloudmini/pull/9
PR: -
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add a safe Cloudmini provider V3 error fixture path for Billing evidence and collect redacted rate-limit runtime evidence if it can be deployed to non-production safely.

## Scope

- Work in a clean provider worktree from `origin/main`, not the dirty `/opt/proxy-cloudmini` checkout.
- Add a side-effect-free V3 fixture path that can return `RATE_LIMITED` without tripping the shared limiter.
- Keep the fixture disabled outside explicitly non-production/test configuration.
- If the fixture is deployed and reachable on the approved dev provider, run the guarded Billing smoke from T258 and record redacted runtime evidence.
- Do not create/delete proxies, do not change production servers, and do not print raw provider payloads, provider credentials, tokens, cookies, DSNs, proxy credentials, or raw provider IDs.

## Acceptance Criteria

- Provider fixture is side-effect-free and guarded for non-production/test use.
- Billing rate-limit evidence is either safely collected or remains blocked with a concrete deployment blocker.
- Relevant provider and Billing checks pass.
- PR/CI/merge workflow is followed for changed repos.

## Notes

- The provider checkout at `/opt/proxy-cloudmini` had pre-existing uncommitted changes; this task uses `/tmp/proxy-cloudmini-t259` to avoid modifying them.
- Provider 5xx and cancel/delete rejected fixtures are out of scope unless they share the same safe fixture mechanism without expanding risk.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from Billing `origin/main`; provider worktree created at `/tmp/proxy-cloudmini-t259` from provider `origin/main`.
- 2026-05-18: Added provider V3 rate-limit fixture route guarded by `VPM_BILLING_ERROR_FIXTURES_ENABLED` and `X-Cloudmini-Error-Fixture: rate_limited`.
- 2026-05-18: Updated Billing rate-limit runner to send the provider-required fixture header.
- 2026-05-18: Provider PR #9 merged with fixture code and runbook docs.
- 2026-05-18: Did not run live Billing rate-limit evidence because no local `vpm-manager`/Cloudmini manager service is running on this host; dev provider deployment to `https://cz.resvn.net/` still needs the merged provider code plus `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes`.
- 2026-05-18: Provider validation passed for PR #9: targeted API tests, manager/agent build with `-buildvcs=false`, and `git diff --check`; provider full `go test -buildvcs=false ./...` remains blocked by pre-existing root `test_list.go` and `cmd/test_ssh` compile issues.
- 2026-05-18: Billing validation passed: `make fmt`, `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`, `make test`, `make build`, `go run ./cmd/taskguard`, `git diff --check`, and targeted diff secret-value scan excluding fake test fixtures.
