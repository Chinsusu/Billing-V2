# T020 - Catalog API handlers

Status: IN_PROGRESS
Owner: Codex
Branch: feat/catalog-api-handlers
PR: -
Risk: API/catalog
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add catalog HTTP handlers on top of the catalog service layer so future app wiring can expose catalog APIs without handlers calling stores directly.

## Scope

- Add catalog handler routes for admin catalog setup, reseller clone/list, and client catalog list skeletons.
- Add request parsing, response DTO mapping, validation error mapping, and safe client response shape.
- Add app-level option wiring so a catalog route registrar can be attached without forcing DB wiring in default `NewAPI`.
- Add focused handler tests.
- Out of scope: auth/RBAC enforcement, production tenant context middleware, DB connection wiring, audit/outbox writes, seed data, checkout, order, invoice, or provisioning.

## Acceptance Criteria

- Handlers use catalog service interface, not store directly.
- Client catalog response does not expose reseller cost.
- Tenant catalog routes reject missing tenant context before service call.
- App option wiring can register catalog routes in tests.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Header-based tenant context is only an adapter for this skeleton. A future auth/tenant middleware task must replace it before public production use.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T020`.
