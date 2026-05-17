# 70 - Launch Evidence Completion Packet

**Date:** 2026-05-17
**Scope:** Single completion packet for the remaining launch blockers before reconsidering the pilot Go/No-Go decision.  
**Decision:** NO-GO until every required evidence section below is complete, redacted, reviewed, and signed off.

## Purpose

This packet is the final evidence checklist for the work that cannot be proven by repository code or local/dev smokes alone.

The repository currently has strong local/dev evidence for core implementation, a T241 single-owner launch assignment to `Admin`, and T242 target-server staging-equivalent backup/restore evidence. It still does not contain full real provider account proof, staging/full E2E proof, production notification delivery proof, or final target-environment security/finance/QA sign-off.

Do not change the pilot decision to GO or CONDITIONAL GO by filling this packet with assumptions. Every row needs actual evidence or an explicit owner-approved exception.

## Redaction Boundary

Never commit or paste:

- raw DSNs, database passwords, provider API keys, bearer tokens, SMTP passwords, Telegram bot tokens, private keys, cookies, or authorization headers;
- dump files, raw provider request/response payloads, customer data, service credentials, reset tokens, or credential reveal output;
- production account IDs or production customer identifiers unless a security owner explicitly approves a redacted reference.

Use display IDs, redacted placeholders, dates, command names, check counts, and owner names instead.

## Completion Matrix

| Gate | Current repo status | Required completion evidence | Required owner sign-off |
|---|---|---|---|
| Real provider sandbox | Blocked for broader pilot. T199 proves local fake provider behavior; doc 66/T208 defines the provider evidence packet. T213 records Cloudmini V3 API version and non-production base URL. T214 recorded the earlier edge/gateway HTTP `403` blocker; T215 documents the provider-owner unblock; T216 records a successful 2026-05-16 read-only rerun through the public hostname using bearer, `X-API-Key`, and `X-ACCESS-CODE` from a local dev credential source. T217 adds multi-endpoint runtime support. T218 defines a controlled pilot approval packet with a redacted `ipv4_dc` mapping candidate and quota/cleanup guardrails. T219 adds guarded non-production catalog mapping tooling. T220 applied the mapping on the approved Billing dev runtime env and T221 read-only evidence passed with readiness `ready` for the pilot `cloudmini_v3` source. T228 ran one controlled dev Billing-path create/delete pilot with encrypted credential storage and same-session cleanup. T229 adds repo-side fail-closed handling for non-usable Cloudmini statuses and lifecycle-worker provider cleanup before service termination. T230 deployed and build-tested that hardening on the approved test server without provider mutations. T231 proves non-mutating worker registry activation with the real Cloudmini adapter and protected dev credential path. T232 attempted owner-approved dev mutating/lifecycle activation; Billing reached Cloudmini create but manual-reviewed provider status `creating`, then same-session direct V3 cleanup succeeded. T233 adds bounded post-create status polling and the target rerun passed with one active service, encrypted credential storage, lifecycle-worker provider cleanup, provider final `404`, and worker restore. T241 assigns the owner roles to `Admin`, but approved shared credential storage, live duplicate/timeout evidence, full error examples, usable-status owner sign-off, target auth/RBAC/finance/security sign-off, or broader owner approval is still not recorded. | Approved sandbox account, shared secret-store path, owner-approved quota/cost limit, SKU/location sign-off, live timeout/idempotency behavior, redacted error examples, cleanup owner, edge/gateway access approval record, Cloudmini usable-status semantics, and target/staging owner sign-off. | Provider Owner, Engineering Lead, Ops Lead, Security Owner |
| Shared staging backup/restore | Partial with target staging-equivalent pass. T203 proves local restore. T242 proves a target-server staging-equivalent clean source/restore drill with checksum, restore, smoke, and cleanup evidence. The long-lived target app DB was not used as pass evidence because prior dev/test smoke mutations make strict seed-baseline `dev-db` smoke unsuitable. | Admin/Ops/QA/Security acceptance of the staging-equivalent scope, or an additional approved clean shared staging app snapshot restore with suitable post-restore smoke. | Ops Lead, QA Lead, Security Owner |
| Staging/full E2E | Partial. T204 proves local/dev full gate with fake provider. | Approved staging or signed staging-equivalent run covering auth/RBAC, top-up approval, checkout, wallet payment, provisioning boundary, service activation, T206 renewal, lifecycle jobs, frontend smoke, and audit checks. | QA Lead, Engineering Lead, Product Owner |
| Notification delivery or fallback | Partial. T200 provides local notification foundation only. T222 defines the manual fallback owner/SLA/evidence packet, but no production delivery proof or approved fallback drill is recorded. | Production SMTP/Telegram delivery proof for launch-critical events, or a manual fallback with owner, SLA, escalation path, and sample redacted notification records. | Ops Lead, Support Owner, Security Owner |
| Launch-day owners | Assigned with single-owner risk. T241 records the user-provided assignment that `Admin` owns Product, Engineering, QA, Ops, Finance, Security, Support, and Provider launch-day roles. | A launch window, escalation path, and explicit acceptance that one person owns all role decisions for the selected launch scope. | Product Owner, Engineering Lead |
| Target-environment verification | Partial. T230 proves the hardened backend/worker/frontend code deploys and builds on the approved test server, target services are active, `/healthz`, `/readyz`, and frontend return HTTP `200`, and protected Cloudmini dev credentials stay outside git. T231 proves non-mutating Cloudmini registry activation. T232 proves the target can reach a real Cloudmini create path and safely cleanup by fallback. T233 proves one target Cloudmini lifecycle-worker cleanup activation. T235 proves target top-up review create/approve/reject via HTTP API on the approved test server: approve top-up display `10003` posted ledger display `10005` and audit display `10015`; reject top-up display `10004` posted no ledger and audit display `10016`; wallet delta was `111`; provider side effects were `none`. T236 proves target API session/RBAC behavior on the approved test server: client cookie-only `/client/catalog` passed without `X-Actor-*`, unsatisfied platform admin session was blocked with `auth.2fa_required`, invalid session was blocked with `auth.session_invalid`, missing actor was blocked with `auth.actor_required`, cross-tenant mismatch was blocked with `tenant.context_mismatch`, and three low-permission RBAC checks were blocked with `auth.permission_denied`. T237 proves target credential reveal audit/redaction behavior: the seeded client session revealed one encrypted dev/test fixture for service display `43001`, response headers were `no-store`, audit display `10017` was recorded, reveal rate-limit state was recorded, and smoke output excluded plaintext credentials, encrypted payloads, raw credential IDs, session tokens, cookies, DSNs, provider payloads, and provider credentials. T238 proved target finance reconciliation read paths and found one dev/test wallet projection mismatch. T239 fixed the seed baseline, repaired the target dev/test projection from posted ledger source-of-truth with audit display `10018`, and reran target finance reconciliation with `balanced` status and zero mismatch counts. T240 removed cloudflared token flag exposure from process argv by switching to `--token-file`, verified restricted secret file modes, and confirmed target local/domain HTTP `200`. External browser auth, Finance owner sign-off, Security owner sign-off, and broader owner sign-off still remain missing. | Finance owner review/sign-off for the balanced T239 evidence, Security owner review/sign-off for T240 secret/key evidence, external browser auth evidence if required for launch, and owner sign-off. | Security Owner, Finance Lead, QA Lead |

Any missing required sign-off keeps the launch decision at NO-GO.

## Evidence Packet

Fill one packet per launch candidate. Store only redacted evidence in git.

```text
Launch candidate ID:
Date/time UTC:
Pilot scope:
Environment:
Evidence collector:
Final reviewer:
Decision requested: GO / CONDITIONAL GO / NO-GO
```

### 1. Real Provider Sandbox

```text
Provider:
Provider owner:
Sandbox account reference: redacted
Credential storage path: redacted secret-store reference only
Credential scope:
Quota/cost limit:
Provider support contact:
Billing plan code:
Provider SKU:
Sandbox location:
Timeout policy:
Idempotency level:
Cleanup owner:
Real pilot run ID:
Run result:
Redacted evidence link/reference:
Provider owner sign-off:
Security owner sign-off:
```

Pass criteria:

- Provider account and credentials are sandbox-only and stored outside git.
- Billing plan maps to an explicit provider SKU/location.
- Duplicate create and timeout-after-send behavior are documented and tested.
- Pilot run creates at most one provider resource and cleanup is recorded.
- No raw provider secret, credential, or payload appears in logs, PRs, tasks, or docs.

### 2. Shared Staging Backup/Restore

```text
Drill ID:
T242-target-20260517T134247Z
Source classification:
Temporary target-server staging-equivalent seed DB, no production data.
Target classification:
Temporary target-server staging-equivalent restore DB, approved to overwrite.
Target overwrite approval:
Bounded to billing_t242_restore_20260517134247 only.
Backup artifact path: redacted non-repo path
/tmp/billing-t242-backup-restore/billing-billing_t242_source_20260517134247-20260517T134248Z.dump
Backup checksum:
be364dcbd3b434402f89bfbfef941d66e96c04e3d88e4d7ef70b91d9b4f0c0e2
Restore command:
bash scripts/backup_restore_drill.sh --run with redacted target-server DSNs
Restore result:
PASS; pg_restore completed and the drill reported backup/restore passed.
Smoke command:
go run ./cmd/smoke -dsn "$BILLING_RESTORE_TARGET_DSN" -timeout 90s dev-db
Smoke result:
PASS on restored target.
Migration count:
Source smoke applied 25 migrations; restored target applied 0 new migrations and reported 25 schema migration rows.
Smoke check count:
20 checks on source and 20 checks on restored target.
Cleanup/retention decision:
Dump/checksum files deleted after evidence capture; temporary source and restore DBs dropped.
Ops sign-off:
Admin assigned by T241; final evidence acceptance pending.
QA sign-off:
Admin assigned by T241; final evidence acceptance pending.
Security sign-off:
Admin assigned by T241; final evidence acceptance pending.
```

Pass criteria:

- Source and target are approved non-production or staging-equivalent databases.
- Restore target overwrite is approved before running destructive restore.
- Restored target passes `dev-db` smoke.
- Backup artifact retention or deletion owner is recorded.

### 3. Staging Or Staging-Equivalent Full E2E

```text
Gate ID:
Environment:
DB/API classification:
Provider mode: fake/manual/real sandbox
Frontend target:
Backend result:
DB smoke result:
API/RBAC smoke result:
Billing mutation result:
Renewal path result:
Lifecycle job result:
Frontend smoke result:
Audit/redaction result:
Exception requested: yes/no
Exception owner and reason:
QA sign-off:
Engineering sign-off:
Product sign-off:
```

Pass criteria:

- T206 renewal path is included with wallet debit, invoice/payment records, lifecycle renewal, and audit evidence.
- RBAC negative checks and cross-tenant attempts fail.
- Credential reveal remains masked by default and audited when revealed.
- Real provider work is excluded unless section 1 is complete.
- Any staging-equivalent exception names the owner, reason, limits, expiry date, and residual risk.

### 4. Notification Delivery Or Manual Fallback

```text
Delivery mode: SMTP / Telegram / dashboard / manual fallback
Launch-critical events covered:
Credential/secret storage path: redacted secret-store reference only
Successful delivery evidence:
Failure/retry evidence:
Manual fallback owner:
Manual fallback SLA:
Escalation path:
Support owner sign-off:
Ops sign-off:
Security sign-off:
```

Pass criteria:

- At least top-up status, provisioning failure/manual review, service lifecycle, password reset, and support/abuse critical events have delivery or fallback coverage.
- Failure mode and retry/manual fallback are tested or explicitly accepted.
- Notification payloads are redacted and contain no credentials or reset tokens.
- If using manual fallback, complete `docs/03_execution_operations_launch/72_Notification_Delivery_And_Manual_Fallback_Runbook.md` with named owner, SLA, escalation channel, sampled events, and redacted evidence.

### 5. Launch-Day Owners

```text
Product Owner:
Admin
Engineering Lead:
Admin
QA Lead:
Admin
Ops Lead:
Admin
Finance Lead:
Admin
Security Owner:
Admin
Support Owner:
Admin
Provider Owner:
Admin
Escalation channel:
Admin direct launch channel; single-person escalation accepted by user statement on 2026-05-17.
Launch window:
Not approved until remaining P0 evidence gates are complete.
Owner availability confirmed:
Yes for owner assignment by user statement on 2026-05-17: "1 mình tao cân hết. Admin".
Single-owner risk:
Accepted for owner assignment only. This does not waive missing staging, backup/restore, notification/fallback, provider duplicate/timeout/error, shared secret-store, or final sign-off evidence.
```

Pass criteria:

- Every role has a named human owner before launch.
- Each owner has accepted their launch-day responsibility.
- Escalation channel and launch window are recorded.

### 6. Target-Environment Verification

```text
Auth/session target check:
T236 PASS on the approved test server local API. Client seed login set an HttpOnly `billing_session`; cookie-only `/client/catalog` passed without `X-Actor-*` dev helper headers. Invalid session returned `auth.session_invalid`. Missing actor returned `auth.actor_required`.
Admin 2FA enrollment/enforcement: T236 PASS for enforcement on the approved test server local API. Platform staff login returned an unsatisfied 2FA session, and cookie-only admin provider-readiness access returned `auth.2fa_required`. Enrollment of named production admins is still missing.
Credential reveal audit access:
T237 PASS on the approved test server local API. The smoke created/refreshed one encrypted dev/test credential fixture for service display 43001, logged in as the seeded client, revealed through the client credential reveal route with only the HttpOnly session cookie, verified `Cache-Control: no-store` and `Pragma: no-cache`, verified `last_revealed_by` and rate-limit state, and recorded audit display 10017 for `credential.revealed`. The evidence excludes plaintext credentials, encrypted payloads, raw credential IDs, session tokens, cookies, DSNs, provider payloads, and provider credentials.
Top-up review target check: T235 PASS on the approved test server. Approve top-up display 10003 posted ledger display 10005 and audit display 10015. Reject top-up display 10004 posted no ledger and audit display 10016. Wallet delta was 111 minor units. Provider side effects were none.
Secret/key handling review:
T240 PASS on the approved test server for target dev/test secret handling evidence. `/opt/Billing/.env.dev` is mode 640 with owner `root:billing-svc`; `/opt/cred-cloudmini-dev.env` is mode 600 with owner `root:root`; `/etc/cloudflared/tunnel.token` is mode 600 with owner `root:root`. Cloudflared was changed from token flag usage to `--token-file /etc/cloudflared/tunnel.token`, restarted, and verified active with `cloudflared_token_in_argv=no`. Reachability returned HTTP 200 for `http://localhost:3000`, `https://billing.resvn.net`, `https://reseller.resvn.net`, and `https://client.resvn.net`. No raw token, DSN, API key, cookie, credential, or file contents were printed or recorded. Security owner sign-off and approved production/shared secret-store evidence are still missing.
Finance reconciliation owner run:
T238 evidence captured on the approved test server local API/DB. The read-only smoke selected transaction display 51001, invoice display 44001, wallet display 41001, and ledger display 50002; payment reconciliation list/detail returned matching public display evidence; daily reconciliation for 2026-04-23 initially returned `mismatched` with one wallet mismatch. T239 traced the root cause to dev/test wallet projection drift from an inconsistent seed baseline and later smoke runs, fixed the seed baseline, repaired wallet display 41001 from posted ledger source-of-truth with audit display 10018, and reran the target smoke. The rerun returned `balanced` with wallets checked 2, wallet mismatches 0, invoices checked 1, invoice mismatches 0, payments checked 1, and duplicate payment references 0. The repair inserted no ledger rows, updated no posted ledger rows, called no money/provider mutation routes, and printed no secrets. Finance owner review/sign-off is still missing.
Cross-tenant negative check: T236 PASS on the approved test server local API. A mismatched actor tenant was denied with `tenant.context_mismatch`, and low-permission RBAC checks returned `auth.permission_denied`.
Support coverage check:
Residual risks:
Security sign-off:
Finance sign-off:
QA sign-off:
```

Pass criteria:

- Target auth/session configuration is verified outside dev-only actor headers.
- Admin 2FA is enrolled and enforced in the target environment.
- Credential reveal audit is visible to authorized operators and redacted elsewhere.
- Finance owner runs or reviews reconciliation evidence.
- Cross-tenant negative tests fail safely.

## Final Decision Rule

After all sections are complete:

1. Re-run `docs/03_execution_operations_launch/69_Pilot_Go_No_Go_Record.md`.
2. Keep decision NO-GO if any P0 section is missing, unreviewed, or based on assumptions.
3. Use CONDITIONAL GO only for non-P0 exceptions with a named owner, mitigation, expiry date, and rollback path.
4. Use GO only when every P0 gate has passing evidence and required owner sign-off.

Until then, this repository remains launch-ready for local/dev validation only, not for external private beta, pilot launch, or real-provider production-like provisioning.
