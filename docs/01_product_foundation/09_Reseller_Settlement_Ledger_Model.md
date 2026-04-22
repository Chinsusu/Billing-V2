# Tài liệu 09 - Reseller Settlement & Ledger Model

## 1. Mục tiêu tài liệu
Tài liệu này khóa mô hình tiền giữa Platform, Reseller và Client trong phase 1. Đây là phần P0 vì nếu sai settlement, hệ thống có thể cấp tài nguyên thật nhưng không thu được tiền thật.

Tài liệu này chưa mô tả code. Nó mô tả rule nghiệp vụ, ledger, flow và acceptance criteria để đội dev/ops triển khai thống nhất.

## 2. Nguyên tắc lõi
### 2.1. Platform chỉ tin reseller wallet
Client wallet trong tenant reseller là quan hệ tiền giữa client và reseller. Platform không xem đó là tiền đã thanh toán cho platform.

Platform chỉ cấp tài nguyên thật khi:
- Client wallet đủ tiền bán lẻ trong tenant.
- Reseller wallet đủ reseller cost phải trả cho platform.
- Hai debit liên quan được ghi ledger thành công.

### 2.2. Không có ledger thì giao dịch không tồn tại
Mọi thay đổi số dư phải có ledger entry. Không sửa ledger cũ. Nếu sai, tạo adjustment entry.

### 2.3. Snapshot tại thời điểm mua
Order/service phải lưu:
- selling price snapshot
- reseller cost snapshot
- plan policy snapshot
- refund policy snapshot
- billing cycle snapshot
- fx snapshot nếu có

Không dùng giá hiện tại để xử lý order cũ.

## 3. Các loại ví
### 3.1. Platform/Admin retail client wallet
Ví của client mua trực tiếp từ Admin storefront. Khi mua, tiền client trả trực tiếp cho platform.

### 3.2. Reseller wallet
Ví settlement của reseller với platform. Reseller nạp trước, platform cấp hàng sau.

### 3.3. Client wallet trong tenant reseller
Ví của client cuối trong storefront reseller. Reseller có thể xác nhận top-up cho client theo policy tenant, nhưng platform không coi đó là tiền platform đã nhận.

## 4. Các mô hình giao dịch
### 4.1. Admin retail purchase
Khi client thuộc Admin mua plan giá 20 USD:

| Bước | Ví/Tài khoản | Amount | Ý nghĩa |
|---|---:|---:|---|
| 1 | Client wallet | -20 | Client trả tiền mua dịch vụ |
| 2 | Platform revenue ledger | +20 | Doanh thu retail platform |
| 3 | Cost snapshot | theo source | Dùng cho profit report, không nhất thiết là wallet movement |

Rule:
- Nếu client wallet không đủ tiền, không tạo provisioning job.
- Nếu provision fail chắc chắn trước khi tạo tài nguyên, refund 100% về client wallet.

### 4.2. Reseller client purchase
Client A thuộc Reseller R mua VPS.

Ví dụ:
- Selling price: 20 USD.
- Reseller cost: 12 USD.
- Gross profit reseller: 8 USD.

| Bước | Ví/Tài khoản | Amount | Ý nghĩa |
|---|---:|---:|---|
| 1 | Client wallet trong tenant reseller | -20 | Client trả tiền cho reseller |
| 2 | Reseller revenue internal ledger | +20 | Doanh thu nội bộ reseller |
| 3 | Reseller wallet | -12 | Reseller trả platform giá nhập |
| 4 | Platform reseller revenue ledger | +12 | Doanh thu platform từ reseller |
| 5 | Reseller profit report | +8 | Profit tính toán từ ledger/snapshot |

Rule khóa:
- Nếu client wallet đủ 20 nhưng reseller wallet không đủ 12, order không được provision.
- Nếu reseller tự approve top-up ảo cho client nhưng reseller wallet không có tiền thật, platform vẫn không cấp hàng.
- Profit reseller không được tính bằng số dư ví đơn thuần; phải tính từ ledger.

### 4.3. Reseller tự mua cho chính mình
Nếu reseller mua dịch vụ để dùng nội bộ:
- Có thể xem reseller như client đặc biệt của tenant.
- Vẫn phải debit reseller wallet theo reseller cost hoặc theo policy platform.
- Nếu cần tracking profit nội bộ, tạo order với `buyer_type = reseller_internal`.

## 5. Ledger entry types đề xuất
### 5.1. Top-up
- `topup.request.created`
- `topup.credit.client`
- `topup.credit.reseller`
- `topup.rejected`

### 5.2. Purchase/Renew
- `purchase.client_wallet.debit`
- `purchase.reseller_wallet.debit`
- `purchase.platform_revenue.credit`
- `renewal.client_wallet.debit`
- `renewal.reseller_wallet.debit`

### 5.3. Refund/Adjustment
- `refund.client_wallet.credit`
- `refund.reseller_wallet.credit`
- `adjustment.client_wallet.credit`
- `adjustment.client_wallet.debit`
- `adjustment.reseller_wallet.credit`
- `adjustment.reseller_wallet.debit`

### 5.4. Paid action
- `paid_action.client_wallet.debit`
- `paid_action.reseller_wallet.debit`
- `paid_action.refund.credit`

## 6. Ledger fields tối thiểu
Mỗi ledger entry cần có:
- ledger_entry_id
- tenant_id
- wallet_id
- wallet_owner_type: platform_client, reseller, reseller_client
- actor_id hoặc system actor
- direction: credit/debit
- amount
- currency
- fx_snapshot
- entry_type
- reference_type
- reference_id
- order_id nếu có
- service_id nếu có
- idempotency_key
- balance_before
- balance_after
- status: pending, posted, voided_by_adjustment
- reason
- created_at

Rule:
- `posted` ledger không sửa.
- Nếu sai, tạo entry adjustment có reference tới entry gốc.

## 7. Top-up state machine
Top-up request có trạng thái:
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
- Manual approve phải có actor, timestamp, proof/reference.
- Reject phải có reason.
- Expired không được tự credit.

## 8. Checkout flow chuẩn cho reseller client
1. Client chọn plan.
2. Hệ thống lấy tenant context từ domain/session.
3. Hệ thống lấy tenant plan version và price snapshot.
4. Kiểm tra plan active, source active, stock available.
5. Kiểm tra client wallet >= selling price.
6. Kiểm tra reseller wallet >= reseller cost.
7. Tạo reservation atomic.
8. Debit client wallet và ghi ledger.
9. Debit reseller wallet và ghi ledger.
10. Tạo order paid.
11. Tạo provisioning job.
12. Nếu service active, report revenue/profit được tính từ ledger và snapshot.

## 9. Rollback và refund
### 9.1. Reserve fail
Nếu reserve fail trước debit:
- Không trừ ví.
- Order cancelled/out_of_stock.

### 9.2. Debit fail
Nếu một trong hai debit fail:
- Không tạo provisioning job.
- Nếu đã debit một ví nhưng debit ví kia fail, phải reverse bằng refund/adjustment trong cùng control flow.
- Ghi audit.

### 9.3. Provisioning fail chắc chắn trước khi tạo resource
- Release reservation.
- Refund client wallet 100% charge liên quan.
- Refund reseller wallet 100% reseller cost debit liên quan.
- Order/service chuyển failed/refunded theo lifecycle.

### 9.4. Provider partial success
Nếu không chắc provider đã tạo resource hay chưa:
- Không retry mù.
- Chuyển provisioning job sang manual_review.
- Không refund tự động cho tới khi xác minh.

### 9.5. Cancel giữa kỳ
Nếu plan cho phép cancel:
- Client refund = 80% x giá trị thời gian còn lại theo selling price snapshot.
- Reseller wallet refund = 80% x giá trị thời gian còn lại theo reseller cost snapshot nếu platform reclaim được nguồn hoặc policy cho phép.
- Nếu nguồn upstream không refund cho platform, reseller cost refund có thể bằng 0; policy này phải hiển thị rõ cho reseller.

## 10. Công thức tài chính
### 10.1. Reseller gross profit theo order
`reseller_gross_profit = selling_price_snapshot - reseller_cost_snapshot - reseller_refund_impact`

### 10.2. Refund theo ngày
`remaining_days = max(0, floor((term_end - now) theo ngày))`

`daily_value = original_price / total_cycle_days`

`refund_amount = 0.8 x remaining_days x daily_value`

Rounding rule:
- USD: làm tròn 2 chữ số thập phân.
- VND: làm tròn theo cấu hình tenant/payment method.

### 10.3. Renew trong grace
Mặc định:
`new_term_end = old_term_end + cycle`

Không cộng từ `now` trừ khi plan policy ghi rõ.

## 11. Controls bắt buộc
- Không âm ví trừ khi có credit limit rõ ràng.
- Không refund vượt số tiền đã debit theo cùng reference.
- Không debit hai lần cùng `order_id + wallet_id + charge_type`.
- Mọi manual adjustment phải có reason và audit.
- Các report tài chính tính từ ledger posted.
- Không dùng wallet balance cache để tính doanh thu.

## 12. Acceptance criteria
- Client reseller không thể provision nếu reseller wallet thiếu tiền.
- Double click checkout không tạo hai charge.
- Top-up approved hai lần không credit hai lần.
- Refund không vượt charge gốc.
- Order cũ vẫn dùng price snapshot cũ sau khi Admin đổi giá.
- Report reseller profit khớp selling price trừ reseller cost theo ledger.
- Mọi adjustment có reason và audit.
