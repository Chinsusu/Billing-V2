# T231 - Cloudmini runtime activation preflight

Status: DONE
Owner: Codex
Branch: codex/t231-cloudmini-runtime-activation
PR: https://github.com/Chinsusu/Billing-V2/pull/495
Risk: provider/provisioning/lifecycle/credential/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Add and run a non-mutating worker runtime activation check that validates the Cloudmini provider registry can boot from the approved test-server credential/config path without claiming jobs or calling Cloudmini mutating routes.

## Scope

- Add a worker command that builds the provider registry from runtime env and reports a redacted activation summary.
- Ensure the command does not claim jobs, run lifecycle transitions, call provider create/delete/action routes, or print secrets/provider-private identifiers.
- Deploy the command to the approved test server and run it with the local Cloudmini dev credential file.
- Record redacted evidence and remaining blockers in the Cloudmini launch/provider docs.

## Acceptance Criteria

- `cmd/worker` has a non-mutating provider registry activation check with tests.
- Target server evidence shows the check passes with `PROVIDER_DEFAULT_MODE=cloudmini_v3` and the approved local credential path.
- Evidence states no provider API mutation, Billing provisioning mutation, raw DSN, token, raw group ID, raw provider payload, or proxy credential was printed.
- Broader pilot remains blocked until a separate owner-approved mutating/lifecycle cleanup window exists.

## Notes

- This task is not approval to run a provisioning or lifecycle cleanup job.
- Keep the always-on worker service in fake mode unless a later task explicitly approves a bounded Cloudmini worker window.

## Agent Log

- 2026-05-17: Task created and claimed by Codex from latest `origin/main` on branch `codex/t231-cloudmini-runtime-activation`.
- 2026-05-17: Added non-mutating `cmd/worker provider-registry-check` command and tests. The command builds the provider registry from env without DB access, job claims, or provider API calls, and prints only redacted activation evidence.
- 2026-05-17: Deployed to the approved test server and ran the check with `APP_ENV=dev`, `PROVIDER_DEFAULT_MODE=cloudmini_v3`, `.env.dev`, and `/opt/cred-cloudmini-dev.env`. Result passed with real Cloudmini adapter, one source mapping, no provider API calls, no mutating routes, no job claims, and no secrets printed.
- 2026-05-17: Validation passed: target `go test ./cmd/worker`; target `go build -o bin/worker ./cmd/worker`; local `go test ./...`; local `go build -buildvcs=false -o bin/worker ./cmd/worker`; `go run ./cmd/taskguard`; `git diff --check`; changed-file secret scan found only existing documentation text about avoiding `?token=` query credentials. Local plain `go build -o bin/worker ./cmd/worker` was blocked by a stray `/tmp/.git` VCS stamping issue, so compile was verified with `-buildvcs=false` locally and plain build on the target server.
- 2026-05-17: Opened PR https://github.com/Chinsusu/Billing-V2/pull/495.
- 2026-05-17: PR https://github.com/Chinsusu/Billing-V2/pull/495 merged after CI passed.
