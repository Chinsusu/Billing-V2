# T204 - Full E2E quality gate

Status: REVIEW
Owner: Codex
Branch: codex/t204-full-e2e-quality-gate
PR: https://github.com/Chinsusu/Billing-V2/pull/439
Risk: CI quality gates, database smoke, billing flow, provisioning, and release safety
Created: 2026-05-13
Updated: 2026-05-14

## Summary

Create a repeatable full end-to-end quality gate for the launch-critical billing and provisioning flow.

## Scope

- Run or wire local/CI parity for DB migration, seed, API, worker, frontend build, and smoke flows.
- Include `dev-db`, `dev-api`, and `dev-billing` smoke where prerequisites are available.
- Document any environment prerequisites and blocked checks.
- Do not point any smoke to production services.

## Acceptance Criteria

- A single documented validation path proves top-up, checkout, payment, job creation, worker provisioning, service activation, and frontend smoke.
- CI or documented local equivalent reports clear pass/fail output.
- Missing prerequisites are explicit blockers with remediation steps.
- Relevant validation and CI pass.

## Notes

- This task may depend on T189-T202 for full launch coverage, but it can start by documenting and automating current coverage.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-14: Codex claimed task on `codex/t204-full-e2e-quality-gate`.
- 2026-05-14: Added full E2E quality gate script, Makefile target, validation matrix entry, and launch runbook.
- 2026-05-14: First local gate run failed because unit tests inherited `DB_DSN`; fixed wrapper to unset `DB_DSN` for CI-parity backend checks and keep it only for smoke commands.
- 2026-05-14: Full local gate exposed and fixed catalog admin tenant context, smoke job JSON UUID cast, and jobs claim RETURNING ambiguity.
- 2026-05-14: Updated gate logging to redact `DB_DSN`, then reran the full local gate on temporary DB `billing_t204_e2e_20260514032311`: backend checks, dev-db, dev-api, dev-billing worker fulfillment, frontend install/audit/lint/build, and admin browser smoke passed.
- 2026-05-14: Dropped temporary local gate databases after evidence capture.
- 2026-05-14: Opened PR #439 and moved task to REVIEW.
