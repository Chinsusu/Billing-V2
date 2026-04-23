# T039 - Invoice domain schema

Status: REVIEW
Owner: Codex
Branch: feat/invoice-domain-schema
PR: -
Risk: invoice/migration
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add invoice and invoice item domain models with PostgreSQL schema so paid orders can later generate account-facing invoices.

## Scope

- Add invoice and invoice item domain models with validation.
- Add migration files for invoices and invoice items.
- Keep UUID primary keys and numeric display IDs for records shown in the UI.
- Add migration validation and focused domain tests.
- Out of scope: invoice creation workflow, PDF export, payment gateway integration, or frontend views.

## Acceptance Criteria

- Invoice records include tenant, buyer account, status, currency, subtotal, tax, discount, and total fields.
- Invoice item records can reference order items when available.
- Numeric display IDs are generated for UI-visible invoice records.
- Migration validation passes and rollback notes are clear.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task can start after order status and billing status paths are stable.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
- 2026-04-23: Claimed by Codex from latest `origin/main` in `/tmp/Billing-T039`.
- 2026-04-23: Added invoice domain models, invoice tables, rollback script. Validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`.
