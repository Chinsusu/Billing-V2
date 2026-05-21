# 78 - Notification Telegram Evidence Appendix

**Date:** 2026-05-21
**Scope:** Redacted Telegram notification delivery evidence moved out of the operational runbook to keep the runbook concise and below file-size limits.
**Parent runbook:** `docs/03_execution_operations_launch/72_Notification_Delivery_And_Manual_Fallback_Runbook.md`

## Redaction Boundary

Evidence in this appendix may include command names, environment labels, safe display IDs, worker summary counts, status codes, redacted error codes, and file metadata.

Evidence must not include raw `payload_redacted`, customer data, DB DSNs, provider payloads, cookies, credentials, reset tokens, Telegram token values, chat ID values, raw Telegram request/response bodies, tenant UUIDs, job UUIDs, notification UUIDs, or command lines containing secrets.

## T279 Telegram Preflight Evidence

```text
Evidence ID:
T279-telegram-preflight-20260520
Date/time UTC:
2026-05-20T00:20Z
Environment:
Selected host with APP_ENV=staging.
Evidence collector:
Codex
Secret/config path:
/etc/billing/secrets/telegram.env on the selected host; values redacted.
Secret-file metadata:
mode 600, owner root:root.
Config keys verified:
TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID, TELEGRAM_API_BASE_URL.
Command:
notification-telegram-preflight
Result:
PASS
Preflight fields:
telegram_api_called=yes
message_payload_redacted=yes
secrets_printed=no
Process argv secret check:
0 matches excluding the checker process.
Payload/customer data:
No notification payload, customer data, DSN, provider payload, cookie, credential, token value, chat ID value, or raw Telegram response body was printed or recorded.
Open exceptions:
This is selected-host redacted preflight proof, not queued launch-critical event delivery proof. Broader production notification approval still requires owner-approved scope, runtime worker activation plan, queued event delivery evidence or a signed exception, and failure/retry handling evidence.
Decision:
PASS for selected-host Telegram preflight reachability and redaction boundary.
```

## T280 Queued Telegram Delivery Drill Evidence

```text
Evidence ID:
T280-queued-telegram-delivery-20260520
Date/time UTC:
2026-05-20T00:51Z
Environment:
Selected host with APP_ENV=staging.
Evidence collector:
Codex
Secret/config path:
/etc/billing/secrets/telegram.env on the selected host; values redacted.
Secret-file metadata:
mode 600, owner root:root.
Config keys verified:
TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID, TELEGRAM_API_BASE_URL.
DB access:
Protected service-file handling from /etc/billing/secrets/billing-api.env; DB_DSN value was not printed or passed on the worker argv.
Pre-run queue state:
claimable_telegram_jobs=0
claimable_generic_notification_jobs=0
Notification created:
display_id=10000
event_type=service.lifecycle
template_key=t280.telegram.queued_drill
channel=telegram
recipient_group=ops
correlation_id_present=yes
Job created:
display_id=10000
job_type=notification.deliver.telegram
status_before=queued
Worker command:
notification-telegram-once with worker-id t280-telegram-drill, batch-size 1, timeout 60s.
Worker result:
claimed=1
succeeded=1
retried=0
manual_review=0
terminal_failed=0
cancelled=0
Post-run DB state:
notification_status=sent
notification_sent_at_present=yes
notification_error_code_present=no
job_status=succeeded
job_attempt_count=1
attempt_rows_succeeded=1
claimable_telegram_jobs_after=0
claimable_generic_notification_jobs_after=0
Process argv secret checks:
0 Telegram token/chat ID matches before and after the worker, excluding the checker process.
Payload/customer data:
No raw payload_redacted, customer data, DB DSN, provider payload, cookie, credential, reset token, Telegram token value, chat ID value, raw Telegram request/response body, tenant UUID, job UUID, notification UUID, or command line was printed or recorded.
Artifact handling:
The dev/test notification and job remain in the selected database as sent/succeeded drill evidence with safe display IDs and template/event labels.
Open exceptions:
This is one selected-host queued success drill, not a failure/retry drill by itself. T281 separately proves controlled failure/retry classification using a local fake API endpoint.
Decision:
PASS for one selected-host queued Telegram delivery through the real worker; not enough by itself for broader production primary-path approval.
```

## T281 Telegram Failure/Retry Drill Evidence

```text
Evidence ID:
T281-telegram-failure-retry-20260520
Date/time UTC:
2026-05-20T01:20Z
Environment:
Selected host with APP_ENV=staging, selected DB, and local fake Telegram API endpoint on 127.0.0.1.
Evidence collector:
Codex
Secret/config path:
/etc/billing/secrets/telegram.env on the selected host; values redacted.
Secret-file metadata:
mode 600, owner root:root.
Config keys verified:
TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID, TELEGRAM_API_BASE_URL.
DB access:
Protected service-file handling from /etc/billing/secrets/billing-api.env; DB_DSN value was not printed or passed on worker argv.
Telegram API boundary:
No real Telegram API call was made. The worker used fake token/chat values and TELEGRAM_API_BASE_URL pointed to a local fake endpoint.
Pre-run queue state:
claimable_telegram_jobs=0
claimable_generic_notification_jobs=0
Retryable notification created:
display_id=10001
event_type=service.lifecycle
template_key=t281.telegram.retryable_drill
channel=telegram
recipient_group=ops
Retryable job created:
display_id=10001
job_type=notification.deliver.telegram
status_before=queued
Retryable fake API response:
HTTP 500
Retryable worker command:
notification-telegram-once with worker-id t281-telegram-retryable, batch-size 1, timeout 60s.
Retryable worker result:
claimed=1
succeeded=0
retried=1
manual_review=0
terminal_failed=0
cancelled=0
Retryable pre-cleanup DB state:
notification_status=failed
notification_error_code=telegram_http_500
job_status=failed_retryable
job_error_code=telegram_http_500
job_attempt_count=1
attempt_rows_failed_retryable=1
Terminal notification created:
display_id=10002
event_type=service.lifecycle
template_key=t281.telegram.terminal_drill
channel=telegram
recipient_group=ops
Terminal job created:
display_id=10002
job_type=notification.deliver.telegram
status_before=queued
Terminal fake API response:
HTTP 400
Terminal worker command:
notification-telegram-once with worker-id t281-telegram-terminal, batch-size 1, timeout 60s.
Terminal worker result:
claimed=1
succeeded=0
retried=0
manual_review=0
terminal_failed=1
cancelled=0
Terminal post-run DB state:
notification_status=failed
notification_error_code=telegram_http_400
job_status=failed_terminal
job_error_code=telegram_http_400
job_attempt_count=1
attempt_rows_failed_terminal=1
Cleanup result:
retryable_cleanup_cancelled_jobs=1
retryable_cleanup_cancelled_notifications=1
retryable_post_cleanup_notification_status=cancelled
retryable_post_cleanup_job_status=cancelled
claimable_telegram_jobs_after_cleanup=0
claimable_generic_notification_jobs_after_cleanup=0
Fake API request counts:
fake_api_status_500_calls=1
fake_api_status_400_calls=1
Process argv secret checks:
0 Telegram token/chat ID/DB_DSN matches before and after the worker, excluding the checker process.
Payload/customer data:
No raw payload_redacted, customer data, DB DSN, provider payload, cookie, credential, reset token, Telegram token value, chat ID value, raw Telegram request/response body, tenant UUID, job UUID, notification UUID, or command line was printed or recorded.
Artifact handling:
The retryable dev/test artifact was cancelled after evidence capture so it cannot be claimed again. The terminal dev/test artifact remains failed_terminal as evidence and is not claimable.
Open exceptions:
This proves worker classification and cleanup hygiene with a controlled local fake API endpoint. It does not prove real Telegram outage behavior, production worker always-on operations, production support monitoring, or owner approval for Telegram as the sole broader production notification path.
Decision:
PASS for selected-host controlled Telegram worker failure classification: HTTP 500 maps to retryable, HTTP 400 maps to terminal, and retryable drill cleanup leaves no claimable Telegram jobs.
```
