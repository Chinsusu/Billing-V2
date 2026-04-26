# T151 - Humanize transaction type labels

Status: DONE
Owner: Codex
Branch: codex/t151-humanize-transaction-type-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/333
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable transaction type labels instead of raw backend keys.

## Scope

- Add a shared payment display label helper for transaction types.
- Use the helper in admin, reseller, and client transaction tables.
- Move the payment method label helper into the shared payment label module if needed.
- Keep raw API values unchanged.
- Update admin browser smoke to verify a readable transaction type label.

## Acceptance Criteria

- Transaction type cells show readable labels such as Charge, Top-up, Refund, Purchase, and Service renewal.
- Raw type keys such as `topup`, `purchase.client_wallet.debit`, and `renewal.client_wallet.debit` are not rendered by transaction tables.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change billing APIs or mock source values.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T150 was marked done; starting transaction type label cleanup.
- 2026-04-26: Added shared payment label helpers, applied transaction type labels across admin/reseller/client transaction tables, and updated admin smoke coverage.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/333.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/333 into main; marking task done.
