# T284 - UAT runbook

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t284-uat-runbook
PR: -
Risk: docs, QA/UAT workflow
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Add a UAT runbook for client, reseller, and admin flows in the selected non-production pilot environment.

## Scope

- In scope: UAT entry/exit criteria, portal matrix, execution sequence, bug severity, evidence packet, and sign-off checklist.
- In scope: update docs index and task board metadata.
- Out of scope: running UAT evidence, mutating provider resources, changing product behavior, or production launch approval.

## Acceptance Criteria

- UAT doc clearly separates test/non-production UAT from production launch.
- Client, reseller, and admin flows each have positive and negative checks.
- P0 areas cover money, tenant/RBAC, credential, provider/provisioning, notification/fallback, and audit.
- Evidence/sign-off packet avoids secrets, customer data, raw provider payloads, DSNs, cookies, and credentials.
- `taskguard`, docs secret scan, line-count check, and `git diff --check` pass.

## Notes

- Follow-up task should run UAT evidence against the selected test environment after this runbook is merged.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t284-uat-runbook`.
- 2026-05-21: Added `docs/03_execution_operations_launch/79_UAT_Client_Reseller_Admin_Runbook.md` and linked it from `docs/00_README.md`.
- 2026-05-21: Validation passed: `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, docs secret scan, and line-count check for touched files.
