# T110 - Frontend sensitive text guard

Status: REVIEW
Owner: Codex
Branch: codex/t110-frontend-sensitive-text-guard
PR: pending
Risk: frontend/security
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a frontend guard that prevents sensitive/internal backend field names from being rendered or introduced in admin UI copy and mock data.

## Scope

- Scan frontend source for high-risk field names such as `payload_json`, `capability_profile`, `provider_account_id`, `secret`, `raw_response`, and credential/token variants.
- Allow narrow exceptions only where type definitions or explicit redaction tests require them.
- Add an npm script or documented command.
- Keep output clear enough for agents to fix quickly.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Guard fails on newly introduced sensitive/internal UI strings outside approved allowlist locations.
- Existing `npm run lint`, `npm run build`, and backend validation pass.
- Documentation tells frontend agents to run the guard before PR.

## Notes

- This should complement, not replace, backend redaction tests.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
- 2026-04-25: Codex claimed the task; adding a static frontend sensitive text guard.
- 2026-04-25: Added `npm run check:sensitive-text`, CI coverage, narrow allowlist markers for browser redaction tests, and removed direct provider account field use from the admin provider screen. Validation passed: `npm run check:sensitive-text`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `npm run smoke:admin`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker ./cmd/contractguard`.
