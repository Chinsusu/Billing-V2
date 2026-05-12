# T201 - Support and abuse basic backend

Status: TODO
Owner: -
Branch: codex/t201-support-abuse-basic-backend
PR: -
Risk: support operations, abuse workflow, tenant isolation, service suspension, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add basic backend support and abuse control records needed for MVP operations.

## Scope

- Add minimal support ticket or support case records if current backend does not provide them.
- Add basic abuse flag/case workflow with evidence notes and service/account references.
- Enforce tenant/RBAC access and audit all sensitive support/abuse actions.
- Do not build a full custom ticket system beyond MVP needs.

## Acceptance Criteria

- Admin/reseller/client access follows tenant and permission boundaries.
- Abuse actions can record reason/evidence and trigger supported manual suspension path when applicable.
- Tests cover allowed/denied access and audit behavior.
- Relevant backend validation and CI pass.

## Notes

- Stop and ask if support data retention or abuse takedown policy is unclear.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
