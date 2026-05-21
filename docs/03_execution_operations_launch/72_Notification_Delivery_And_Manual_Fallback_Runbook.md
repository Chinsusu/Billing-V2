# 72 - Notification Delivery And Manual Fallback Runbook

**Date:** 2026-05-16
**Scope:** Launch-critical notification delivery proof or approved manual fallback evidence.
**Decision:** Manual fallback is launch-ready for an owner-approved pilot after T244. T279 adds selected-host Telegram preflight proof, T280 proves one queued Telegram delivery, and T281 proves controlled retryable/terminal Telegram worker classification. Broader production notification delivery still needs scope-specific owner approval before using Telegram as a sole primary path.

## Current State

T200 added the notification foundation:

- notification schema and module service;
- redacted launch-critical event builders;
- local/dev delivery runner;
- tests for notification creation and secret redaction.

T277 exposes that local/dev runner through `cmd/worker notification-local-once` and `cmd/worker notification-local-loop`. These commands are intentionally limited to `APP_ENV=local` or `APP_ENV=dev` because local delivery only marks notification jobs as sent after the local handler runs; it does not send SMTP, Telegram, webhook, or customer-facing external messages.

This is not production delivery proof. Before GO, launch evidence must show either:

- production SMTP/Telegram delivery proof for launch-critical events; or
- an approved manual fallback with owner, SLA, escalation path, and redacted evidence.

T244 records the second path for the current pilot scope: Admin owns Support, Ops, and Security fallback decisions; fallback messages use the Admin direct launch channel; sample events are redacted dev/test evidence references; and no external production delivery channel is claimed. T279 later proves the selected-host Telegram channel can receive a redacted preflight message, but it does not replace the selected-pilot manual fallback packet for historical pilot-window evidence. T280 later proves one selected-host queued Telegram notification can be delivered by the real worker. T281 proves the worker classifies a controlled retryable Telegram API failure and a controlled terminal Telegram API failure correctly using a local fake API endpoint, without calling Telegram or printing secrets.

## Local/Dev Delivery Worker

Use these commands only for local/dev notification job plumbing checks:

```bash
APP_ENV=dev go run ./cmd/worker notification-local-once -dsn <redacted-dsn> -worker-id notification-local-1
APP_ENV=dev go run ./cmd/worker notification-local-loop -dsn <redacted-dsn> -worker-id notification-local-1
```

Required evidence for local/dev checks:

- command used and environment name;
- job summary counts only: claimed, succeeded, retried, manual review, terminal failed, cancelled;
- no raw `payload_redacted`, credentials, tokens, DSNs, SMTP secrets, Telegram tokens, provider payloads, cookies, or customer data.

Do not use these commands as proof of automated production notification delivery. Production SMTP/Telegram proof still requires real channel configuration, secret owner approval, a safe delivery target, redacted delivered-message evidence, failure/retry evidence, and owner sign-off.

## Telegram Delivery Worker

T278 adds a Telegram delivery path for queued `telegram` notifications. It uses a channel-specific job type so the Telegram worker does not claim dashboard/email notification jobs.

Required secret/config keys:

```text
TELEGRAM_BOT_TOKEN
TELEGRAM_CHAT_ID
TELEGRAM_API_BASE_URL optional; defaults to https://api.telegram.org
BILLING_TELEGRAM_DELIVERY_PRODUCTION_APPROVED=yes required only for APP_ENV=production
```

Commands:

```bash
APP_ENV=staging go run ./cmd/worker notification-telegram-preflight
APP_ENV=staging go run ./cmd/worker notification-telegram-once -dsn <redacted-dsn> -worker-id notification-telegram-1
APP_ENV=staging go run ./cmd/worker notification-telegram-loop -dsn <redacted-dsn> -worker-id notification-telegram-1
```

Evidence may include only:

- command name, environment, and timestamp;
- preflight result fields: `telegram_api_called=yes`, `message_payload_redacted=yes`, `secrets_printed=no`;
- worker summary counts;
- notification display ID, event type, template key, and correlation ID if safe.

Evidence must not include bot token, chat ID, raw Telegram request/response body, raw `payload_redacted`, credentials, reset tokens, DSNs, provider payloads, cookies, or customer data.

## Telegram Evidence References

Detailed selected-host Telegram evidence is tracked in `docs/03_execution_operations_launch/78_Notification_Telegram_Evidence_Appendix.md`:

- T279 selected-host Telegram preflight evidence.
- T280 queued Telegram delivery drill evidence.
- T281 controlled retryable/terminal Telegram worker classification evidence.

Keep future raw evidence packets in that appendix or a successor appendix so this runbook remains focused on operational procedure.

## Launch-Critical Events

Treat these as launch-critical during pilot:

| Event | Audience | Fallback priority | Required response |
|---|---|---:|---|
| `auth.password_reset` | client/reseller/admin | P0 | Notify account owner through approved channel only; never expose reset token material. |
| `wallet.topup.approved` | client/reseller | P1 | Notify customer before they retry payment or contact support. |
| `wallet.topup.rejected` | client/reseller | P1 | Notify customer with a safe review note and support path. |
| `provisioning.failed` | support/ops/provider | P0 | Stop blind retry; verify provider state. |
| `provisioning.manual_review` | support/ops/provider | P0 | Assign manual review owner and SLA. |
| `service.lifecycle` | client/reseller/support | P1 | Notify customer of lifecycle transition; do not send credentials in the message. |
| `service.expiring` | client/reseller | P1 | Notify customer before expiry window closes. |
| `service.expired` | client/reseller/support | P1 | Notify customer and support of grace/suspend policy. |
| `service.suspended` | client/reseller/support | P0 | Notify customer and support with reason and appeal path. |
| `service.terminated` | client/reseller/support | P0 | Notify customer and support; verify no residual billable provider resource. |

Provider-down and abuse events are also launch-critical when they are emitted or manually tracked. If they are not present in `notifications`, record the external alert or support case reference in the fallback evidence packet.

## Manual Fallback Preconditions

Manual fallback may be used only when production delivery proof is unavailable or degraded.

Required before declaring fallback ready:

```text
Fallback ID:
Environment:
Support owner:
Ops owner:
Security owner:
Escalation channel reference:
Support coverage window:
P0 acknowledgement SLA:
P0 customer-contact SLA:
P1 customer-contact SLA:
Evidence storage reference:
Reviewer sign-off:
```

Minimum pilot SLA:

- `P0 acknowledgement SLA`: 15 minutes.
- `P0 customer-contact SLA`: 30 minutes when customer-facing.
- `P1 customer-contact SLA`: 4 business hours.
- `Support coverage window`: entire pilot window plus 2 hours after pilot close.

If any owner or SLA field is unassigned, fallback is not ready and launch remains `NO-GO`.

## Fallback Workflow

1. Identify launch-critical notification records from the app, admin screen, or a read-only DB query.
2. Confirm the event is still relevant and has not already been delivered through a production channel.
3. Assign the event to the support or ops owner based on priority.
4. Send a manual message through the approved fallback channel.
5. Record redacted evidence with display IDs, event type, channel used, timestamp, owner, and result.
6. Escalate if the SLA is missed or the message cannot be delivered.

Read-only DB inspection may use this shape when an approved target DB access path exists:

```sql
SELECT
  display_id,
  event_type,
  channel,
  priority,
  recipient_group,
  reference_type,
  status,
  last_error_code,
  created_at,
  updated_at
FROM notifications
WHERE event_type IN (
  'wallet.topup.approved',
  'wallet.topup.rejected',
  'auth.password_reset',
  'provisioning.failed',
  'provisioning.manual_review',
  'service.lifecycle',
  'service.expiring',
  'service.expired',
  'service.suspended',
  'service.terminated'
)
ORDER BY created_at DESC
LIMIT 100;
```

Do not include `payload_redacted` in pasted evidence unless Security explicitly approves it for the launch packet. Even redacted payload can contain sensitive business context.

## Safe Message Rules

Manual fallback messages must not include:

- plaintext service credentials;
- provider API keys, tokens, auth headers, or raw provider payloads;
- database DSNs;
- raw payment proof URLs;
- private abuse evidence;
- backend UUIDs as the user-facing label when a display ID is available.

Use display IDs and safe summaries:

```text
Order #<display_id>
Service #<display_id>
Top-up #<display_id>
Notification #<display_id>
Status: <safe status>
Next step: <safe action>
Support reference: <correlation_id if already visible to support>
```

## Manual Fallback Evidence Packet

Store one packet per launch candidate or fallback drill:

```text
Fallback ID:
Date/time UTC:
Environment:
Evidence collector:
Support owner:
Ops owner:
Security owner:
Coverage window:
Escalation channel reference:
P0 acknowledgement SLA:
P0 customer-contact SLA:
P1 customer-contact SLA:
Events sampled:
Delivery channel used:
Delivery result:
SLA result:
Redacted evidence reference:
Open exceptions:
Support owner sign-off:
Ops owner sign-off:
Security owner sign-off:
Decision: PASS / FAIL
```

Pass criteria:

- All owner fields are named.
- All SLA fields are filled and accepted.
- At least one launch-critical customer-facing event and one ops-facing event are sampled, or an owner-approved exception explains why no event exists.
- Evidence uses display IDs and safe summaries only.
- No credential, token, DSN, raw provider payload, raw payment proof, or private abuse evidence appears in the packet.
- Missed SLA or failed manual delivery keeps the launch decision at `NO-GO`.

## Escalation

Escalate immediately to Security and Ops if:

- any fallback message exposes a secret, credential, raw provider payload, or private abuse evidence;
- a P0 event is not acknowledged within SLA;
- provider down/manual review notifications cannot reach the provider or ops owner;
- customer-facing service suspended/terminated notification cannot be delivered;
- support coverage is unavailable during the pilot window.

Escalation result must be recorded in the evidence packet before GO can be reconsidered.

## T244 Manual Fallback Evidence

```text
Fallback ID:
T244-manual-fallback-20260517
Date/time UTC:
2026-05-17T14:30Z
Environment:
Approved target test server evidence plus owner-approved pilot fallback procedure.
Evidence collector:
Codex
Support owner:
Admin
Ops owner:
Admin
Security owner:
Admin
Coverage window:
Entire approved pilot window plus 2 hours after pilot close; no launch window is approved yet.
Escalation channel reference:
Admin direct launch channel.
P0 acknowledgement SLA:
15 minutes.
P0 customer-contact SLA:
30 minutes for customer-facing P0 events.
P1 customer-contact SLA:
4 business hours.
Events sampled:
Customer-facing P1 top-up approved event from T235: top-up display 10003, ledger display 10005, audit display 10015.
Ops-facing P0 provisioning manual-review event from T232: Cloudmini create reached provider status creating and Billing moved the job to manual review before same-session cleanup.
Delivery channel used:
Manual fallback through Admin direct launch channel. No SMTP, Telegram, or automated production delivery channel was used.
Delivery result:
PASS for manual fallback drill review. No external customer message was sent because the sampled events are dev/test evidence.
SLA result:
PASS for owner-approved SLA values; live SLA measurement starts only when an approved pilot window starts.
Redacted evidence reference:
T235 and T232 task logs plus this T244 packet; no payload, credential, token, DSN, provider ID, provider payload, cookie, reset token, or customer data is recorded.
Open exceptions:
Production SMTP delivery and queued Telegram launch-critical event delivery remain unproven. If Admin is unavailable or misses SLA, pilot must pause and launch-critical events stay in manual review.
Support owner sign-off:
Admin, by T241 owner assignment and T244 fallback acceptance.
Ops owner sign-off:
Admin, by T241 owner assignment and T244 fallback acceptance.
Security owner sign-off:
Admin, by T241 owner assignment and T244 fallback acceptance.
Decision:
PASS for manual fallback readiness; T279 separately proves selected-host Telegram preflight reachability, T280 proves one selected-host queued Telegram delivery, and T281 proves controlled retryable/terminal Telegram worker classification. Broader primary-path approval still requires scope-specific owner approval.
```

Safe message samples approved for manual fallback:

```text
Top-up #10003 was approved and wallet credit was recorded. If the balance is not visible, contact support through the approved pilot channel.
Provisioning for the pilot service needs manual review before retry. Do not retry provider create until provider state is verified.
```

These samples intentionally omit credentials, reset tokens, raw provider references, raw payment proof, backend UUIDs, and private support details.
