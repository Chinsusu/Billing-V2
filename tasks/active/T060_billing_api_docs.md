# T060 - Billing API docs

Status: DONE
Owner: Codex
Branch: docs/billing-api-contracts
PR: #135
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
- 2026-04-23: Codex claimed the task and started a billing API operational reference based on live handler/filter code.
- 2026-04-23: PR #135 opened with the new billing API operational reference and AGENTS link update.
- 2026-04-23: Validation for the docs branch was `git diff --check`; CI passed before merge.
- 2026-04-23: PR #135 merged into `main` with commit `178239ba667cef6fbfb0e90d24978b7a0bf7cdfd`.
