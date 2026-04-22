# Tài liệu 01 - Product Scope & Business Model

## 1. Mục tiêu tài liệu
Tài liệu này là nguồn sự thật ở mức business cho toàn bộ dự án. Mục tiêu là chốt:
- Phạm vi sản phẩm
- Mô hình kinh doanh
- Đối tượng người dùng
- Nguồn cung
- Kênh bán hàng
- Ranh giới phase 1

## 2. Executive summary
Hệ thống là một nền tảng bán và cho thuê VPS/Proxy theo mô hình hybrid multi-provider, multi-tenant, hỗ trợ:
- Admin bán lẻ trực tiếp
- Reseller bán lại qua storefront white-label riêng
- Client mua, gia hạn và quản lý dịch vụ qua portal

Nền tảng kết hợp:
- Nguồn tự sở hữu: Proxmox, proxy-manager
- Nguồn upstream qua API: OVH, Hetzner, Smarthost, proxy-cheap và các nguồn bổ sung sau này

## 3. Bài toán kinh doanh cần giải quyết
Dự án cần giải quyết đồng thời 5 bài toán:
1. Bán nhiều loại VPS/Proxy từ nhiều nguồn khác nhau nhưng cho người mua một trải nghiệm thống nhất.
2. Chuẩn hóa provisioning và inventory cho nhiều provider có API và rule khác nhau.
3. Tạo mô hình reseller white-label để mở rộng kênh bán hàng.
4. Chuẩn hóa wallet, billing cycle, expiry, renew, refund và ledger.
5. Tách rõ boundary dữ liệu giữa Admin, Reseller và Client để tránh lẫn tenant.

## 4. Mục tiêu sản phẩm
### 4.1. Mục tiêu ngắn hạn
- Có thể bán thật và vận hành thật ở phase 1
- Hỗ trợ nhiều provider thông qua adapter layer chuẩn
- Có storefront riêng cho reseller
- Có dashboard và audit cơ bản cho vận hành

### 4.2. Mục tiêu trung hạn
- Giảm thao tác thủ công của Admin
- Tăng tỷ lệ provisioning thành công
- Cho phép mở rộng thêm provider/plan mà không phải đập lại business model

### 4.3. Mục tiêu dài hạn
- Trở thành nền tảng phân phối VPS/Proxy có khả năng white-label mạnh
- Hỗ trợ nhiều tenant, nhiều thương hiệu, nhiều region trên một codebase

## 5. Nhóm người dùng
### 5.1. Admin
Admin là chủ sở hữu platform. Admin:
- Quản lý provider, source, inventory
- Quản lý catalog gốc
- Quản lý giá gốc, giá nhập reseller, tỷ giá, policy
- Xác nhận thanh toán thủ công khi cần
- Xử lý provisioning fail thủ công
- Bán lẻ trực tiếp cho client của Admin

### 5.2. Reseller
Reseller là tenant bán lại theo mô hình white-label. Reseller:
- Có domain, login page, logo, theme, email template riêng
- Có ví riêng, nạp trước bán sau
- Clone catalog từ Admin
- Tự set giá
- Tự quản lý client, service, report và staff trong tenant của mình

### 5.3. Client
Client là người mua cuối. Client:
- Thuộc đúng một seller: Admin hoặc một Reseller
- Đăng ký tự do
- Đăng nhập bằng email
- Không có sub-user
- Nạp tiền trước rồi mới mua
- Quản lý các action dịch vụ tùy capability của provider/plan

## 6. Mô hình kinh doanh
### 6.1. Hybrid supply model
#### Nguồn tự sở hữu
- VPS self-host qua Proxmox
- Proxy self-host qua proxy-manager

#### Nguồn upstream
- OVH
- Hetzner
- Smarthost
- proxy-cheap
- Các provider sẽ bổ sung sau

### 6.2. Kênh bán hàng
- Admin retail storefront
- Reseller storefront white-label

### 6.3. Quy tắc ownership
- Một client chỉ thuộc đúng một seller
- Không có sub-reseller trong phase 1
- Reseller có thể có staff nội bộ, hiện tại dùng một role staff chung

## 7. Danh mục sản phẩm
### 7.1. VPS
- Fixed plan là chính
- Một số gói Proxmox hỗ trợ custom
- Có Linux và Windows
- Chu kỳ ngày hoặc tháng tùy plan/provider
- Không có addon như backup, snapshot, anti-DDoS, license, managed support trong phase 1

### 7.2. Proxy
- Datacenter
- ISP
- Residential
- Mobile
- IPv4 / IPv6
- Static / Dynamic
- Bán theo IP, bandwidth hoặc thời gian
- Giai đoạn đầu thiên về tháng hoặc ngày tùy provider

## 8. Các quyết định sản phẩm đã khóa
### 8.1. Billing cycle
Hệ thống hỗ trợ:
- Day
- Week
- Month

Rule “month” phải hiển thị đúng theo từng provider/plan:
- Có plan 30 ngày
- Có plan theo calendar month

### 8.2. Auto renew
- Mặc định tắt

### 8.3. Expiry
- Hết hạn suspend ngay
- Grace period 3 ngày
- Sau 3 ngày thì xóa service

### 8.4. Cancel và refund
- Chỉ một số plan self-host từ Proxmox và proxy-manager được cancel giữa kỳ
- Refund = 80% giá trị phần thời gian còn lại tính theo ngày
- Các plan khác không cho cancel giữa kỳ
- Cancellation policy phải gắn theo product/plan

### 8.5. Change IP
- Chỉ một số VPS Proxmox được đổi IP
- Proxy tùy provider/plan
- Mỗi lần đổi IP có tính phí

### 8.6. Capability masking
- Provider nào không hỗ trợ action nào thì action đó ẩn khỏi UI

## 9. Mô hình storefront
### 9.1. Admin storefront
Admin có storefront bán lẻ riêng cho client của mình.

### 9.2. Reseller storefront
Mỗi reseller có thể có:
- Domain riêng
- Brand riêng
- Theme riêng
- Login riêng
- Email template riêng
- Contact/Telegram riêng
- Bảng giá riêng

## 10. Phase 1
### 10.1. Trong phạm vi
- Auth theo tenant
- Registration
- Wallet/top-up
- Manual payment verification
- Product catalog
- Order flow
- Reserve stock 5 phút
- Provisioning flow
- Service management cơ bản
- White-label reseller
- Audit log
- Dashboard cơ bản
- Exchange rate thủ công
- Invoice request thủ công
- Contact support qua Telegram/link

### 10.2. Ngoài phạm vi
- Internal ticketing đầy đủ
- KYC
- Anti-fraud automation
- Sub-reseller
- Client sub-user
- Public API cho client
- Mobile app
- E-invoice tự động

## 11. Ràng buộc nghiệp vụ chính
- Mỗi provider có billing rule khác nhau
- Không phải provider nào cũng có cùng action
- Không auto failover giữa provider
- Provisioning fail sẽ alert Admin để xử lý thủ công
- Hệ thống dùng tỷ giá USD/VND nhập tay
- Support qua Telegram thay vì ticket nội bộ

## 12. KPI phase 1
### KPI kinh doanh
- Doanh thu
- Lợi nhuận
- MRR
- ARPU
- Doanh thu retail Admin
- Doanh thu theo reseller

### KPI vận hành
- Tỷ lệ provision thành công
- Tỷ lệ out-of-stock
- Tỷ lệ renew
- Tỷ lệ expire không gia hạn
- Tỷ lệ lỗi API theo provider

## 13. Rủi ro chính
### Rủi ro kỹ thuật
- API upstream thay đổi
- Sync tồn kho lệch gây oversell
- Mapping capability sai

### Rủi ro nghiệp vụ
- Sai lệch billing cycle
- Sai logic refund
- Tỷ giá thủ công gây chênh lệch

### Rủi ro vận hành
- Backlog provisioning fail
- Không có ticket nội bộ
- Không có anti-fraud phase 1

## 14. Tiêu chí hoàn tất tài liệu
Tài liệu này hoàn tất khi:
- Khóa được phase 1 và ngoài phạm vi
- Khóa được supply model, kênh bán hàng và ownership
- Khóa được các chính sách nền như billing, expiry, cancel/refund và white-label

## 15. Bản vá v1.1 - Các khóa P0 bổ sung trước khi giao dev
Các điểm dưới đây được xem là điều kiện sống còn của phase 1. Nếu chưa khóa, không nên chuyển sang build production.

### 15.1. Reseller settlement là P0
Khi client thuộc reseller mua dịch vụ, hệ thống phải tách rõ hai lớp tiền:
- **Client wallet**: ví nội bộ giữa client và seller/reseller.
- **Reseller wallet**: ví settlement giữa reseller và platform.

Platform chỉ được cấp tài nguyên thật khi reseller wallet đủ tiền giá nhập/reseller cost. Client đã nạp tiền cho reseller không đồng nghĩa platform đã nhận được tiền thật.

Rule khóa:
- Client của Admin mua hàng: debit client wallet, ghi nhận doanh thu retail platform.
- Client của Reseller mua hàng: debit client wallet theo selling price, đồng thời debit reseller wallet theo reseller cost.
- Nếu reseller wallet không đủ reseller cost, order không được provision dù client wallet đủ tiền.
- Profit của reseller = selling price snapshot - reseller cost snapshot - refund/adjustment liên quan.

Chi tiết nằm ở `09_Reseller_Settlement_Ledger_Model.md`.

### 15.2. Tenant isolation là P0
Mọi dữ liệu client/order/service/wallet/ledger/audit phải bị ép scope theo tenant ở tầng backend, không chỉ ở UI. Không được tin `tenant_id` gửi từ request body.

Rule khóa:
- Tenant context lấy từ domain/session/token đã xác thực.
- Mọi query nghiệp vụ phải lọc theo tenant context.
- Admin emergency access chỉ dùng khi có reason và phải audit.
- Credential/access info không được lộ cross-tenant trong mọi trường hợp.

Chi tiết nằm ở `10_Tenant_Security_Access_Control_Spec.md`.

### 15.3. Provisioning safety là P0
Dự án này không được retry mù. Provider timeout có thể là partial success. Nếu hệ thống gọi tạo VPS/proxy, provider đã cấp tài nguyên nhưng API timeout, retry mù có thể tạo trùng tài nguyên và mất tiền.

Rule khóa:
- Mỗi order có idempotency key.
- Mỗi provisioning job lưu external request id/provider response nếu có.
- Lỗi không rõ đã tạo tài nguyên hay chưa phải vào manual review.
- Một order chỉ được attach một service chính thức.

Chi tiết nằm ở `11_Provisioning_Idempotency_And_Inventory_Locking.md`.

### 15.4. Credential và abuse control là P0
Credential VPS/proxy là tài sản nhạy cảm. Abuse control tối thiểu là bắt buộc dù anti-fraud automation chưa làm.

Rule khóa:
- Encrypt credential at rest.
- Không log plaintext credential.
- Mọi lần reveal credential phải audit.
- Phase 1 phải có manual abuse flag, suspend reason, evidence log và blacklist cơ bản.

Chi tiết nằm ở `13_Abuse_Fraud_Operational_Policy_Phase1.md`.

## 16. MVP đề xuất sau khi vá scope
Để không build quá rộng ngay từ đầu, phase 1 nên tách thành các lớp triển khai:

### MVP 0 - Nền tài chính và tenant
- Auth theo tenant.
- Admin tenant và ít nhất 1 reseller tenant mẫu.
- Wallet + immutable ledger.
- Manual top-up approval.
- Master catalog và tenant catalog clone.
- Audit log cơ bản.

### MVP 1 - Bán được thật
- Checkout bằng wallet.
- Reservation atomic.
- Provisioning queue.
- 1 VPS source tự sở hữu hoặc 1 provider ổn định.
- 1 proxy source mẫu.
- Service detail, renew, expiry/suspend/terminate.

### MVP 2 - White-label reseller đủ dùng
- Reseller storefront.
- Domain/brand/theme/contact.
- Tenant pricing riêng.
- Client management.
- Reseller settlement wallet.
- Tenant report cơ bản.

### MVP 3 - Mở rộng provider
- Bổ sung adapter OVH/Hetzner/Smarthost/proxy-cheap.
- Provider health.
- Partial success handling nâng cao.
- Inventory reconciliation.
