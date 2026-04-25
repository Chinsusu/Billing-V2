# T120 - Public display ID API policy

Status: REVIEW
Owner: Codex
Branch: codex/t120-public-display-id-api-policy
PR: https://github.com/Chinsusu/Billing-V2/pull/271
Risk: API/docs
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Document the API rule for numeric public display IDs versus backend UUID references before adding more route behavior.

## Scope

- Define which fields are safe for frontend labels for accounts, tenants, providers, orders, services, invoices, transactions, jobs, top-ups, and audit logs.
- Define when backend UUIDs may be used for API paths or actions but must not be shown as UI labels.
- Add naming guidance for query filters that use public display IDs.
- Link the policy from frontend and API operational docs.
- Keep edited docs under 500 lines.

## Acceptance Criteria

- Agents can tell which ID to display in UI and PR descriptions.
- API docs distinguish action/internal references from user-facing labels.
- Related frontend/API docs link to the policy.
- Docs validation commands pass.

## Notes

- This task should keep terminology simple: "public ID" for numeric display IDs, "backend reference" for UUIDs.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
- 2026-04-25: Codex claimed the task; documenting public ID versus backend reference rules before backend/filter work.
- 2026-04-25: Opened PR https://github.com/Chinsusu/Billing-V2/pull/271 after task board, contract, error code, and diff checks passed.
