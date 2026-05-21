# 70 - Launch Evidence Completion Packet

**Date:** 2026-05-19
**Scope:** Single completion packet for the selected bounded non-production pilot Go/No-Go decision.
**Decision:** GO for the selected bounded non-production pilot only. NO-GO for production, broader private beta, broader provider scope, production customer data, and real-provider production-like provisioning outside the approved selected Cloudmini scope.

## Purpose

This packet is the final evidence checklist for the work that cannot be proven by repository code or local/dev smokes alone.

The repository currently has strong local/dev evidence for core implementation, a T241 single-owner launch assignment to `Admin`, T242 target-server staging-equivalent backup/restore evidence, T243 target-server staging-equivalent full E2E evidence including T206 renewal, T244 owner-approved notification manual fallback evidence, T245 Admin sign-off for the completed target evidence gates, T247 target dev/test secret-file metadata recheck, T252 Cloudmini usable-status owner sign-off, T253 Cloudmini cleanup owner/procedure evidence, T254 selected-host self-managed secret-store proof after owner-confirmed key rotation, T255 source-inspection/runbook evidence for Cloudmini provider-controlled error cases, T256 permission-denied runtime evidence, T257 out-of-capacity runtime evidence, T258 Billing runner support for a future rate-limit fixture, T259 merged provider-side rate-limit fixture code, T260 rate-limit runtime evidence, T261 provider 5xx runtime evidence, T262 cancel/delete rejected runtime evidence, T263 selected Cloudmini provider-owner approval, T264 selected-pilot manual notification fallback approval, T266 selected-domain auth/RBAC smoke evidence, T267 protected systemd runtime evidence, T268 selected launch-window, single-owner acceptance, and target Admin 2FA enrollment/enforcement evidence, T269 selected launch-window finance reconciliation, T270 selected launch-window ops health/support coverage evidence, T271/T272/T273/T274/T275 selected support-window checkpoint evidence, T276 final support-window closeout evidence, T278 Telegram delivery worker implementation, T279 selected-host Telegram preflight evidence, T280 one-message queued Telegram delivery evidence, T281 controlled Telegram failure/retry classification evidence, T290 domain-aware target auth smoke support, and T291 selected test-server split-domain target auth smoke deploy evidence. It does not approve provider scope beyond the selected bounded non-production Cloudmini pilot. Broader production notification delivery still needs scope-specific owner approval before using Telegram as a sole primary path.

Do not broaden the pilot decision by filling this packet with assumptions. Every broader row needs actual evidence or an explicit owner-approved exception.

## Redaction Boundary

Never commit or paste:

- raw DSNs, database passwords, provider API keys, bearer tokens, SMTP passwords, Telegram bot tokens, private keys, cookies, or authorization headers;
- dump files, raw provider request/response payloads, customer data, service credentials, reset tokens, or credential reveal output;
- production account IDs or production customer identifiers unless a security owner explicitly approves a redacted reference.

Use display IDs, redacted placeholders, dates, command names, check counts, and owner names instead.

## Completion Matrix

| Gate | Current repo status | Required completion evidence | Required owner sign-off |
|---|---|---|---|
| Real provider sandbox | Ready for selected bounded non-production Cloudmini pilot. Docs 66 and 77 record Cloudmini V3 intake, read-only reachability, source mapping, controlled create/delete pilot, fail-closed status handling, lifecycle cleanup, target deploy/build evidence, registry activation, bounded status polling, duplicate-create and timeout-after-send evidence, redacted provider error examples, usable-status sign-off, cleanup owner/procedure, selected-host secret-store proof, and T263 provider-owner approval. | Additional provider accounts, additional secret paths, quota above one active resource, production customer data, or unbounded production-like provisioning require a new owner-approved packet and repeated secret-store proof. | Admin as Provider Owner, Engineering Lead, Ops Lead, and Security Owner for the selected scope |
| Shared staging backup/restore | Pass for target staging-equivalent scope. T203 proves local restore. T242 proves a target-server staging-equivalent clean source/restore drill with checksum, restore, smoke, and cleanup evidence. T245 records Admin/Ops/QA/Security acceptance of the staging-equivalent scope. The long-lived target app DB was not used as pass evidence because prior dev/test smoke mutations make strict seed-baseline `dev-db` smoke unsuitable. | Use the T242/T245 staging-equivalent exception for pilot, or run an additional approved clean shared staging app snapshot restore if the launch scope rejects that exception. | Ops Lead, QA Lead, Security Owner |
| Staging/full E2E | Pass for target staging-equivalent scope. T204 proves local/dev full gate with fake provider. T243 extends the gate to include T206 renewal and passes it on the approved test server using a temporary target DB, local API, fake-provider fulfillment, and mocked frontend browser smoke. T245 records Admin/QA/Engineering/Product acceptance of this staging-equivalent scope. | Real provider work remains excluded unless the provider gate is complete. External browser/auth-session evidence remains separate if required for launch scope. | QA Lead, Engineering Lead, Product Owner |
| Notification delivery or fallback | Accepted for selected bounded pilot via manual fallback; selected-host Telegram preflight passed after T279; one selected-host queued Telegram delivery passed after T280; controlled retryable/terminal classification passed after T281. T200 provides local notification foundation, T222 defines the fallback packet, T244 records Admin-owned manual fallback SLA/escalation with redacted customer-facing and ops-facing sample events, T264 selects that fallback as the notification path for the selected pilot, and T270/T271/T272/T273/T274/T275/T276 confirm selected support-window coverage with no launch-critical notification payloads read. T278 adds Telegram delivery commands. T279 records selected-host Telegram preflight PASS with protected secret-file metadata, redacted payload, no secrets printed, and zero process argv secret matches excluding the checker. T280 records one queued Telegram worker delivery PASS for notification `10000` and job `10000`. T281 records local-fake-API classification PASS: HTTP 500 produced `failed_retryable`/`retried=1`, HTTP 400 produced `failed_terminal`/`terminal_failed=1`, cleanup cancelled the retryable artifact, and claimable notification jobs returned to `0`. | Use T244/T264 manual fallback for the historical selected-pilot window. Use T279/T280/T281 as selected-host Telegram reachability, redaction, one-message queued delivery, and controlled failure classification evidence. Before broader production notification approval, capture scope-specific owner approval for Telegram primary-path operation. | Admin as Ops Lead, Support Owner, and Security Owner for selected-pilot fallback and selected-host Telegram evidence |
| Launch-day owners | Complete for selected bounded non-production pilot with single-owner risk. T241 records the user-provided assignment that `Admin` owns Product, Engineering, QA, Ops, Finance, Security, Support, and Provider launch-day roles. T268 records the selected launch window, Admin direct escalation path, and explicit acceptance that one person owns all role decisions for the selected launch scope. | Repeat owner assignment and risk acceptance if the launch scope, window, coverage, owner, provider account, quota, customer-data classification, or notification path changes. | Admin as Product Owner and Engineering Lead |
| Target-environment verification | Complete for selected bounded non-production pilot. T230-T254 prove target deploy/build, health checks, non-mutating registry activation, one Cloudmini lifecycle cleanup activation, target top-up review, local API auth/RBAC, credential reveal audit/redaction, balanced finance reconciliation, cloudflared token-file handling, and selected-host secret-store proof. T266 recovered the selected target route and passed `dev-target-auth-rbac` against `https://billing.resvn.net/backend` with no provider or money mutation routes called. T267 promotes that runtime to `billing-api` and `billing-frontend` systemd services using protected env files outside git. T268 enrolls the named selected target Admin 2FA, verifies a 2FA-satisfied admin route, and reruns domain auth/RBAC smoke successfully. T269 completes the selected launch-window finance reconciliation smoke with daily status `balanced` and zero mismatch counts. T270 confirms launch-window service health, selected domain health/readiness, process command-line secret-pattern checks, and protected secret/token file metadata. T271, T272, T273, T274, T275, and T276 repeat those checks and rerun finance reconciliation successfully; T276 is the final support-window closeout after `2026-05-19 22:00 Asia/Ho_Chi_Minh`. T290 adds separate client/admin base URL support for target auth smoke. T291 deploys the merged smoke runner to the selected test server and records remote binary PASS through `client.resvn.net/backend` and `billing.resvn.net/backend`, proving split-domain public auth/RBAC without provider or money mutation routes. | Repeat T254/T267/T268/T270/T271/T272/T273/T274/T275/T276 and T291-style split-domain auth proof for any new host/path, service configuration, admin user, or production scope before launch use. Continue daily finance reconciliation and pause on any mismatch. | Admin as Security Owner, Finance Lead, and QA Lead |

Any missing required sign-off keeps the requested scope at NO-GO.

## Evidence Packet

Fill one packet per launch candidate. Store only redacted evidence in git.

```text
Launch candidate ID:
selected-bounded-nonprod-pilot-2026-05-19
Date/time UTC:
2026-05-19 11:00-13:00 UTC
Pilot scope:
Selected bounded non-production pilot only: approved test-server/domain runtime, one active Cloudmini dev/test resource maximum, manual notification fallback, no production customer data, no broader private beta, and no production-like provider provisioning.
Environment:
Selected target dev/test runtime behind billing.resvn.net, client.resvn.net, and reseller.resvn.net.
Evidence collector:
Codex
Final reviewer:
Admin
Decision requested:
GO for selected bounded non-production pilot only; NO-GO for production, broader private beta, broader provider scope, production customer data, and real-provider production-like provisioning outside the approved selected Cloudmini scope.
```

### 1. Real Provider Sandbox

```text
Provider:
Cloudmini V3
Provider owner:
Admin for the selected bounded non-production Cloudmini pilot scope.
Sandbox account reference: redacted
Credential storage path:
/etc/billing/secrets/cloudmini.env on the selected host; values redacted.
Secret-store evidence:
T254 selected-host proof passed after owner-confirmed API key rotation; repeat proof for any new host/path.
Credential scope:
Non-production Cloudmini V3 access for the approved dev/test account.
Quota/cost limit:
One active test resource, no parallel mutating calls, single-dev-resource exposure.
Provider support contact:
Admin direct provider-support owner; upstream provider support channel is kept outside git as a redacted account reference.
Billing plan code:
proxy-static-10gb-monthly
Provider SKU:
Cloudmini `ipv4_dc` group reference redacted.
Sandbox location:
Approved target dev/test provider account.
Timeout policy:
Timeout-after-send maps to `PROVIDER_TIMEOUT_REQUEST_KNOWN` and `manual_review_required`.
Idempotency level:
Duplicate-create with the same idempotency key produced one distinct redacted resource and cleanup success.
Cleanup owner:
Admin for the approved dev/test run and selected pilot same-session cleanup, source disable, and incident/follow-up ownership if cleanup cannot be confirmed.
Real pilot run ID:
T249-duplicate-20260518T032613Z; T249-timeout-20260518T032823Z.
Run result:
PASS for T249 duplicate/timeout, T250 safe error examples, T256 permission-denied evidence, T257 out-of-capacity evidence, T260 rate-limit evidence, T261 provider 5xx evidence, T262 cancel/delete rejected evidence, and T263 provider-owner approval in approved target dev/test.
Redacted evidence link/reference:
docs/03_execution_operations_launch/73_Cloudmini_Idempotency_Evidence_Runbook.md#t249-target-devtest-evidence
T250 safe error evidence is recorded in doc 66. T255 provider-controlled error plan, T256 permission-denied runtime evidence, T257 out-of-capacity runtime evidence, T260 rate-limit runtime evidence, T261 provider 5xx runtime evidence, and T262 cancel/delete rejected runtime evidence are recorded in doc 77.
Usable-status semantics:
T252 approved `running`, `active`, `ready`, and `available` as the only Cloudmini statuses that may activate a Billing service; all other statuses, timeout, and credential-missing outcomes remain fail-closed/manual-review.
Provider owner sign-off:
Admin for the selected bounded non-production Cloudmini pilot scope, including one-resource quota, edge/header policy, SKU mapping, support ownership, usable-status semantics, and cleanup procedure.
Security owner sign-off:
Admin for selected pilot usable-status semantics, cleanup redaction boundary, T254 selected-host secret-store proof, and T263 no-secret provider approval boundary.
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
Admin accepted the T242 staging-equivalent backup/restore scope in T245.
QA sign-off:
Admin accepted the T242 staging-equivalent backup/restore scope in T245.
Security sign-off:
Admin accepted the T242 staging-equivalent backup/restore scope in T245.
```

Pass criteria:

- Source and target are approved non-production or staging-equivalent databases.
- Restore target overwrite is approved before running destructive restore.
- Restored target passes `dev-db` smoke.
- Backup artifact retention or deletion owner is recorded.

### 3. Staging Or Staging-Equivalent Full E2E

```text
Gate ID:
T243-target-20260517T140625Z
Environment:
Approved target test server, staging-equivalent dev run.
DB/API classification:
Temporary DB billing_t243_e2e_20260517140625 and local API http://127.0.0.1:18083; no production data.
Provider mode: fake/manual/real sandbox
fake provider for provisioning worker fulfillment; no real provider routes called.
Frontend target:
Next standalone browser smoke at http://127.0.0.1:3120 in demo portal mode.
Backend result:
PASS: taskguard, make test, contract guard, error code guard, and build.
DB smoke result:
PASS: dev-db applied 25 migrations and passed 20 checks.
API/RBAC smoke result:
PASS: dev-api passed 35 checks including RBAC negative checks.
Billing mutation result:
PASS: top-up, checkout, wallet payment, provisioning job, fake-provider worker fulfillment, and active service verified.
Renewal path result:
PASS: service display 10000 renewed; renewal invoice 10002 paid; renewal transaction 10001 posted; renewal ledger 10002 recorded; service term increased.
Lifecycle job result:
PASS for provisioning worker lifecycle boundary in fake-provider mode; real-provider lifecycle remains covered by provider-specific evidence only.
Frontend smoke result:
PASS: npm ci, audit, sensitive-text, lint, build, and admin browser smoke.
Audit/redaction result:
PASS: wallet.topup.approved, invoice.wallet_paid, service.renewed, and renewal invoice.wallet_paid audit checks passed; evidence omits raw DSNs, tokens, provider payloads, service credentials, and dump files.
Exception requested: yes/no
yes, staging-equivalent scope because the target run uses a temporary dev DB and fake provider instead of shared staging with real provider.
Exception owner and reason:
Admin assigned by T241; reason is to prove full E2E and T206 renewal safely without real provider mutation or production data.
QA sign-off:
Admin accepted the T243 staging-equivalent full E2E and renewal scope in T245.
Engineering sign-off:
Admin accepted the T243 staging-equivalent full E2E and renewal scope in T245.
Product sign-off:
Admin accepted the T243 staging-equivalent full E2E and renewal scope in T245.
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
manual fallback for the selected pilot window; Telegram selected-host preflight after T279; one queued Telegram delivery after T280; controlled retryable/terminal worker classification after T281.
Launch-critical events covered:
Top-up status, provisioning failure/manual review, service lifecycle, password reset, support/abuse critical events by manual fallback procedure.
Credential/secret storage path: redacted secret-store reference only
Manual fallback used no notification credential. T279 Telegram preflight used `/etc/billing/secrets/telegram.env` on the selected host; values were redacted, file mode was `600`, owner was `root:root`, and only key names were recorded.
Successful delivery evidence:
T244 manual fallback drill review passed; T264 selects manual fallback as the notification path for the selected bounded pilot; sampled customer-facing top-up approved event from T235 and ops-facing Cloudmini manual-review event from T232. T270 checked selected support-window coverage at 2026-05-19 18:25 Asia/Ho_Chi_Minh, T271 repeated the checkpoint at 2026-05-19 18:55 Asia/Ho_Chi_Minh, T272 repeated it at 2026-05-19 19:36 Asia/Ho_Chi_Minh, T273 repeated it at 2026-05-19 20:01 Asia/Ho_Chi_Minh, T274 repeated it at 2026-05-19 20:14 Asia/Ho_Chi_Minh, T275 repeated it at 2026-05-19 20:44 Asia/Ho_Chi_Minh, and T276 completed final support-window closeout at 2026-05-19 22:06 Asia/Ho_Chi_Minh. All seven runs read only notification counts: launch-critical notification total `0`; no notification payloads, customer data, DSNs, tokens, provider payloads, or credentials were read or recorded. T279 ran `notification-telegram-preflight` with `APP_ENV=staging` on 2026-05-20 after the owner corrected the channel target: result `PASS`, `telegram_api_called=yes`, `message_payload_redacted=yes`, `secrets_printed=no`, and process argv secret check returned `0` matches excluding the checker. T280 ran `notification-telegram-once` with worker `t280-telegram-drill`, batch-size `1`, and timeout `60s` against one dev/test queued Telegram notification: worker result `claimed=1 succeeded=1 retried=0 manual_review=0 terminal_failed=0 cancelled=0`; notification display `10000` status `sent`; job display `10000` status `succeeded`; one succeeded attempt row; post-run claimable Telegram and generic notification jobs `0`.
Failure/retry evidence:
If Admin misses SLA or cannot deliver a fallback message, pilot pauses and the related event remains in manual review until Security/Ops review. T281 ran `notification-telegram-once` against a local fake Telegram API endpoint, not real Telegram: HTTP 500 produced worker `claimed=1 succeeded=0 retried=1 manual_review=0 terminal_failed=0 cancelled=0`, notification `10001` failed with `telegram_http_500`, job `10001` `failed_retryable`, and one failed_retryable attempt row; HTTP 400 produced worker `claimed=1 succeeded=0 retried=0 manual_review=0 terminal_failed=1 cancelled=0`, notification `10002` failed with `telegram_http_400`, job `10002` `failed_terminal`, and one failed_terminal attempt row. The retryable artifact was cancelled after evidence capture, post-cleanup claimable Telegram/generic notification jobs were `0`, fake API received exactly one 500 and one 400 call, and argv checks found `0` Telegram token/chat ID/DB_DSN matches before and after the run. Broader production approval still requires scope-specific owner approval before Telegram becomes the sole primary path.
Manual fallback owner:
Admin
Manual fallback SLA:
P0 acknowledgement 15 minutes; P0 customer contact 30 minutes; P1 customer contact 4 business hours.
Escalation path:
Admin direct launch channel; single-owner escalation accepted for this fallback scope.
Support owner sign-off:
Admin
Ops sign-off:
Admin
Security sign-off:
Admin
```

Pass criteria:

- At least top-up status, provisioning failure/manual review, service lifecycle, password reset, and support/abuse critical events have delivery or fallback coverage.
- Failure mode and retry/manual fallback are tested or explicitly accepted.
- Notification payloads are redacted and contain no credentials or reset tokens.
- If using manual fallback, complete `docs/03_execution_operations_launch/72_Notification_Delivery_And_Manual_Fallback_Runbook.md` with named owner, SLA, escalation channel, sampled events, and redacted evidence.
- If using Telegram as the primary broader notification path, include selected-host secret metadata, a redacted preflight PASS, queued launch-critical event evidence or signed exception, and failure/retry evidence.

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
Admin direct launch channel; single-person escalation accepted by user statement on 2026-05-17 and selected-scope GO packet on 2026-05-19.
Launch window:
2026-05-19 18:00-20:00 Asia/Ho_Chi_Minh (2026-05-19 11:00-13:00 UTC). Support/Ops coverage for the selected pilot window plus two hours: 2026-05-19 18:00-22:00 Asia/Ho_Chi_Minh.
Owner availability confirmed:
Yes for owner assignment by user statement on 2026-05-17: "1 mình tao cân hết. Admin"; selected launch window and single-owner risk accepted in T268 on 2026-05-19.
Single-owner risk:
Accepted for owner assignment, T244/T264 selected-pilot notification manual fallback, T245 target evidence sign-off, T249 target dev/test duplicate/timeout evidence, T252 usable-status semantics, T253 cleanup owner/procedure, T254 selected-host secret-store proof, T256 permission-denied evidence, T257 out-of-capacity evidence, T260 rate-limit evidence, T261 provider 5xx evidence, T262 cancel/delete rejected evidence, T263 selected provider-owner approval, and T268 selected launch-window/Admin 2FA GO packet. This does not waive automated production notification delivery evidence for broader launch or approve provider scope beyond the selected bounded non-production Cloudmini pilot.
```

Pass criteria:

- Every role has a named human owner before launch.
- Each owner has accepted their launch-day responsibility.
- Escalation channel and launch window are recorded.

### 6. Target-Environment Verification

```text
Auth/session target check:
T236 PASS on the approved test server local API. Client seed login set an HttpOnly `billing_session`; cookie-only `/client/catalog` passed without `X-Actor-*` dev helper headers. Invalid session returned `auth.session_invalid`. Missing actor returned `auth.actor_required`.
T265 domain check:
Initial result was BLOCKED: `APP_ENV=dev go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 20s dev-target-auth-rbac` failed before auth assertions with `target auth login expected HTTP 200, got 500`. T266 recovered the selected dev/staging-equivalent target route by applying 25 migrations and idempotent dev seed to the empty `billing_smoke` database, starting the Billing API on `127.0.0.1:8080`, and serving frontend `/backend` proxy on `3000` using root-only `/run` env files. Follow-up health checks returned HTTP `200` for local API health, local frontend root, local `/backend/healthz`, domain root, and domain `/backend/healthz`. `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac` then passed: client session cookie-only access, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three low-permission RBAC denials. No raw cookies, session tokens, passwords, DSNs, provider payloads, or credentials were printed, and the smoke reported no provider or money mutation routes called.
T267 runtime service check:
PASS. `/opt/Billing` was updated to `origin/main` merge commit `87a6584`, `bin/api` and the Next.js standalone frontend bundle were built, and `billing-api.service` plus `billing-frontend.service` were installed, enabled, and active. Both services run as `billing-svc`; API `ExecStart` is `/opt/Billing/bin/api`; frontend `ExecStart` is `/usr/bin/node /opt/Billing/frontend/.next/standalone/server.js`; service command lines contain no raw DSN, token, cookie, password, provider payload, or credential. Protected env file metadata: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac` passed after service promotion with no provider or money mutation routes called. No raw env file contents, DSNs, tokens, cookies, provider payloads, or credentials were printed.
T291 domain-aware target auth check:
PASS. T290 added separate client/admin base URL support to `dev-target-auth-rbac`; T291 deployed the merged code to the selected test server, rebuilt `/opt/Billing/bin/smoke`, verified local API/frontend health and public `billing`, `client`, and `reseller` health endpoints, and ran the remote smoke binary with `client_base_url=https://client.resvn.net/backend` and `admin_base_url=https://billing.resvn.net/backend`. The run passed cookie-only client catalog access, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three RBAC denials. It reported `domain_aware_base_urls=pass`, `provider_mutation_routes_called=no`, and `money_mutation_routes_called=no`; no raw cookies, session tokens, passwords, DSNs, provider payloads, or credentials were printed.
T270 launch-window ops health check:
PASS at 2026-05-19 18:25 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded.
T271 support-window checkpoint:
PASS at 2026-05-19 18:55 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This is an in-window checkpoint, not final support-window closeout.
T272 support-window checkpoint:
PASS at 2026-05-19 19:36 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This is an in-window checkpoint, not final support-window closeout.
T273 launch-window end checkpoint:
PASS at 2026-05-19 20:01 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This records the selected launch-window boundary, not final support-window closeout.
T274 support-extension checkpoint:
PASS at 2026-05-19 20:14 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This is a support-extension checkpoint, not final support-window closeout.
T275 support-extension checkpoint:
PASS at 2026-05-19 20:44 Asia/Ho_Chi_Minh. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This is a support-extension checkpoint, not final support-window closeout.
T276 final support-window closeout:
PASS at 2026-05-19 22:06 Asia/Ho_Chi_Minh, after the approved selected support window ended. `billing-api`, `billing-frontend`, and `cloudflared` were active and enabled with main PIDs present. Secret-pattern checks on process command lines reported `none` for DSN, token, password, credential, bearer, Cloudmini token, and cloudflared token patterns without printing command lines. Protected metadata remained restrictive: `/etc/billing/secrets` mode `700` owner `root:root`; `/etc/billing/secrets/billing-api.env` mode `600` owner `root:root`; `/etc/billing/secrets/billing-frontend.env` mode `600` owner `root:root`; `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Domain checks returned HTTP `200` for `https://billing.resvn.net/`, `https://billing.resvn.net/backend/healthz`, `https://billing.resvn.net/backend/readyz`, `https://client.resvn.net/`, and `https://reseller.resvn.net/`. Read-only target finance reconciliation remained `balanced` with wallets/invoices/payments checked `1/1/1`, zero wallet/invoice/duplicate-payment mismatches, and no money or provider mutation routes called. Launch-critical notification total was `0`; no payload was read. No file contents, command lines, DSNs, tokens, cookies, provider payloads, credentials, or customer data were printed or recorded. This closes the selected bounded non-production pilot support window.
Admin 2FA enrollment/enforcement:
T236 PASS for enforcement on the approved test server local API. Platform staff login returned an unsatisfied 2FA session, and cookie-only admin provider-readiness access returned `auth.2fa_required`.
T268 PASS for selected target Admin enrollment and enforcement. Before enrollment, redacted target metadata showed `admin@local.billing|platform_staff|active|disabled|totp_not_enabled`. The selected target API setup returned HTTP `201`, verify returned HTTP `200`, a 2FA-satisfied admin provider-readiness route returned HTTP `200`, and post-enrollment metadata showed `admin@local.billing|platform_staff|active|enabled|totp_enabled`. The temporary enrollment script was removed. `APP_ENV=dev GOFLAGS=-buildvcs=false go run ./cmd/smoke -base-url https://billing.resvn.net/backend -timeout 60s dev-target-auth-rbac` passed after enrollment: client session cookie-only access, admin 2FA gate, invalid session denial, missing actor denial, tenant mismatch denial, and three low-permission RBAC denials. No TOTP secret, TOTP code, cookie, session token, password, DSN, provider payload, or credential was printed or recorded.
Credential reveal audit access:
T237 PASS on the approved test server local API. The smoke created/refreshed one encrypted dev/test credential fixture for service display 43001, logged in as the seeded client, revealed through the client credential reveal route with only the HttpOnly session cookie, verified `Cache-Control: no-store` and `Pragma: no-cache`, verified `last_revealed_by` and rate-limit state, and recorded audit display 10017 for `credential.revealed`. The evidence excludes plaintext credentials, encrypted payloads, raw credential IDs, session tokens, cookies, DSNs, provider payloads, and provider credentials.
Top-up review target check: T235 PASS on the approved test server. Approve top-up display 10003 posted ledger display 10005 and audit display 10015. Reject top-up display 10004 posted no ledger and audit display 10016. Wallet delta was 111 minor units. Provider side effects were none.
Secret/key handling review:
T240 PASS on the approved test server for target dev/test secret handling evidence. `/opt/Billing/.env.dev` is mode 640 with owner `root:billing-svc`; `/opt/cred-cloudmini-dev.env` is mode 600 with owner `root:root`; `/etc/cloudflared/tunnel.token` is mode 600 with owner `root:root`. Cloudflared was changed from token flag usage to `--token-file /etc/cloudflared/tunnel.token`, restarted, and verified active with no token present in process arguments. Reachability returned HTTP 200 for `http://localhost:3000`, `https://billing.resvn.net`, `https://reseller.resvn.net`, and `https://client.resvn.net`. T247 rechecked the same target metadata on 2026-05-18: the three files kept the same restrictive modes and owners, cloudflared remained active with `--token-file` and no token in process arguments, and `billing-api` plus `billing-worker` were active. T254 recorded owner-confirmed API key rotation, canonical Cloudmini path `/etc/billing/secrets/cloudmini.env` mode 600 owner `root:root`, directory mode 700 owner `root:root`, required Cloudmini keys present, `DB_DSN` absent, and no exact `--token` arg in the running cloudflared process. T267 recorded long-lived selected-runtime env metadata for `/etc/billing/secrets/billing-api.env` and `/etc/billing/secrets/billing-frontend.env`, both mode `600` owner `root:root`, with API/frontend systemd command lines free of raw secrets. No raw token, DSN, API key, cookie, credential, provider payload, or file contents were printed or recorded. T245 records Admin Security sign-off for the earlier target dev/test evidence; T254/T267 record selected-host secret-store proof only.
Finance reconciliation owner run:
T238 evidence captured on the approved test server local API/DB. The read-only smoke selected transaction display 51001, invoice display 44001, wallet display 41001, and ledger display 50002; payment reconciliation list/detail returned matching public display evidence; daily reconciliation for 2026-04-23 initially returned `mismatched` with one wallet mismatch. T239 traced the root cause to dev/test wallet projection drift from an inconsistent seed baseline and later smoke runs, fixed the seed baseline, repaired wallet display 41001 from posted ledger source-of-truth with audit display 10018, and reran the target smoke. The rerun returned `balanced` with wallets checked 2, wallet mismatches 0, invoices checked 1, invoice mismatches 0, payments checked 1, and duplicate payment references 0. The repair inserted no ledger rows, updated no posted ledger rows, called no money/provider mutation routes, and printed no secrets. T245 records Admin Finance sign-off for this evidence.
T269 selected launch-window finance reconciliation:
PASS at 2026-05-19 18:11 Asia/Ho_Chi_Minh using the selected dev/test API and database. `APP_ENV=dev API_BASE_URL=https://billing.resvn.net/backend GOFLAGS=-buildvcs=false go run ./cmd/smoke -timeout 90s dev-target-finance-reconciliation` returned transaction display `51001`, invoice display `44001`, wallet display `41001`, ledger display `50002`, daily date `2026-04-23`, daily status `balanced`, wallets checked `1`, invoices checked `1`, payments checked `1`, wallet mismatches `0`, invoice mismatches `0`, duplicate payment references `0`, `money_mutation_routes_called=no`, and `provider_mutation_routes_called=no`. Output intentionally excluded raw transaction IDs, invoice IDs, wallet IDs, ledger IDs, actor IDs, session tokens, cookies, DSNs, provider payloads, and credentials.
Cross-tenant negative check: T236 PASS on the approved test server local API. A mismatched actor tenant was denied with `tenant.context_mismatch`, and low-permission RBAC checks returned `auth.permission_denied`.
Domain cross-tenant negative check: T266 PASS through `https://billing.resvn.net/backend`; the domain smoke reached the assertion and denied tenant mismatch without dev actor headers. T291 repeated the selected public-domain path with separate client/admin base URLs and denied the tenant mismatch through `client.resvn.net/backend`.
Support coverage check:
T244 PASS for owner-approved manual fallback readiness. T264 accepts this as the selected bounded pilot notification path. Admin owns Support/Ops/Security fallback decisions, coverage is the approved pilot window plus 2 hours, P0 acknowledgement SLA is 15 minutes, P0 customer-contact SLA is 30 minutes, and P1 customer-contact SLA is 4 business hours. Evidence samples include T235 top-up display 10003 and T232 provisioning manual-review evidence. T270 checked the selected support window at 2026-05-19 18:25 Asia/Ho_Chi_Minh, T271 repeated the checkpoint at 2026-05-19 18:55 Asia/Ho_Chi_Minh, T272 repeated it at 2026-05-19 19:36 Asia/Ho_Chi_Minh, T273 repeated it at 2026-05-19 20:01 Asia/Ho_Chi_Minh, T274 repeated it at 2026-05-19 20:14 Asia/Ho_Chi_Minh, T275 repeated it at 2026-05-19 20:44 Asia/Ho_Chi_Minh, and T276 completed final closeout at 2026-05-19 22:06 Asia/Ho_Chi_Minh; launch-critical notification total was `0` in all seven checks, so no live fallback delivery was required. Production SMTP/Telegram delivery remains unproven for broader launch.
Residual risks:
Single-person support/ops/security coverage for the selected support window; no automated production notification delivery for broader launch; day-one and final closeout finance reconciliation passed for the selected launch window, but daily reconciliation must continue and launch must pause on any future mismatch, health/readiness failure, secret exposure signal, or support SLA breach.
Security sign-off:
Admin for T237/T240/T247 target dev/test evidence, T244/T264 selected-pilot manual fallback scope, T254/T267/T270/T271/T272/T273/T274/T275/T276 selected-host secret-store proof, and T268 selected target Admin 2FA enrollment/enforcement; repeat the proof for any new host/path, admin user, or launch scope before launch use.
Finance sign-off:
Admin for T239 balanced target reconciliation evidence, T269 selected launch-window finance reconciliation, and T271/T272/T273/T274/T275/T276 support-window finance reruns. Daily reconciliation must continue during pilot operation.
Ops sign-off:
Admin for T267 protected systemd runtime evidence, T270 selected launch-window ops health evidence, T271/T272/T273/T274/T275 in-window support checkpoint evidence, and T276 final support-window closeout evidence. Pause the selected pilot on health/readiness failure, secret exposure signal, or unavailable support coverage.
QA sign-off:
Admin for T242/T243/T245 target staging-equivalent evidence scope.
```

Pass criteria:

- Target auth/session configuration is verified outside dev-only actor headers.
- Admin 2FA is enrolled and enforced in the target environment.
- Credential reveal audit is visible to authorized operators and redacted elsewhere.
- Finance owner runs or reviews reconciliation evidence.
- Cross-tenant negative tests fail safely.

## Final Decision Rule

For any requested scope:

1. Re-run `docs/03_execution_operations_launch/69_Pilot_Go_No_Go_Record.md`.
2. Keep that scope at NO-GO if any P0 section is missing, unreviewed, or based on assumptions.
3. Use CONDITIONAL GO only for non-P0 exceptions with a named owner, mitigation, expiry date, and rollback path.
4. Use GO only when every P0 gate has passing evidence and required owner sign-off.

T268 satisfies this rule for the selected bounded non-production pilot only. Production, broader private beta, broader provider scope, production customer data, and real-provider production-like provisioning outside the approved selected Cloudmini scope remain NO-GO until separately proven and signed off.
