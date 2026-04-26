# T152 - Humanize admin type labels

Status: DONE
Owner: Codex
Branch: codex/t152-humanize-admin-type-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/335
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable admin type labels instead of raw backend keys.

## Scope

- Add shared display label helpers for account, tenant, provider source, and product types.
- Use the helpers in admin account, customer, provider, and product rows.
- Keep filter values and raw API values unchanged.
- Update admin browser smoke coverage for a visible provider source type label.

## Acceptance Criteria

- Admin type badges show readable labels such as Reseller owner, Provider webhook, Self-hosted, Manual pool, VPS Linux, and VPS Windows.
- Raw keys such as `reseller_owner`, `provider_webhook`, `self-host`, and `vps-linux` are not rendered in the targeted admin rows.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is display-only; query parameters and API payloads keep the original keys.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T151 was marked done; starting admin type label cleanup.
- 2026-04-26: Added shared admin display label helpers and applied them to account, tenant, provider source, product, inventory, and risk labels.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/335.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/335 into main; marking task done.
