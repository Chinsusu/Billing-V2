# T192 - Credential encrypted storage

Status: TODO
Owner: -
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

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
