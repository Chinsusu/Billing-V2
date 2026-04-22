# 22 - Deployment DevOps And Environment Runbook

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa runbook triển khai và vận hành môi trường cho platform VPS/Proxy.

Đây không phải tài liệu chọn framework/cloud cụ thể. Mục tiêu là khóa:
- Environment cần có.
- Secret quản lý thế nào.
- Database/queue/worker triển khai thế nào.
- Backup/restore cần gì.
- Monitoring/alert cần gì.
- Release/rollback ra sao.
- Checklist trước production.

---

## 2. Environment model

### 2.1 Local development

Dành cho dev build/test chức năng.

Cho phép:
```text
mock provider
manual provider
sandbox email
fake payment proof
seed tenant data
```

Không cho phép:
```text
provider production API key
production database copy chưa mask dữ liệu
real credential plaintext
```

### 2.2 Staging

Dành cho QA/UAT.

Yêu cầu:
```text
gần giống production
có queue/worker/scheduler thật
có test provider/sandbox provider
có test email/telegram sandbox
có backup/restore test
```

Staging phải test được:
- top-up flow.
- checkout.
- provisioning success/fail/timeout.
- tenant isolation.
- credential reveal.
- expiry/suspend cron.

### 2.3 Production

Dành cho khách thật.

Yêu cầu:
```text
secret management chuẩn
backup tự động
monitoring/alert
audit retention
rate limit
2FA admin
access control chặt
rollback plan
```

---

## 3. Logical deployment units

Không bắt buộc microservice phase 1. Khuyến nghị:
```text
modular monolith API
frontend portals
worker process
scheduler/cron process
database
queue/cache
object storage
log/monitoring stack
secret manager
```

### 3.1 Frontend

Có thể deploy một app hoặc tách:
```text
admin portal
reseller portal
client storefront/portal
```

Requirements:
- tenant branding/domain support.
- no secret in frontend.
- environment-specific API endpoint.
- CSP/security headers recommended.

### 3.2 Backend API

Responsibilities:
- auth.
- tenant context.
- RBAC.
- catalog/order/wallet/service APIs.
- audit.
- create queue jobs.

Backend không nên giữ long-running provider call trong request user.

### 3.3 Worker

Responsibilities:
- provisioning.
- provider sync.
- service actions.
- notification send.
- export/report long-running tasks.

Worker phải chạy cùng version với backend API hoặc có compatibility contract.

### 3.4 Scheduler/Cron

Responsibilities:
- reservation expiry.
- service expiry.
- suspension/termination.
- renewal reminder.
- provider health.
- inventory sync.
- retention.

Scheduler cần lock để tránh nhiều instance chạy cùng job một lúc.

---

## 4. Environment variables / config groups

Không ghi secret thật trong tài liệu/repo.

Config groups:
```text
APP_ENV
APP_URL
ADMIN_URL
DATABASE_URL
QUEUE_URL
CACHE_URL
OBJECT_STORAGE_CONFIG
EMAIL_PROVIDER_CONFIG
TELEGRAM_BOT_CONFIG
ENCRYPTION_KEY_REFERENCE
JWT/SESSION_SECRET_REFERENCE
PROVIDER_SECRET_REFERENCE
PAYMENT_METHOD_CONFIG
RATE_LIMIT_CONFIG
LOGGING_CONFIG
```

### 4.1 Secret classification

Critical:
```text
database password
session/JWT secret
encryption master key
provider API key/secret
email provider secret
telegram bot token
object storage secret
```

Sensitive:
```text
payment instructions internal refs
support webhook
backup storage credential
```

Public:
```text
public app URL
brand assets
feature flags non-sensitive
```

Rule:
```text
critical secrets never committed.
critical secrets rotate-able.
critical secrets not visible to normal developers in production.
```

---

## 5. Secret management

### 5.1 Required controls

- Provider credentials encrypted at rest.
- Service credentials encrypted at rest.
- Master encryption key stored outside database.
- Secret rotation plan.
- Access to secret manager audited.
- Staging/prod secrets separated.

### 5.2 Encryption key rotation logic

Minimum design:
```text
secret_version stored with encrypted payload
new credentials encrypted with current key version
old credentials decryptable until rotated
rotation job can re-encrypt if needed
```

### 5.3 Never log

```text
Authorization header
session token
provider API key
root password
proxy password
encryption key
payment proof private URL if sensitive
```

---

## 6. Database operations

### 6.1 Migration rule

Every schema change needs:
```text
migration description
backward compatibility note
rollback note
data migration risk
tested on staging
```

Financial/tenant tables require extra review:
```text
wallets
wallet_ledger_entries
orders
reservations
services
service_credentials
audit_logs
provider_resource_mappings
```

### 6.2 Migration safety

Before production migration:
- backup completed.
- migration tested on staging copy.
- long-running migration estimated.
- rollback or forward-fix plan.
- maintenance window if needed.

### 6.3 Data integrity checks

Periodic checks:
```text
wallet balance cache equals ledger sum
orders paid have ledger entries
services active have provider mapping
reservations allocated have service
provider mappings unique
credential rows encrypted
ledger entries immutable
```

---

## 7. Queue/worker deployment

### 7.1 Worker scaling

Workers can scale horizontally if:
```text
job locking prevents duplicate processing
idempotency enforced
provider rate limits respected
```

### 7.2 Job locking

Each job should be claimed by one worker:
```text
queued -> running with atomic lock
heartbeat/timeout for stuck jobs
stuck running job recovery policy
```

### 7.3 Stuck job policy

If job running too long:
```text
mark as stale
do not blindly rerun provider create
move to manual_review if provider action may have been sent
```

### 7.4 Provider rate limit

Provider-specific worker throttles:
```text
max concurrent jobs per source
min delay/backoff on 429
source maintenance mode
```

---

## 8. Object storage

Used for:
```text
payment proof attachments
abuse evidence files
private export files
possibly redacted provider raw response
brand assets
```

Rules:
- Private files require signed/authorized access.
- Payment proofs and abuse evidence not public.
- Brand assets can be public.
- Retention policy by file type.
- Virus/malware scan if attachments from users.

---

## 9. Logging and monitoring

### 9.1 Application logs

Logs should include:
```text
timestamp
level
request_id
correlation_id
tenant_id when safe
actor_id when safe
module
action
error_code
```

No plaintext secret.

### 9.2 Metrics

P0 metrics:
```text
api error rate
api latency
login failures
checkout success/fail rate
wallet ledger posting failures
provisioning success/fail/manual_review
queue backlog by job type
provider health
provider latency/error
reservation expiry count
service expiry/suspend/terminate count
notification failure rate
credential reveal count
ledger adjustment count
```

### 9.3 Alerts

Critical alerts:
```text
database down
queue down
worker not running
provider down repeated
provisioning manual_review spike
wallet ledger posting error
checkout failure spike
backup failed
credential reveal spike
admin login brute-force spike
```

Warning alerts:
```text
reseller low balance
source stock low
top-up pending too long
provider degraded
email notification failure
```

---

## 10. Backup and restore

### 10.1 Backup scope

Must back up:
```text
database
object storage private files
encryption key references/secret manager backup procedure
configuration snapshots
```

Database backup is useless if encryption keys are lost.

### 10.2 Backup frequency

Suggested:
```text
database: automated daily + point-in-time if possible
object storage: daily/incremental
config/secret references: whenever changed
```

### 10.3 Restore test

At least before production and periodically:
```text
restore staging from backup
verify users/tenants/orders/ledger/services
verify encrypted credential can decrypt with key reference
verify audit logs accessible
verify provider mappings intact
```

### 10.4 Restore acceptance

Restore is valid only if:
- Ledger balances reconcile.
- Services still map to provider resources.
- Credentials decrypt.
- Audit logs preserved.
- Tenant domains/config present.
- App can boot with restored config.

---

## 11. Release process

### 11.1 Release checklist

Before release:
- changelog written.
- migrations reviewed.
- staging tests passed.
- QA P0 passed.
- backup completed.
- rollback plan documented.
- monitoring ready.
- support/admin notified if needed.

### 11.2 Deployment order

Typical order:
```text
1. backup
2. run compatible DB migration
3. deploy backend API
4. deploy worker
5. deploy scheduler
6. deploy frontend
7. run smoke tests
8. monitor
```

If breaking change:
```text
use maintenance window or two-phase migration
```

### 11.3 Rollback

Rollback plan must define:
- app rollback.
- worker rollback.
- migration rollback/forward-fix.
- queue compatibility.
- jobs in progress handling.
- provider calls already sent.

Important:
```text
Không rollback bừa nếu migration đã thay đổi financial ledger semantics.
Trong hệ thống tiền, forward-fix đôi khi an toàn hơn rollback.
```

---

## 12. Incident response

### 12.1 Incident severity

| Severity | Example | Response |
|---|---|---|
| SEV1 | wallet/ledger wrong, credential leak, tenant data leak | immediate freeze affected functions |
| SEV2 | provider outage, provisioning failing widely | disable source, manual review |
| SEV3 | notification delayed, report wrong | fix within normal SLA |
| SEV4 | UI cosmetic | backlog |

### 12.2 Immediate actions

For wallet/ledger issue:
```text
pause checkout if needed
pause top-up approval if needed
snapshot affected ledger/orders
investigate correlation_id
create adjustment only after root cause known
```

For credential leak:
```text
disable reveal if needed
rotate affected credentials if possible
notify impacted users per policy
audit access logs
```

For provider duplicate resource:
```text
stop retry for source
list provider resources by time/correlation metadata
link or terminate duplicates carefully
adjust ledger/refund if needed
```

For tenant data leak:
```text
disable affected endpoint
audit access
notify affected tenant per policy/legal need
fix tenant guard
add regression test
```

---

## 13. Security hardening checklist

Minimum production:
- Admin 2FA required.
- Reseller owner 2FA enabled/default.
- Strong password policy.
- Rate limit login/register/checkout/reveal.
- CSP/security headers.
- HTTPS only.
- Secure cookies/session.
- Provider API IP allowlist if possible.
- Secrets not in repo/logs.
- DB access restricted.
- Audit critical actions.
- Backup encrypted.
- Least privilege staff roles.

---

## 14. Go-live checklist

Do not go production until:
- Tenant isolation QA passed.
- Wallet/ledger QA passed.
- Checkout/reservation QA passed.
- Provisioning success/fail/timeout QA passed.
- Credential security QA passed.
- Backup/restore tested.
- Monitoring/alert configured.
- Admin 2FA active.
- Provider source tested with small live resource.
- Abuse/manual suspend SOP ready.
- Support flow ready.
- Terms/AUP visible on storefront.

---

## 15. Runbook acceptance criteria

Deployment/runbook đạt khi:
- Có environment local/staging/prod rõ.
- Secret không commit, không log.
- Worker/scheduler có deployment riêng.
- Backup/restore test được.
- Monitoring/alert bao phủ money/provisioning/provider/credential.
- Release có checklist và rollback/forward-fix plan.
- Incident response có bước pause/freeze cho chức năng nguy hiểm.
- Production không chạy nếu chưa pass P0 QA.

Câu nền: **production không phải nơi để chứng minh dev chạy được; production là nơi chứng minh hệ thống vẫn đúng khi mọi thứ xung quanh bắt đầu sai.**
