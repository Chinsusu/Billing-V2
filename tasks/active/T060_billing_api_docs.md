# T060 - Billing API docs

Status: TODO
Owner: -
Branch: docs/billing-api-contracts
PR: -
Risk: docs/API
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Document the billing API routes and operational query parameters so agents and frontend work from a stable contract.

## Scope

- Document order, service, invoice, wallet, top-up, transaction, reconciliation, and audit routes.
- Include mutation payloads, idempotency requirements, and common error responses.
- Include display ID and amount filter query params from T056.
- Link the docs from the agent/workflow guidance if useful.

## Acceptance Criteria

- Docs describe path params, query params, request body, response shape, and auth expectations.
- Docs avoid over-specialized wording and stay clear for future agents.
- No code behavior changes unless needed for doc accuracy.

## Notes

- Keep docs compact and operational; do not turn this into generated OpenAPI unless it fits the repo style.

## Agent Log

- 2026-04-23: Task created to stabilize API contract knowledge.
