# CHANGELOG_TECHNICAL_HANDOFF_v1_2

## Mục tiêu bản v1.2
Bản v1.2 bổ sung bộ Technical Build Handoff Package để chuyển bộ blueprint/vá lỗi v1.1 thành tài liệu dev/backend/frontend/QA/DevOps có thể dùng để build.

Bản này vẫn giữ nguyên nguyên tắc:
- Chưa code.
- Không chọn framework cụ thể.
- Ưu tiên khóa behavior, data contract, permission, state, queue, QA và vận hành.

## File mới được thêm

### `14_System_Architecture_Blueprint.md`
Bổ sung kiến trúc tổng thể:
- portal layer.
- API/auth/tenant/RBAC.
- core modules.
- financial core.
- security core.
- queue/worker/cron.
- provider adapter layer.
- observability và fail-safe principle.

### `15_Database_Schema_And_ERD.md`
Bổ sung data contract:
- tenants/domains/users/roles/permissions.
- catalog/product/plan/source.
- wallet/ledger/top-up.
- order/item/reservation/service/credential.
- provider/provisioning/resource mapping.
- audit/risk/abuse/notification.
- index, constraint, enum và acceptance criteria.

### `16_API_Contract_And_Permission_Spec.md`
Bổ sung API contract:
- conventions.
- auth/tenant/catalog/wallet/order/service/provider/audit/risk APIs.
- role/permission.
- validation.
- error codes.
- idempotency.
- rate limit.
- audit action.

### `17_RBAC_Permission_Matrix.md`
Bổ sung quyền chi tiết:
- Platform Super Admin.
- Platform Staff.
- Finance Agent.
- Support Agent.
- Provisioning Operator.
- Reseller Owner.
- Reseller Staff.
- Client.
- Read-only Auditor.
- permission matrix và critical controls.

### `18_Provider_Adapter_Technical_Spec.md`
Bổ sung adapter contract:
- capability profile.
- operation result.
- provision/status/suspend/terminate/renew/reset/reinstall/change IP.
- retry safety.
- error normalization.
- idempotency.
- credential handling.
- manual provider.
- provider onboarding checklist.

### `19_Worker_Queue_And_Cron_Jobs_Spec.md`
Bổ sung job nền:
- provisioning_worker.
- provider_sync_worker.
- service_action_worker.
- notification_worker.
- reservation_expiry_job.
- service_expiry/suspension/termination jobs.
- provider health/inventory sync.
- manual review queue.
- monitoring metrics.

### `20_UI_Wireflow_And_Screen_Spec.md`
Bổ sung screen spec:
- Client Portal.
- Reseller Portal.
- Admin Portal.
- service detail/credential reveal.
- checkout flow.
- wallet/top-up.
- catalog/pricing.
- admin provisioning/manual review.
- common status/timeline/confirmation/error states.

### `21_QA_Test_Cases_And_Acceptance_Plan.md`
Bổ sung test plan:
- tenant isolation.
- RBAC.
- wallet/ledger.
- checkout/reservation.
- provisioning.
- service lifecycle.
- credential security.
- catalog/pricing.
- abuse/risk.
- notification/report.
- deployment smoke tests.

### `22_Deployment_DevOps_And_Environment_Runbook.md`
Bổ sung runbook:
- local/staging/production.
- secret management.
- DB migration.
- queue/worker.
- logging/monitoring/alert.
- backup/restore.
- release/rollback.
- incident response.
- go-live checklist.

### `23_Notification_Email_Telegram_Template_Spec.md`
Bổ sung notification:
- event naming.
- channels.
- priority/dedupe.
- client/reseller/admin templates.
- payload security.
- timing matrix.
- audit/notification record.

## File được cập nhật
- `00_README.md`: cập nhật mô tả v1.2, danh sách tài liệu và cách đọc.
- `MANIFEST.txt`: cập nhật danh sách file.
- `VPS_Proxy_Project_Master_Document.md`: tổng hợp lại toàn bộ tài liệu `00–23`.

## Điểm khóa thêm ở bản v1.2
1. Dev có schema logic đủ để thiết kế DB.
2. Backend có API behavior và error code rõ.
3. Frontend có wireflow/screen/action/state rõ.
4. QA có test case để nghiệm thu P0.
5. DevOps có checklist production, backup, monitoring, rollback.
6. Provider adapter có retry/idempotency/partial success behavior rõ.
7. RBAC không còn nằm mơ hồ trong tenant doc.
8. Notification được xem là hệ thống vận hành, không phải phần trang trí.

## Khuyến nghị bước tiếp theo
Sau v1.2, nếu chuẩn bị thuê dev hoặc triển khai sprint, nên làm tiếp:
- Backlog/Sprint plan theo milestone.
- User stories + acceptance criteria theo từng sprint.
- Data migration checklist khi chọn framework/database thật.
- Provider-specific adapter spec cho provider đầu tiên, ví dụ Proxmox hoặc một proxy upstream cụ thể.
