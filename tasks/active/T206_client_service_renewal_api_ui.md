# T206 - Client service renewal API and UI action

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t206-client-service-renewal-api-ui
PR: -
Risk: wallet debit, invoice/payment, service lifecycle, tenant isolation, RBAC, and audit
Created: 2026-05-13
Updated: 2026-05-14

## Summary

Add a production-safe client service renewal path once the required backend contract is explicit.

## Scope

- Add or confirm a client-facing service renewal API contract for expiring or expiry-suspended services.
- Ensure renewal uses wallet/invoice/payment rules, idempotency, service lifecycle transitions, and audit.
- Wire the client dashboard/service renewal CTA to the production API after the backend path is ready.
- Keep unsupported renew actions hidden or routed to safe checkout/support paths until implemented.

## Acceptance Criteria

- Client renewal cannot bypass wallet, invoice, tenant, or service lifecycle checks.
- Renewal action is idempotent and audited.
- Frontend only shows the direct renewal CTA when the API/capability supports it.
- Tests cover allowed, denied, insufficient balance, and invalid service-status cases.

## Notes

- Created from T202 audit: `ClientDashboard` had a static `Renew now` CTA, but there is no explicit client renew HTTP endpoint yet.

## Agent Log

- 2026-05-13: Task created by Codex during T202 frontend production integration audit.
- 2026-05-14: Codex claimed task on `codex/t206-client-service-renewal-api-ui`.
