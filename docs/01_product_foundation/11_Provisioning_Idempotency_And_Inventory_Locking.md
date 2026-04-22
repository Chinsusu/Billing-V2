# Tài liệu 11 - Provisioning Idempotency & Inventory Locking

## 1. Mục tiêu tài liệu
Tài liệu này khóa cách reserve stock, debit wallet, tạo provisioning job, retry provider và xử lý partial success.

Mục tiêu: không oversell, không cấp trùng VPS/proxy, không mất tiền do retry sai.

## 2. Nguyên tắc lõi
- Reservation phải atomic.
- Debit wallet và ledger phải idempotent.
- Gọi provider là bước bất đồng bộ, không giữ DB transaction mở.
- Timeout sau khi gửi request tới provider là tình huống nguy hiểm, không retry mù.
- Một order chỉ attach một service chính thức.

## 3. Inventory counters
Mỗi source nên theo dõi:
- capacity
- reserved_count
- allocated_count
- unavailable_count nếu có
- stock_state: available, low_stock, out_of_stock, disabled, unknown

Available tính theo:
`available = capacity - reserved_count - allocated_count - unavailable_count`

Rule:
- Không dùng read-then-write thường cho reservation.
- Phải có atomic update/row lock/transaction guard.

## 4. Reservation lifecycle
Trạng thái:
- pending_reserve
- reserved
- reservation_expired
- reservation_released
- allocated

Rule:
- Reservation mặc định hết hạn sau 5 phút.
- Hết hạn phải release reserved_count.
- Khi provisioning bắt đầu hoặc resource được gắn chắc chắn, reservation chuyển allocated.
- Nếu order cancel trước provisioning, reservation released.

## 5. Checkout transaction boundary
### 5.1. Trong DB transaction
- Validate plan/source active.
- Reserve inventory atomic.
- Validate wallet balances.
- Debit wallet(s).
- Ghi ledger entries.
- Tạo order paid.
- Tạo provisioning job queued.

### 5.2. Ngoài DB transaction
- Worker lấy job queued.
- Adapter gọi provider.
- Provider response được ghi lại.
- Service instance được tạo/cập nhật.

Rule:
- Không gọi provider trong transaction đang khóa wallet/inventory.

## 6. Idempotency key
Mỗi flow cần key riêng:
- checkout idempotency key
- wallet debit idempotency key
- reservation idempotency key
- provisioning job idempotency key
- provider request id nếu provider hỗ trợ

Format có thể do dev quyết, nhưng phải đủ để ngăn double click/double retry tạo giao dịch trùng.

## 7. Provisioning job fields
- job_id
- tenant_id
- order_id
- reservation_id
- service_id nếu đã tạo
- provider_id
- source_id
- idempotency_key
- external_request_id
- external_resource_id
- status: queued, provisioning, provisioned, failed, manual_review, cancelled
- attempt_count
- max_attempts
- retry_safety_level: safe, unsafe, unknown
- last_error_code
- last_error_summary
- last_provider_response_summary
- next_retry_at
- created_at
- updated_at

## 8. Retry matrix
| Tình huống | retry_safety_level | Hành động |
|---|---|---|
| Provider reject do invalid payload | safe/failed | Không retry cho tới khi sửa payload |
| Provider hết hàng rõ ràng | safe/failed | Release/refund hoặc chọn source khác thủ công |
| Network fail trước khi request gửi | safe | Retry được |
| Timeout sau khi request đã gửi | unknown | Manual review |
| Provider trả external_request_id nhưng chưa rõ resource | unknown | Poll/status check trước khi retry |
| Provider trả resource_id nhưng credential thiếu | unsafe | Không tạo resource mới; recover credential/manual |
| Webhook lặp | safe | Idempotent update |
| DB lưu service fail sau provider success | unsafe | Manual recovery, không retry create |

## 9. Partial success handling
Partial success là khi provider có thể đã tạo resource nhưng hệ thống chưa hoàn tất.

Ví dụ:
- Provider timeout sau create request.
- Provider trả resource id nhưng credential response lỗi.
- Provider tạo VPS nhưng callback/webhook không tới.
- DB crash sau khi provider success.

Rule:
- Chuyển job manual_review.
- Gắn external_request_id/resource_id nếu có.
- Không refund tự động trước khi xác minh.
- Không retry create resource mới.
- Ops phải có màn hình resolve: attach existing resource, fetch credential, terminate duplicate, refund.

## 10. Service creation rule
Service instance chỉ chuyển active khi:
- Có provider/source mapping.
- Có external_resource_id hoặc bằng chứng resource cấp thành công.
- Credential/access info đã lưu an toàn hoặc có handoff thủ công.
- Capability snapshot đã ghi.
- Billing term đã ghi.

## 11. Credential persistence
Nếu provider trả credential một lần duy nhất:
- Lưu encrypted ngay khi nhận.
- Nếu lưu thất bại, job manual_review ngay.
- Không gọi lại create để lấy credential mới.

## 12. Inventory reconciliation
Cần job đối soát:
- source capacity từ provider
- resource tồn tại ở provider nhưng không có trong DB
- DB active nhưng provider missing
- reservation expired nhưng chưa release
- allocated_count lệch với service active/suspended

Reconciliation nên tạo alert/manual review, không tự sửa tài chính nếu chưa có rule rõ.

## 13. No auto failover phase 1
Không tự chuyển sang provider khác khi source fail/out-of-stock vì:
- giá vốn khác
- capability khác
- billing cycle khác
- refund/change IP policy khác

Chỉ cho Admin chọn source khác thủ công nếu user chấp nhận hoặc policy đã rõ.

## 14. Acceptance criteria
- 10 checkout đồng thời cho 1 slot chỉ có 1 reservation thành công.
- User double click checkout không bị debit hai lần.
- Provider timeout sau create request không tạo retry mù.
- Webhook lặp không tạo service trùng.
- Reservation expired tự release stock.
- Service active luôn có capability snapshot và billing snapshot.
- Credential không bị log plaintext.
