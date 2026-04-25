# T119 - Go validation package scope hygiene

Status: DONE
Owner: Codex
Branch: codex/t119-go-validation-package-scope-hygiene
PR: https://github.com/Chinsusu/Billing-V2/pull/269
Risk: tooling/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Keep Go validation commands focused on Billing source packages instead of accidentally traversing frontend dependency folders.

## Scope

- Review current `go test ./...` and build command behavior with `frontend/node_modules` present.
- Add a small repo-local helper, Makefile target, or documented command path that excludes frontend dependency folders.
- Update CI or docs only if needed to make the intended command unambiguous.
- Keep the approach portable for Windows hosts without `make`.
- Do not change application runtime behavior.

## Acceptance Criteria

- Local Go validation no longer reports packages from `frontend/node_modules`.
- Makefile and documentation point to the same package-scope approach.
- CI still runs backend tests/builds successfully.
- Task guard and diff check pass.

## Notes

- This task exists because full Go validation currently sees a package under frontend dependencies on this workspace.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
- 2026-04-25: Codex claimed the task; adding a repo-scoped Go package selector for validation commands.
- 2026-04-25: Added `cmd/gopackages`, wired Makefile `fmt/test` through it, and updated validation docs; local PowerShell package-scoped tests passed without `frontend/node_modules`.
- 2026-04-25: Opened PR https://github.com/Chinsusu/Billing-V2/pull/269 for review and CI.
- 2026-04-25: PR https://github.com/Chinsusu/Billing-V2/pull/269 merged; T119 is done.
