# T124 - Admin read model related display IDs

Status: TODO
Owner: -
Branch: codex/t124-admin-read-model-related-display-ids
PR: -
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
