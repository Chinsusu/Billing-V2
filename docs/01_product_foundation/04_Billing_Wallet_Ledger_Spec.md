# Tài liệu 04 - Billing, Wallet & Ledger Specification

## 1. Mục tiêu tài liệu
Tài liệu này định nghĩa:
- Ví nội bộ
- Nạp tiền
- Thanh toán
- Gia hạn
- Expiry
- Refund
- Tỷ giá
- Ledger

## 2. Nguyên tắc tài chính nền tảng
- Wallet-first
- Nạp tiền trước rồi mới mua
- Mọi thay đổi số dư phải có ledger
- USD là đơn vị chính
- VND là đơn vị quy đổi theo tỷ giá nhập tay

## 3. Mô hình ví
### 3.1. Ví của Reseller
- Dùng để mua nguồn hàng và bán lại
- Nạp trước bán sau

### 3.2. Ví của Client
- Dùng để mua, gia hạn, chi trả các action trả phí như change IP

### 3.3. Các số dư cần có
- Available balance
- Locked balance
- Pending top-up amount nếu cần

## 4. Phương thức nạp tiền
- Chuyển khoản tay
- PayPal
- USDT manual
- USDT auto

## 5. Luồng top-up
1. User tạo yêu cầu nạp tiền
2. Hệ thống sinh mã tham chiếu hoặc hướng dẫn thanh toán
3. Thanh toán được xác nhận tự động hoặc thủ công
4. Số dư ví được cộng
5. Ledger ghi nhận giao dịch
6. Audit ghi nhận actor và phương thức

## 6. Luồng mua hàng
1. User chọn plan
2. Hệ thống kiểm tra stock và reserve tối đa 5 phút
3. Hệ thống kiểm tra available balance
4. Trừ ví
5. Tạo order ở trạng thái paid
6. Chuyển sang queue provisioning

## 7. Luồng gia hạn
- Auto-renew mặc định tắt
- Người dùng hoặc seller chủ động renew
- Khi renew:
  - Kiểm tra policy plan
  - Kiểm tra số dư
  - Trừ ví
  - Cập nhật service term
  - Ghi ledger

## 8. Billing cycle
- Day
- Week
- Month

Rule:
- “Month” phải được lưu đúng theo plan/provider
- Không dùng duy nhất một công thức 30 ngày cho mọi plan

## 9. Expiry & grace
- Hết hạn -> suspend ngay
- Grace period 3 ngày
- Hết grace -> terminate/xóa service

Đề xuất:
- Nên có ít nhất 2 mốc thông báo:
  - Sắp hết hạn
  - Đã suspend
  - Sắp xóa hẳn

## 10. Refund
### 10.1. Điều kiện
Chỉ áp dụng cho plan có cờ cho phép cancel giữa kỳ.

### 10.2. Công thức
Refund = 80% x giá trị phần thời gian còn lại theo ngày

### 10.3. Ghi nhận
- Trả về ví
- Ghi ledger loại refund
- Ghi audit actor
- Có lý do refund

## 11. Ledger model
Các loại bút toán tối thiểu:
- topup_credit
- manual_adjustment_credit
- manual_adjustment_debit
- purchase_charge
- renewal_charge
- refund_credit
- paid_action_charge
- currency_conversion_snapshot

## 12. Tỷ giá USD/VND
- USD là đơn vị chuẩn
- VND dùng để hiển thị và giao dịch theo cấu hình
- Tỷ giá do Admin cập nhật tay

Đề xuất:
- Lưu tỷ giá snapshot tại thời điểm tạo giao dịch
- Không recalc giao dịch cũ khi tỷ giá thay đổi

## 13. Invoice & VAT
- Khách có thể yêu cầu hóa đơn
- Phase 1 xử lý thủ công
- Chưa có API hóa đơn tự động

Đề xuất trạng thái:
- not_requested
- requested
- processing
- issued
- rejected

## 14. Payment verification
### 14.1. Manual bank transfer
- Cần màn hình xác nhận
- Có trạng thái chờ xác nhận
- Cần mã đối chiếu / tham chiếu giao dịch

### 14.2. PayPal / USDT
- Có thể auto hoặc semi-auto tùy mức tích hợp thực tế
- Nếu không chắc chắn, phải có luồng manual fallback

## 15. Những rule cần chốt ở DB/API
- Không được trừ ví 2 lần cho cùng một order
- Reservation hết hạn phải release đúng
- Refund không được vượt quá mức cho phép
- Ledger phải immutable ở mức nghiệp vụ; nếu sửa phải dùng adjustment entry
- Mọi transaction tài chính phải có actor, time, amount, currency, fx_snapshot

## 16. Chỉ số cần theo dõi
- Tổng top-up
- Tỷ lệ payment verified
- Tổng charge
- Tổng refund
- Wallet balance theo tenant
- Nợ grace period
- Tỷ lệ renew

## 17. Bản vá v1.1 - Settlement và ledger guard
File này cần được đọc cùng `09_Reseller_Settlement_Ledger_Model.md`. Các rule dưới đây là bản vá bắt buộc.

### 17.1. Hai lớp ví trong mô hình reseller
Trong tenant reseller luôn tồn tại hai lớp tiền:
- **Client wallet**: tiền nội bộ của client với reseller.
- **Reseller wallet**: tiền thật/reserved settlement giữa reseller và platform.

Khi client của reseller mua hàng, không được chỉ trừ client wallet. Platform chỉ provision nếu reseller wallet đủ reseller cost.

### 17.2. Flow mua hàng reseller-client
Luồng chuẩn:
1. Client chọn plan trong storefront reseller.
2. Hệ thống kiểm tra plan active, stock available, capability/policy snapshot.
3. Hệ thống kiểm tra client wallet >= selling price.
4. Hệ thống kiểm tra reseller wallet >= reseller cost.
5. Reserve inventory atomic.
6. Debit client wallet theo selling price.
7. Debit reseller wallet theo reseller cost.
8. Tạo ledger entries tương ứng.
9. Tạo provisioning job.
10. Nếu provisioning fail chắc chắn, rollback/refund theo rule.

### 17.3. Ledger entry không được đứng một mình
Mỗi thay đổi số dư phải có:
- wallet_id
- tenant_id
- actor_id hoặc system actor
- amount
- currency
- direction
- entry_type
- reference_type/reference_id
- idempotency_key
- balance_before
- balance_after
- created_at
- fx_snapshot nếu có

Không cho ghi số dư mới nếu không ghi ledger thành công.

### 17.4. Top-up state machine
Top-up request nên có trạng thái:
- created
- pending_payment
- pending_verification
- approved
- rejected
- expired
- cancelled

Rule:
- Chỉ `approved` mới credit wallet.
- Một top-up request chỉ được approve một lần.
- Reject không được xóa request; phải lưu reason.
- Manual approve phải audit actor.

### 17.5. Refund theo reseller
Refund cho client của reseller có hai tầng:
- Client nhận refund theo selling price snapshot và refund policy snapshot.
- Reseller nhận hoặc không nhận adjustment ở reseller wallet tùy policy platform đối với reseller cost.

Rule mặc định đề xuất:
- Nếu provisioning fail trước khi tạo tài nguyên thật: hoàn 100% client charge và hoàn 100% reseller cost debit.
- Nếu cancel giữa kỳ theo policy: client refund = 80% phần thời gian còn lại theo selling price; reseller cost refund = 80% phần thời gian còn lại theo reseller cost nếu nguồn platform cũng được reclaim.
- Nếu abuse/vi phạm AUP: refund có thể bằng 0 theo policy và phải có reason/audit.

### 17.6. Renew term rule
Khi renew service còn active:
- `new_term_end = old_term_end + cycle`.

Khi renew trong grace:
- Mặc định vẫn cộng từ `old_term_end`, không cộng từ `now`, để tránh khách nhận thêm ngày miễn phí ngoài policy.
- Nếu muốn chính sách thân thiện hơn, phải cấu hình rõ ở plan policy.

### 17.7. Financial invariants
- Không âm ví trừ khi Admin bật credit line có limit rõ ràng.
- Không trừ ví hai lần cho cùng `order_id + charge_type`.
- Không refund vượt quá số tiền đã charge theo cùng reference.
- Không sửa ledger cũ; mọi sửa sai dùng adjustment entry.
- Mọi report tài chính phải tính từ ledger, không tính từ số dư cache.
