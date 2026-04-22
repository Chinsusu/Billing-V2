# CHANGELOG_FIXES - Bản vá tài liệu v1.1

## Mục tiêu bản vá
Bản v1.1 vá các lỗ hổng P0 trong bộ tài liệu gốc để biến blueprint thành build-spec rõ hơn cho dev/ops, vẫn giữ phạm vi “chưa code”.

## Các lỗi/thiếu sót đã vá
1. Thêm mô hình reseller settlement: client wallet, reseller wallet, platform revenue, reseller profit.
2. Khóa rule không provision nếu reseller wallet thiếu reseller cost.
3. Bổ sung tenant enforcement: không tin tenant_id từ body, backend bắt buộc scope theo tenant context.
4. Bổ sung emergency access có reason và audit.
5. Chuyển 2FA Admin thành P0 phase 1; Reseller Owner bật mặc định/khuyến nghị bắt buộc.
6. Bổ sung catalog versioning, snapshot, propagation rule và margin guard.
7. Bổ sung top-up state machine và financial invariants.
8. Bổ sung idempotency key, partial success, retry safety và no retry mù.
9. Bổ sung atomic inventory reservation, reserved/allocated counters và expiry release.
10. Bổ sung credential security: encrypt, masked reveal, audit, redaction.
11. Bổ sung renew/cancel/refund guard và calendar month rule.
12. Bổ sung acceptance criteria cho flow/API/data model.
13. Bổ sung manual abuse/fraud controls phase 1.
14. Bổ sung audit naming, redaction và correlation_id xuyên flow.
15. Bổ sung report formula cho Admin và Reseller.

## File mới được thêm
- `09_Reseller_Settlement_Ledger_Model.md`
- `10_Tenant_Security_Access_Control_Spec.md`
- `11_Provisioning_Idempotency_And_Inventory_Locking.md`
- `12_API_Data_Model_Acceptance_Criteria.md`
- `13_Abuse_Fraud_Operational_Policy_Phase1.md`

## File cũ đã được vá
- `00_README.md`
- `01_Product_Scope_Business_Model.md`
- `02_Tenant_Model_Role_Architecture.md`
- `03_Product_Catalog_Pricing_Rules.md`
- `04_Billing_Wallet_Ledger_Spec.md`
- `05_Provisioning_Provider_Adapter_Spec.md`
- `06_Order_Service_Lifecycle_State_Machine.md`
- `07_Portal_Functional_Spec.md`
- `08_Audit_Reports_Operational_Control.md`

## Gợi ý dùng với dev
Đưa dev đọc theo thứ tự:
1. `00_README.md`
2. `CHANGELOG_FIXES.md`
3. `09`, `10`, `11` trước vì đây là ba tài liệu P0 nhất.
4. `12` để chuyển sang backlog/API/data model.
5. `13` để khóa vận hành rủi ro trước production.

## Bổ sung bản v1.2 - Technical Handoff
Bản v1.2 không thay đổi nguyên tắc vá lỗi P0 của v1.1, mà bổ sung bộ tài liệu `14–23` để dev/backend/frontend/QA/DevOps có thể chuyển blueprint thành kế hoạch build cụ thể hơn.

Các nhóm tài liệu mới:
- Architecture
- Database schema/ERD
- API contract
- RBAC matrix
- Provider adapter contract
- Worker/queue/cron
- UI wireflow
- QA acceptance plan
- DevOps runbook
- Notification templates
