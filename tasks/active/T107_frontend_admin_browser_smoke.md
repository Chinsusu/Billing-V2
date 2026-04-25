# T107 - Frontend admin browser smoke

Status: REVIEW
Owner: Codex
Branch: codex/t107-frontend-admin-browser-smoke
PR: https://github.com/Chinsusu/Billing-V2/pull/243
Risk: frontend/QA
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a repeatable browser smoke for the admin frontend shell so critical admin screens are checked after frontend or API-contract changes.

## Scope

- Cover admin navigation and at least these screens: overview, provisioning queue, provider readiness, top-up verification, and audit logs.
- Use local/mock-safe data or intercepted API responses; do not require production backend access.
- Verify that key display IDs and readiness/job status labels render.
- Verify that sensitive/internal field names such as `payload_json`, `capability_profile`, `provider_account_id`, `secret`, and `raw_response` do not appear in the rendered UI.
- Add a clear command or npm script that agents can run locally.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Browser smoke can be run from a clean checkout after `npm ci`.
- Smoke fails when a required screen cannot be reached or sensitive/internal text is visible.
- `npm run lint`, `npm run build`, and backend validation commands pass.
- Docs mention when agents should run the browser smoke.

## Notes

- Prefer a lightweight harness that does not require real provider credentials.
- If adding a browser dependency is too heavy for CI, keep CI optional and document the local command clearly.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
- 2026-04-25: Codex claimed the task on `codex/t107-frontend-admin-browser-smoke`.
- 2026-04-25: Added `npm run smoke:admin`, a Playwright browser smoke that starts Next locally, intercepts admin API calls, checks admin overview/provisioning/provider/top-up/audit screens, and blocks sensitive/internal text in rendered UI.
- 2026-04-25: Validation passed: `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `npm run smoke:admin`, `go test ./...`, and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`.
- 2026-04-25: Opened PR #243 and moved the task to review.
