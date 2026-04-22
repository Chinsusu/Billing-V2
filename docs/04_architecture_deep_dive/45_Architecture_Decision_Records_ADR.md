# 45 - Architecture Decision Records ADR

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Architecture decisions, rationale, consequences, future review triggers  
Related docs: 14, 22, 24, 26, 37, 38, 39, 40, 41, 42, 43, 44

---

## 1. Mục tiêu tài liệu

Tài liệu này ghi lại các quyết định kiến trúc quan trọng để team không phải tranh luận lại từ đầu khi triển khai.

ADR không thay thế spec chi tiết. ADR trả lời:

```text
Quyết định là gì?
Vì sao chọn?
Đánh đổi là gì?
Khi nào xem lại?
```

Format mỗi ADR:

```text
Status
Context
Decision
Consequences
Review trigger
```

Status allowed:

```text
Proposed
Accepted
Superseded
Rejected
Deprecated
```

---

## ADR-001 - Use Go for core backend

Status: Proposed

### Context

Hệ thống cần API, worker, scheduler, provider clients, timeout/retry, concurrency, deployment gọn.

### Decision

Backend lõi dùng Go.

Frontend vẫn có thể dùng TypeScript/React/Next.js hoặc stack tương đương.

### Consequences

Ưu điểm:

```text
runtime gọn
concurrency tốt cho provider I/O
dễ deploy nhiều process từ cùng codebase
phù hợp worker/scheduler
```

Đánh đổi:

```text
team cần discipline về module boundary
không có framework full-stack magic
frontend phải tách rõ API contract
```

### Review trigger

Xem lại nếu:

```text
team không có năng lực Go
backend scope đổi sang hệ sinh thái khác đã có sẵn
compliance/integration bắt buộc runtime khác
```

---

## ADR-002 - Modular monolith before microservices

Status: Proposed

### Context

Wallet, ledger, order, inventory, provisioning, service lifecycle liên quan chặt. Tách microservice sớm sẽ tạo distributed transaction/saga trước khi domain ổn định.

### Decision

Một Go codebase modular monolith cho MVP. API/worker/scheduler/migrate là process riêng nhưng cùng codebase và cùng PostgreSQL primary.

### Consequences

Ưu điểm:

```text
transaction rõ
deploy đơn giản
debug nhanh
module boundary vẫn kiểm soát được
```

Đánh đổi:

```text
cần tránh module gọi vòng tròn
scale runtime theo process nhưng codebase chung
team phải giữ discipline dependency direction
```

### Review trigger

Chỉ cân nhắc tách service khi:

```text
module có tải vượt xa phần còn lại
team riêng cần deploy độc lập
boundary nghiệp vụ đã ổn định
compliance bắt buộc tách
```

---

## ADR-003 - PostgreSQL as source of truth

Status: Proposed

### Context

Money, inventory, service lifecycle, tenant data và provider state cần consistency và audit trail.

### Decision

PostgreSQL là source of truth. Redis/cache/queue/provider API là phụ trợ.

### Consequences

Ưu điểm:

```text
transaction mạnh
constraint/index hỗ trợ invariant
reconciliation có nguồn sự thật rõ
giảm distributed consistency trong MVP
```

Đánh đổi:

```text
phải thiết kế query/index tốt
job/outbox trên PostgreSQL cần monitoring
report nặng có thể cần read replica/materialized view sau này
```

### Review trigger

Xem lại khi:

```text
volume vượt khả năng primary
reporting ảnh hưởng transaction path
cần regional data architecture
```

---

## ADR-004 - Immutable ledger for all wallet mutations

Status: Proposed

### Context

Nền tảng wallet-first, reseller settlement, refund/adjustment. Sai ledger là lỗi P0.

### Decision

Mọi thay đổi số dư phải có ledger entry append-only. Không update/delete posted ledger. Wallet balance cache chỉ là projection.

### Consequences

Ưu điểm:

```text
finance trace rõ
reconciliation khả thi
refund/reversal minh bạch
tranh chấp dùng snapshot/lịch sử thật
```

Đánh đổi:

```text
schema và transaction phức tạp hơn CRUD balance
phải có adjustment/reversal thay vì sửa trực tiếp
report cần hiểu ledger type
```

### Review trigger

Không nên xem lại nguyên tắc này. Chỉ xem lại implementation projection khi volume tăng.

---

## ADR-005 - Wallet-first and settlement-first checkout

Status: Proposed

### Context

Client reseller có ví client nhưng platform chỉ tin reseller wallet để trả cost hạ tầng.

### Decision

Checkout chỉ tạo provisioning job khi:

```text
client wallet đủ selling_price
reseller wallet đủ reseller_cost nếu tenant là reseller
ledger debit đã tạo trong transaction
reservation đã tạo
```

### Consequences

Ưu điểm:

```text
không provision khi chưa có tiền thật ở lớp platform
reseller gross profit tính rõ
tránh postpaid/credit rủi ro trong MVP
```

Đánh đổi:

```text
reseller phải nạp trước
client có thể đủ tiền nhưng order bị chặn do reseller wallet thiếu
UI/support phải giải thích rõ
```

### Review trigger

Xem lại nếu có chính sách credit line/postpaid chính thức với risk controls.

---

## ADR-006 - PostgreSQL-backed outbox and jobs first

Status: Proposed

### Context

MVP cần async provisioning/notification/expiry nhưng chưa cần vận hành queue phức tạp.

### Decision

Dùng PostgreSQL-backed outbox + jobs table trước. Queue ngoài thêm sau nếu có tải thật. Outbox vẫn giữ ngay cả khi có queue ngoài.

### Consequences

Ưu điểm:

```text
business state và event/job atomic trong cùng transaction
vận hành đơn giản
debug bằng DB dễ
phù hợp MVP
```

Đánh đổi:

```text
PostgreSQL chịu thêm workload queue
cần index/retention/monitoring queue depth
không phù hợp fanout cực lớn nếu scale mạnh
```

### Review trigger

Thêm queue ngoài khi:

```text
outbox/job depth tăng liên tục
worker throughput bị DB queue bottleneck
cần consumer độc lập nhiều service
notification/report fanout lớn
```

---

## ADR-007 - Provider operations use normalized error taxonomy

Status: Proposed

### Context

Provider lỗi đa dạng và không đáng tin. Raw error không đủ để quyết định retry/manual review.

### Decision

Mọi adapter operation trả normalized result gồm status, error_code, retry_safety, external ids, redacted message.

### Consequences

Ưu điểm:

```text
worker retry an toàn
manual review rõ
UI/API error ổn định
provider mới có contract kiểm thử
```

Đánh đổi:

```text
adapter implementation mất công map lỗi
provider-specific edge case phải cập nhật taxonomy
```

### Review trigger

Cập nhật taxonomy khi onboard provider mới hoặc gặp failure mode chưa map.

---

## ADR-008 - No blind retry after uncertain provider create

Status: Proposed

### Context

Provider có thể timeout sau khi đã tạo resource. Retry create có thể tạo VPS/proxy trùng.

### Decision

Nếu create/provision result unknown hoặc unsafe, worker không tự retry create. Chuyển manual review hoặc status lookup an toàn.

### Consequences

Ưu điểm:

```text
tránh duplicate resource
giữ provider cost an toàn
operator có điểm xử lý rõ
```

Đánh đổi:

```text
một số order chậm hơn vì manual review
cần UI/support giải thích pending/manual review
```

### Review trigger

Chỉ nới nếu provider có idempotency guarantee mạnh đã test thực tế.

---

## ADR-009 - Tenant context is server-side only

Status: Proposed

### Context

Multi-tenant white-label có rủi ro client tự gửi tenant_id/reseller_id để đọc dữ liệu tenant khác.

### Decision

Tenant context resolve từ domain/session/server-side membership. Client/reseller API không tin tenant_id từ body/query.

### Consequences

Ưu điểm:

```text
giảm cross-tenant leak
repository/API rõ boundary
phù hợp white-label domain
```

Đánh đổi:

```text
local/dev phải seed domain mapping
admin cross-tenant flow cần route riêng và audit
```

### Review trigger

Không xem lại nguyên tắc. Có thể cải tiến implementation bằng RLS hoặc identity model global sau.

---

## ADR-010 - Backend RBAC enforcement, frontend hiding is UX only

Status: Proposed

### Context

Frontend có thể bị bypass bằng direct API call.

### Decision

Backend bắt buộc kiểm tra permission, resource ownership, tenant scope, capability, status policy. Frontend chỉ ẩn action để UX tốt hơn.

### Consequences

Ưu điểm:

```text
API an toàn trước bypass
test quyền rõ
permission matrix có giá trị thật
```

Đánh đổi:

```text
handler/service phải khai báo permission
UI và API cần đồng bộ error/capability
```

### Review trigger

Không xem lại nguyên tắc.

---

## ADR-011 - Credential reveal is explicit and audited

Status: Proposed

### Context

Service credentials là secret thật. Hiển thị mặc định trên service detail hoặc log/audit là rủi ro P0.

### Decision

Credential lưu encrypted at rest, service detail chỉ trả masked hint. Plaintext chỉ trả qua reveal endpoint có tenant/permission/rate limit/audit.

### Consequences

Ưu điểm:

```text
giảm lộ credential
truy vết ai đã xem secret
phù hợp support/security operation
```

Đánh đổi:

```text
UX thêm một bước reveal
cần key management và redaction test
```

### Review trigger

Không xem lại nguyên tắc. Có thể thay đổi UX reveal timeout/2FA policy.

---

## ADR-012 - Capability snapshot controls service actions

Status: Proposed

### Context

Provider/source capability có thể đổi sau khi service đã bán.

### Decision

Order/service lưu capability_snapshot tại thời điểm mua. UI/API/worker dựa vào snapshot và current source safety status để cho phép action.

### Consequences

Ưu điểm:

```text
tránh action sai với service cũ
tranh chấp dựa trên snapshot
UI/API nhất quán
```

Đánh đổi:

```text
cần snapshot JSONB đầy đủ
capability migration cần policy riêng
```

### Review trigger

Xem lại khi có feature migration capability hoặc upgrade plan.

---

## ADR-013 - No auto failover in MVP provisioning

Status: Proposed

### Context

Auto failover giữa provider/source đụng đến giá, stock, policy, location, credential, provider terms và refund.

### Decision

MVP không auto failover provider khi provisioning fail. Fail/unknown đi theo retry safety/manual review/refund policy.

### Consequences

Ưu điểm:

```text
giảm lỗi cấp sai sản phẩm
giữ snapshot/order rõ
support dễ giải thích
```

Đánh đổi:

```text
ít tự động hơn
cần operator xử lý manual review
```

### Review trigger

Chỉ thêm failover khi có source equivalence model, pricing policy, inventory policy và QA riêng.

---

## ADR-014 - Observability and audit are launch requirements

Status: Proposed

### Context

Không thể vận hành tiền/tenant/provider/credential nếu không có trace, audit, metrics và alert.

### Decision

MVP phải có structured logs, correlation_id, audit events, worker/provider metrics, P0 alerts, backup/restore monitoring.

### Consequences

Ưu điểm:

```text
incident response khả thi
QA/launch gate đo được
support điều tra order stuck
finance reconcile rõ
```

Đánh đổi:

```text
thêm công build instrumentation
cần redaction discipline
```

### Review trigger

Không bỏ. Chỉ tinh chỉnh tooling/threshold.

---

## ADR-015 - Fail closed for money, tenant, credential, provisioning

Status: Proposed

### Context

Một số lỗi nếu "cho qua" sẽ tạo thiệt hại nghiêm trọng.

### Decision

Các vùng sau phải fail closed:

```text
wallet/ledger mutation
tenant-scoped data read/write
credential reveal
provider create after unknown result
admin critical action without 2FA/reason
```

### Consequences

Ưu điểm:

```text
an toàn hơn khi hệ thống không chắc
giảm rủi ro tiền/credential/tenant
```

Đánh đổi:

```text
có thể tăng pending/manual review
support cần SOP tốt
```

### Review trigger

Không xem lại nguyên tắc. Có thể tối ưu tooling để giảm false positive.

---

## 2. ADR maintenance rule

Khi có thay đổi kiến trúc lớn:

```text
thêm ADR mới
không sửa lịch sử để che quyết định cũ
đánh dấu Superseded nếu thay thế
link tài liệu/spec liên quan
ghi rõ migration/rollback nếu có
```

Examples requiring ADR:

```text
switch database
split microservice
add external queue
introduce postpaid/credit line
enable auto failover
change identity model to global account
change encryption/key management model
```

---

## 3. Open decisions

Các quyết định cần chốt khi bắt đầu implementation:

```text
Frontend stack cụ thể: Next.js/React hay framework khác.
Session strategy: server-side session hay signed token.
Secret manager cụ thể theo hạ tầng deploy.
Migration tool cụ thể cho Go/PostgreSQL.
Provider đầu tiên cho VPS.
Provider/manual source đầu tiên cho proxy.
Exact retention policy theo pháp lý/vận hành.
Whether to use PostgreSQL RLS in addition to app-level tenant guard.
```

---

## 4. Acceptance criteria

ADR package đạt khi:

```text
Architecture decisions are documented with rationale and consequences.
Docs 37-44 align with ADR decisions.
Open decisions are tracked before implementation.
Any future major architecture change adds or supersedes an ADR.
Team can explain why MVP is modular monolith, PostgreSQL-first, outbox/job-based, wallet-first, tenant-first.
```

---

## 5. Tóm tắt

```text
ADR giữ ký ức kiến trúc cho team.
Quyết định quan trọng phải có lý do, đánh đổi, và điều kiện xem lại.
MVP ưu tiên correctness và operability hơn cảm giác scale sớm.
```
