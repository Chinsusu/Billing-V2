# T165 - Demo top-up actor labels

Status: REVIEW
Owner: Codex
Branch: codex/t165-demo-topup-actor-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/361
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up top-up demo fallback requester labels so operators do not see raw wallet actor keys when the live top-up API is unavailable.

## Scope

- Humanize demo top-up actor labels.
- Add browser smoke coverage for the top-up API fallback path.

## Acceptance Criteria

- Demo fallback shows `Reseller Wallet` and `Client Wallet`.
- Raw keys such as `reseller_wallet`, `client_wallet`, and `pending_verification` are not visible.
- Admin browser smoke covers the fallback path.

## Notes

- Live top-up rows already use API view model labels.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Opened PR #361 after lint, sensitive-text guard, build, admin smoke, taskguard, and diff check passed.
