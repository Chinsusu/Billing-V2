# Tài liệu 07 - Portal Functional Specification

## 1. Mục tiêu tài liệu
Tài liệu này mô tả:
- Portal Admin
- Portal Reseller
- Portal Client
- Các module chính
- Những action hiển thị theo role và capability

## 2. Nguyên tắc UI/UX phase 1
- Ưu tiên dễ vận hành hơn là quá nhiều tính năng
- Capability không hỗ trợ phải ẩn
- Out-of-stock phải disable/gray-out rõ
- Billing cycle phải hiển thị đúng theo plan/provider
- Chính sách cancel/refund cần hiện rõ ngay ở detail page

## 3. Admin portal
### 3.1. Dashboard
- Revenue
- Profit
- Inventory overview
- Provider status
- Provision failures gần đây
- MRR / ARPU

### 3.2. Provider management
- Danh sách provider
- Trạng thái kết nối
- Last sync
- Recent errors
- Enable/disable provider
- Adapter config

### 3.3. Master catalog
- Danh sách sản phẩm gốc
- Plan/variant
- Mapping source
- Base pricing
- Retail pricing
- Reseller cost

### 3.4. Retail client management
- Danh sách client của Admin
- Wallet
- Orders
- Services
- Manual actions khi cần

### 3.5. Reseller management
- Danh sách reseller
- White-label settings overview
- Wallet
- Status
- Catalog scope
- Financial summary

### 3.6. Payment & ledger
- Top-up requests
- Verification queue
- Wallet adjustments
- Refund actions
- Ledger search

### 3.7. Provisioning queue
- Pending provisioning
- Failed provisioning
- Manual review
- Retry / rollback / cancel

### 3.8. Reports
- Revenue by period
- Profit by period
- Revenue by reseller
- Active services
- Expiry risk
- Provider issues

### 3.9. Audit log
- Login logs
- Price changes
- Refund logs
- Service deletion
- IP changes

## 4. Reseller portal
### 4.1. Dashboard
- Revenue tenant
- Profit tenant
- Active clients
- Active services
- Wallet balance

### 4.2. White-label settings
- Logo
- Brand name
- Theme
- Domain mapping
- Email template
- Telegram/contact

### 4.3. Catalog clone & pricing
- Enable/disable plan trong tenant
- Set selling price
- Hiển thị policy plan

### 4.4. Client management
- Danh sách client
- Wallet của client
- Orders
- Services
- Manual suspend theo quyền

### 4.5. Wallet & transactions
- Balance hiện tại
- Lịch sử top-up
- Lịch sử charge/refund tenant
- Sổ giao dịch cơ bản

### 4.6. Staff management
- Danh sách staff
- Tạo/sửa/xóa staff
- Gán role staff chung
- Ghi chú: nếu phase 1 không có granular permission thì đây là management đơn giản

### 4.7. Reports
- Revenue
- Profit
- Client growth
- Service growth
- Expiring services

### 4.8. Tenant audit
- Login
- Price changes
- Refunds
- Service suspension
- White-label changes

## 5. Client portal
### 5.1. Dashboard
- Wallet balance
- Active services
- Upcoming expiry
- Usage summary

### 5.2. Auth & profile
- Register
- Login
- Change password
- Optional 2FA placeholder cho phase sau

### 5.3. Wallet
- Top-up
- Payment instructions
- Transaction history

### 5.4. Catalog & checkout
- Browse plans
- Xem detail plan
- Xem billing cycle
- Xem cancellation/refund policy
- Checkout bằng số dư ví

### 5.5. Service list
- Danh sách service
- Trạng thái
- Ngày hết hạn
- Hành động theo capability

### 5.6. VPS service detail
- Start
- Stop
- Reboot
- Reinstall OS
- Change password
- VNC/console
- Change IP nếu plan cho phép
- Usage realtime
- Billing info

### 5.7. Proxy service detail
- Endpoint/access info
- Basic status
- Billing info
- Contact support
- Không có nhiều self-service action ở phase 1

### 5.8. Renew & cancel
- Renew thủ công
- Cancel chỉ xuất hiện nếu plan cho phép
- Nếu plan không cho cancel giữa kỳ thì chỉ hiển thị policy tương ứng

## 6. Điều kiện hiển thị theo capability
Ví dụ:
- Nếu `supports_console = false` -> ẩn VNC/Console
- Nếu `allow_mid_cycle_cancel = false` -> ẩn nút cancel
- Nếu `supports_change_ip = false` -> ẩn action đổi IP
- Nếu `stock_state = out_of_stock` -> disable mua mới

## 7. Contact support
Phase 1 không có module ticket nội bộ.
Portal chỉ cần:
- Link Telegram hoặc contact method
- Hiển thị contact theo tenant
- Admin và Reseller tách contact riêng

## 8. Các flow UX cần có wireframe sau
- Register/Login
- Top-up wallet
- Browse catalog
- Checkout
- Service detail VPS
- Service detail Proxy
- Renew
- Cancel/refund eligible flow
- Payment verification admin
- Failed provisioning admin
- White-label config reseller

## 9. Bản vá v1.1 - Functional spec format bắt buộc cho từng flow
Danh sách module chưa đủ để giao dev. Mỗi flow P0 cần được viết theo format dưới đây.

### 9.1. Template cho flow
Mỗi flow phải có:
- Entry point.
- Actor/role được phép.
- Preconditions.
- Main steps.
- Validation.
- Error states.
- Audit event.
- Notification nếu có.
- Acceptance criteria.

### 9.2. Flow checkout client thuộc reseller
Preconditions:
- User đã đăng nhập.
- Tenant active.
- Plan active trong tenant catalog.
- Source còn stock.
- Client wallet đủ selling price.
- Reseller wallet đủ reseller cost.
- Plan không bị margin block.

Main steps:
1. Client chọn plan.
2. UI hiển thị billing cycle, refund policy, capability và giá snapshot.
3. Client xác nhận checkout.
4. Hệ thống reserve stock.
5. Hệ thống debit client wallet và reseller wallet.
6. Hệ thống tạo order paid.
7. Hệ thống tạo provisioning job.
8. UI chuyển sang màn hình order/provisioning status.

Error states:
- Insufficient client balance.
- Insufficient reseller balance.
- Out of stock.
- Plan disabled.
- Provider/source disabled.
- Reservation expired.
- Provisioning failed/manual review.

Audit events:
- order.created
- wallet.client.debited
- wallet.reseller.debited
- reservation.created
- provisioning.job.created

Acceptance criteria:
- Không có provisioning job nếu thiếu một trong hai debit.
- Không trừ tiền hai lần khi user double click.
- User nhìn thấy trạng thái rõ: pending/provisioning/active/failed.

### 9.3. Flow top-up manual
Preconditions:
- User thuộc đúng tenant.
- Payment method enabled.

Main steps:
1. User tạo top-up request.
2. Hệ thống sinh mã tham chiếu.
3. User chuyển khoản/gửi bằng chứng nếu cần.
4. Admin/Reseller có quyền mở verification queue.
5. Người verify approve/reject.
6. Nếu approve, hệ thống credit wallet và ghi ledger.

Acceptance criteria:
- Một top-up chỉ approve được một lần.
- Reject phải có reason.
- Amount/currency/fx snapshot được lưu.

### 9.4. Flow reveal credential
Preconditions:
- User có quyền xem service.
- Service thuộc tenant context.
- Service đã provisioned.

Main steps:
1. UI hiển thị credential masked.
2. User bấm reveal.
3. Nếu role nhạy cảm, yêu cầu re-auth/2FA theo policy.
4. Hệ thống ghi audit.
5. UI hiển thị credential trong thời gian ngắn.

Acceptance criteria:
- Credential không xuất hiện trong page source/log/audit plaintext ngoài cơ chế hiển thị được kiểm soát.
- Reveal cross-tenant bị chặn.

### 9.5. Flow pricing update của reseller
Preconditions:
- Actor là Reseller Owner hoặc staff có quyền pricing.
- Plan thuộc tenant clone.

Validation:
- Selling price phải lớn hơn 0.
- Nếu selling price < reseller cost, hiển thị cảnh báo margin âm.
- Nếu policy block negative margin, không cho save hoặc auto-disable plan.

Audit events:
- pricing.tenant_plan.updated

Acceptance criteria:
- Order cũ giữ price snapshot cũ.
- Storefront hiển thị giá mới chỉ sau khi update thành công.
