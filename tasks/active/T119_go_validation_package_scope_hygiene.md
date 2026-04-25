# T119 - Go validation package scope hygiene

Status: TODO
Owner: -
Branch: codex/t119-go-validation-package-scope-hygiene
PR: -
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
