# 23 - Notification Email Telegram Template Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa notification events, channels và template logic cho nền tảng VPS/Proxy.

Notification không chỉ để “gửi email cho đẹp”. Nó giúp:
- giảm support.
- giảm tranh chấp tiền.
- cảnh báo trước khi mất dịch vụ.
- báo lỗi provisioning/provider sớm.
- nhắc reseller nạp ví settlement.
- lưu dấu vận hành.

---

## 2. Notification channels

### 2.1 Dashboard notification

Dùng cho:
```text
mọi user đã đăng nhập
client/reseller/admin inbox
```

Ưu điểm:
- không phụ thuộc email.
- giữ lịch sử trong hệ thống.
- có thể gắn link đến order/service/top-up.

### 2.2 Email

Dùng cho:
```text
registration/verification
top-up status
order/service status
expiry reminder
abuse warning
support update
```

Không nên gửi plaintext credential qua email mặc định.

### 2.3 Telegram/Admin alert

Dùng cho:
```text
provider down
provisioning failed/manual review
reseller low balance
abuse critical
queue backlog
backup failed
```

Telegram cho admin/reseller team nên chứa correlation_id và link admin/reseller portal.

### 2.4 Webhook optional

Phase sau hoặc reseller advanced:
```text
service activated
service expiring
top-up approved
order failed
```

---

## 3. Notification event naming

Format:
```text
module.event
```

Examples:
```text
auth.email_verification
wallet.topup.submitted
wallet.topup.approved
wallet.topup.rejected
wallet.reseller_low_balance
order.created
order.failed
order.manual_review
service.activated
service.expiring
service.expired
service.suspended
service.terminated
credential.revealed_alert_optional
provisioning.failed
provisioning.manual_review
provider.down
provider.recovered
abuse.warning
abuse.suspended
support.ticket_created
support.ticket_replied
```

---

## 4. Notification payload rules

Payload must include:
```text
tenant_id
recipient_user_id or recipient_group
template_key
channel
reference_type
reference_id
correlation_id
dedupe_key optional
priority
```

Payload must not include:
```text
plaintext password
provider API key
private token
full payment proof URL if public
sensitive abuse evidence
```

---

## 5. Priority levels

| Priority | Meaning | Examples |
|---|---|---|
| Low | Informational | catalog update, general support note |
| Normal | User flow status | top-up submitted, order created |
| High | Action needed | service expiring, top-up rejected |
| Critical | Operational risk | provider down, abuse critical, backup failed |

Critical notifications should go to dashboard + admin Telegram/email depending config.

---

## 6. Dedupe policy

Some events must not spam.

Dedupe examples:
```text
service_expiring:{service_id}:{window}
reseller_low_balance:{tenant_id}:{date}
provider_down:{source_id}:{hour}
queue_backlog:{job_type}:{hour}
topup_pending:{topup_request_id}:{day}
```

Critical provider alert can repeat after threshold:
```text
first alert immediately
repeat every N minutes/hours while unresolved
send recovered notification
```

---

## 7. Template variables

Common variables:
```text
tenant_name
brand_name
user_name
order_number
service_name
service_id_short
plan_name
amount
currency
wallet_balance
term_end_at
days_remaining
support_link
dashboard_link
correlation_id
```

Admin variables:
```text
source_name
provider_type
job_id
error_code
manual_review_reason
queue_backlog_count
tenant_name
reseller_name
```

---

## 8. Client templates

### 8.1 Email verification

Event:
```text
auth.email_verification
```

Subject:
```text
Verify your account for {brand_name}
```

Body:
```text
Hi {user_name},

Please verify your email address to activate your account on {brand_name}.

Verification link:
{verification_link}

If you did not create this account, you can ignore this message.
```

Channels:
```text
email
```

### 8.2 Top-up submitted

Event:
```text
wallet.topup.submitted
```

Subject:
```text
Top-up request received: {amount} {currency}
```

Body:
```text
Hi {user_name},

We received your wallet top-up request for {amount} {currency}.
Your request is now waiting for review.

Reference:
{topup_reference}

You will be notified once it is approved or rejected.
```

Channels:
```text
email optional
dashboard
```

### 8.3 Top-up approved

Event:
```text
wallet.topup.approved
```

Subject:
```text
Wallet top-up approved
```

Body:
```text
Hi {user_name},

Your top-up of {amount} {currency} has been approved.
Your updated wallet balance is {wallet_balance} {currency}.

You can now purchase or renew services from your dashboard.
```

Channels:
```text
email
dashboard
```

### 8.4 Top-up rejected

Event:
```text
wallet.topup.rejected
```

Subject:
```text
Wallet top-up could not be approved
```

Body:
```text
Hi {user_name},

Your top-up request for {amount} {currency} could not be approved.

Reason:
{review_note}

Please check your payment details or contact support if you believe this is a mistake.
```

Channels:
```text
email
dashboard
```

### 8.5 Order created / provisioning queued

Event:
```text
order.created
```

Subject:
```text
Order {order_number} has been created
```

Body:
```text
Hi {user_name},

Your order {order_number} has been created and is being processed.

Plan:
{plan_name}

Status:
{order_status}

You can track the order from your dashboard.
```

Channels:
```text
dashboard
email optional
```

### 8.6 Service activated

Event:
```text
service.activated
```

Subject:
```text
Your service is active: {service_name}
```

Body:
```text
Hi {user_name},

Your service is now active.

Service:
{service_name}

Plan:
{plan_name}

Expiry:
{term_end_at}

For security, login to your dashboard to reveal credentials.
```

Important:
```text
Do not include plaintext password by default.
```

Channels:
```text
email
dashboard
```

### 8.7 Provisioning failed

Event:
```text
order.failed
```

Subject:
```text
Order {order_number} could not be provisioned
```

Body:
```text
Hi {user_name},

We could not provision your order {order_number}.

Status:
{failure_status}

If payment was already captured from your wallet, a refund/reversal will be handled according to the platform policy.

Reference:
{correlation_id}
```

Channels:
```text
email
dashboard
```

### 8.8 Service expiring reminder

Event:
```text
service.expiring
```

Subject:
```text
Your service expires in {days_remaining} day(s)
```

Body:
```text
Hi {user_name},

Your service {service_name} will expire on {term_end_at}.

Please renew before expiry to avoid suspension or termination according to the service policy.

Renew here:
{service_link}
```

Channels:
```text
email
dashboard
```

Dedupe:
```text
service_expiring:{service_id}:{days_remaining_window}
```

### 8.9 Service expired

Event:
```text
service.expired
```

Subject:
```text
Your service has expired: {service_name}
```

Body:
```text
Hi {user_name},

Your service {service_name} expired on {term_end_at}.

Depending on the product policy, the service may enter a grace period before suspension or termination.
Please renew as soon as possible if you want to keep using it.
```

### 8.10 Service suspended

Event:
```text
service.suspended
```

Subject:
```text
Your service has been suspended
```

Body:
```text
Hi {user_name},

Your service {service_name} has been suspended.

Reason:
{suspension_reason}

Please contact support or resolve the related issue to request reactivation.
```

### 8.11 Service terminated

Event:
```text
service.terminated
```

Subject:
```text
Your service has been terminated
```

Body:
```text
Hi {user_name},

Your service {service_name} has been terminated.

Reason:
{termination_reason}

Terminated services may not be recoverable. Please contact support if you need clarification.
```

### 8.12 Abuse warning

Event:
```text
abuse.warning
```

Subject:
```text
Important notice about your service
```

Body:
```text
Hi {user_name},

We received an abuse or policy notice related to your service {service_name}.

Issue:
{abuse_case_type}

Please review and resolve this immediately. Continued violations may lead to suspension or termination.
```

---

## 9. Reseller templates

### 9.1 Reseller wallet top-up submitted

Event:
```text
wallet.reseller_topup.submitted
```

Recipient:
```text
reseller owner/staff
platform finance/admin alert optional
```

Body:
```text
Your reseller wallet top-up request for {amount} {currency} has been submitted and is waiting for platform review.
```

### 9.2 Reseller top-up approved

Body:
```text
Your reseller settlement wallet has been credited with {amount} {currency}.
New client orders can continue to provision as long as your balance covers reseller cost.
```

### 9.3 Reseller low balance

Event:
```text
wallet.reseller_low_balance
```

Subject:
```text
Your reseller wallet balance is low
```

Body:
```text
Hi {reseller_name},

Your reseller settlement wallet balance is currently {wallet_balance} {currency}.

New client orders may not be provisioned if your balance is lower than the platform reseller cost.

Please top up your reseller wallet to avoid checkout/provisioning interruption.
```

Channels:
```text
email
dashboard
telegram optional
```

Dedupe:
```text
reseller_low_balance:{tenant_id}:{date}
```

### 9.4 Client top-up pending review

Event:
```text
wallet.client_topup.pending_review
```

Body:
```text
A client top-up request is waiting for review.

Client:
{client_name}

Amount:
{amount} {currency}

Reference:
{topup_reference}
```

### 9.5 Reseller plan margin risk

Event:
```text
catalog.margin_risk
```

Body:
```text
One or more plans in your catalog may have low or negative margin.

Plan:
{plan_name}

Selling price:
{selling_price} {currency}

Reseller cost:
{reseller_cost} {currency}

Please update your selling price or disable the plan.
```

---

## 10. Admin templates

### 10.1 Provider down

Event:
```text
provider.down
```

Subject:
```text
Provider/source down: {source_name}
```

Telegram/dashboard body:
```text
Provider/source health check failed.

Source:
{source_name}

Provider:
{provider_type}

Error:
{error_code}

Last check:
{last_health_check_at}

Impact:
New provisioning through this source may fail or be paused.

Correlation:
{correlation_id}
```

Channels:
```text
telegram
dashboard
email optional
```

### 10.2 Provider recovered

Event:
```text
provider.recovered
```

Body:
```text
Provider/source has recovered.

Source:
{source_name}

Previous status:
{previous_status}

Current status:
healthy
```

### 10.3 Provisioning manual review

Event:
```text
provisioning.manual_review
```

Body:
```text
A provisioning job requires manual review.

Job:
{job_id}

Tenant:
{tenant_name}

Order:
{order_number}

Source:
{source_name}

Reason:
{manual_review_reason}

Retry safety:
{retry_safety}

Correlation:
{correlation_id}
```

### 10.4 Provisioning failed

Event:
```text
provisioning.failed
```

Body:
```text
Provisioning failed.

Job:
{job_id}

Order:
{order_number}

Source:
{source_name}

Error:
{error_code}

Attempts:
{attempt_count}/{max_attempts}

Next action:
{recommended_action}
```

### 10.5 Queue backlog

Event:
```text
system.queue_backlog
```

Body:
```text
Queue backlog threshold exceeded.

Job type:
{job_type}

Backlog:
{queue_backlog_count}

Oldest job age:
{oldest_job_age}

Please check worker health and provider status.
```

### 10.6 Backup failed

Event:
```text
system.backup_failed
```

Body:
```text
Production backup failed.

Backup type:
{backup_type}

Environment:
{environment}

Error:
{error_code}

Immediate action required.
```

### 10.7 Ledger adjustment created

Event:
```text
wallet.adjustment.created
```

Body:
```text
A wallet ledger adjustment was created.

Actor:
{actor_name}

Wallet:
{wallet_reference}

Amount:
{amount} {currency}

Direction:
{direction}

Reason:
{reason}

Correlation:
{correlation_id}
```

Critical audit/notification for finance/admin.

### 10.8 Credential reveal spike

Event:
```text
security.credential_reveal_spike
```

Body:
```text
Credential reveal activity is higher than normal.

Tenant:
{tenant_name}

Actor:
{actor_name}

Count:
{count}

Window:
{time_window}

Please review audit logs.
```

---

## 11. Notification timing matrix

| Event | Client | Reseller | Admin |
|---|---|---|---|
| top-up submitted | immediate dashboard/email | if client top-up: reseller alert | optional |
| top-up approved | immediate | optional | no |
| order created | immediate | optional tenant order alert | no |
| service activated | immediate | optional | no |
| provisioning failed | client if final failed | reseller if tenant client affected | immediate |
| manual review | maybe pending status only | optional | immediate |
| service expiring | 7/3/1 days | reseller summary optional | no |
| service expired | immediate | optional | no |
| service suspended | immediate | optional | if abuse/admin action |
| provider down | no | if affects reseller maybe optional | immediate |
| reseller low balance | no | immediate/daily dedupe | optional |
| abuse warning | immediate | reseller owner | admin abuse queue |

---

## 12. Template localization

Phase 1 can support one default language. Recommended structure:
```text
template_key
language
subject
body
channel
variables_schema
```

If platform targets international reseller/client:
```text
en
vi
```

Tenant can override safe text:
```text
brand greeting
footer
support link
```

Tenant should not override security/legal core wording unless approved.

---

## 13. Notification audit

Important notification events should create audit or notification record:
```text
notification.queued
notification.sent
notification.failed
```

For financial/security:
```text
top-up approved
ledger adjustment
credential reveal alert
provider down
abuse warning
```

Store:
```text
recipient
channel
template_key
reference_id
status
sent_at
error redacted
correlation_id
```

---

## 14. Security rules

Do not send:
```text
root password
proxy password
provider secret
private abuse evidence
payment proof file without authorized link
session/auth token
```

Use links:
```text
login to dashboard to reveal credentials
view top-up request
view service
view admin job
```

Links should expire or require auth.

---

## 15. Notification acceptance criteria

Notification spec đạt khi:
- Mỗi lifecycle tiền/order/service có notification hợp lý.
- Credential không gửi plaintext mặc định.
- Admin nhận provider/provisioning/queue/backup critical alerts.
- Reseller nhận low balance và pending client top-up.
- Expiry reminders có dedupe.
- Payload redacted.
- Notification records trace được bằng correlation_id.
- Templates có biến rõ, không phụ thuộc text hardcode rải rác.
- Critical failure có “next action” trong message.

Câu nền: **notification tốt là hệ thần kinh của platform: nó không làm thay việc vận hành, nhưng nó báo đau đúng chỗ trước khi cơ thể bị thương nặng.**
