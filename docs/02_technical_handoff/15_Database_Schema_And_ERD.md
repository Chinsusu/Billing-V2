# 15 - Database Schema And ERD

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa schema logic cho nền tảng VPS/Proxy. Đây là **data contract**, không phải migration code.

Mục tiêu:
- Khóa các bảng chính.
- Khóa quan hệ tenant/user/order/wallet/service/provider.
- Khóa các field bắt buộc, constraint, index, encryption, audit.
- Giảm rủi ro sai tiền, sai tenant, sai stock, cấp trùng tài nguyên.
- Làm nền cho API, worker, QA và reporting.

---

## 2. Nguyên tắc database bắt buộc

### 2.1 Tenant scope

Mọi bảng chứa dữ liệu thuộc khách/tenant phải có:

```text
tenant_id
created_at
updated_at
```

Các query của reseller/client/staff phải scope theo `tenant_id`.

Không dùng `tenant_id` từ request body để ghi DB. Tenant phải đến từ request context server-side.

### 2.2 Immutable financial ledger

Bảng ledger không được update/delete sau khi entry ở trạng thái `posted`.

Nếu sai, tạo entry điều chỉnh:
```text
adjustment_credit
adjustment_debit
refund
reversal
```

### 2.3 Snapshot cho giao dịch

Order/service phải lưu snapshot để không bị ảnh hưởng bởi thay đổi giá/policy tương lai.

Snapshot tối thiểu:
```text
product_snapshot
plan_snapshot
price_snapshot
billing_cycle_snapshot
capability_snapshot
provider_source_snapshot
reseller_cost_snapshot
fx_snapshot
policy_snapshot
```

### 2.4 Credential encryption

Mọi credential phải encrypt at rest.

Không lưu plaintext trong:
```text
audit_logs
worker_logs
provider_requests raw response
notifications
support tickets public notes
```

### 2.5 Idempotency uniqueness

Các thao tác có khả năng retry phải có unique key:
```text
checkout_idempotency_key
provisioning_idempotency_key
provider_request_idempotency_key
payment_webhook_idempotency_key
```

---

## 3. ERD logic tổng quan

```text
tenants
  ├── tenant_domains
  ├── users
  │     ├── user_roles
  │     ├── wallets
  │     │     └── wallet_ledger_entries
  │     ├── topup_requests
  │     └── orders
  │           ├── order_items
  │           ├── reservations
  │           └── services
  │                 ├── service_credentials
  │                 └── service_lifecycle_events
  ├── tenant_products
  │     └── tenant_plans
  ├── notifications
  ├── audit_logs
  ├── risk_flags
  └── abuse_cases

master_products
  └── master_plans
        └── plan_sources
              └── provider_sources
                    ├── provider_accounts
                    ├── provider_inventory
                    └── provider_resource_mappings

provisioning_jobs
  └── provider_requests
```

---

## 4. Core tenancy tables

### 4.1 `tenants`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| tenant_id | id | yes | Primary key |
| parent_tenant_id | id/null | no | Null với platform/root, set với reseller hierarchy nếu dùng |
| tenant_type | enum | yes | `platform`, `reseller`, `direct_store` |
| name | string | yes | Tên hiển thị |
| slug | string | yes | Unique |
| status | enum | yes | `active`, `suspended`, `disabled`, `pending_setup` |
| default_currency | string | yes | Ví dụ USD, VND |
| timezone | string | yes | Mặc định Asia/Ho_Chi_Minh nếu không có |
| owner_user_id | id/null | no | Owner chính |
| branding_settings | json | no | Logo, màu, footer, support links |
| billing_settings | json | no | Minimum balance, reseller threshold |
| risk_settings | json | no | Manual review thresholds |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(slug)
tenant_type in allowed enum
status in allowed enum
```

Index:
```text
tenant_type
status
owner_user_id
```

### 4.2 `tenant_domains`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| domain_id | id | yes | Primary key |
| tenant_id | id | yes | FK tenants |
| domain | string | yes | Lowercase normalized |
| domain_type | enum | yes | `system_subdomain`, `custom_domain` |
| verification_status | enum | yes | `pending`, `verified`, `failed`, `disabled` |
| verification_token_hash | string | no | Không lưu token plaintext nếu không cần |
| tls_status | enum | yes | `pending`, `active`, `failed`, `expired` |
| is_primary | bool | yes |  |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(domain)
only one primary domain per tenant
```

Security:
```text
domain -> tenant_id mapping là nguồn xác định tenant context cho storefront.
```

---

## 5. User, role, permission tables

### 5.1 `users`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| user_id | id | yes | Primary key |
| tenant_id | id | yes | FK tenants |
| email | string | yes | Normalize lowercase |
| email_verified_at | datetime/null | no |  |
| password_hash | string | yes | Không bao giờ log |
| full_name | string | no |  |
| user_type | enum | yes | `platform_staff`, `reseller_staff`, `client` |
| status | enum | yes | `active`, `suspended`, `disabled`, `pending_verification` |
| two_factor_status | enum | yes | `required`, `enabled`, `disabled` |
| last_login_at | datetime/null | no |  |
| failed_login_count | number | yes | rate limit |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(tenant_id, email)
```

Important:
```text
Cùng một email có thể tồn tại ở nhiều tenant nếu business cho phép.
Nếu muốn global identity, cần bảng identities riêng. Phase 1 có thể giữ email unique theo tenant để đơn giản.
```

### 5.2 `roles`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| role_id | id | yes | Primary key |
| tenant_id | id/null | no | Null = system role, set = custom tenant role |
| role_key | string | yes | `platform_super_admin`, `reseller_owner`, `client` |
| name | string | yes |  |
| is_system | bool | yes |  |
| created_at / updated_at | datetime | yes |  |

### 5.3 `permissions`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| permission_id | id | yes |  |
| permission_key | string | yes | Ví dụ `wallet.topup.approve` |
| module | string | yes | `wallet`, `catalog`, `service` |
| risk_level | enum | yes | `low`, `medium`, `high`, `critical` |

### 5.4 `role_permissions`

```text
role_id
permission_id
created_at
```

Constraint:
```text
unique(role_id, permission_id)
```

### 5.5 `user_roles`

```text
user_id
tenant_id
role_id
created_at
```

Constraint:
```text
unique(user_id, tenant_id, role_id)
```

---

## 6. Catalog tables

### 6.1 `master_products`

| Field | Required | Note |
|---|---:|---|
| product_id | yes | Primary key |
| product_type | yes | `vps`, `proxy`, `service_addon` |
| name | yes | Master name |
| description | no |  |
| status | yes | `draft`, `active`, `disabled`, `archived` |
| display_order | no |  |
| created_by | yes | Admin user |
| created_at / updated_at | yes |  |

### 6.2 `master_plans`

| Field | Required | Note |
|---|---:|---|
| plan_id | yes | Primary key |
| product_id | yes | FK master_products |
| plan_code | yes | Internal SKU |
| name | yes |  |
| specs | yes | json: CPU/RAM/SSD/location/bandwidth/IP |
| billing_cycle_type | yes | `day`, `month_30d`, `calendar_month`, `custom` |
| billing_cycle_value | yes | number |
| base_cost | yes | platform internal cost |
| suggested_price | yes |  |
| reseller_min_price | no | margin guard |
| status | yes | `active`, `disabled`, `archived` |
| version | yes | increment khi đổi giá/spec/policy |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(plan_code, version)
```

### 6.3 `provider_sources`

| Field | Required | Note |
|---|---:|---|
| source_id | yes | Primary key |
| source_type | yes | `proxmox`, `ovh`, `hetzner`, `manual`, `proxy_upstream`, `custom_api` |
| name | yes |  |
| provider_account_id | no | FK nếu có |
| location | no | country/region/datacenter |
| status | yes | `active`, `disabled`, `maintenance`, `out_of_stock` |
| capability_profile | yes | json |
| inventory_mode | yes | `finite_stock`, `provider_live`, `manual_unlimited`, `preloaded_list` |
| risk_level | yes | `low`, `medium`, `high` |
| created_at / updated_at | yes |  |

### 6.4 `plan_sources`

| Field | Required | Note |
|---|---:|---|
| plan_source_id | yes |  |
| plan_id | yes | FK master_plans |
| source_id | yes | FK provider_sources |
| priority | yes | Provider priority |
| cost_override | no | Nếu source có cost riêng |
| capacity_policy | yes | json |
| status | yes | `active`, `disabled` |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(plan_id, source_id)
```

### 6.5 `tenant_products`

| Field | Required | Note |
|---|---:|---|
| tenant_product_id | yes | Primary key |
| tenant_id | yes | FK tenants |
| master_product_id | yes | FK master_products |
| name_override | no |  |
| description_override | no |  |
| status | yes | `active`, `hidden`, `disabled` |
| clone_version | yes | version của master lúc clone/sync |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(tenant_id, master_product_id)
```

### 6.6 `tenant_plans`

| Field | Required | Note |
|---|---:|---|
| tenant_plan_id | yes |  |
| tenant_id | yes |  |
| tenant_product_id | yes |  |
| master_plan_id | yes |  |
| selling_price | yes | Giá bán cho client tenant |
| reseller_cost | yes | Giá platform thu reseller |
| currency | yes |  |
| margin_policy | yes | json |
| visibility | yes | `public`, `hidden`, `private` |
| status | yes | `active`, `disabled`, `margin_risk`, `archived` |
| source_policy_snapshot | yes | json |
| plan_snapshot | yes | json |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(tenant_id, master_plan_id)
selling_price >= 0
reseller_cost >= 0
```

Index:
```text
tenant_id, status, visibility
```

---

## 7. Wallet and billing tables

### 7.1 `wallets`

| Field | Required | Note |
|---|---:|---|
| wallet_id | yes |  |
| tenant_id | yes | Owner tenant |
| owner_type | yes | `tenant`, `user`, `reseller_settlement`, `platform` |
| owner_id | yes | user_id hoặc tenant_id |
| currency | yes |  |
| status | yes | `active`, `frozen`, `closed` |
| available_balance_cache | yes | Cache, không là source of truth |
| locked_balance_cache | yes | Nếu dùng lock/pending |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(owner_type, owner_id, currency)
```

Important:
```text
available_balance_cache chỉ để đọc nhanh.
Source of truth là wallet_ledger_entries.
```

### 7.2 `wallet_ledger_entries`

| Field | Required | Note |
|---|---:|---|
| ledger_id | yes | Primary key |
| wallet_id | yes | FK wallets |
| tenant_id | yes | Scope |
| direction | yes | `credit`, `debit` |
| amount | yes | positive decimal |
| currency | yes |  |
| entry_type | yes | `topup`, `purchase`, `reseller_cost`, `refund`, `adjustment`, `reversal`, `commission`, `lock`, `unlock` |
| status | yes | `posted`, `voided_by_reversal` |
| balance_after | yes | Balance sau entry |
| reference_type | yes | `order`, `topup_request`, `service`, `manual_adjustment`, `settlement` |
| reference_id | yes |  |
| idempotency_key | yes | Unique theo wallet/action |
| created_by | no | actor |
| reason | no | required với adjustment |
| correlation_id | yes | trace |
| created_at | yes | immutable |

Constraints:
```text
unique(wallet_id, idempotency_key)
amount > 0
posted entry cannot be updated/deleted
```

Index:
```text
wallet_id, created_at
tenant_id, reference_type, reference_id
correlation_id
```

### 7.3 `topup_requests`

| Field | Required | Note |
|---|---:|---|
| topup_request_id | yes |  |
| tenant_id | yes |  |
| wallet_id | yes |  |
| requested_by | yes | user |
| amount | yes |  |
| currency | yes |  |
| payment_method | yes | `bank_transfer`, `crypto`, `manual`, `other` |
| payment_reference | no |  |
| proof_attachment_id | no | private file |
| status | yes | `draft`, `submitted`, `under_review`, `approved`, `rejected`, `expired`, `cancelled` |
| reviewed_by | no |  |
| reviewed_at | no |  |
| review_note | no |  |
| ledger_id | no | set when approved |
| created_at / updated_at | yes |  |

Constraint:
```text
approved topup must have ledger_id
```

---

## 8. Order, reservation, service tables

### 8.1 `orders`

| Field | Required | Note |
|---|---:|---|
| order_id | yes | Primary key |
| tenant_id | yes | Tenant selling to client |
| client_user_id | yes | Buyer |
| seller_tenant_id | yes | Reseller tenant hoặc platform tenant |
| order_number | yes | Human readable unique |
| order_type | yes | `new_service`, `renewal`, `upgrade`, `addon` |
| order_status | yes | `draft`, `pending_payment`, `paid`, `provisioning`, `active`, `failed`, `manual_review`, `cancelled`, `refunded`, `expired` |
| payment_status | yes | `unpaid`, `paid`, `partially_refunded`, `refunded`, `failed` |
| provisioning_status | yes | `not_started`, `queued`, `running`, `success`, `failed`, `partial_success`, `manual_review` |
| subtotal | yes |  |
| discount_amount | yes |  |
| total_amount | yes | selling price total |
| currency | yes |  |
| client_wallet_debit_ledger_id | no |  |
| reseller_wallet_debit_ledger_id | no |  |
| idempotency_key | yes | Unique per buyer checkout |
| correlation_id | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(tenant_id, order_number)
unique(client_user_id, idempotency_key)
```

### 8.2 `order_items`

| Field | Required | Note |
|---|---:|---|
| order_item_id | yes |  |
| order_id | yes | FK orders |
| tenant_id | yes |  |
| tenant_plan_id | yes |  |
| master_plan_id | yes |  |
| quantity | yes | phase 1 thường = 1 |
| unit_price | yes | selling price snapshot |
| reseller_unit_cost | yes | reseller cost snapshot |
| billing_cycle_snapshot | yes | json |
| product_snapshot | yes | json |
| plan_snapshot | yes | json |
| capability_snapshot | yes | json |
| provider_source_snapshot | yes | json |
| policy_snapshot | yes | json |
| created_at | yes |  |

### 8.3 `reservations`

| Field | Required | Note |
|---|---:|---|
| reservation_id | yes |  |
| tenant_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| source_id | yes | provider source |
| quantity | yes |  |
| status | yes | `reserved`, `allocated`, `released`, `expired`, `failed` |
| reserved_at | yes |  |
| expires_at | yes | thường now + 5 phút |
| allocated_at | no |  |
| released_at | no |  |
| idempotency_key | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(source_id, idempotency_key)
allocated reservation cannot expire/release again
```

### 8.4 `services`

| Field | Required | Note |
|---|---:|---|
| service_id | yes |  |
| tenant_id | yes |  |
| client_user_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| reservation_id | yes |  |
| service_type | yes | `vps`, `proxy` |
| service_status | yes | `pending`, `active`, `suspended`, `expired`, `terminated`, `cancelled`, `failed` |
| billing_status | yes | `paid`, `due_soon`, `overdue`, `grace`, `cancelled` |
| suspension_reason | no | required if suspended |
| source_id | yes | provider source |
| external_resource_id | no | Provider resource id |
| service_identifier | no | IP/hostname/proxy endpoint |
| term_start_at | yes |  |
| term_end_at | yes |  |
| grace_until_at | no |  |
| terminate_after_at | no |  |
| capability_snapshot | yes | json |
| plan_snapshot | yes | json |
| price_snapshot | yes | json |
| created_at / updated_at | yes |  |

Index:
```text
tenant_id, client_user_id, service_status
term_end_at
source_id, external_resource_id
```

### 8.5 `service_credentials`

| Field | Required | Note |
|---|---:|---|
| credential_id | yes |  |
| service_id | yes | FK services |
| tenant_id | yes |  |
| credential_type | yes | `vps_root`, `proxy_auth`, `ssh_key`, `console_url`, `api_token` |
| encrypted_payload | yes | encrypted JSON |
| secret_version | yes | key rotation support |
| masked_hint | no | ví dụ `root / ********` |
| last_revealed_at | no |  |
| last_revealed_by | no | user |
| status | yes | `active`, `rotated`, `revoked` |
| created_at / updated_at | yes |  |

Security:
```text
Không lưu plaintext ở bất kỳ field nào.
Mỗi reveal phải ghi audit credential.revealed.
```

### 8.6 `service_lifecycle_events`

| Field | Required | Note |
|---|---:|---|
| event_id | yes |  |
| tenant_id | yes |  |
| service_id | yes |  |
| event_type | yes | `created`, `activated`, `renewed`, `expired`, `suspended`, `unsuspended`, `terminated`, `cancelled`, `failed` |
| from_status | no |  |
| to_status | no |  |
| reason | no | required with suspend/terminate |
| actor_id | no | user or system |
| job_id | no | if worker/cron |
| correlation_id | yes |  |
| created_at | yes |  |

---

## 9. Provider and provisioning tables

### 9.1 `provider_accounts`

| Field | Required | Note |
|---|---:|---|
| provider_account_id | yes |  |
| name | yes |  |
| provider_type | yes | `proxmox`, `ovh`, `hetzner`, `manual`, `proxy_upstream` |
| status | yes | `active`, `disabled`, `maintenance` |
| encrypted_credentials | yes | Provider API credentials |
| allowed_ips | no | API allowlist if used |
| health_status | yes | `unknown`, `healthy`, `degraded`, `down` |
| last_health_check_at | no |  |
| created_at / updated_at | yes |  |

Security:
```text
provider credentials encrypt at rest.
Không trả về qua API thông thường.
```

### 9.2 `provider_inventory`

| Field | Required | Note |
|---|---:|---|
| inventory_id | yes |  |
| source_id | yes |  |
| capacity_total | no | null nếu live/unlimited |
| reserved_count | yes |  |
| allocated_count | yes |  |
| available_count_cache | yes | derived/cache |
| last_synced_at | no |  |
| status | yes | `active`, `out_of_stock`, `sync_failed`, `disabled` |
| updated_at | yes |  |

Rule:
```text
available = capacity_total - reserved_count - allocated_count
Nếu capacity_total null và inventory_mode là provider_live thì phải check live trước reserve.
```

### 9.3 `provisioning_jobs`

| Field | Required | Note |
|---|---:|---|
| job_id | yes |  |
| tenant_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| reservation_id | yes |  |
| source_id | yes |  |
| job_type | yes | `provision`, `renew`, `suspend`, `unsuspend`, `terminate`, `sync`, `reset_password`, `change_ip` |
| status | yes | `queued`, `running`, `success`, `failed`, `partial_success`, `manual_review`, `cancelled` |
| idempotency_key | yes | unique |
| attempt_count | yes |  |
| max_attempts | yes |  |
| next_retry_at | no |  |
| last_error_code | no |  |
| last_error_message_redacted | no |  |
| manual_review_reason | no |  |
| correlation_id | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(source_id, job_type, idempotency_key)
```

### 9.4 `provider_requests`

| Field | Required | Note |
|---|---:|---|
| provider_request_id | yes |  |
| job_id | yes | FK provisioning_jobs |
| source_id | yes |  |
| request_type | yes | same as job/action |
| external_request_id | no | nếu provider có |
| external_resource_id | no | nếu biết |
| status | yes | `sent`, `success`, `failed`, `timeout`, `unknown`, `manual_review` |
| retry_safety | yes | `safe_retry`, `unsafe_retry`, `do_not_retry`, `manual_review_required` |
| request_summary_redacted | no | không chứa secret |
| response_summary_redacted | no | không chứa secret |
| error_code | no |  |
| sent_at | yes |  |
| received_at | no |  |
| correlation_id | yes |  |

### 9.5 `provider_resource_mappings`

| Field | Required | Note |
|---|---:|---|
| mapping_id | yes |  |
| tenant_id | yes |  |
| service_id | yes |  |
| source_id | yes |  |
| external_resource_id | yes |  |
| external_status | no |  |
| last_synced_at | no |  |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(source_id, external_resource_id)
```

---

## 10. Audit, risk, abuse, notification tables

### 10.1 `audit_logs`

| Field | Required | Note |
|---|---:|---|
| audit_id | yes |  |
| tenant_id | no | null nếu platform-level |
| actor_id | no | null nếu system |
| actor_type | yes | `user`, `system`, `worker`, `provider_webhook` |
| action | yes | e.g. `wallet.topup.approved` |
| target_type | yes |  |
| target_id | yes |  |
| before_snapshot_redacted | no | json |
| after_snapshot_redacted | no | json |
| metadata_redacted | no | json |
| ip_address | no |  |
| user_agent | no |  |
| correlation_id | yes |  |
| created_at | yes | immutable |

Index:
```text
tenant_id, created_at
actor_id, created_at
target_type, target_id
correlation_id
action
```

### 10.2 `risk_flags`

| Field | Required | Note |
|---|---:|---|
| risk_flag_id | yes |  |
| tenant_id | yes |  |
| user_id | no |  |
| service_id | no |  |
| order_id | no |  |
| flag_type | yes | `new_account_high_value`, `payment_mismatch`, `abuse_history`, `manual_blacklist`, `provider_risk` |
| severity | yes | `low`, `medium`, `high`, `critical` |
| status | yes | `open`, `reviewing`, `cleared`, `confirmed` |
| note | no |  |
| created_by | no |  |
| created_at / updated_at | yes |  |

### 10.3 `abuse_cases`

| Field | Required | Note |
|---|---:|---|
| abuse_case_id | yes |  |
| tenant_id | yes |  |
| service_id | no |  |
| user_id | no |  |
| provider_source_id | no |  |
| case_type | yes | `spam`, `scan`, `bruteforce`, `copyright`, `fraud`, `aup_violation`, `other` |
| severity | yes | `low`, `medium`, `high`, `critical` |
| status | yes | `open`, `investigating`, `warning_sent`, `suspended`, `resolved`, `closed` |
| evidence_private | no | json/file refs |
| action_taken | no |  |
| created_at / updated_at | yes |  |

### 10.4 `notifications`

| Field | Required | Note |
|---|---:|---|
| notification_id | yes |  |
| tenant_id | yes |  |
| recipient_user_id | no |  |
| channel | yes | `email`, `dashboard`, `telegram`, `webhook` |
| template_key | yes |  |
| status | yes | `queued`, `sent`, `failed`, `cancelled` |
| payload_redacted | yes | json |
| reference_type | no |  |
| reference_id | no |  |
| sent_at | no |  |
| created_at / updated_at | yes |  |

---

## 11. Recommended enums

### 11.1 Order status

```text
draft
pending_payment
paid
provisioning
active
failed
manual_review
cancelled
refunded
expired
```

### 11.2 Provisioning status/job status

```text
not_started
queued
running
success
failed
partial_success
manual_review
cancelled
```

### 11.3 Service status

```text
pending
active
suspended
expired
terminated
cancelled
failed
```

### 11.4 Reservation status

```text
reserved
allocated
released
expired
failed
```

### 11.5 Wallet ledger entry type

```text
topup
purchase
reseller_cost
refund
adjustment
reversal
commission
lock
unlock
```

---

## 12. Critical indexes

P0 indexes:
```text
users(tenant_id, email)
tenant_domains(domain)
tenant_plans(tenant_id, status, visibility)
wallets(owner_type, owner_id, currency)
wallet_ledger_entries(wallet_id, created_at)
wallet_ledger_entries(correlation_id)
orders(tenant_id, client_user_id, created_at)
orders(tenant_id, order_number)
reservations(source_id, status, expires_at)
services(tenant_id, client_user_id, service_status)
services(term_end_at, service_status)
provisioning_jobs(status, next_retry_at)
provisioning_jobs(correlation_id)
provider_resource_mappings(source_id, external_resource_id)
audit_logs(tenant_id, created_at)
audit_logs(correlation_id)
```

---

## 13. Data retention

Suggested phase 1:
```text
wallet_ledger_entries: keep forever
orders/services: keep forever or legal retention
audit_logs financial/security: keep long-term
provider_request detailed redacted body: keep 90-180 days
notification payload: keep 30-90 days
temporary payment proof files: follow finance retention policy
```

Không xóa cứng dữ liệu tài chính. Dùng status/archival.

---

## 14. Acceptance criteria database

Schema được xem là đạt khi:

- Mọi bảng tenant-owned có `tenant_id`.
- Ledger posted không thể sửa/xóa ở tầng application.
- Credential/provider secret không có plaintext field.
- Checkout có thể trace bằng `correlation_id`.
- Provisioning có unique idempotency key.
- Provider external resource có unique mapping.
- Reservation không thể allocate/release hai lần.
- Service/order lưu đủ snapshot để xử lý tranh chấp sau này.
- Có index cho các query P0: wallet, order, service, provisioning, audit, expiry.
- Admin emergency access để lại audit đầy đủ.

Câu nền: **database không chỉ lưu dữ liệu; nó là hàng rào chống con người, provider và bug làm sai tiền.**
