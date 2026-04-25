# T108 - Smoke RBAC negative checks

Status: TODO
Owner: -
Branch: codex/t108-smoke-rbac-negative-checks
PR: -
Risk: backend/RBAC/QA
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Extend smoke coverage with safe negative RBAC checks for admin catalog readiness, job reads, and recovery actions.

## Scope

- Add smoke checks that verify denied access returns stable auth/permission errors.
- Cover at least one admin catalog readiness route and one job route.
- Do not depend on production credentials or production DSNs.
- Keep response assertions focused on envelope, code, status, and redaction.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Smoke command fails if a low-permission actor can read provider readiness or job data unexpectedly.
- Smoke command fails if a denied response leaks provider credentials, raw payloads, or internal job payload names.
- `go test ./...` and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker` pass.
- Frontend validation commands pass if frontend files are touched.

## Notes

- Follow the permission map in `docs/05_development_standards/56_Billing_API_Operational_Reference.md`.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
