# T192 - Credential encrypted storage

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t192-credential-encrypted-storage
PR: -
Risk: credential security, database migration, provider provisioning, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add encrypted-at-rest credential storage for provisioned service credentials.

## Scope

- Design and add the credential storage model and migration required for provisioned services.
- Store credentials encrypted at rest using repository-approved config and secret rules.
- Ensure provider/provisioning paths store only encrypted credential payloads.
- Do not add reveal UI/API in this task; T193 owns reveal behavior.

## Acceptance Criteria

- Credential plaintext is not persisted in database columns, logs, audit, or tests.
- Migration has rollback notes and data impact documented.
- Tests cover encryption/decryption boundaries with safe fixtures.
- Migration validation, Go tests, secret/sensitive guards, and CI pass.

## Notes

- Stop and ask before choosing key management behavior if docs are insufficient.
- Migration data impact: creates new `service_credentials` table/enum types and adds a tenant-consistency unique constraint on `service_instances`; no backfill or existing data mutation.
- Rollback plan: `migrations/rollback/T192_service_credentials_down.sql` drops the new table and enum types for clean/dev rollback only.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Claimed by Codex on branch `codex/t192-credential-encrypted-storage`.
- 2026-05-13: Added service credential migration/store, shared AES-GCM cipher package, encrypted provider credential envelope helper, and provisioning worker storage path without reveal API/UI.
