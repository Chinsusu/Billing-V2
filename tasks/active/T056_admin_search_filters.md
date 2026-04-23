# T056 - Admin search filters

Status: IN_PROGRESS
Owner: Codex
Branch: feat/admin-search-filters
PR: -
Risk: API/search/admin
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add practical admin search filters for operational billing lists so support can find records by numeric display IDs and common customer/payment attributes.

## Scope

- Add display ID filters where the APIs already expose FE-visible records.
- Add customer/account filters where the schema supports tenant-safe lookup.
- Add amount range filters for payment or reconciliation lists if practical.
- Keep tenant and RBAC checks intact for every filtered lookup.
- Document filter query parameters in handler tests or API notes.

## Acceptance Criteria

- Admin list APIs support the agreed filters with tests.
- Filters never bypass tenant scoping.
- Invalid filter input returns a clear API error.
- Backend quality gates pass.

## Notes

- UUID remains the internal identifier; display ID is for UI/support lookup.

## Agent Log

- 2026-04-23: Task created from operational admin needs.
- 2026-04-23: Codex claimed task and started admin filter implementation from latest `origin/main`.
