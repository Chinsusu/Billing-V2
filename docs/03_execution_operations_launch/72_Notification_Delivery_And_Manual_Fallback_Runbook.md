# 72 - Notification Delivery And Manual Fallback Runbook

**Date:** 2026-05-16
**Scope:** Launch-critical notification delivery proof or approved manual fallback evidence.
**Decision:** Not launch-ready until production delivery proof or the manual fallback packet below is completed with named owners.

## Current State

T200 added the notification foundation:

- notification schema and module service;
- redacted launch-critical event builders;
- local/dev delivery runner;
- tests for notification creation and secret redaction.

This is not production delivery proof. Before GO, launch evidence must show either:

- production SMTP/Telegram delivery proof for launch-critical events; or
- an approved manual fallback with owner, SLA, escalation path, and redacted evidence.

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
