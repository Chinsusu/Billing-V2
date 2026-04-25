# T133 - Reseller remove backend reference fallbacks

Status: REVIEW
Owner: Codex
Branch: codex/t133-reseller-remove-backend-fallbacks
PR: https://github.com/Chinsusu/Billing-V2/pull/297
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Remove reseller frontend fallbacks that join wallet/customer labels by backend user IDs now that related public IDs are available.

## Scope

- Use `buyer_display_id` and `account_display_id` only for reseller billing customer labels.
- Use `owner_display_id` only for reseller wallet-to-customer display lookups.
- Keep backend IDs only in action bodies or internal join paths that still require them.

## Acceptance Criteria

- Reseller billing labels do not depend on `buyer_user_id` or `account_user_id` fallbacks.
- Reseller client/dashboard wallet display lookup does not depend on `owner_id`.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- If a related public ID is missing, UI should show the public ID fallback or `-`, not a backend reference.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T132 was marked done; starting reseller backend-reference fallback cleanup.
- 2026-04-26: Removed reseller backend-ID label fallbacks for billing customer labels and wallet owner display lookups.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/297 for review.
