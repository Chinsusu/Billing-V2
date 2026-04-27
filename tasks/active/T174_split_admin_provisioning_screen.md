# T174 - Split admin provisioning screen

Status: REVIEW
Owner: Codex
Branch: codex/t174-split-admin-provisioning
PR: https://github.com/Chinsusu/Billing-V2/pull/379
Risk: frontend provisioning UI organization
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Reduce `frontend/src/modules/admin/screens/AdminProvisioning.tsx` file-size risk by moving stable table/filter/detail helpers into focused modules without changing provisioning UI behavior.

## Scope

- Split selected helpers or presentational pieces out of `AdminProvisioning.tsx`.
- Keep the admin provisioning screen behavior and smoke coverage unchanged.
- Do not change API contracts, provisioning business rules, or backend code.

## Acceptance Criteria

- `AdminProvisioning.tsx` is safely below the 500-line file limit.
- Any new files stay below 500 lines.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- Follow-up after T173; `AdminProvisioning.tsx` was 467 lines and close to the repository file-size limit.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Moved provisioning queue table and recovery controls into `AdminProvisioningQueueTable`; local gates pass.
- 2026-04-27: Opened PR #379 and moved task to REVIEW.
