# T122 - Admin frontend public ID filters

Status: TODO
Owner: -
Branch: codex/t122-admin-frontend-public-id-filters
PR: -
Risk: frontend/API
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Update admin frontend filters and labels to prefer numeric public IDs after backend filter support is ready.

## Scope

- Review admin filters for invoices, transactions, logs, services, top-ups, provisioning, and accounts.
- Replace visible raw field placeholders such as account or actor backend references with public ID wording and query usage where supported.
- Reuse the shared frontend API view-model boundary from T117 where it fits.
- Keep layout changes small and files under 500 lines.
- Do not add fake backend routes in frontend code.

## Acceptance Criteria

- Admin filter labels and placeholders use public ID language.
- Frontend queries use public display ID filters where backend support exists.
- Sensitive backend references do not appear in user-facing labels.
- Frontend lint, build, sensitive-text guard, and admin smoke pass.

## Notes

- This task should start after T121 or explicitly document any backend filter still missing.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
