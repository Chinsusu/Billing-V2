# Tài liệu 08 - Audit, Reports & Operational Control

## 1. Mục tiêu tài liệu
Tài liệu này định nghĩa:
- Audit log
- Báo cáo
- Chỉ số vận hành
- Control panel cho Admin/Reseller
- SOP tối thiểu

## 2. Audit log scope
Các hành động nhạy cảm cần audit:
- Login
- Đổi giá
- Xóa service
- Đổi IP
- Refund
- Payment verification
- Manual provisioning action
- White-label config change
- Domain mapping change
- Staff management change

## 3. Trường dữ liệu audit tối thiểu
- audit_id
- actor_type
- actor_id
- tenant_id
- resource_type
- resource_id
- action
- before_snapshot hoặc before_summary
- after_snapshot hoặc after_summary
- timestamp
- source_ip / context

## 4. Retention
Rule đã chốt:
- Audit log lưu không giới hạn

Đề xuất:
- Dù retention logic là vô hạn, vẫn cần chiến lược archive hoặc cold storage khi dữ liệu lớn.

## 5. Báo cáo cho Admin
### 5.1. Kinh doanh
- Revenue theo ngày/tuần/tháng
- Gross profit
- Revenue retail
- Revenue theo reseller
- MRR
- ARPU

### 5.2. Vận hành
- Provision success/fail
- Provider status
- Out-of-stock by source/plan
- Expiring services
- Refund volume

## 6. Báo cáo cho Reseller
- Revenue tenant
- Profit tenant
- Active clients
- Active services
- Wallet balance
- Expiring services
- Top-selling plans

## 7. View cho Client
- Wallet balance
- Active services
- Upcoming expiry
- Usage summary
- Transaction history

## 8. Operational controls cho Admin
- Enable/disable provider
- Enable/disable plan
- Update exchange rate
- Verify payment
- Manual wallet adjustment
- Retry provisioning
- Cancel failed order
- Process refund
- Force suspend/terminate khi cần

## 9. Operational controls cho Reseller
- Enable/disable plan trong tenant
- Update selling price
- Suspend service của client thuộc tenant
- Quản lý staff
- Quản lý white-label config
- Quản lý contact/support info

## 10. SOP tối thiểu
### 10.1. SOP xác nhận thanh toán manual
1. Mở payment verification queue
2. Đối chiếu mã tham chiếu
3. Kiểm tra số tiền và currency
4. Approve/Reject
5. Ghi ledger và audit

### 10.2. SOP provisioning fail
1. Xác định source/provider lỗi
2. Xem payload và external response
3. Retry nếu phù hợp
4. Nếu không retry được thì hủy hoặc refund
5. Ghi audit đầy đủ

### 10.3. SOP out-of-stock
1. Disable/gray-out plan hoặc source
2. Kiểm tra mapping source
3. Cập nhật availability
4. Thông báo nội bộ nếu cần

### 10.4. SOP cancel/refund
1. Kiểm tra plan có cho cancel giữa kỳ không
2. Tính số ngày còn lại
3. Tính refund 80%
4. Trả tiền vào ví
5. Cleanup service theo policy
6. Ghi ledger và audit

### 10.5. SOP abuse thủ công
Phase 1 chưa có anti-fraud tự động, nên abuse xử lý thủ công:
- Suspend service nếu cần
- Terminate nếu vi phạm nặng
- Ghi lý do ở audit/log

## 11. Provider health monitoring
### 11.1. Phase 1
- Connected/Disconnected
- Last sync
- Recent error
- Recent provision failures

### 11.2. Phase 2
- Health score tổng hợp
- Error rate
- Latency
- Freshness inventory

## 12. Những điểm cần chuẩn hóa sớm
- Naming chuẩn cho audit action
- Khoảng thời gian mặc định cho report
- Logic tính profit
- Cách tách retail Admin với doanh thu reseller trong dashboard tổng
- Quy tắc export CSV/XLSX cho báo cáo ở phase sau

## 13. Bản vá v1.1 - Audit naming, redaction và traceability
Audit không chỉ để xem lịch sử. Audit là công cụ cứu tranh chấp, cứu tiền và cứu dữ liệu khi có sự cố.

### 13.1. Trường audit bổ sung
Ngoài các trường tối thiểu đã có, audit nên thêm:
- request_id
- correlation_id
- actor_role
- tenant_context_source
- user_agent
- reason
- risk_level
- redaction_version

`correlation_id` phải đi xuyên từ order -> wallet -> ledger -> reservation -> provisioning -> service.

### 13.2. Redaction rule
Không ghi plaintext các trường sau vào audit/log:
- password
- API key
- provider token
- VPS password
- proxy username/password nếu nhạy cảm
- private key
- payment secret

Trường nhạy cảm chỉ được lưu dạng masked hoặc hash/reference.

### 13.3. Naming chuẩn cho audit action
Action nên dùng format `domain.object.action`.

Ví dụ P0:
- auth.login.success
- auth.login.failed
- tenant.emergency.read
- tenant.emergency.write
- wallet.topup.created
- wallet.topup.approved
- wallet.adjustment.created
- wallet.client.debited
- wallet.reseller.debited
- ledger.entry.created
- order.created
- order.paid
- order.refunded
- reservation.created
- reservation.expired
- provisioning.job.created
- provisioning.job.failed
- provisioning.job.manual_review
- service.activated
- service.suspended
- service.terminated
- pricing.master_plan.updated
- pricing.tenant_plan.updated
- credential.revealed
- abuse.flag.created
- abuse.service.suspended

### 13.4. Report formula khóa
Admin report:
- Platform retail revenue = tổng purchase/renewal của client thuộc Admin tenant.
- Platform reseller revenue = tổng reseller cost debit từ reseller wallet.
- Platform gross revenue = retail revenue + reseller revenue.
- Refund = tổng refund credit theo reference.
- Net revenue = gross revenue - refund.

Reseller report:
- Reseller revenue = tổng charge từ client wallet theo selling price.
- Reseller cost = tổng debit từ reseller wallet theo reseller cost.
- Reseller gross profit = reseller revenue - reseller cost - reseller refund impact.

Rule: report tài chính phải dựa trên ledger entries đã posted, không dựa vào order status đơn lẻ.

### 13.5. SOP bổ sung khi khách báo “bị trừ tiền nhưng chưa có service”
1. Tìm theo correlation_id hoặc order_id.
2. Kiểm tra ledger debit của client.
3. Nếu reseller flow, kiểm tra reseller wallet debit.
4. Kiểm tra reservation status.
5. Kiểm tra provisioning job status và provider response.
6. Nếu provider timeout/partial success, chuyển manual review.
7. Nếu fail chắc chắn, refund theo policy.
8. Ghi audit kết luận và action xử lý.
