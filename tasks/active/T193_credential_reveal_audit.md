# T193 - Credential reveal audit

Status: TODO
Owner: -
Branch: codex/t193-credential-reveal-audit
PR: -
Risk: credential security, tenant isolation, RBAC, rate limiting, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add controlled credential reveal behavior with masking, authorization, rate limiting, and audit.

## Scope

- Add backend reveal API for service credentials using encrypted storage from T192.
- Enforce tenant/RBAC permission and reveal rate limits.
- Audit every reveal without logging plaintext credential values.
- Add frontend masking/reveal behavior for the relevant client/admin/reseller service detail paths.

## Acceptance Criteria

- Credentials are masked by default and reveal only after authorized action.
- Cross-tenant and unauthorized reveal attempts are denied and tested.
- Reveal audit includes actor, tenant, service, reason/context where required, and no secret plaintext.
- Relevant backend/frontend validation and CI pass.

## Notes

- This task depends on T192.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
