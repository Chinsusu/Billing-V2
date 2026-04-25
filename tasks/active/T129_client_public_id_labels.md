# T129 - Client public ID labels

Status: REVIEW
Owner: Codex
Branch: codex/t129-client-public-id-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/289
Risk: frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Clean up client-facing billing screens so visible labels use public numeric IDs and do not fall back to backend references.

## Scope

- Use related `order_display_id` and `invoice_display_id` directly in client invoice and transaction tables.
- Hide backend service references such as `provider_source_id` and `tenant_plan_id` from visible client labels.
- Use public source labels or safe `not shown` text for client service source/region fields.
- Fix visible mojibake text in the client wallet summary.

## Acceptance Criteria

- Client invoice and transaction reference columns use public IDs from response fields when available.
- Client services no longer show backend provider source or plan references as user-facing labels.
- Client wallet summary renders ASCII-safe readable text.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- Backend references may remain in memory for joins/actions, but should not be visible labels.

## Agent Log

- 2026-04-25: Codex created and claimed the task after T128 was marked done; starting client public ID label cleanup.
- 2026-04-25: Codex updated client dashboard, services, invoices, transactions, checkout invoice list, and wallet labels to use public IDs or safe hidden-reference text. Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, and `npm --prefix frontend run build`.
- 2026-04-25: Opened PR https://github.com/Chinsusu/Billing-V2/pull/289 for review and CI.
