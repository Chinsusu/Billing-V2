# Tài liệu 03 - Product Catalog & Pricing Rules

## 1. Mục tiêu tài liệu
Tài liệu này định nghĩa:
- Mô hình product/plan/source
- Catalog VPS/Proxy
- Capability theo plan/provider
- Các lớp giá
- Chính sách cancel/refund/change IP

## 2. Nguyên tắc mô hình sản phẩm
Hệ thống nên tách 3 lớp:
1. Product hiển thị trên storefront
2. Plan/variant kinh doanh
3. Source thực tế từ provider

Cách tách này giúp:
- Giao diện bán hàng nhất quán
- Dễ thay đổi nguồn ở backend
- Dễ quản lý capability và tồn kho

## 3. Danh mục VPS
### 3.1. Fixed plan
Ví dụ:
- VPS 2C/4G/80G
- VPS 4C/8G/160G

### 3.2. Custom plan cho Proxmox
Có 2 mode:
- Khách tự kéo CPU/RAM/Disk trong khoảng min-max
- Admin định nghĩa vài option mở rộng sẵn

Đề xuất:
- Hai mode nên cùng map về một loại plan “customizable”
- UI chỉ khác ở cách chọn cấu hình

### 3.3. Hệ điều hành
- Linux
- Windows

### 3.4. Addon
Phase 1 không làm:
- Backup
- Snapshot
- Anti-DDoS
- License
- Managed support

## 4. Danh mục Proxy
### 4.1. Loại proxy
- Datacenter
- ISP
- Residential
- Mobile

### 4.2. Tính chất
- IPv4 / IPv6
- Static / Dynamic

### 4.3. Đơn vị bán
- Theo IP
- Theo bandwidth
- Theo thời gian

## 5. Product model đề xuất
### 5.1. Master product
Ví dụ:
- VPS Singapore
- Proxy Datacenter US
- Proxy Residential Global

### 5.2. Plan
Ví dụ:
- VPS 2C4G monthly
- VPS 2C4G daily
- Proxy DC 10 IP monthly
- Proxy Residential 50 GB monthly

### 5.3. Source mapping
Ví dụ:
- Master plan “VPS 2C4G SG” có thể map tới:
  - Proxmox local SG
  - OVH SG equivalent
  - Hetzner SG equivalent
- Nhưng phase 1 không auto failover

## 6. Capability flags
### 6.1. VPS capability
- start
- stop
- reboot
- reinstall
- change password
- VNC/console
- change IP
- realtime usage

### 6.2. Proxy capability
Phase 1 chủ yếu:
- view access info
- xem basic status
- một số plan có change IP trả phí nếu provider hỗ trợ

### 6.3. Rule hiển thị
- Capability không hỗ trợ sẽ ẩn khỏi UI
- Capability nên được quyết định ở cấp source/provider và tổng hợp lên plan

## 7. Billing cycle ở cấp plan
- Day
- Week
- Month

Rule:
- Không ép chuẩn hóa “month”
- Mỗi plan phải mang theo cycle_type và cycle_definition
- Ví dụ:
  - `month_30d`
  - `calendar_month`
  - `day_1`
  - `week_1`

## 8. Pricing layers
### 8.1. Giá gốc của Admin
Là giá nền để tính toán margin và giá nhập.

### 8.2. Giá nhập của Reseller
Có thể theo:
- Fixed amount
- Discount %
- Bảng giá riêng từng plan

### 8.3. Giá bán của Reseller
Reseller được tự set trong tenant.

### 8.4. Giá bán lẻ trực tiếp của Admin
Áp dụng cho client mua trực tiếp từ Admin storefront.

## 9. Rule custom pricing cho Proxmox
Đề xuất:
- Định nghĩa giá đơn vị cho CPU/RAM/Disk/Bandwidth nếu cần
- Hoặc dùng package tier nội bộ để tránh quá nhiều tổ hợp
- Cần có ngưỡng min/max
- Cần có bước validate tài nguyên khả dụng trước khi thanh toán

## 10. Change IP policy
### 10.1. VPS
- Chỉ một số plan Proxmox hỗ trợ
- Mỗi lần đổi IP có phí
- Cần định nghĩa số lần tối đa hoặc policy manual review nếu cần

### 10.2. Proxy
- Tùy provider/plan
- Có thể là manual, semi-auto hoặc auto
- Phase 1 nên mô hình hóa như một paid action có thể bị tắt hoàn toàn

## 11. Cancel & refund policy
### 11.1. Các plan được phép cancel
- Một số plan self-host từ Proxmox
- Một số plan self-host từ proxy-manager

### 11.2. Công thức
Refund = 80% x số ngày còn lại x đơn giá ngày

### 11.3. Các plan không được cancel
- Tất cả các plan không có cờ `allow_mid_cycle_cancel`

### 11.4. Cách hiển thị ở storefront
Mỗi plan nên có:
- Can cancel mid-cycle: Yes/No
- Refund policy summary
- Change IP fee nếu có

## 12. Out-of-stock và availability
- Hết hàng thì disable hoặc gray-out
- Không auto đổi sang provider khác
- Không auto thay đổi capability theo hướng gây bất ngờ cho khách

## 13. Những trường dữ liệu business nên có
### Product
- product_id
- seller_scope
- type (vps/proxy)
- family
- location
- display_name
- visibility

### Plan
- plan_id
- product_id
- cycle_type
- cycle_definition
- os_support
- pricing_mode
- allow_mid_cycle_cancel
- refund_policy
- change_ip_allowed
- change_ip_fee
- capability_flags

### Source mapping
- source_id
- provider_id
- plan_id
- inventory_key
- capacity
- stock_state
- capability_override

## 14. Những điểm cần giữ nhất quán xuyên toàn hệ thống
- Plan là đơn vị bán
- Source là đơn vị cấp hàng
- Chính sách cancel/refund/change IP phải nằm ở cấp plan
- Capability phải có thể bị override theo source

## 15. Bản vá v1.1 - Versioning, snapshot và margin guard
Catalog không chỉ là bảng hiển thị. Catalog là nguồn của tiền, quyền sử dụng, policy và tranh chấp sau này. Vì vậy mọi thay đổi quan trọng phải có version hoặc snapshot.

### 15.1. Master plan versioning
Khi Admin sửa các trường ảnh hưởng tới tiền hoặc quyền sử dụng, hệ thống nên tạo version mới thay vì ghi đè mù:
- reseller cost
- retail price
- billing cycle
- refund policy
- change IP fee
- capability flags
- source mapping chính
- visibility/availability policy

Order/service đã mua phải giữ snapshot của version cũ.

### 15.2. Tenant catalog clone propagation
Khi master plan thay đổi, tenant clone không nên tự bị thay đổi im lặng. Rule đề xuất:
- Master plan thay đổi tạo version mới.
- Tenant clone giữ snapshot cũ cho tới khi sync/accept version mới.
- Nếu version mới làm `selling_price < reseller_cost`, hệ thống cảnh báo margin âm.
- Có thể auto-disable plan trong tenant nếu policy `block_negative_margin = true`.

### 15.3. Snapshot bắt buộc tại thời điểm mua
Mỗi order/service nên lưu các snapshot sau:
- product snapshot
- plan snapshot
- price snapshot
- reseller cost snapshot nếu qua reseller
- billing cycle snapshot
- refund policy snapshot
- capability snapshot
- source/provider snapshot
- fx snapshot nếu giao dịch hiển thị/thu bằng VND

Rule: xử lý tranh chấp, refund, renew, audit order cũ phải dựa trên snapshot lúc mua; không dùng giá/policy hiện tại để áp ngược.

### 15.4. Margin guard
Reseller được tự set giá nhưng không được gây rủi ro platform.

Rule đề xuất:
- Platform không cấp hàng nếu reseller wallet không đủ reseller cost.
- Nếu selling price thấp hơn reseller cost, UI phải cảnh báo rõ.
- Nếu policy yêu cầu, plan bị disable cho tới khi reseller tăng giá hoặc Admin cho phép bán lỗ có kiểm soát.
- Dashboard reseller phải hiển thị gross margin theo plan.

### 15.5. Billing cycle cụ thể hơn
Với `calendar_month`, dùng logic cộng tháng theo lịch và clamp ngày cuối tháng.
Ví dụ:
- 2026-01-31 + 1 calendar month = 2026-02-28.
- 2026-03-31 + 1 calendar month = 2026-04-30.

Với `month_30d`, term_end = start_at + 30 ngày.
