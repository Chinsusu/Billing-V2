# Tài liệu 05 - Provisioning & Provider Adapter Specification

## 1. Mục tiêu tài liệu
Tài liệu này định nghĩa:
- Adapter layer cho nhiều provider
- Inventory model
- Reservation
- Provisioning flow
- Failure handling
- Sync strategy

## 2. Nguyên tắc thiết kế
- Một giao diện chuẩn cho provider
- Capability khác nhau giữa provider phải được mô hình hóa, không hard-code
- Provisioning fail phải có đường lui rõ ràng
- Phase 1 ưu tiên tính kiểm soát hơn là tự động hóa tuyệt đối

## 3. Danh sách provider ban đầu
### 3.1. VPS
- Proxmox
- OVH
- Hetzner
- Smarthost

### 3.2. Proxy
- proxy-manager
- proxy-cheap

## 4. Provider adapter layer
Mỗi provider adapter nên hỗ trợ một tập chuẩn các hàm nghiệp vụ:
- inventory_sync()
- provision()
- suspend()
- terminate()
- reboot()
- reinstall()
- change_password()
- fetch_console()
- fetch_usage()
- change_ip()

Không bắt buộc provider nào cũng hỗ trợ toàn bộ. Action không hỗ trợ phải được đánh dấu capability = false.

## 5. Inventory model
Inventory cần được theo dõi theo ít nhất các chiều:
- provider
- location
- protocol
- quality
- resource bucket / pool
- availability

### 5.1. Source
Source là đơn vị nguồn hàng thực tế. Ví dụ:
- Proxmox node/pool
- OVH plan mapping
- Hetzner plan mapping
- proxy-manager pool
- proxy-cheap pool

## 6. Reservation
- Khi user bắt đầu thanh toán/mua, hệ thống reserve stock tối đa 5 phút
- Nếu thanh toán không hoàn thành hoặc order bị hủy, reservation bị release
- Nếu provisioning bắt đầu, reservation chuyển thành allocation

## 7. Provisioning flow chuẩn
1. Order paid
2. Reserve source
3. Push vào provisioning queue
4. Adapter gọi API provider
5. Nhận kết quả
6. Tạo service instance
7. Lưu credential/info
8. Chuyển trạng thái active
9. Notify user

## 8. Failure handling
### 8.1. Khi provision fail
- Chuyển order/service sang failed hoặc manual_review
- Alert Admin
- Cho phép xử lý thủ công

### 8.2. Các cách xử lý thủ công
- Retry provisioning
- Chọn source khác thủ công
- Refund/rollback
- Hủy order

## 9. Không auto failover ở phase 1
Rule đã chốt:
- Khi source/provider hết hàng, không tự nhảy sang provider khác
- UI sẽ disable/gray-out plan hoặc source tương ứng

Lý do:
- Tránh lệch capability
- Tránh lệch giá vốn
- Tránh cấp sai nguồn không đúng kỳ vọng

## 10. Sync strategy
Không ép tất cả provider theo một cơ chế duy nhất. Cho phép:
- Realtime
- Webhook
- Cron

Nhưng cần một model chuẩn để lưu:
- last_sync_at
- sync_status
- recent_error
- availability_state

## 11. Provider health monitoring
### Phase 1
- Connection status
- Last sync
- Recent errors
- Provision success/failure count

### Phase 2
- Health score tổng hợp
- API latency
- Error rate window
- Inventory freshness score

## 12. UI capability masking
- Action nào không được provider hỗ trợ thì không render trên UI
- Capability phải có thể override theo source hoặc plan cụ thể
- Không dựa hoàn toàn vào loại sản phẩm chung

## 13. Dữ liệu business nên lưu cho service instance
- provider_id
- source_id
- external_resource_id
- provision_payload_snapshot
- capability_snapshot
- credentials/access_info
- billing_cycle_snapshot
- lifecycle_status
- last_sync_status

## 14. Queue và xử lý bất đồng bộ
Đề xuất:
- Có provisioning queue riêng
- Có sync queue riêng
- Có alerting channel cho failed jobs

Nếu phase 1 chưa có queue đầy đủ, vẫn nên thiết kế theo tư duy queue để dễ nâng cấp.

## 15. Những lỗi phổ biến cần tính trước
- Provider timeout
- Provider trả partial success
- Provider cấp tài nguyên xong nhưng callback lỗi
- Credential tạo thành công nhưng chưa lưu DB
- Stock thực tế khác stock sync
- Action idempotency không chắc chắn

## 16. Rule idempotency tối thiểu
- Mỗi order chỉ được attach một service instance chính thức
- Retry phải có idempotency key
- Webhook hoặc callback từ provider phải xử lý lặp an toàn

## 17. Bản vá v1.1 - Idempotency, locking và credential safety
File này cần đọc cùng `11_Provisioning_Idempotency_And_Inventory_Locking.md`.

### 17.1. Transaction boundary đề xuất
Một checkout thành công không đồng nghĩa provider đã cấp xong. Flow phải tách rõ các boundary:
1. Tạo order/reservation trong DB transaction.
2. Debit wallet và ghi ledger trong DB transaction.
3. Tạo provisioning job với idempotency key.
4. Worker/adapter gọi provider bên ngoài transaction DB.
5. Kết quả provider được ghi về DB theo trạng thái an toàn.

Không giữ DB transaction mở trong lúc gọi API provider.

### 17.2. Idempotency key
Mỗi provisioning job phải có:
- `idempotency_key`
- `order_id`
- `reservation_id`
- `source_id`
- `attempt_count`
- `external_request_id` nếu provider trả về
- `external_resource_id` nếu provider đã tạo tài nguyên
- `last_provider_response_summary`
- `retry_safety_level`

### 17.3. Retry safety
Không phải lỗi nào cũng retry.

| Tình huống | Hành động |
|---|---|
| Provider reject rõ ràng trước khi tạo resource | Có thể retry/chọn source khác/refund |
| Timeout trước khi request gửi đi | Có thể retry nếu chắc chắn request chưa tới provider |
| Timeout sau khi request đã gửi | Manual review, không retry mù |
| Provider trả partial success | Manual review |
| Provider tạo resource nhưng credential lưu lỗi | Không tạo resource mới; khôi phục/lưu lại credential nếu có thể |
| Webhook/callback lặp | Idempotent update, không tạo service trùng |

### 17.4. Credential safety
Credential/access info phải được xử lý như secret:
- Encrypt at rest.
- Không log plaintext.
- Không lưu trong audit before/after snapshot dạng thô.
- UI chỉ hiển thị masked by default.
- Reveal credential phải ghi audit action `credential.revealed`.

### 17.5. Service uniqueness
Một order chỉ được có một service chính thức ở trạng thái active/suspended/terminated theo dòng đời của nó. Nếu provider cấp nhầm nhiều tài nguyên do sự cố, tài nguyên phụ phải vào bảng exception/manual review, không attach im lặng vào order.

### 17.6. Reconciliation
Cần có job đối soát định kỳ:
- DB nói service active nhưng provider không thấy resource.
- Provider có resource nhưng DB chưa có service.
- Stock cached khác stock provider.
- Credential missing hoặc external_resource_id missing.

Kết quả reconciliation không tự sửa dữ liệu tài chính; chỉ tạo alert/manual review trừ khi rule đã được khóa rõ.
