# T139 - Demo account service audit prefixes

Status: DONE
Owner: Codex
Branch: codex/t139-demo-entity-prefixes
PR: https://github.com/Chinsusu/Billing-V2/pull/309
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Align remaining visible demo and fallback account, service, and audit log IDs with the live public ID prefixes.

## Scope

- Replace demo customer and reseller-client `C-*` / `RC-*` account IDs with `ACC-*`.
- Replace demo service IDs that still use `prx-*`, `vps-*`, `bw-*`, or `svc-*` values with `SVC-*`.
- Replace demo audit log IDs that still use `LOG-*` values with `AUD-*`.
- Update visible demo audit targets that reference service/provider IDs to public prefixes.
- Keep request IDs, product SKUs, plan codes, and backend-only correlation IDs unchanged.

## Acceptance Criteria

- Demo account rows show `ACC-*` public IDs.
- Demo service rows and ledger references show `SVC-*` public IDs.
- Demo audit log rows show `AUD-*` public IDs.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- This task only changes mock data used when live API data is unavailable.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T138 was marked done; starting remaining demo public ID cleanup.
- 2026-04-26: Replaced visible demo account, service, audit log, and provider-source references with `ACC-*`, `SVC-*`, `AUD-*`, and `SRC-*` public prefixes.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/309 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/309 merged into `main`; marking task done.
