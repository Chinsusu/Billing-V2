# T150 - Humanize payment method labels

Status: REVIEW
Owner: Codex
Branch: codex/t150-humanize-payment-method-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/331
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable payment method labels instead of raw backend keys.

## Scope

- Add a shared payment method label helper.
- Use the helper in admin, reseller, and client top-up displays.
- Use the helper for admin transaction reconciliation provider labels where applicable.
- Keep raw API values unchanged.
- Update admin browser smoke to verify the displayed label.

## Acceptance Criteria

- Payment method cells show readable labels such as Bank transfer, Wallet, Crypto, VietQR, and USDT.
- Raw method keys such as `bank_transfer` are not rendered by top-up tables.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change payment or top-up APIs.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T149 was marked done; starting payment method label cleanup.
- 2026-04-26: Added a shared payment method label helper, applied it to top-up/transaction displays, and updated admin smoke coverage for the visible label.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/331.
