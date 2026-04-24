# T087 - Provisioning job read API

Status: TODO
Owner: -
Branch: codex/t087-provisioning-job-read-api
PR: -
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expose read APIs for provisioning jobs so admin/reseller screens and operations can inspect paid order fulfillment without direct database access.

## Scope

- Work mainly in `internal/modules/jobs/**/*`, `cmd/api/**/*`, and API docs.
- Add list/detail reads for `jobs` with tenant scoping and permission middleware.
- Include filters for `display_id`, `job_type`, `status`, `reference_type`, `reference_id`, `source_id`, and `limit`.
- Keep responses numeric-display-ID first and avoid exposing raw provider secrets.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin can list and get jobs for the effective tenant.
- Reseller can list and get jobs for the effective tenant.
- Responses include job display ID, type, reference, source, status, attempt counts, next attempt, last redacted error, manual review reason, and timestamps.
- API docs list the new routes and filters.
- Backend and frontend validation commands pass.

## Notes

- Start with read-only routes. Do not add retry/cancel/manual-review mutations in this task.
- This unlocks frontend fulfillment visibility without direct DB reads.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
