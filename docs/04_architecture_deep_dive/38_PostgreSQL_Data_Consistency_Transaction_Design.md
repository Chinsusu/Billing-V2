# 38 - PostgreSQL Data Consistency And Transaction Design

**Version:** v1.4.1 Architecture Deep Dive Hotfix  
**Status:** Architecture decision draft for backend/dev/DBA review  
**Scope:** PostgreSQL consistency model, transaction boundaries, ledger integrity, inventory locking, provisioning idempotency, outbox/jobs, tenant-safe data access  
**Related docs:** 04, 09, 11, 15, 19, 21, 37, 39, 40, 41, 42, 43, 44, 45

---

## 0. Executive summary

Dự án VPS/Proxy này không phải web CRUD bình thường. Nó là hệ thống điều phối 4 thứ nhạy cảm cùng lúc:

```text
1. Money       - ví, ledger, reseller settlement, refund, adjustment
2. Resource    - VPS/proxy stock, reservation, allocation, provider resource
3. Identity    - tenant, reseller, client, RBAC, credential visibility
4. State       - order, provisioning, lifecycle, billing, abuse, audit
```

Vì vậy PostgreSQL phải được xem là **source of truth**. Cache, queue, provider API, webhook, Redis, worker memory đều chỉ là phụ trợ.

Các quyết định khóa:

```text
PostgreSQL is the source of truth.
Financial mutations must be transaction-first.
Ledger entries are immutable.
Inventory reservation must be atomic.
Provider provisioning must be idempotent and recoverable.
Outbox/job records must be created in the same DB transaction as business state.
Every tenant-owned query must be tenant-scoped.
Every sensitive state transition must leave audit evidence.
```

Câu gốc để team nhớ:

```text
Không có transaction thì không có sự thật.
Không có ledger thì không có tiền.
Không có tenant scope thì không có an toàn.
Không có idempotency thì provisioning chỉ là đánh bạc.
```

---

## 1. Why PostgreSQL is the system of record

PostgreSQL chịu trách nhiệm lưu trạng thái chính của toàn hệ thống:

```text
tenants
users
roles / permissions
tenant_domains
catalog products / plans / sources
tenant catalog clones
wallets
wallet_ledger_entries
topup_requests
orders
order_items
reservations
services
service_credentials metadata
provider_sources
provider_accounts
provisioning_jobs
provider_requests
provider_resource_mappings
audit_logs
risk_flags
abuse_cases
notifications
outbox_events
job_claims
```

Redis có thể dùng cho:

```text
rate limit
short-lived cache
session cache nếu cần
distributed lock phụ trợ
```

Queue ngoài có thể dùng cho:

```text
async fanout
email/telegram notification
worker job dispatch khi scale
```

Provider API là nguồn trạng thái bên ngoài, nhưng không phải nguồn sự thật nội bộ. Nếu provider và PostgreSQL lệch nhau, hệ thống phải chuyển qua `sync_required` hoặc `manual_review`, không tự sửa bằng phỏng đoán.

Rule:

```text
If PostgreSQL and cache disagree, PostgreSQL wins.
If PostgreSQL and provider disagree, trigger reconciliation/manual review.
If ledger and wallet projection disagree, ledger wins but operation must freeze until reviewed.
```

---

## 2. Data consistency classes

Không phải dữ liệu nào cũng cần cùng mức consistency. Phân loại đúng giúp tránh over-engineer nhưng vẫn bảo vệ chỗ sống còn.

### 2.1 Class A - Financial state

Bao gồm:

```text
wallets
wallet_ledger_entries
topup_requests
refunds
adjustments
reseller_settlement_records
platform_revenue_records
```

Yêu cầu:

```text
Strong consistency.
Single transaction per mutation.
Immutable ledger.
Idempotency key bắt buộc với operation có thể retry.
Audit bắt buộc.
Không eventual consistency trong checkout/debit/approval/refund.
```

Không được:

```text
- update wallet balance không có ledger entry.
- sửa ledger entry cũ để fix sai.
- xử lý checkout dựa trên balance cache không transactional.
- approve top-up bằng 2 request đồng thời mà credit 2 lần.
```

### 2.2 Class B - Inventory and reservation state

Bao gồm:

```text
source_capacity
reservations
allocation counters
provider_source availability
```

Yêu cầu:

```text
Atomic reservation.
Row lock hoặc conditional atomic update.
Reservation state machine rõ.
Expired cleanup idempotent.
```

Không được:

```text
- read available stock rồi update sau mà không lock.
- allocate reservation đã expired/released.
- decrement reserved_count nhiều lần khi cleanup job chạy lặp.
```

### 2.3 Class C - Provisioning state

Bao gồm:

```text
provisioning_jobs
provider_requests
provider_response_records
provider_resource_mappings
```

Yêu cầu:

```text
Idempotency.
Explicit retry policy.
Unknown/partial success state.
Recoverable mapping.
No blind retry after dangerous timeout.
```

Provider provisioning là vùng nguy hiểm nhất vì có thể xảy ra:

```text
DB says pending, provider already created VPS.
Provider response timeout, but resource exists.
Worker crashed after provider success before DB commit.
Webhook duplicated.
Retry creates second resource.
```

### 2.4 Class D - Tenant-owned business state

Bao gồm:

```text
orders
services
client records
support tickets
notifications
tenant settings
```

Yêu cầu:

```text
tenant_id bắt buộc.
Index bắt đầu bằng tenant_id ở bảng lớn.
Query theo resource_id + tenant_id.
Composite FK cho bảng quan trọng nếu khả thi.
```

### 2.5 Class E - Secrets and credentials

Bao gồm:

```text
provider API keys
service credentials
proxy credentials
2FA secrets
webhook secrets
```

Yêu cầu:

```text
Encrypt at rest.
Redact in logs/audit/snapshots.
Reveal action must be audited.
No plaintext in normal service response.
```

---

## 3. Core database invariants

Đây là những điều kiện luôn phải đúng. Nếu vi phạm, coi như incident.

### 3.1 Financial invariants

```text
Every wallet balance mutation has at least one ledger entry.
Ledger entries are append-only.
Wallet projection equals SUM(ledger entries), except during transaction before commit.
Top-up approval can credit once only.
Refund never deletes or modifies original debit.
Adjustment must reference reason, actor, and approval trail.
Client order cannot be provisioned if payment/settlement rule is not satisfied.
```

### 3.2 Reseller settlement invariants

```text
Client wallet and reseller wallet are different accounting layers.
Client balance controlled by reseller does not prove platform has received money.
Platform only trusts reseller wallet for infrastructure cost.
If client belongs to reseller tenant, provisioning requires reseller wallet >= reseller_cost unless explicit credit policy exists.
Every reseller-client order stores selling_price_snapshot and reseller_cost_snapshot.
```

### 3.3 Inventory invariants

```text
allocated_count cannot exceed capacity.
reserved_count cannot be negative.
available = capacity - reserved_count - allocated_count.
Only reserved reservation can become allocated.
Expired/released/cancelled reservation cannot become allocated.
```

### 3.4 Provisioning invariants

```text
Each provisioning job has unique idempotency key.
Each provider request has request status and retry classification.
Provider resource mapping must be unique per external_resource_id/provider.
Unknown provider result must not be treated as terminal failure.
Service cannot become active without allocated reservation or explicit manual allocation reason.
```

### 3.5 Tenant invariants

```text
Tenant-owned rows must have tenant_id.
Child rows must not point to parent rows in another tenant.
Client/reseller API never trusts tenant_id from request body.
Platform admin cross-tenant access uses explicit admin route and audit.
```

---

## 4. Recommended PostgreSQL schema patterns

### 4.1 IDs

Use UUID primary keys for internal references, foreign keys, idempotency, and security-sensitive operations.

Rows shown to FE/admin/support must also have a numeric `display_id` generated by a table-specific sequence. This is the human-facing number shown in lists and detail pages for account, service, provider, order, invoice, transaction, audit log, job, and similar records.

```text
UUID id: stable internal identity and FK target.
display_id: short numeric identity for UI, support, search, and screenshots.
```

Recommended:

```text
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
display_id BIGINT NOT NULL UNIQUE DEFAULT nextval('<table>_display_id_seq')
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NULL for mutable records
```

Rules:

- Do not use `display_id` as a foreign key or idempotency key.
- Do not authorize by `display_id` alone. Resolve it under tenant/context scope to the UUID row first.
- `display_id` starts at 10000 per table unless a task explicitly documents another start point.
- Existing rows are backfilled in stable creation order.
- Pure junction/system tables such as role-permission mappings do not need `display_id` unless they become first-class UI records.

### 4.2 Tenant-owned tables

Every tenant-owned table should include:

```sql
tenant_id UUID NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ
```

Query pattern:

```sql
SELECT *
FROM services
WHERE id = $1
  AND tenant_id = $2;
```

Bad pattern:

```sql
SELECT * FROM services WHERE id = $1;
```

### 4.3 Composite integrity for critical tenant tables

For high-risk relationships, use composite uniqueness to prevent cross-tenant references.

Example:

```sql
ALTER TABLE wallets
ADD CONSTRAINT wallets_id_tenant_unique UNIQUE (id, tenant_id);

ALTER TABLE wallet_ledger_entries
ADD CONSTRAINT ledger_wallet_tenant_fk
FOREIGN KEY (wallet_id, tenant_id)
REFERENCES wallets (id, tenant_id);
```

This prevents a ledger entry in tenant A from accidentally referencing wallet of tenant B.

### 4.4 JSONB snapshots

Use JSONB snapshots for historical proof:

```text
product_snapshot
plan_snapshot
price_snapshot
billing_cycle_snapshot
provider_source_snapshot
capability_snapshot
refund_policy_snapshot
terms_snapshot
fx_snapshot nếu multi-currency
```

But keep report-critical fields as normal columns:

```text
selling_price
reseller_cost
currency
billing_cycle_type
term_start
term_end
provider_id
source_id
```

Reason:

```text
JSONB is good for evidence.
Columns are better for reporting, indexing, constraints, and reconciliation.
```

---

## 5. Transaction isolation strategy

Default isolation should be:

```text
READ COMMITTED
```

For most operations, `READ COMMITTED` plus explicit row locks/conditional updates is enough.

Use stronger isolation only when required:

```text
REPEATABLE READ: complex report snapshot or reconciliation batch.
SERIALIZABLE: rare high-risk batch operations where explicit locking is not enough.
```

Do not blindly set all transactions to SERIALIZABLE. It can reduce throughput and cause retry complexity.

Golden rule:

```text
Use explicit locks for specific rows that must not race.
Use unique constraints for idempotency.
Use reconciliation to detect impossible states.
```

---

## 6. Transaction boundaries by business flow

### 6.1 Checkout transaction

Checkout is the most important transaction in the system.

Target result:

```text
Either everything required to provision exists after commit,
or nothing irreversible happened in DB.
```

Recommended transaction:

```text
BEGIN
  1. Load actor + tenant context.
  2. Validate plan is sellable for tenant.
  3. Lock required wallet rows.
  4. Lock or atomically reserve inventory.
  5. Create order and order_items with snapshots.
  6. Create reservation status=reserved.
  7. Create wallet ledger debit entries.
  8. Update wallet balance projections.
  9. Create provisioning_job with idempotency_key.
  10. Insert outbox events.
  11. Insert audit records.
COMMIT
```

For reseller-owned client order:

```text
Client wallet debit: selling_price
Reseller wallet debit: reseller_cost
Platform revenue recognition basis: reseller_cost
Reseller margin snapshot: selling_price - reseller_cost
```

Critical guard:

```text
If client wallet insufficient: reject before reservation finalization.
If reseller wallet insufficient: reject before provisioning_job creation.
If stock insufficient: reject without ledger debit.
```

Implementation note:

```text
Wallet rows should be locked in deterministic order to reduce deadlock.
For example: lock reseller wallet first, then client wallet, sorted by wallet_id if multiple.
```

### 6.2 Manual top-up approval transaction

```text
BEGIN
  SELECT topup_request FOR UPDATE
  verify status = pending
  verify payment_reference not already approved if applicable
  SELECT wallet FOR UPDATE
  insert ledger credit with unique idempotency key
  update wallets.balance = balance + amount
  update topup_request status = approved
  insert audit log
  insert notification outbox
COMMIT
```

Idempotency rules:

```text
Same topup_request_id cannot create two credit entries.
Same payment_reference cannot approve two topups unless explicitly allowed and audited.
```

Suggested unique constraint:

```sql
CREATE UNIQUE INDEX uq_ledger_topup_credit
ON wallet_ledger_entries (wallet_id, reference_type, reference_id, entry_type)
WHERE entry_type = 'topup_credit';
```

### 6.3 Top-up rejection transaction

```text
BEGIN
  SELECT topup_request FOR UPDATE
  verify status = pending
  update status = rejected
  set reject_reason
  insert audit log
  insert notification outbox
COMMIT
```

No ledger entry should be created for rejected top-up, unless a separate fee/adjustment policy exists.

### 6.4 Refund transaction

Refund must never mutate original debit.

```text
BEGIN
  lock refund_request/order/service as needed
  calculate refundable amount from original order snapshots
  verify amount <= refundable_remaining
  SELECT wallet FOR UPDATE
  insert ledger credit/refund entry
  update wallet balance projection
  update refund_request status
  update order/service billing state if needed
  insert audit log
  insert notification outbox
COMMIT
```

Important:

```text
Original purchase debit remains unchanged.
Refund is a new ledger entry pointing to original order/service.
Partial refund must track remaining refundable amount.
```

### 6.5 Renewal transaction

Renewal should extend service term only if billing succeeds.

```text
BEGIN
  lock service row
  verify service renewable
  lock wallet rows
  debit wallet(s)
  insert ledger entries
  calculate new_term_end from old term policy
  update service term_end
  create renewal record
  insert audit/outbox
COMMIT
```

Term rule should be deterministic:

```text
If active or in grace: new_term_end = old_term_end + billing_cycle.
If expired beyond grace and policy allows restore: new_term_end = now + billing_cycle or old_term_end + cycle depending product policy.
```

This must be product policy, not developer guess.

### 6.6 Provisioning success transaction

Worker after provider confirms success:

```text
BEGIN
  SELECT provisioning_job FOR UPDATE
  verify status in processable states
  upsert provider_request/result
  upsert provider_resource_mapping
  lock reservation
  verify reservation status = reserved
  transition reservation reserved -> allocated
  update source_capacity counters
  create service record
  store encrypted credential payload or credential reference
  update order/provisioning status
  insert lifecycle event
  insert audit log
  insert notification outbox
COMMIT
```

If provider success but DB transaction fails, worker must not call provider create again blindly. It should:

```text
record/recover provider external_resource_id if available
retry DB finalization using same idempotency key
or move to manual_review if state uncertain
```

### 6.7 Provisioning terminal failure transaction

Only use terminal failure if adapter is certain provider did not create resource.

```text
BEGIN
  lock provisioning_job
  lock reservation
  mark provisioning_job = failed_terminal
  release reservation
  update capacity reserved_count
  refund wallet if debit was already committed and policy says auto-refund
  update order status
  audit + notification
COMMIT
```

Do not run this flow after dangerous timeout unless provider status is checked.

### 6.8 Provisioning unknown/manual review transaction

Use when provider result is unclear.

```text
BEGIN
  lock provisioning_job
  set status = manual_review
  save last_provider_response / error classification
  keep reservation depending policy
  audit
  alert admin/support
COMMIT
```

Policy choice:

```text
If unknown after provider create request, keep reservation temporarily.
Do not auto-release if resource may exist.
```

---

## 7. Wallet and ledger design

### 7.1 Ledger is append-only

`wallet_ledger_entries` must be treated as append-only.

Allowed:

```text
INSERT new debit/credit/adjustment/refund entries.
```

Not allowed:

```text
UPDATE amount.
UPDATE direction.
DELETE row.
Rewrite historical metadata to hide mistake.
```

If metadata must be corrected, create separate correction/audit note, not silent mutation.

### 7.2 Suggested ledger fields

```text
id UUID PK
tenant_id UUID NOT NULL
wallet_id UUID NOT NULL
actor_id UUID NULL
entry_type TEXT NOT NULL
reference_type TEXT NOT NULL
reference_id UUID NOT NULL
direction TEXT NOT NULL -- debit/credit
amount NUMERIC(18, 6) NOT NULL
currency TEXT NOT NULL
balance_after NUMERIC(18, 6) NULL
idempotency_key TEXT NOT NULL
correlation_id TEXT NOT NULL
metadata_json JSONB NOT NULL DEFAULT '{}'
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Suggested constraints:

```sql
ALTER TABLE wallet_ledger_entries
ADD CONSTRAINT chk_ledger_amount_positive CHECK (amount > 0);

ALTER TABLE wallet_ledger_entries
ADD CONSTRAINT chk_ledger_direction CHECK (direction IN ('debit', 'credit'));

CREATE UNIQUE INDEX uq_ledger_idempotency
ON wallet_ledger_entries (tenant_id, wallet_id, entry_type, idempotency_key);
```

### 7.3 Balance projection

`wallets.balance` can be stored for fast reads, but it is projection only.

Rule:

```text
wallets.balance update must be in same transaction as ledger insert.
No ledger insert, no balance update.
No balance update, ledger mutation transaction fails.
```

Rebuild formula:

```text
wallet balance = SUM(credits) - SUM(debits)
```

### 7.4 Negative balance policy

Default MVP policy:

```text
No negative balance.
```

Exception only if future credit-line feature exists:

```text
wallet.credit_limit
wallet.available_balance = balance + credit_limit
```

But phase 1 should avoid credit line unless business deliberately accepts risk.

---

## 8. Reseller settlement data consistency

The platform must not confuse client wallet with reseller settlement wallet.

### 8.1 Two-layer money model

```text
Client wallet:
  Internal balance between reseller tenant and its client.

Reseller wallet:
  Real platform-facing balance used to pay infrastructure cost.
```

Checkout by reseller client requires:

```text
client_wallet.balance >= selling_price
reseller_wallet.balance >= reseller_cost
```

Unless explicit credit policy exists, which is out of MVP scope.

### 8.2 Transaction example

```text
BEGIN
  lock client_wallet
  lock reseller_wallet
  verify client_wallet >= selling_price
  verify reseller_wallet >= reseller_cost
  debit client_wallet selling_price
  debit reseller_wallet reseller_cost
  create order with selling_price_snapshot + reseller_cost_snapshot
  create platform revenue record or reportable ledger metadata
  create provisioning_job
COMMIT
```

### 8.3 Reporting truth

Platform revenue should be based on reseller cost charged to reseller, not client selling price.

```text
platform_gross_revenue = SUM(reseller_cost debits)
reseller_gross_profit = SUM(client_selling_price - reseller_cost)
```

Do not calculate old reseller profit from current catalog price.

---

## 9. Inventory reservation locking

### 9.1 Problem

If stock has only 1 slot and 10 users checkout concurrently, only one can reserve.

### 9.2 Safe approach A - row lock

```text
BEGIN
  SELECT * FROM source_capacity
  WHERE source_id = $source_id
  FOR UPDATE;

  available = capacity - reserved_count - allocated_count

  if available >= qty:
    update reserved_count += qty
    insert reservation
  else:
    fail OUT_OF_STOCK
COMMIT
```

### 9.3 Safe approach B - conditional update

```sql
UPDATE source_capacity
SET reserved_count = reserved_count + $1
WHERE source_id = $2
  AND capacity - reserved_count - allocated_count >= $1;
```

Then check affected rows:

```text
1 row affected = reserved
0 rows affected = out_of_stock
```

### 9.4 Reservation state machine

Allowed states:

```text
reserved
allocated
released
expired
cancelled
```

Allowed transitions:

```text
reserved -> allocated
reserved -> released
reserved -> expired
reserved -> cancelled
```

Forbidden transitions:

```text
allocated -> reserved
expired -> allocated
released -> allocated
cancelled -> allocated
```

### 9.5 Expiry cleanup idempotency

Cleanup job:

```text
BEGIN
  select reservation FOR UPDATE
  if status = reserved and expires_at < now:
    update status = expired
    decrement reserved_count
    audit/outbox if needed
  else:
    no-op
COMMIT
```

Running the same cleanup twice must not decrement stock twice.

---

## 10. Provisioning job consistency

### 10.1 Job table requirements

Minimum fields:

```text
id
tenant_id
order_id
order_item_id
reservation_id
provider_id
source_id
idempotency_key
status
attempt_count
max_attempts
next_attempt_at
locked_by
locked_until
last_error_code
last_error_message redacted
last_provider_status
correlation_id
created_at
updated_at
```

Status values:

```text
queued
claimed
running
succeeded
failed_retryable
failed_terminal
manual_review
cancelled
```

### 10.2 Worker claim pattern

A Postgres-backed worker can claim jobs with lock-safe update.

Example conceptual pattern:

```sql
WITH picked AS (
  SELECT id
  FROM provisioning_jobs
  WHERE status IN ('queued', 'failed_retryable')
    AND next_attempt_at <= now()
  ORDER BY created_at ASC
  FOR UPDATE SKIP LOCKED
  LIMIT 10
)
UPDATE provisioning_jobs j
SET status = 'claimed',
    locked_by = $worker_id,
    locked_until = now() + interval '2 minutes'
FROM picked
WHERE j.id = picked.id
RETURNING j.*;
```

Important:

```text
Claiming a job is not the same as completing it.
If worker dies, locked_until allows re-claim.
Worker must use idempotency_key before doing dangerous provider action.
```

### 10.3 Provider request records

Each external provider call should be tracked:

```text
provider_requests:
  id
  tenant_id
  provisioning_job_id
  provider_id
  operation
  idempotency_key
  external_request_id nullable
  external_resource_id nullable
  request_payload_hash
  response_payload_redacted
  status
  error_classification
  created_at
  completed_at
```

Do not store plaintext provider secret or service credential in response payload.

### 10.4 Unknown result rule

Dangerous timeout example:

```text
Worker sends create VPS request.
Provider connection times out after request leaves system.
```

Correct handling:

```text
status = manual_review or unknown_pending_sync
run provider lookup if possible
never blindly retry create call
```

Wrong handling:

```text
mark failed_retryable and call create again immediately
```

That can create duplicate VPS/proxy resources.

---

## 11. Outbox pattern

### 11.1 Why outbox exists

Problem A:

```text
DB commit succeeds but queue publish fails.
=> order paid but no worker/event.
```

Problem B:

```text
Queue publish succeeds but DB rollback happens.
=> worker processes non-existent order.
```

Outbox solves this by storing event in DB transaction.

### 11.2 Outbox fields

```text
id
tenant_id nullable for platform events
aggregate_type
aggregate_id
event_type
payload_json redacted
status pending/processing/published/failed
dedupe_key
attempt_count
next_attempt_at
correlation_id
created_at
published_at
```

### 11.3 Example checkout transaction with outbox

```text
BEGIN
  create order
  reserve inventory
  debit ledger
  create provisioning_job
  insert outbox_event('order.paid')
  insert outbox_event('provisioning.job.created')
COMMIT
```

Outbox dispatcher later publishes/executes.

### 11.4 Idempotent consumers

Every consumer of outbox events must be idempotent.

```text
Receiving same event twice must not duplicate ledger, service, notification, or provider action.
```

Use:

```text
dedupe_key
unique constraint
processed_events table if needed
```

---

## 12. Idempotency design

### 12.1 Operations requiring idempotency

```text
checkout
manual top-up approval
payment webhook if future gateway exists
refund
renew
provision
suspend/unsuspend/terminate
reset password
change IP
provider webhook
notification send if user-visible duplicate is harmful
```

### 12.2 Idempotency key scope

Idempotency keys must include scope:

```text
tenant_id
actor_id or system actor
operation_type
business_reference
```

Suggested constraint:

```sql
CREATE UNIQUE INDEX uq_idempotency_operation
ON idempotency_keys (tenant_id, operation_type, idempotency_key);
```

Or per business table:

```sql
CREATE UNIQUE INDEX uq_provisioning_idempotency
ON provisioning_jobs (provider_id, idempotency_key);
```

### 12.3 Response replay

For API idempotency, store minimal response summary:

```text
status_code
response_body redacted
resource_id
created_at
expires_at
```

On duplicate request:

```text
Return same result if operation already completed.
Return processing state if still in progress.
Reject if same idempotency key used with different payload hash.
```

---

## 13. Tenant data integrity

### 13.1 Application-level tenant guard

All repositories must accept tenant context explicitly for tenant-owned data.

Preferred function shape:

```go
GetService(ctx, tenantID, serviceID)
```

Avoid:

```go
GetService(ctx, serviceID)
```

### 13.2 Database-level tenant protection

For tables with high risk, use composite FK and possibly Row Level Security later.

MVP minimum:

```text
tenant_id on all tenant-owned tables.
repository tests for tenant isolation.
code review rule: no tenant-owned query without tenant_id.
```

Advanced defense-in-depth:

```text
PostgreSQL RLS using app.current_tenant_id setting.
Separate DB role for app runtime.
No superuser connection from application.
```

### 13.3 Platform admin access

Platform admin may cross tenant only through explicit routes and policies.

Required audit fields:

```text
actor_id
actor_role
target_tenant_id
resource_type
resource_id
action
reason nullable but required for sensitive operations
request_id
created_at
```

---

## 14. Constraints and indexes

### 14.1 Financial constraints

```sql
ALTER TABLE wallets
ADD CONSTRAINT chk_wallet_balance_not_null CHECK (balance IS NOT NULL);

ALTER TABLE wallet_ledger_entries
ADD CONSTRAINT chk_ledger_amount_positive CHECK (amount > 0);

CREATE INDEX idx_ledger_wallet_created
ON wallet_ledger_entries (wallet_id, created_at DESC);

CREATE INDEX idx_ledger_tenant_reference
ON wallet_ledger_entries (tenant_id, reference_type, reference_id);
```

### 14.2 Order/service indexes

```sql
CREATE INDEX idx_orders_tenant_created
ON orders (tenant_id, created_at DESC);

CREATE INDEX idx_orders_tenant_status
ON orders (tenant_id, order_status, created_at DESC);

CREATE INDEX idx_services_tenant_user_status
ON services (tenant_id, user_id, service_status);

CREATE INDEX idx_services_expiry
ON services (service_status, term_end)
WHERE service_status IN ('active', 'suspended', 'grace');
```

### 14.3 Job indexes

```sql
CREATE INDEX idx_jobs_claimable
ON provisioning_jobs (status, next_attempt_at, created_at)
WHERE status IN ('queued', 'failed_retryable');

CREATE INDEX idx_jobs_locked_until
ON provisioning_jobs (locked_until)
WHERE status IN ('claimed', 'running');
```

### 14.4 Outbox indexes

```sql
CREATE INDEX idx_outbox_pending
ON outbox_events (status, next_attempt_at, created_at)
WHERE status IN ('pending', 'failed');

CREATE UNIQUE INDEX uq_outbox_dedupe
ON outbox_events (dedupe_key)
WHERE dedupe_key IS NOT NULL;
```

---

## 15. Deadlock prevention

Deadlocks happen when two transactions lock resources in different order.

Rule:

```text
Always lock shared resources in deterministic order.
```

Suggested order:

```text
1. tenant/account config if needed
2. wallets sorted by wallet_id
3. inventory/source rows sorted by source_id
4. order/service rows
5. provisioning_job rows
6. audit/outbox append rows
```

Example:

```text
If checkout needs client wallet and reseller wallet:
lock lower wallet_id first, then higher wallet_id.
```

Keep transactions short:

```text
Do not call provider API inside DB transaction.
Do not send email inside DB transaction.
Do not run slow report inside checkout transaction.
```

---

## 16. What must never happen inside a DB transaction

Never do these inside a critical DB transaction:

```text
Call provider API.
Send email/Telegram.
Call payment gateway.
Perform long file export.
Wait on external network.
Run unbounded query.
Prompt human review.
```

Why:

```text
External calls make transactions long.
Long transactions hold locks.
Held locks cause contention/deadlocks/timeouts.
Timeouts create ambiguous state.
```

Correct pattern:

```text
Transaction creates state + outbox/job.
Worker handles external side effect after commit.
Worker writes result back transactionally.
```

---

## 17. Migrations and schema evolution

### 17.1 Migration discipline

```text
All schema changes through migration files.
No manual production DB edits except incident runbook.
Migration tested on staging copy if possible.
Backward-compatible rollout preferred.
```

### 17.2 Safe migration pattern

For adding required field:

```text
1. Add nullable column.
2. Deploy code that writes it.
3. Backfill old rows.
4. Add NOT NULL constraint.
5. Add index/unique constraint concurrently if table is large.
```

For renaming field:

```text
1. Add new column.
2. Dual-write old + new.
3. Backfill.
4. Read from new.
5. Stop writing old.
6. Drop old in later release.
```

### 17.3 Migration risk labels

```text
P0 migration: touches ledger/wallet/order/service/provisioning mapping.
P1 migration: touches tenant/RBAC/security/credential.
P2 migration: normal catalog/UI/report.
```

P0 migrations require:

```text
backup verified
rollback plan
staging test
reconciliation after deploy
```

---

## 18. Reconciliation jobs

Reconciliation is not optional. It is the early warning system.

### 18.1 Daily finance reconciliation

Check:

```text
wallets.balance == SUM(ledger entries)
ledger entries with missing wallet
ledger entries with invalid reference
paid orders without ledger debit
ledger debit without order/refund/topup reference
reseller wallet negative
manual adjustment without approved reason
```

### 18.2 Order/provisioning reconciliation

Check:

```text
paid orders without provisioning_job
provisioning_job succeeded without service
service active without allocated reservation
allocated reservation without service
provider_resource_mapping without service
manual_review jobs older than threshold
```

### 18.3 Inventory reconciliation

Check:

```text
capacity >= reserved_count + allocated_count
reserved_count equals count of active reserved reservations
allocated_count equals count of allocated reservations/services
expired reservations still counted as reserved
```

### 18.4 Tenant isolation audit

Check:

```text
child rows where child.tenant_id != parent.tenant_id
wallet ledger entries pointing across tenants
orders referencing plans/sources not visible to tenant
services referencing users from another tenant
```

### 18.5 Credential hygiene check

Check:

```text
plaintext-looking credential in audit/log payload fields
credential reveal without audit event
encrypted credential without key version
```

---

## 19. Reporting rules

Reports must be based on historical truth, not current config.

Use:

```text
ledger entries
order snapshots
service lifecycle events
provider request records
```

Do not use:

```text
current plan price to calculate past revenue
current provider cost to calculate old margin
current reseller config to judge old client order
```

Core formulas:

```text
Client wallet balance = credits - debits
Reseller wallet balance = topups + refunds + adjustments - reseller_cost_charges
Platform gross revenue = SUM(reseller_cost charged)
Reseller gross profit = SUM(client_selling_price - reseller_cost)
Refund total = SUM(refund ledger credits)
Manual adjustment total = SUM(adjustment entries by direction)
```

---

## 20. Failure mode handling

### 20.1 DB transaction fails before commit

Result:

```text
No business mutation should be visible.
No provider call should have happened.
No success notification should be sent.
```

Action:

```text
Return safe error.
Log internal error with correlation_id.
No financial adjustment required if no commit.
```

### 20.2 DB commit succeeds but worker fails

Result:

```text
Business state exists.
Job/outbox exists.
Worker can retry.
```

Action:

```text
Retry job according to policy.
Alert if job exceeds threshold.
No duplicate ledger/provisioning due to idempotency.
```

### 20.3 Provider success but DB finalization fails

Dangerous case.

Action:

```text
Do not call create again.
Use same idempotency key.
Try to finalize DB using external_resource_id.
If uncertain, manual_review.
```

### 20.4 Ledger mismatch detected

Action:

```text
P0 finance incident.
Freeze affected wallet operations if needed.
Export ledger trail.
Identify mutation source.
Correct only via approved adjustment.
Postmortem required.
```

### 20.5 Tenant mismatch detected

Action:

```text
P0 security incident.
Disable affected route/query if needed.
Audit affected access.
Patch query/constraint/test.
Notify according to policy if data exposure occurred.
```

---

## 21. Testing requirements

### 21.1 Concurrent checkout test

Scenario:

```text
capacity = 1
10 concurrent checkout requests
```

Expected:

```text
1 success
9 out_of_stock or rejected
reserved_count/allocated_count never exceeds capacity
only one provisioning_job created
ledger only for successful checkout
```

### 21.2 Double top-up approval test

Scenario:

```text
same admin double-clicks approve
or two finance admins approve same pending top-up concurrently
```

Expected:

```text
only one ledger credit
wallet balance credited once
topup_request approved once
second request returns idempotent success or conflict
```

### 21.3 Reseller insufficient wallet test

Scenario:

```text
client wallet sufficient
reseller wallet insufficient
```

Expected:

```text
no provisioning_job
no provider call
no reseller wallet debit
client wallet not debited unless policy explicitly creates pending unpaid order
clear error: INSUFFICIENT_RESELLER_BALANCE
```

### 21.4 Provider timeout after create test

Scenario:

```text
provider create request times out after request sent
```

Expected:

```text
job moves to manual_review/unknown
no blind retry create
provider lookup attempted if supported
audit and admin alert emitted
```

### 21.5 Tenant cross-access test

Scenario:

```text
tenant A user requests tenant B service id
```

Expected:

```text
403 or 404
no credential returned
audit/security event recorded if suspicious
```

### 21.6 Ledger reconciliation test

Scenario:

```text
simulate multiple topups, purchases, refunds, adjustments
```

Expected:

```text
wallet balance projection equals SUM ledger
report formulas match snapshots
no old ledger mutation required
```

---

## 22. Database access rules for Go backend

### 22.1 Repository rules

```text
Every repository method for tenant-owned data receives tenant_id.
Every financial method receives transaction handle.
No raw DB access from HTTP handler for money/resource state changes.
```

Preferred flow:

```text
handler -> application service -> domain service -> repository within transaction
```

### 22.2 Transaction helper

Use one transaction wrapper pattern:

```go
err := db.WithTx(ctx, func(tx Tx) error {
    // lock rows
    // insert ledger
    // update wallet projection
    // create outbox
    return nil
})
```

Rules:

```text
No provider call inside WithTx.
No email send inside WithTx.
No goroutine launched inside WithTx.
Context timeout should be reasonable.
```

### 22.3 SQL style

Prefer explicit SQL for critical operations.

Good:

```text
sqlc/pgx with handwritten SQL and generated typed methods.
```

Risky:

```text
ORM query magic that may forget tenant_id or hide locking behavior.
```

---

## 23. Operational dashboards

DB consistency should be visible to operations.

Dashboard widgets:

```text
wallet reconciliation status
orders paid without provisioning job
manual_review provisioning jobs
queue depth
reservation expired cleanup count
provider_resource_mapping conflicts
tenant isolation anomaly count
ledger adjustment count today
negative wallet count
```

P0 alerts:

```text
ledger mismatch
capacity counter mismatch
cross-tenant FK anomaly
duplicate provider resource mapping
credential plaintext detection
queue stuck beyond threshold
```

---

## 24. Acceptance checklist

A build is not acceptable unless:

```text
[ ] Checkout order + reservation + ledger + provisioning_job are committed atomically.
[ ] No provider API call happens inside checkout DB transaction.
[ ] Wallet ledger entries are immutable by application path.
[ ] Wallet balance projection reconciliation job exists.
[ ] Top-up approval is concurrency-safe and idempotent.
[ ] Refund creates new ledger entry, never mutates old debit.
[ ] Reservation cannot oversell under concurrent test.
[ ] Reservation cleanup is idempotent.
[ ] Provisioning job has unique idempotency key.
[ ] Dangerous provider timeout goes to unknown/manual_review, not blind retry.
[ ] Outbox event is inserted in same transaction as business state.
[ ] Tenant-owned tables have tenant_id.
[ ] Critical queries include tenant_id in WHERE clause.
[ ] Composite FK or equivalent test prevents cross-tenant child references.
[ ] Price/cost/policy/capability snapshots are stored on order/service.
[ ] Reconciliation checks run daily at minimum.
[ ] P0 anomaly alert exists for ledger mismatch.
```

---

## 25. Final architecture rule

The PostgreSQL layer is not just storage. It is the contract that keeps the business sane.

```text
Transaction protects cause and effect.
Ledger protects money memory.
Lock protects scarce resources.
Idempotency protects retry chaos.
Snapshot protects historical truth.
Tenant scope protects customer trust.
Outbox protects async reliability.
Reconciliation protects the founder from invisible drift.
```

Nếu team chỉ nhớ một câu:

```text
Mọi thứ có thể retry, trừ tiền sai, tenant leak, và provider cấp trùng tài nguyên.
Thiết kế database phải ngăn ba thứ đó trước khi chúng xảy ra.
```
