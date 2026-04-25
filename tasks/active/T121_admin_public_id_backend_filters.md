# T121 - Admin public ID backend filters

Status: TODO
Owner: -
Branch: codex/t121-admin-public-id-backend-filters
PR: -
Risk: API/DB
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add or normalize backend admin filters so operators can search common records by numeric public IDs instead of UUIDs.

## Scope

- Review existing admin list filters for accounts, orders, services, invoices, transactions, jobs, top-ups, providers, and audit logs.
- Add focused public display ID filters where missing and where the data model already has `display_id`.
- Keep raw UUID filters available only when needed for internal action paths or precise support diagnostics.
- Update tests and API operational docs for new filter names.
- Avoid broad read-model rewrites in this task.

## Acceptance Criteria

- Important admin list endpoints accept documented numeric public ID filters.
- Existing UUID-based action paths still work.
- Tests cover at least one success and one invalid public ID filter case.
- Contract and error-code guards pass when relevant.

## Notes

- Prefer query names that are obvious to operators, such as `display_id`, `order_display_id`, or `account_display_id`.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
