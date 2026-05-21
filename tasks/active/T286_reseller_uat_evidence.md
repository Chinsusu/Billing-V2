# T286 - Reseller UAT evidence

Status: REVIEW
Owner: Codex
Branch: codex/t286-reseller-uat-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/604
Risk: UAT, tenant/RBAC, wallet, credential, provisioning, audit
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Capture redacted reseller-portal UAT evidence against the selected non-production test environment.

## Scope

- In scope: reseller login/session, reseller catalog/customer/service/invoice/transaction/wallet/top-up/job read scopes, reseller browser portal scope, negative admin/client access checks, cleanup/status evidence, and redacted evidence docs.
- In scope: update the UAT evidence docs index and task board metadata.
- Out of scope: client UAT already covered by T285, admin UAT, production launch approval, production customer data, broad private beta, or unbounded real-provider provisioning.

## Acceptance Criteria

- Evidence uses only the selected non-production/test runtime.
- Evidence records command names, timestamps, public display IDs, status counts, and redacted outcomes only.
- No raw DSNs, passwords, cookies, session tokens, provider credentials, provider payloads, Telegram tokens, TOTP values, or service credentials are committed.
- Reseller UAT result is clearly marked PASS, FAIL, or BLOCKED with residual risk.
- Any unrun/manual-only UAT checks are reported explicitly and converted into follow-up tasks if needed.
- Required local validation passes before PR: `taskguard`, `git diff --check`, docs secret scan, and touched-file line-count check.

## Notes

- This task does not approve production. It can only support selected non-production UAT continuation.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t286-reseller-uat-evidence`.
- 2026-05-21: Verified target health for `billing.resvn.net`, `client.resvn.net`, `reseller.resvn.net`, local API, and local frontend.
- 2026-05-21: Reseller UAT automated evidence passed: browser login/logout scope check, reseller session API read checks, admin/client negative access checks, reseller UI navigation, and finance reconciliation read smoke.
- 2026-05-21: Added `docs/03_execution_operations_launch/81_Reseller_UAT_Evidence.md` and linked the evidence doc from docs index files.
- 2026-05-21: Opened PR https://github.com/Chinsusu/Billing-V2/pull/604 and moved task to REVIEW.
