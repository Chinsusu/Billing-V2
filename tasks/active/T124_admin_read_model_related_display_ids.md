# T124 - Admin read model related display IDs

Status: DONE
Owner: Codex
Branch: codex/t124-admin-read-model-related-display-ids
PR: https://github.com/Chinsusu/Billing-V2/pull/279
Risk: API/frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Expose safe related public IDs in admin read models so the UI does not need to show backend references for linked records.

## Scope

- Review admin invoice, transaction, service, provisioning job, top-up, and audit response models.
- Add small related display ID fields where the backend can join or derive them safely.
- Keep response changes backward compatible.
- Update frontend API types and view models for the new fields.
- Add tests for any new response fields.

## Acceptance Criteria

- Admin UI can show linked record public IDs for at least one high-value flow without UUID labels.
- API docs name the new related display ID fields.
- Existing clients remain compatible.
- Backend and frontend validation commands pass.

## Notes

- This is intentionally separate from T121 so filter behavior and response enrichment can be reviewed independently.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
- 2026-04-25: Codex claimed the task after T123 merged; reviewing backend read models and frontend view-model gaps for related public IDs.
- 2026-04-25: Added related display IDs to invoice, transaction, service, top-up, job, and audit read responses; updated admin UI view models to prefer public IDs.
- 2026-04-25: Validation passed: module Go tests, full `go test`, Go command build, frontend lint/build/sensitive-text/admin smoke, taskguard, and diff whitespace check.
- 2026-04-25: Opened PR https://github.com/Chinsusu/Billing-V2/pull/279 for review.
- 2026-04-25: PR https://github.com/Chinsusu/Billing-V2/pull/279 merged into `main`.
