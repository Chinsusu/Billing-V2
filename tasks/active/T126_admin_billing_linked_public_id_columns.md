# T126 - Admin billing linked public ID columns

Status: DONE
Owner: Codex
Branch: codex/t126-admin-billing-linked-public-id-columns
PR: https://github.com/Chinsusu/Billing-V2/pull/283
Risk: frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Show related public ID labels in admin billing tables so operators can trace invoices and transactions without seeing backend references.

## Scope

- Add linked order/account public labels to the admin invoice table.
- Add linked account/order/invoice public labels to the admin transaction table.
- Keep raw backend UUIDs hidden from admin billing screens.
- Extend admin browser smoke assertions for the new visible labels.

## Acceptance Criteria

- Admin invoice rows show invoice ID, account public ID, and order public ID when the API returns them.
- Admin transaction rows show transaction ID, account public ID, order public ID, and invoice public ID when the API returns them.
- Demo fallback remains usable when live API is unavailable.
- Existing frontend quality gates and taskguard pass.

## Notes

- This follows T124/T125; backend fields already exist and smoke coverage already checks they are returned.

## Agent Log

- 2026-04-25: Codex created and claimed the task after T125 merged; starting admin billing table display updates.
- 2026-04-25: Added linked public ID columns to admin invoices and transactions, with browser smoke assertions for live labels.
- 2026-04-25: Validation passed: frontend smoke, lint, build, sensitive-text guard, taskguard, and diff check.
- 2026-04-25: Opened PR #283 for review.
- 2026-04-25: PR https://github.com/Chinsusu/Billing-V2/pull/283 merged into `main`.
