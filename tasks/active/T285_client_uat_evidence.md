# T285 - Client UAT evidence

Status: REVIEW
Owner: Codex
Branch: codex/t285-client-uat-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/602
Risk: UAT, tenant/RBAC, wallet, credential, provisioning, audit
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Capture redacted client-portal UAT evidence against the selected non-production test environment.

## Scope

- In scope: client login/session, client catalog, wallet/top-up/checkout/provisioning/renewal smoke, credential reveal safety, client RBAC negative checks, finance/audit support evidence, and redacted evidence docs.
- In scope: update the UAT evidence docs index and task board metadata.
- Out of scope: reseller UAT, admin UAT, production launch approval, production customer data, broad private beta, or unbounded real-provider provisioning.

## Acceptance Criteria

- Evidence uses only the selected non-production/test runtime.
- Evidence records command names, timestamps, public display IDs, status counts, and redacted outcomes only.
- No raw DSNs, passwords, cookies, session tokens, provider credentials, provider payloads, Telegram tokens, TOTP values, or service credentials are committed.
- Client UAT result is clearly marked PASS, FAIL, or BLOCKED with residual risk.
- Any unrun/manual-only UAT checks are reported explicitly and converted into follow-up tasks if needed.
- Required local validation passes before PR: `taskguard`, `git diff --check`, docs secret scan, and touched-file line-count check.

## Notes

- This task does not approve production. It can only support selected non-production UAT continuation.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t285-client-uat-evidence`.
- 2026-05-21: Verified target health for `billing.resvn.net`, `client.resvn.net`, `reseller.resvn.net`, local API, and local frontend.
- 2026-05-21: Found and fixed selected target DB domain mapping gap for `billing.resvn.net`, `client.resvn.net`, and `reseller.resvn.net` using verified/active `tenant_domains` rows with public display IDs only.
- 2026-05-21: Client UAT automated evidence passed: auth/RBAC smoke, billing mutation smoke, credential reveal smoke, finance reconciliation smoke, and client browser login/logout scope check.
- 2026-05-21: Added `docs/03_execution_operations_launch/80_Client_UAT_Evidence.md`, updated the UAT runbook entry criteria, and linked the evidence doc from docs index files.
- 2026-05-21: Opened PR https://github.com/Chinsusu/Billing-V2/pull/602 and moved task to REVIEW.
