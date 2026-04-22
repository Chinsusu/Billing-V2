# Tài liệu 02 - Tenant Model & Role Architecture

## 1. Mục tiêu tài liệu
Tài liệu này định nghĩa:
- Mô hình multi-tenant
- Ownership dữ liệu
- Vai trò hệ thống
- Boundary giữa Admin, Reseller, Staff và Client
- Quy tắc phân quyền mức business

## 2. Mô hình tenant
Hệ thống dùng một codebase multi-tenant.

### 2.1. Tenant loại 1: Admin tenant
Admin tenant là tenant lõi của hệ thống, có quyền quản lý platform và đồng thời có retail storefront riêng.

### 2.2. Tenant loại 2: Reseller tenant
Mỗi reseller là một tenant riêng, có white-label riêng nhưng dùng chung platform core.

### 2.3. Client ownership
Client không phải tenant độc lập. Client thuộc đúng một seller:
- Hoặc Admin
- Hoặc một Reseller

## 3. Các vai trò chính
### 3.1. Admin
Phạm vi:
- Platform-wide operations
- Provider management
- Catalog management
- Pricing base
- Exchange rate
- Payment verification
- Refund handling
- Manual provisioning queue
- Reports ở mức toàn hệ thống

### 3.2. Reseller Owner
Phạm vi:
- Quản lý storefront white-label
- Clone catalog từ Admin
- Set giá riêng
- Quản lý client của tenant
- Quản lý order/service của tenant
- Quản lý staff của tenant
- Quản lý ví tenant
- Xem report/audit của tenant

### 3.3. Reseller Staff
Hiện tại dùng một role staff chung.
Đề xuất scope mặc định:
- Xem dashboard tenant
- Xem client và service
- Hỗ trợ tác vụ vận hành hằng ngày
- Không được thay đổi policy nền hoặc white-label nhạy cảm trừ khi Owner cấp phép

### 3.4. Client
Phạm vi:
- Đăng ký/đăng nhập
- Nạp tiền
- Mua hàng
- Gia hạn
- Quản lý dịch vụ trong phạm vi capability cho phép
- Xem lịch sử giao dịch và thông tin billing

## 4. White-label model
Mỗi reseller tenant có thể có:
- Domain riêng
- Trang login riêng
- Theme riêng
- Logo và brand name riêng
- Email template riêng
- Contact/Telegram riêng
- Timezone riêng
- Ngôn ngữ hiển thị riêng trong phạm vi hỗ trợ

## 5. Ownership dữ liệu
### 5.1. Dữ liệu platform-level
Thuộc platform core:
- Danh sách provider
- Adapter config
- Exchange rate
- Catalog gốc
- Chính sách lõi
- Report tổng
- Audit platform

### 5.2. Dữ liệu tenant-level
Thuộc Admin tenant hoặc từng Reseller tenant:
- White-label config
- Catalog clone
- Giá bán
- Client list
- Wallet tenant
- Payment records của tenant
- Audit của tenant

### 5.3. Dữ liệu client-level
- Hồ sơ tài khoản
- Ví
- Đơn hàng
- Dịch vụ
- Lịch sử gia hạn
- Yêu cầu hóa đơn

### 5.4. Dữ liệu service-level
- Product/plan đã mua
- Nguồn provision thực tế
- Trạng thái vòng đời
- Các credential/info cấp dịch vụ
- Các action log liên quan tới service

## 6. Boundary giữa Admin và Reseller
### 6.1. Rule đã chốt
- Admin không vận hành reseller như một user reseller bình thường
- Reseller tự quản lý client và service trong tenant của mình
- Reseller có ví riêng và trách nhiệm tài chính riêng

### 6.2. Giả định/đề xuất cần chốt thêm
Để phục vụ audit, dispute, abuse và khôi phục sự cố, nên có chính sách:
- Admin có emergency read-only access ở mức cần thiết
- Mọi lần truy cập khẩn cấp đều bị audit

Nếu không áp dụng chính sách này, cần chấp nhận rủi ro khó xử lý sự cố tenant.

## 7. Permission matrix mức cao
| Module | Admin | Reseller Owner | Reseller Staff | Client |
|---|---|---|---|---|
| Provider management | Full | No | No | No |
| Master catalog | Full | No | No | No |
| Catalog clone | View | Full | Limited | View storefront only |
| Pricing | Base pricing | Tenant pricing | Limited | No |
| Wallet | Full platform scope | Tenant scope | Limited | Own wallet |
| Payment verification | Full | Tenant scope nếu có | Limited | No |
| Order management | Full platform scope | Tenant scope | Limited | Own orders |
| Service management | Full platform scope | Tenant scope | Limited | Own services |
| White-label settings | No/Own Admin brand | Full | Limited | No |
| Reports | Full platform scope | Tenant scope | Limited | Own summary |
| Audit logs | Full platform scope | Tenant scope | Limited | No |

## 8. Permission matrix chi tiết gợi ý
### 8.1. Admin
- Create/update/disable provider
- Create/update master plan
- Set base price/reseller cost
- Verify payment
- Manual credit/debit wallet
- Force suspend/terminate service
- Approve refund
- Update exchange rate
- View global reports
- Access manual provisioning queue

### 8.2. Reseller Owner
- Clone/enable/disable catalog items trong tenant
- Set tenant selling price
- Create/edit staff
- Suspend service của client thuộc tenant
- View tenant audit
- View wallet, ledger, revenue, profit
- Configure Telegram/contact, theme, domain mapping

### 8.3. Reseller Staff
- Scope nên cấu hình bằng cờ quyền trong phase sau
- Phase 1 có thể dùng role staff chung với một preset mặc định
- Các quyền nhạy cảm nên mặc định tắt: refund, đổi giá, đổi white-label, đổi domain

### 8.4. Client
- View catalog
- Top up wallet
- Create order
- Renew
- Cancel nếu plan cho phép
- Thực hiện các action service được hỗ trợ
- Xem usage và lịch sử thanh toán của chính mình

## 9. Quy tắc cô lập tenant
- Mọi dữ liệu tenant phải có tenant_id hoặc seller scope tương đương
- Không cho phép join dữ liệu chéo tenant ở tầng nghiệp vụ
- Các dashboard và báo cáo phải luôn chạy theo context tenant
- Các action admin ở mức tenant phải lưu audit riêng

## 10. Danh tính và xác thực
### 10.1. Client auth
- Email là định danh đăng nhập
- 2FA không bắt buộc
- Không social login ở phase 1

### 10.2. Staff/Admin auth
Giả định/đề xuất:
- 2FA bắt buộc cho Admin trong phase 1
- 2FA bật mặc định/khuyến nghị bắt buộc cho Reseller Owner trong phase 1
- Nếu có ngoại lệ tạm thời, ngoại lệ phải có expiry date, reason và audit
- Client 2FA có thể để optional ở phase sau, nhưng phải có chính sách mật khẩu mạnh và audit login

## 11. Multi-language, timezone, currency
- Ngôn ngữ hỗ trợ: Việt / Anh
- Timezone có thể cấu hình theo tenant
- USD là đơn vị chính
- VND là đơn vị quy đổi theo tỷ giá nhập tay

## 12. Những câu hỏi nên khóa ở vòng tiếp theo
- Admin có emergency read-only access hay không
- Role staff chung có những quyền gì mặc định
- Có cho reseller tự map domain qua quy trình self-service hay không
- Có cho reseller tự xuất thủ công invoice từ tenant của mình hay chỉ request lên Admin

## 13. Bản vá v1.1 - Tenant enforcement bắt buộc
Các rule dưới đây là rule kỹ thuật-nghiệp vụ bắt buộc để tránh lộ dữ liệu chéo tenant.

### 13.1. Nguồn xác định tenant context
Tenant context phải được xác định theo thứ tự ưu tiên:
1. Domain hoặc custom domain đang truy cập.
2. Session/token đã xác thực.
3. Actor role và seller scope.

Không được lấy tenant context từ request body cho các tài nguyên nhạy cảm như order, wallet, service, credential, ledger.

### 13.2. Rule kiểm tra ownership
Mọi thao tác đọc/ghi với tài nguyên tenant phải thỏa một trong hai điều kiện:
- `resource.tenant_id` khớp với tenant context hiện tại.
- Actor là platform admin đang dùng emergency access có reason và audit.

Các tài nguyên cần guard nghiêm ngặt:
- users/clients
- wallets
- top-up requests
- orders
- reservations
- provisioning jobs
- services
- credentials/access info
- ledger entries
- invoices
- audit events

### 13.3. Composite ownership đề xuất
Ở mức data model, các bảng nghiệp vụ quan trọng nên có khóa/unique constraint hoặc rule truy vấn dựa trên cặp:
- `tenant_id + user_id`
- `tenant_id + wallet_id`
- `tenant_id + order_id`
- `tenant_id + service_id`
- `tenant_id + ledger_entry_id`

Mục tiêu không phải làm DB phức tạp, mà là ép dev không thể vô tình query nhầm tenant.

### 13.4. Emergency access
Admin chỉ được truy cập dữ liệu tenant ngoài phạm vi thông thường khi:
- Có reason bắt buộc.
- Có ticket/reference nội bộ hoặc note vận hành.
- Mặc định read-only.
- Nếu có thao tác write, phải ghi audit riêng với action `tenant.emergency.write`.

### 13.5. Domain mapping và white-label
Custom domain phải có quy trình xác minh ownership trước khi active:
- Reseller khai báo domain.
- Hệ thống sinh verification token.
- Reseller thêm DNS TXT/CNAME theo hướng dẫn.
- Admin hoặc hệ thống xác minh.
- Domain mới chuyển sang active sau khi TLS/certificate sẵn sàng.

Không cho một domain active ở nhiều tenant cùng lúc.

### 13.6. 2FA và session security
- Admin bắt buộc bật 2FA.
- Reseller Owner bật mặc định; nếu tắt phải có policy và audit.
- Reseller Staff nên có 2FA nếu được phép thao tác wallet, pricing, service action.
- Session admin/reseller owner nên ngắn hơn client session.
- Các action nhạy cảm nên yêu cầu re-auth hoặc step-up auth: refund, manual adjustment, reveal credential, domain change, staff permission change.
