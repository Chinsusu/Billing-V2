# T195 - Daily reconciliation report

Status: TODO
Owner: -
Branch: codex/t195-daily-reconciliation-report
PR: -
Risk: finance reconciliation, wallet, ledger, invoice, and payment reporting
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add the daily reconciliation report needed for pilot operations.

## Scope

- Produce a backend reconciliation report for wallet balances, ledger entries, invoices, and payment transactions.
- Flag mismatches and duplicate payment references.
- Expose a safe admin read path or operator command for the report.
- Do not mutate financial records in this task.

## Acceptance Criteria

- Report identifies balanced and mismatched states with deterministic output.
- Tests cover clean data, wallet mismatch, invoice/payment mismatch, and duplicate references.
- Output uses public display IDs where shown to operators.
- Relevant backend validation and CI pass.

## Notes

- Coordinate with T194 if refund/adjustment entries affect reconciliation categories.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
