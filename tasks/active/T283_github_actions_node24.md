# T283 - GitHub Actions Node 24 runtime

Status: DONE
Owner: Codex
Branch: codex/t283-actions-node24
PR: https://github.com/Chinsusu/Billing-V2/pull/598
Risk: CI workflow
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Upgrade official GitHub Actions wrappers used by CI to Node 24-compatible major versions so CI no longer emits Node 20 action-runtime deprecation warnings.

## Scope

- In scope: update `.github/workflows/ci.yml` action versions for `actions/checkout`, `actions/setup-go`, and `actions/setup-node`.
- In scope: keep frontend `node-version: "20"` unchanged because the warning is about action runtime, not the app toolchain runtime.
- Out of scope: frontend dependency upgrades, Node toolchain migration, Go version changes, or CI job behavior changes.

## Acceptance Criteria

- CI workflow uses Node 24-compatible official action versions.
- Local workflow YAML review, `taskguard`, `git diff --check`, and required local equivalents pass.
- GitHub CI passes without the previous Node 20 action-runtime deprecation annotations for official actions.

## Notes

- Official action release notes indicate `actions/checkout@v5`, `actions/setup-go@v6`, and `actions/setup-node@v5` use Node 24 and require runner v2.327.1 or newer.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t283-actions-node24`.
- 2026-05-21: Updated CI action wrappers to `actions/checkout@v5`, `actions/setup-go@v6`, and `actions/setup-node@v5`; kept frontend `node-version: "20"` unchanged.
- 2026-05-21: Validation passed: `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, workflow action version check, line-count check, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false make contract-guard`, `GOFLAGS=-buildvcs=false make error-code-guard`, `GOFLAGS=-buildvcs=false make build`, `npm ci`, `npm audit --omit=dev`, `npm run check:sensitive-text`, `npm run lint`, `npm run build`, and `npm run smoke:admin:ci`.
- 2026-05-21: Opened PR https://github.com/Chinsusu/Billing-V2/pull/598 and moved task to REVIEW.
- 2026-05-21: PR https://github.com/Chinsusu/Billing-V2/pull/598 merged after GitHub checks passed and log scan found no Node 20 action-runtime deprecation annotation; task marked DONE.
