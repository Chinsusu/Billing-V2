# Tài liệu 12 - API, Data Model & Acceptance Criteria

## 1. Mục tiêu tài liệu
Tài liệu này không viết code. Nó mô tả data model nghiệp vụ, API flow và acceptance criteria để đội dev biết cần build gì và test thế nào.

## 2. Shared fields nên có
Các bảng nghiệp vụ quan trọng nên có:
- id
- tenant_id nếu thuộc tenant
- created_at
- updated_at
- created_by
- updated_by
- status
- metadata hoặc notes nếu cần
- request_id/correlation_id với flow tài chính/provisioning

Các bảng tài chính/provisioning nên có thêm:
- idempotency_key
- reference_type/reference_id
- snapshot fields

## 3. Entities P0
### 3.1. Tenant
Fields:
- tenant_id
- tenant_type: admin, reseller
- status: active, suspended, disabled
- brand_name
- default_currency_display
- timezone
- contact_support

### 3.2. Domain mapping
Fields:
- domain_id
- tenant_id
- domain
- verification_status
- tls_status
- status
- verification_token_reference

### 3.3. User
Fields:
- user_id
- tenant_id/seller_scope
- role
- email
- status
- two_factor_enabled
- last_login_at

### 3.4. Wallet
Fields:
- wallet_id
- tenant_id
- owner_type: admin_client, reseller, reseller_client
- owner_id
- currency
- available_balance
- locked_balance
- status

### 3.5. Ledger entry
Fields:
- ledger_entry_id
- tenant_id
- wallet_id
- direction
- amount
- currency
- fx_snapshot
- entry_type
- reference_type/reference_id
- idempotency_key
- balance_before
- balance_after
- status
- reason

### 3.6. Product / Plan / Source
Product fields:
- product_id
- type
- display_name
- location
- visibility

Plan fields:
- plan_id
- product_id
- version
- cycle_type
- cycle_definition
- retail_price
- reseller_cost
- refund_policy
- capability_flags
- status

Source fields:
- source_id
- provider_id
- plan_id
- capacity
- reserved_count
- allocated_count
- stock_state
- capability_override
- status

### 3.7. Tenant plan clone
Fields:
- tenant_plan_id
- tenant_id
- master_plan_id
- master_plan_version
- selling_price
- reseller_cost_snapshot
- enabled
- margin_state
- policy_snapshot

### 3.8. Order
Fields:
- order_id
- tenant_id
- buyer_user_id
- seller_type: admin, reseller
- plan_snapshot
- price_snapshot
- reseller_cost_snapshot
- order_status
- billing_status
- correlation_id

### 3.9. Reservation
Fields:
- reservation_id
- tenant_id
- order_id
- source_id
- status
- quantity
- expires_at

### 3.10. Provisioning job
Fields:
- job_id
- tenant_id
- order_id
- reservation_id
- source_id
- provider_id
- idempotency_key
- external_request_id
- external_resource_id
- status
- attempt_count
- retry_safety_level
- last_error_summary

### 3.11. Service instance
Fields:
- service_id
- tenant_id
- owner_user_id
- order_id
- provider_id
- source_id
- external_resource_id
- service_status
- billing_status
- suspension_reason
- term_start
- term_end
- grace_end
- capability_snapshot
- billing_cycle_snapshot
- credential_reference

### 3.12. Audit event
Fields:
- audit_id
- tenant_id
- actor_id
- actor_role
- action
- resource_type
- resource_id
- before_summary
- after_summary
- reason
- request_id
- correlation_id
- source_ip
- created_at

### 3.13. Risk/abuse flag
Fields:
- risk_flag_id
- tenant_id
- target_type: user, order, service, payment, ip, domain
- target_id
- risk_type
- severity
- status
- evidence_summary
- created_by

## 4. API flow P0 và acceptance criteria
### 4.1. Register/Login
Acceptance:
- User được gắn đúng tenant từ domain context.
- Email trùng trong cùng tenant bị chặn theo policy.
- Login failed bị rate limit.
- Admin login không có 2FA bị chặn nếu policy bắt buộc.

### 4.2. Create top-up request
Acceptance:
- Request thuộc đúng tenant/user.
- Amount > 0.
- Payment method enabled.
- Sinh reference duy nhất.
- Không credit wallet trước khi approved.

### 4.3. Approve top-up
Acceptance:
- Actor có quyền approve.
- Một request chỉ approve một lần.
- Credit wallet và ledger xảy ra cùng một control flow.
- Audit có actor/reason/reference.

### 4.4. Checkout Admin client
Acceptance:
- Plan active.
- Stock available.
- Client wallet đủ tiền.
- Reservation atomic.
- Wallet debit + ledger + order paid + provisioning job được tạo an toàn.
- Double click không double debit.

### 4.5. Checkout Reseller client
Acceptance:
- Client wallet đủ selling price.
- Reseller wallet đủ reseller cost.
- Thiếu một trong hai ví thì không provision.
- Order lưu selling price snapshot và reseller cost snapshot.
- Report reseller profit tính đúng.

### 4.6. Renew service
Acceptance:
- Service active hoặc suspended trong grace.
- Term mới cộng từ old_term_end theo cycle.
- Nếu thuộc reseller, debit cả client wallet và reseller wallet.
- Terminated service không renew.

### 4.7. Cancel/refund
Acceptance:
- Chỉ plan có allow_mid_cycle_cancel mới hiện action.
- Refund dùng policy snapshot.
- Refund không vượt charge gốc.
- Service không còn active sau refund toàn phần/cancel hợp lệ.
- Audit có reason.

### 4.8. Service action
Các action: start, stop, reboot, reinstall, change password, console, change IP.

Acceptance:
- Action chỉ hiện nếu capability_snapshot cho phép.
- Tenant ownership được check backend.
- Paid action debit wallet trước khi gọi provider.
- Provider fail thì refund/rollback theo policy action.
- Action log/audit được ghi.

### 4.9. Reveal credential
Acceptance:
- Credential masked mặc định.
- Reveal yêu cầu quyền hợp lệ.
- Cross-tenant bị chặn.
- Audit `credential.revealed` được ghi.

### 4.10. Reseller pricing update
Acceptance:
- Chỉ Owner/staff có quyền mới update.
- Giá mới không ảnh hưởng order cũ.
- Negative margin bị cảnh báo hoặc block theo policy.
- Audit `pricing.tenant_plan.updated`.

### 4.11. Domain mapping
Acceptance:
- Domain phải verify trước active.
- Không active domain trùng tenant khác.
- Mapping change audit.
- Domain suspended không route vào tenant.

## 5. Error states P0
UI/API cần trả trạng thái rõ cho:
- insufficient_client_balance
- insufficient_reseller_balance
- out_of_stock
- reservation_expired
- plan_disabled
- source_disabled
- provider_unavailable
- permission_denied
- tenant_scope_mismatch
- provisioning_manual_review
- payment_verification_required
- refund_not_allowed

## 6. Acceptance test tối thiểu trước production
- Cross-tenant API tests.
- Double-click checkout tests.
- Concurrent reservation tests.
- Top-up approve idempotency tests.
- Reseller settlement tests.
- Provider timeout/partial success tests.
- Refund boundary tests.
- Credential redaction tests.
- Audit correlation tests.
