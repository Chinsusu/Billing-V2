# 37 - Go Backend Architecture And Module Boundaries

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Backend architecture only. No production code in this document.

---

## 1. Mục tiêu tài liệu

Tài liệu này khóa hướng kiến trúc backend cho nền tảng VPS/Proxy multi-tenant, reseller white-label, wallet-first, provider-adapter-based.

Mục tiêu chính:

- Chọn kiến trúc backend phù hợp với giai đoạn MVP và đường scale sau này.
- Định nghĩa ranh giới module để dev không viết chồng chéo, không phá tiền, không phá tenant.
- Khóa cách tổ chức Go codebase ở mức kiến trúc, chưa phải implementation chi tiết.
- Làm nền cho các tài liệu sau: PostgreSQL transaction, async worker, provider adapter, security, observability, scaling.

Kết luận kiến trúc:

```text
Backend chính: Go
Kiểu kiến trúc: modular monolith trước, microservices sau nếu có tải thật
Database chính: PostgreSQL
Async: outbox + jobs table trước, queue ngoài sau
Deploy process: api / worker / scheduler tách process, cùng codebase
Provider integration: adapter interface + capability mask + normalized error taxonomy
```

Câu chốt:

```text
Không xây một app CRUD bán VPS.
Xây một lõi điều phối tiền, quyền, tài nguyên, trạng thái và provider.
```

---

## 2. Vì sao chọn Go cho backend lõi

Go phù hợp với dự án này vì hệ thống có nhiều I/O network, nhiều worker nền, nhiều timeout/retry, nhiều provider API và cần deploy gọn.

Các phần Go phù hợp:

```text
HTTP API
Provider API client
Queue worker
Scheduler/cron
Webhook receiver
Health check
CLI internal tool
Migration runner wrapper
Admin utility command
```

Go không nên gánh:

```text
Frontend portal phức tạp
Figma/design UI
BI dashboard kéo thả
Data science fraud model giai đoạn sau
```

Frontend vẫn nên là TypeScript/React/Next.js hoặc stack tương đương. Backend Go cung cấp API rõ ràng.

---

## 3. Quyết định kiến trúc: Modular Monolith

### 3.1 Chọn modular monolith trước

Giai đoạn này không nên tách sớm thành microservices. Các domain của dự án liên quan giao dịch chặt với nhau:

```text
Order cần Wallet
Wallet cần Ledger
Order cần Inventory Reservation
Provisioning cần Provider Adapter
Service Lifecycle cần Billing
Audit cần bám mọi action
Tenant context cần đi xuyên toàn bộ hệ thống
```

Nếu tách microservices quá sớm, các transaction quan trọng sẽ biến thành distributed transaction/saga khi chưa có đủ đội ngũ vận hành. Với hệ thống có tiền và provisioning, tách sớm dễ làm tăng lỗi hơn là tăng scale.

Quyết định:

```text
Một codebase Go.
Một PostgreSQL primary database.
Nhiều module nội bộ rõ ràng.
Nhiều process chạy từ cùng codebase: api, worker, scheduler, migrate.
```

### 3.2 Khi nào mới tách microservice

Chỉ cân nhắc tách service khi có ít nhất một điều kiện thật:

```text
Một module có tải vượt xa phần còn lại.
Một team riêng cần deploy độc lập.
Một boundary nghiệp vụ đã rất ổn định.
Một tenant/provider tạo tải cực lớn.
Yêu cầu compliance bắt buộc tách runtime/data.
```

Ứng viên tách sau này:

```text
notification-service
reporting-service
provider-worker-service
public-api-gateway
billing-service, nhưng chỉ sau khi ledger model cực kỳ chín
```

Không tách sớm:

```text
ledger
wallet
order
inventory
provisioning core
```

Các module này phải ở gần nhau trong MVP để giữ transaction rõ.

---

## 4. Runtime process model

Một codebase, nhiều entrypoint:

```text
/cmd/api        HTTP API cho Admin/Reseller/Client portal
/cmd/worker     xử lý jobs/outbox/provider provisioning/notification
/cmd/scheduler  enqueue recurring jobs: expiry, reminder, health check
/cmd/migrate    migration runner hoặc wrapper
/cmd/cli        internal admin tool, optional
```

Triển khai production ban đầu:

```text
api replicas:       2+
worker replicas:    1-3, tăng theo queue depth
scheduler replicas: 1 active leader
postgres:           managed primary + backup
redis:              optional cho cache/rate limit/lock phụ
```

Scheduler cần cơ chế chống chạy trùng:

```text
PostgreSQL advisory lock
hoặc scheduler_locks table
hoặc deployment đảm bảo một replica active
```

Nếu scheduler chạy trùng, các job expiry/suspension/reminder có thể enqueue lặp. Vì vậy mọi scheduled job vẫn phải idempotent.

---

## 5. Cấu trúc thư mục đề xuất

```text
/cmd
  /api
    main.go
  /worker
    main.go
  /scheduler
    main.go
  /migrate
    main.go

/internal
  /app
    api_app.go
    worker_app.go
    scheduler_app.go
    dependencies.go

  /platform
    /config
    /db
    /logger
    /crypto
    /httpserver
    /middleware
    /queue
    /metrics
    /tracing
    /ratelimit
    /email
    /telegram
    /clock

  /modules
    /tenant
    /identity
    /rbac
    /catalog
    /wallet
    /ledger
    /settlement
    /order
    /inventory
    /provisioning
    /provider
    /service
    /audit
    /notification
    /abuse
    /reporting

/migrations
/docs
/scripts
```

Quy tắc:

```text
/cmd chỉ bootstrap process, không chứa business logic.
/internal/platform chứa hạ tầng dùng chung.
/internal/modules chứa nghiệp vụ.
/app wiring dependencies, không chứa logic nghiệp vụ sâu.
```

---

## 6. Module anatomy

Mỗi module nên có cấu trúc gần giống nhau, nhưng không ép máy móc.

Gợi ý:

```text
/internal/modules/order
  handler.go       HTTP handler hoặc transport adapter
  service.go       nghiệp vụ chính
  repository.go    truy vấn DB của module
  model.go         struct/domain model
  policy.go        rule quyền/guard
  events.go        domain events/outbox payload
  errors.go        domain errors
  dto.go           request/response DTO nếu cần
```

Không bắt buộc module nào cũng đủ tất cả file. Module nhỏ có thể gộp. Nhưng các module lớn như wallet, order, provisioning, service, tenant nên giữ rõ.

---

## 7. Ranh giới module bắt buộc

### 7.1 Tenant module

Quản lý:

```text
tenant
reseller tenant
platform tenant metadata
domain mapping
branding config
tenant status
tenant policy
```

Không được quản lý:

```text
wallet balance
order data
service credential
provider secret
```

Tenant module expose các capability:

```text
ResolveTenantByDomain(domain)
GetTenantPolicy(tenant_id)
EnsureTenantActive(tenant_id)
GetTenantBranding(tenant_id)
```

### 7.2 Identity module

Quản lý:

```text
user
password
session
login
2FA
password reset
API key nếu có sau này
```

Không được quyết định business permission phức tạp một mình. Identity chỉ xác định actor là ai. RBAC/policy quyết định actor được làm gì.

### 7.3 RBAC module

Quản lý:

```text
role
permission
role assignment
permission check
emergency access policy
```

Nên expose:

```text
Can(actor, action, resource)
Require(actor, action, resource)
ListPermissions(actor)
```

Không query business data trực tiếp ngoài metadata cần thiết cho policy. Với resource phức tạp, service gọi nên truyền resource owner/tenant/status vào policy.

### 7.4 Catalog module

Quản lý:

```text
product
plan
source
provider_source mapping
tenant catalog clone
tenant pricing override
capability snapshot input
pricing snapshot input
```

Không trực tiếp debit ví hoặc provision.

Catalog module expose:

```text
GetSellablePlan(ctx, tenant_id, plan_id)
BuildOrderSnapshot(plan, tenant_policy)
CheckPlanVisibility(actor, plan)
ValidatePriceOverride(reseller, plan, price)
```

### 7.5 Wallet/Ledger/Settlement modules

Đây là vùng thánh địa.

Quản lý:

```text
wallet
ledger entry
top-up request
manual adjustment
refund/reversal
reseller settlement debit
platform revenue ledger
wallet balance projection
```

Rule cứng:

```text
Không module nào update wallet balance trực tiếp.
Không update/delete ledger entry đã posted.
Không tạo order paid nếu ledger transaction chưa commit.
```

Wallet/Ledger expose method cấp cao:

```text
DebitClientForOrder(...)
DebitResellerCost(...)
ApproveTopup(...)
RefundOrder(...)
CreateAdjustment(...)
GetWalletBalance(...)
```

Không expose:

```text
UpdateBalanceRaw()
DeleteLedgerEntry()
MutatePostedEntry()
```

### 7.6 Order module

Quản lý:

```text
order
order item
checkout state
order status
price snapshot
payment status relation
```

Order không tự gọi provider. Order tạo trạng thái nghiệp vụ và phối hợp với provisioning qua service orchestration/outbox.

### 7.7 Inventory module

Quản lý:

```text
source capacity
reservation
allocation
release
reservation expiry
stock counters
```

Rule:

```text
Reserve phải atomic.
Allocate phải dựa trên reservation hợp lệ.
Release phải idempotent.
Không oversell khi concurrent checkout.
```

### 7.8 Provisioning module

Quản lý:

```text
provisioning_job
provider_request
idempotency_key
attempt
retry/backoff
manual_review state
mapping provider result -> service creation
```

Provisioning module gọi Provider module qua interface, không biết raw implementation cụ thể.

### 7.9 Provider module

Quản lý:

```text
adapter registry
provider account config
provider capability map
provider health
normalized provider error
external resource mapping
```

Provider module không quyết định tiền, không tự tạo service active nếu chưa được provisioning orchestration yêu cầu.

### 7.10 Service module

Quản lý:

```text
service instance
service status
expiry
renewal
suspension
termination
credential reference
lifecycle event
```

Service module không tự tính tiền renew nếu chưa qua Billing/Wallet module.

### 7.11 Audit module

Quản lý:

```text
audit log
security event
financial event reference
action naming
redaction
correlation_id
```

Audit phải được gọi ở mọi action quan trọng. Nhưng không lưu plaintext credential.

---

## 8. Orchestration layer

Một số flow không thuộc riêng module nào. Cần orchestration service để điều phối nhiều module trong cùng transaction.

Ví dụ `CheckoutOrchestrator`:

```text
1. Resolve actor + tenant context
2. Catalog.GetSellablePlan
3. Inventory.Reserve
4. Wallet.DebitClientForOrder
5. Settlement.DebitResellerCost nếu tenant reseller
6. Order.CreatePaidOrder
7. Provisioning.EnqueueJob
8. Audit.Record
9. Commit transaction
```

Ví dụ `RenewalOrchestrator`:

```text
1. Load service with tenant scope
2. Validate renewable
3. Calculate term extension using billing cycle snapshot
4. Debit wallet/settlement
5. Update service term_end
6. Enqueue provider renew if needed
7. Audit.Record
```

Rule:

```text
Orchestrator được gọi nhiều module.
Module service nội bộ không nên gọi vòng tròn lẫn nhau.
Repository không gọi service.
Handler không chứa business transaction dài.
```

---

## 9. Dependency direction

Hướng phụ thuộc mong muốn:

```text
handler -> service/orchestrator -> module services -> repositories/platform
```

Không cho:

```text
repository -> service
platform -> module
module A repository -> module B repository trực tiếp
handler -> repository trực tiếp cho business mutation
```

Một số module nền có thể được module khác gọi:

```text
rbac
audit
notification
```

Nhưng phải tránh việc mọi module đều gọi mọi module. Nếu flow cần nhiều domain, đưa vào orchestrator.

---

## 10. Request lifecycle chuẩn

Flow request portal/API:

```text
1. HTTP request vào API server
2. Request ID middleware
3. Domain mapping middleware resolve tenant context
4. Auth middleware resolve actor/session
5. RBAC middleware hoặc handler-level policy check
6. Rate limit middleware nếu action nhạy cảm
7. Handler parse input
8. Service/orchestrator chạy business logic
9. DB transaction nếu mutation
10. Audit log/outbox event
11. Response chuẩn hóa
```

Context truyền xuyên suốt:

```text
ActorContext:
- request_id
- correlation_id
- actor_id
- actor_type
- tenant_id
- role_ids
- permissions
- domain_context
- is_platform_admin
- emergency_access_id nếu có
```

Không nhận `tenant_id` từ body cho client/reseller. Tenant đến từ domain/session/server context.

---

## 11. Transaction boundary

Transaction không nên mở ở handler. Nên mở ở service/orchestrator.

Các flow bắt buộc transaction:

```text
checkout
wallet debit/credit
reseller settlement
manual top-up approval
refund/reversal
reservation allocate/release
service activation after provider success
renewal payment + term extension
suspension/termination state mutation
```

Không giữ transaction mở trong lúc gọi provider API.

Sai:

```text
Begin DB transaction
Call provider API 30 giây
Update DB
Commit
```

Đúng:

```text
API transaction:
- create order
- debit wallet
- reserve inventory
- enqueue provisioning job
- commit

Worker:
- claim job
- call provider outside critical DB transaction
- short transaction to persist result
```

---

## 12. Error model

Không trả raw provider/database error ra client.

Domain error nên có:

```text
code
message_safe
http_status
retryable
log_level
audit_required
```

Ví dụ code:

```text
TENANT_INACTIVE
PERMISSION_DENIED
PLAN_DISABLED
OUT_OF_STOCK
INSUFFICIENT_CLIENT_BALANCE
INSUFFICIENT_RESELLER_BALANCE
CHECKOUT_IDEMPOTENCY_CONFLICT
PROVISIONING_PENDING_MANUAL_REVIEW
PROVIDER_UNAVAILABLE
CREDENTIAL_REVEAL_DENIED
```

Provider raw error chỉ lưu ở vùng nội bộ đã redact.

---

## 13. Interface placement rule

Trong Go, interface nên đặt gần nơi sử dụng, không đặt tất cả vào một package `interfaces` khổng lồ.

Ví dụ:

```text
provisioning service cần ProviderAdapter interface
=> interface có thể nằm trong provisioning hoặc provider contract package
```

Không tạo package kiểu:

```text
/internal/interfaces
```

vì dễ biến thành bãi rác abstraction.

---

## 14. Repository rule

Repository chịu trách nhiệm SQL/data access của module.

Rule:

```text
Mọi query tenant-owned phải nhận tenant_id hoặc ActorContext.
Không query service/order/wallet bằng id đơn lẻ.
Không trả plaintext secret nếu method không nói rõ.
Không update ledger posted.
```

Ví dụ tên method tốt:

```text
GetServiceForTenant(ctx, tenantID, serviceID)
GetOrderForActor(ctx, actorCtx, orderID)
CreateLedgerEntriesTx(ctx, tx, entries)
ClaimProvisioningJobs(ctx, limit)
```

Tên method nguy hiểm:

```text
GetService(id)
UpdateBalance(walletID, amount)
GetCredential(serviceID)
```

---

## 15. Configuration architecture

Config chia nhóm:

```text
App config:
- environment
- base_url
- log_level
- feature flags

DB config:
- postgres dsn
- migration setting

Security config:
- session secret
- encryption key reference
- password policy
- 2FA required roles

Provider config:
- provider credentials reference
- timeout
- rate limit

Worker config:
- concurrency
- retry limits
- backoff

Notification config:
- SMTP/API
- Telegram bot token
```

Secret không hardcode, không commit. Config object nên validate khi boot.

---

## 16. Testing boundary

Test cần chia:

```text
unit test: policy, price calculation, status transition
repository test: SQL/constraint/tenant scope
integration test: checkout transaction, ledger, reservation
worker test: job claim, retry, manual review
adapter contract test: provider mock
security test: cross-tenant access, credential redaction
```

Mọi bug P0 phải có regression test.

P0 test bắt buộc:

```text
client reseller A không đọc được service reseller B
ledger posted không update/delete được qua service
checkout concurrent không oversell
provider timeout không retry mù
credential reveal luôn audit
```

---

## 17. Anti-patterns cấm trong dự án

```text
Microservices sớm khi chưa có tải thật.
Gọi provider API trực tiếp trong HTTP checkout transaction.
Update wallet balance không qua ledger.
Query dữ liệu tenant bằng id đơn lẻ.
Log plaintext credential.
Retry provisioning timeout như lỗi thường.
Dùng current plan price để xử lý order cũ.
Để frontend quyết định capability action.
Để staff/admin xem credential mà không audit.
Hardcode provider secret trong repo.
Tạo global mutable state cho provider clients mà không timeout/context.
```

---

## 18. Architecture acceptance criteria

Backend architecture được xem là đạt khi:

```text
Có entrypoint api/worker/scheduler riêng.
Tất cả mutation tiền đi qua ledger service.
Tất cả query tenant-owned có tenant scope.
Checkout tạo job/outbox, không gọi provider sync nguy hiểm.
Provider adapter trả normalized result/error.
Credential encrypt + reveal audit.
Request/job log có request_id/correlation_id.
Module boundaries không circular dependency nghiêm trọng.
QA có test cho transaction, tenant, provisioning, credential.
```

---

## 19. Tóm tắt quyết định

```text
Go là backend runtime chính.
Modular monolith là hình thái đúng cho MVP.
PostgreSQL là source of truth.
API, worker, scheduler tách process nhưng chung codebase.
Provider adapter nằm sau provisioning boundary.
Wallet/ledger là vùng bất khả xâm phạm.
Tenant context đi xuyên mọi request/query/job.
Không gọi provider trong transaction dài.
Outbox/jobs là cầu nối an toàn giữa DB và worker.
```

Câu găm sâu:

```text
Code dễ sửa chưa đủ.
Code phải không cho dev vô tình làm sai tiền, sai tenant, sai tài nguyên.
```
