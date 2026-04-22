# Tài liệu 06 - Order & Service Lifecycle State Machine

## 1. Mục tiêu tài liệu
Tài liệu này tách rõ:
- Order lifecycle
- Reservation lifecycle
- Provisioning lifecycle
- Service lifecycle
- Billing status
- Suspension reason

## 2. Nguyên tắc chung
Không nên dùng một cột status duy nhất để gánh mọi ý nghĩa.
Nên tách ít nhất:
- order_status
- reservation_status
- provisioning_status
- service_status
- billing_status
- suspension_reason

## 3. Order lifecycle
### Trạng thái đề xuất
- draft
- pending_payment
- paid
- cancelled
- failed
- refunded

### Chuyển trạng thái chính
- draft -> pending_payment
- pending_payment -> paid
- pending_payment -> cancelled
- paid -> refunded
- paid -> failed (nếu rollback cần tách logic rõ)

## 4. Reservation lifecycle
- pending_reserve
- reserved
- reservation_expired
- reservation_released
- allocated

Rule:
- reserved chỉ tồn tại tối đa 5 phút nếu chưa đi tiếp
- allocated nghĩa là source đã gắn cho provisioning/order

## 5. Provisioning lifecycle
- queued
- provisioning
- provisioned
- failed
- manual_review

## 6. Service lifecycle
- active
- suspended
- expired
- cancelled
- terminated

### Diễn giải
- active: dịch vụ đang chạy/hợp lệ
- suspended: dịch vụ bị ngưng nhưng chưa xóa
- expired: đã hết kỳ sử dụng
- cancelled: đã được hủy theo rule
- terminated: đã xóa hẳn, không khôi phục theo nghiệp vụ chuẩn

## 7. Billing status
- unpaid
- paid
- overdue
- refunded
- partially_refunded

## 8. Suspension reason
- expiry
- manual_admin
- manual_reseller
- abuse
- system_issue

Dù hiện tại user nói suspended vì hết tiền và abuse không khác nhau, hệ thống vẫn nên lưu reason để phục vụ audit và report.

## 9. Luồng mua hàng chuẩn
1. draft
2. pending_payment
3. paid
4. pending_reserve
5. reserved
6. queued
7. provisioning
8. provisioned
9. service active

## 10. Luồng gia hạn
- active -> active với term mới
- nếu service đã suspended do expiry nhưng vẫn trong grace:
  - renew thành công -> active
- nếu đã terminated:
  - mặc định không renew lại, phải tạo order mới

## 11. Luồng expiry
1. term_end reached
2. service_status = suspended
3. billing_status = overdue
4. grace countdown = 3 ngày
5. quá grace -> terminated

## 12. Luồng cancel
### 12.1. Plan cho phép cancel giữa kỳ
- active -> cancelled
- tạo refund entry
- sau bước cleanup -> terminated nếu cần

### 12.2. Plan không cho cancel giữa kỳ
- Client không có action cancel ngay
- Hệ thống chỉ cho “không gia hạn tiếp” nếu có khái niệm đó ở UI

## 13. Luồng provisioning fail
- queued -> provisioning -> failed
- hoặc queued -> manual_review
- Tùy hướng xử lý:
  - retry -> queued
  - refund -> refunded/cancelled
  - terminate flow

## 14. Notifications gợi ý
- order_created
- payment_confirmed
- provisioning_started
- provisioning_failed
- service_activated
- expiry_warning
- service_suspended
- service_terminated
- refund_completed

## 15. Cảnh báo nghiệp vụ
- Không để service_status = active khi billing_status = overdue ngoài policy cho phép
- Không để order refunded nhưng service vẫn active
- Không để terminated service có action operational tiếp tục xuất hiện trên UI

## 16. Bảng mapping ngắn
| Tình huống | order_status | provisioning_status | service_status | billing_status |
|---|---|---|---|---|
| Chưa thanh toán | pending_payment | n/a | n/a | unpaid |
| Đã thanh toán, chưa cấp | paid | queued | n/a | paid |
| Đang cấp | paid | provisioning | n/a | paid |
| Đã cấp xong | paid | provisioned | active | paid |
| Hết hạn trong grace | paid | provisioned | suspended | overdue |
| Đã refund một phần | refunded hoặc paid | provisioned/terminated | cancelled/terminated | partially_refunded |

## 17. Bản vá v1.1 - Transition guard và term calculation
State machine phải có guard để tránh trạng thái đẹp trên giấy nhưng sai tiền/sai tồn kho khi chạy thật.

### 17.1. Guard khi checkout
`pending_payment -> paid` chỉ được phép khi:
- Client wallet debit thành công hoặc payment verified.
- Nếu thuộc reseller: reseller wallet settlement debit thành công.
- Ledger entries đã ghi đầy đủ.
- Reservation còn hiệu lực hoặc reserve được tạo trong cùng transaction.

### 17.2. Guard khi provisioning
`queued -> provisioning` chỉ được phép khi:
- Reservation status = reserved hoặc allocated theo rule.
- Order billing status = paid.
- Idempotency key tồn tại.
- Source/provider đang enabled.

`provisioning -> provisioned` chỉ được phép khi:
- Có external_resource_id hoặc bằng chứng provider đã cấp tài nguyên.
- Credential/access info lưu thành công hoặc có manual secure handoff.
- Service instance tạo thành công.
- Reservation chuyển allocated.

### 17.3. Guard khi refund/cancel
Không cho `paid -> refunded` nếu:
- Service vẫn active và chưa có cancel/terminate action hợp lệ.
- Refund amount vượt allowed amount.
- Actor không có quyền.
- Không có reason.

### 17.4. Renewal calculation
Rule mặc định:
- Service active: `new_term_end = old_term_end + cycle`.
- Service suspended trong grace: `new_term_end = old_term_end + cycle`, sau đó status về active nếu provision/source vẫn còn hợp lệ.
- Service terminated: không renew; tạo order mới.

### 17.5. Calendar month calculation
Với `calendar_month`, cộng tháng theo lịch và clamp ngày cuối tháng.
Với `month_30d`, cộng đúng 30 ngày.

### 17.6. Expiry job idempotency
Expiry/suspend/terminate job phải idempotent:
- Chạy lại nhiều lần không suspend/terminate trùng.
- Không terminate trước khi grace_end.
- Nếu user renew giữa lúc job chạy, job phải đọc lại state mới nhất trước khi terminate.

### 17.7. Invariants bắt buộc
- Không có service active khi order bị refunded toàn phần.
- Không có provisioning job running cho order cancelled/refunded, trừ job cleanup có type riêng.
- Không có reserved inventory quá expiry reservation.
- Không hiện operational action cho service terminated.
- Không cho client tự thao tác service đang manual_review.
