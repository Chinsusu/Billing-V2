# VPS Proxy Project Master Document - v1.4 Package Index And Technical Handoff Base
Generated: 2026-04-22 Asia/Ho_Chi_Minh

This master document keeps the v1.2 technical handoff base from `00` to `23`, plus changelogs and package indexes for later v1.3/v1.4 layers.

---



# ===== 00_README.md =====

# Bộ tài liệu dự án nền tảng VPS/Proxy hybrid multi-tenant - v1.2 Technical Handoff

## Mô tả
Gói này là bản mở rộng technical handoff cho dự án web thuê/bán VPS/Proxy theo mô hình:
- Hybrid multi-provider
- Multi-tenant
- White-label reseller
- Wallet-first
- Admin có bán lẻ trực tiếp
- Reseller có storefront/client/wallet/pricing riêng

Bản v1.1 đã vá các điểm P0 về reseller settlement, tenant isolation, inventory locking, provisioning idempotency, credential security, audit traceability và abuse control phase 1.

Bản v1.2 bổ sung bộ tài liệu `14–23` để chuyển blueprint thành tài liệu bàn giao dev/backend/frontend/QA/DevOps.

## Nguyên tắc bản v1.2
- Chưa code.
- Không ràng buộc framework/ngôn ngữ triển khai.
- Tập trung vào data contract, API behavior, permission, worker, UI wireflow, QA, deployment và notification.
- Mọi rule liên quan tiền, tenant, credential, provider và provisioning được xem là P0.
- Dev không nên tự đoán ở các flow tiền, tenant, provisioning, stock, credential.

## Thành phần gói

### Tài liệu nền đã vá từ bản v1.1
- `01_Product_Scope_Business_Model.md`
- `02_Tenant_Model_Role_Architecture.md`
- `03_Product_Catalog_Pricing_Rules.md`
- `04_Billing_Wallet_Ledger_Spec.md`
- `05_Provisioning_Provider_Adapter_Spec.md`
- `06_Order_Service_Lifecycle_State_Machine.md`
- `07_Portal_Functional_Spec.md`
- `08_Audit_Reports_Operational_Control.md`
- `09_Reseller_Settlement_Ledger_Model.md`
- `10_Tenant_Security_Access_Control_Spec.md`
- `11_Provisioning_Idempotency_And_Inventory_Locking.md`
- `12_API_Data_Model_Acceptance_Criteria.md`
- `13_Abuse_Fraud_Operational_Policy_Phase1.md`

### Tài liệu technical handoff thêm ở bản v1.2
- `14_System_Architecture_Blueprint.md`
- `15_Database_Schema_And_ERD.md`
- `16_API_Contract_And_Permission_Spec.md`
- `17_RBAC_Permission_Matrix.md`
- `18_Provider_Adapter_Technical_Spec.md`
- `19_Worker_Queue_And_Cron_Jobs_Spec.md`
- `20_UI_Wireflow_And_Screen_Spec.md`
- `21_QA_Test_Cases_And_Acceptance_Plan.md`
- `22_Deployment_DevOps_And_Environment_Runbook.md`
- `23_Notification_Email_Telegram_Template_Spec.md`

### Tài liệu tổng hợp và ghi chú
- `VPS_Proxy_Project_Master_Document.md`
- `CHANGELOG_FIXES.md`
- `CHANGELOG_TECHNICAL_HANDOFF_v1_2.md`
- `MANIFEST.txt`

## Cách đọc đề xuất

### Cho founder/product owner
1. `01_Product_Scope_Business_Model.md`
2. `09_Reseller_Settlement_Ledger_Model.md`
3. `14_System_Architecture_Blueprint.md`
4. `20_UI_Wireflow_And_Screen_Spec.md`
5. `21_QA_Test_Cases_And_Acceptance_Plan.md`

### Cho backend/dev lead
1. `14_System_Architecture_Blueprint.md`
2. `15_Database_Schema_And_ERD.md`
3. `16_API_Contract_And_Permission_Spec.md`
4. `17_RBAC_Permission_Matrix.md`
5. `18_Provider_Adapter_Technical_Spec.md`
6. `19_Worker_Queue_And_Cron_Jobs_Spec.md`

### Cho frontend
1. `16_API_Contract_And_Permission_Spec.md`
2. `17_RBAC_Permission_Matrix.md`
3. `20_UI_Wireflow_And_Screen_Spec.md`
4. `23_Notification_Email_Telegram_Template_Spec.md`

### Cho QA
1. `12_API_Data_Model_Acceptance_Criteria.md`
2. `15_Database_Schema_And_ERD.md`
3. `16_API_Contract_And_Permission_Spec.md`
4. `21_QA_Test_Cases_And_Acceptance_Plan.md`

### Cho DevOps/Ops
1. `14_System_Architecture_Blueprint.md`
2. `18_Provider_Adapter_Technical_Spec.md`
3. `19_Worker_Queue_And_Cron_Jobs_Spec.md`
4. `22_Deployment_DevOps_And_Environment_Runbook.md`
5. `23_Notification_Email_Telegram_Template_Spec.md`

## 10 luật nền phải giữ
1. Không provision nếu tiền chưa được debit/lock hợp lệ.
2. Không debit ví nếu không tạo được ledger entry.
3. Không có ledger thì giao dịch không tồn tại.
4. Không có tenant context thì không cho đọc/ghi dữ liệu tenant.
5. Không retry provisioning nếu không biết provider đã tạo tài nguyên hay chưa.
6. Không hiển thị action nếu capability snapshot không cho phép.
7. Không sửa transaction cũ; chỉ tạo adjustment/reversal.
8. Không lưu credential plaintext trong log/audit.
9. Không cho client reseller provision nếu reseller wallet không đủ reseller cost.
10. Không dùng giá/policy hiện tại để xử lý tranh chấp order cũ; dùng snapshot lúc mua.

## Mục tiêu sau bản v1.2
Sau khi đọc xong gói này, team dev phải trả lời được:
- Cần tạo bảng nào và bảng nào bắt buộc có tenant_id.
- API nào cần build và role nào được gọi.
- Checkout, wallet, reseller settlement và provisioning chạy theo flow nào.
- Provider adapter phải trả về gì khi success/fail/timeout/partial success.
- Worker/cron nào phải chạy nền.
- UI cần màn nào và action nào phải ẩn/chặn.
- QA phải test case nào trước khi nghiệm thu.
- Production cần backup, monitoring, secret và rollback như thế nào.

---


# ===== 01_Product_Scope_Business_Model.md =====

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

---


# ===== 02_Tenant_Model_Role_Architecture.md =====

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

---


# ===== 03_Product_Catalog_Pricing_Rules.md =====

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

---


# ===== 04_Billing_Wallet_Ledger_Spec.md =====

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

---


# ===== 05_Provisioning_Provider_Adapter_Spec.md =====

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

---


# ===== 06_Order_Service_Lifecycle_State_Machine.md =====

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

---


# ===== 07_Portal_Functional_Spec.md =====

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

---


# ===== 08_Audit_Reports_Operational_Control.md =====

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

---


# ===== 09_Reseller_Settlement_Ledger_Model.md =====

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

---


# ===== 10_Tenant_Security_Access_Control_Spec.md =====

# Tài liệu 10 - Tenant Security & Access Control Specification

## 1. Mục tiêu tài liệu
Tài liệu này khóa các rule bảo mật tenant, role, domain, credential và audit access cho phase 1.

Mục tiêu lớn nhất: không để client/reseller/staff đọc hoặc thao tác nhầm dữ liệu của tenant khác.

## 2. Nguyên tắc lõi
- Không tin `tenant_id` từ request body.
- Tenant context lấy từ domain/session/token đã xác thực.
- Backend phải enforce tenant scope; UI ẩn không đủ.
- Credential là secret, không phải dữ liệu text bình thường.
- Emergency access của Admin phải có reason và audit.

## 3. Request context chuẩn
Mỗi request sau khi auth nên có context:
- actor_id
- actor_type: admin, reseller_owner, reseller_staff, client, system
- actor_role
- tenant_id
- seller_scope
- domain_context
- permissions
- session_id
- request_id
- correlation_id nếu đang trong flow order/payment/provisioning

Rule:
- API không lấy tenant từ body nếu request đã có tenant context.
- Nếu domain và token conflict tenant, request bị reject và audit security event.

## 4. Domain-to-tenant mapping
### 4.1. Admin storefront
Admin có domain mặc định và có thể có brand riêng.

### 4.2. Reseller storefront
Mỗi custom domain map tới đúng một reseller tenant.

Trạng thái domain:
- requested
- pending_verification
- verified
- tls_pending
- active
- suspended
- removed

### 4.3. Xác minh domain
Flow:
1. Reseller khai báo domain.
2. Hệ thống sinh verification token.
3. Reseller thêm DNS TXT/CNAME.
4. Hệ thống/Admin xác minh.
5. TLS/certificate sẵn sàng.
6. Domain active.

Rule:
- Không cho một domain active ở nhiều tenant.
- Domain suspended không được route vào tenant.
- Domain mapping change phải audit.

## 5. Backend tenant enforcement
Mọi resource tenant-level cần check:
- resource.tenant_id == context.tenant_id
- hoặc actor là platform admin với emergency access hợp lệ.

Tài nguyên bắt buộc scope:
- users/clients
- wallets
- ledger entries
- payment/top-up requests
- catalog clone/tenant plans
- orders
- reservations
- provisioning jobs
- services
- credentials/access info
- invoices
- audit logs
- abuse/risk flags

## 6. Composite ownership rule
Data model nên thiết kế để dev khó query sai:
- wallets: tenant_id + wallet_id
- orders: tenant_id + order_id
- services: tenant_id + service_id
- ledger: tenant_id + ledger_entry_id
- users: tenant_id + user_id nếu user thuộc tenant

Khi lấy service, không chỉ `WHERE service_id = ?`; phải có tenant scope.

## 7. Role và permission phase 1
### 7.1. Platform Admin
Có toàn quyền platform, nhưng khi vào tenant reseller để hỗ trợ nên dùng emergency context.

P0 action:
- provider management
- master catalog
- reseller management
- payment verification platform-level
- refund/adjustment
- force suspend/terminate
- audit/report global

### 7.2. Reseller Owner
Có quyền tenant-level:
- white-label settings
- tenant pricing
- client/service management
- staff management
- tenant wallet/report/audit

### 7.3. Reseller Staff
Phase 1 dùng preset hạn chế:
- view dashboard
- view clients/services
- hỗ trợ thao tác vận hành cơ bản nếu Owner bật

Mặc định không có:
- refund
- manual wallet adjustment
- domain change
- white-label sensitive change
- pricing change
- staff permission change
- credential reveal cho service nhạy cảm nếu chưa bật quyền

### 7.4. Client
Chỉ được xem và thao tác tài nguyên của chính mình trong tenant hiện tại.

## 8. Emergency access
Admin emergency access phải có:
- reason bắt buộc
- time window
- actor_id
- target_tenant_id
- request_id
- audit action

Mặc định read-only. Write access cần action riêng và reason riêng.

Audit action:
- tenant.emergency.read
- tenant.emergency.write

## 9. Credential security
Credential/access info gồm:
- VPS username/password
- proxy endpoint/username/password/token
- VNC/console token
- private key
- provider token/API key

Rule:
- Encrypt at rest.
- Không log plaintext.
- Không lưu plaintext trong audit snapshot.
- UI masked by default.
- Reveal phải audit `credential.revealed`.
- Có thể yêu cầu re-auth/2FA trước khi reveal.
- Download/export credential nếu có phải audit riêng.

## 10. Provider secret management
Provider API key/token không nên lưu thô trong DB.

Rule đề xuất:
- Lưu trong secret vault hoặc encrypted config.
- Chỉ adapter runtime được quyền đọc.
- Rotate secret phải có audit.
- Không expose secret qua admin UI sau khi đã lưu; chỉ cho replace.

## 11. 2FA và session policy
P0:
- Admin bắt buộc 2FA.
- Reseller Owner bật mặc định/khuyến nghị bắt buộc.
- Staff có quyền tiền/credential/pricing nên bắt buộc 2FA.

Session:
- Admin/reseller owner session timeout ngắn hơn client.
- Action nhạy cảm yêu cầu re-auth/step-up auth.
- Login failed nhiều lần phải rate limit.

## 12. Rate limit P0
Áp dụng cho:
- login
- register
- password reset
- top-up request
- checkout
- renew
- change IP
- reinstall
- credential reveal
- domain verification

Rate limit nên có theo IP + actor + tenant.

## 13. Audit và security events
Cần audit:
- login success/failed
- tenant context mismatch
- permission denied
- credential reveal
- domain mapping change
- 2FA enable/disable
- staff permission change
- wallet/refund/adjustment
- emergency access

## 14. Acceptance criteria
- Client tenant A không thể gọi API lấy service tenant B.
- Staff không có quyền pricing không thể update tenant plan price dù biết API endpoint.
- Request body gửi tenant_id khác tenant context bị reject.
- Credential reveal luôn có audit.
- Admin emergency read có reason và log.
- Domain đã active ở tenant A không thể active ở tenant B.
- Provider API key không xuất hiện trong response/log/audit plaintext.

---


# ===== 11_Provisioning_Idempotency_And_Inventory_Locking.md =====

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

---


# ===== 12_API_Data_Model_Acceptance_Criteria.md =====

# Tài liệu 12 - API, Data Model & Acceptance Criteria

## 1. Mục tiêu tài liệu
Tài liệu này không viết code. Nó mô tả data model nghiệp vụ, API flow và acceptance criteria để đội dev biết cần build gì và test thế nào.

## 2. Shared fields nên có
Các bảng nghiệp vụ quan trọng nên có:
- id
- tenant_id nếu thuộc tenant
- created_at
- updated_at
- created_by
- updated_by
- status
- metadata hoặc notes nếu cần
- request_id/correlation_id với flow tài chính/provisioning

Các bảng tài chính/provisioning nên có thêm:
- idempotency_key
- reference_type/reference_id
- snapshot fields

## 3. Entities P0
### 3.1. Tenant
Fields:
- tenant_id
- tenant_type: admin, reseller
- status: active, suspended, disabled
- brand_name
- default_currency_display
- timezone
- contact_support

### 3.2. Domain mapping
Fields:
- domain_id
- tenant_id
- domain
- verification_status
- tls_status
- status
- verification_token_reference

### 3.3. User
Fields:
- user_id
- tenant_id/seller_scope
- role
- email
- status
- two_factor_enabled
- last_login_at

### 3.4. Wallet
Fields:
- wallet_id
- tenant_id
- owner_type: admin_client, reseller, reseller_client
- owner_id
- currency
- available_balance
- locked_balance
- status

### 3.5. Ledger entry
Fields:
- ledger_entry_id
- tenant_id
- wallet_id
- direction
- amount
- currency
- fx_snapshot
- entry_type
- reference_type/reference_id
- idempotency_key
- balance_before
- balance_after
- status
- reason

### 3.6. Product / Plan / Source
Product fields:
- product_id
- type
- display_name
- location
- visibility

Plan fields:
- plan_id
- product_id
- version
- cycle_type
- cycle_definition
- retail_price
- reseller_cost
- refund_policy
- capability_flags
- status

Source fields:
- source_id
- provider_id
- plan_id
- capacity
- reserved_count
- allocated_count
- stock_state
- capability_override
- status

### 3.7. Tenant plan clone
Fields:
- tenant_plan_id
- tenant_id
- master_plan_id
- master_plan_version
- selling_price
- reseller_cost_snapshot
- enabled
- margin_state
- policy_snapshot

### 3.8. Order
Fields:
- order_id
- tenant_id
- buyer_user_id
- seller_type: admin, reseller
- plan_snapshot
- price_snapshot
- reseller_cost_snapshot
- order_status
- billing_status
- correlation_id

### 3.9. Reservation
Fields:
- reservation_id
- tenant_id
- order_id
- source_id
- status
- quantity
- expires_at

### 3.10. Provisioning job
Fields:
- job_id
- tenant_id
- order_id
- reservation_id
- source_id
- provider_id
- idempotency_key
- external_request_id
- external_resource_id
- status
- attempt_count
- retry_safety_level
- last_error_summary

### 3.11. Service instance
Fields:
- service_id
- tenant_id
- owner_user_id
- order_id
- provider_id
- source_id
- external_resource_id
- service_status
- billing_status
- suspension_reason
- term_start
- term_end
- grace_end
- capability_snapshot
- billing_cycle_snapshot
- credential_reference

### 3.12. Audit event
Fields:
- audit_id
- tenant_id
- actor_id
- actor_role
- action
- resource_type
- resource_id
- before_summary
- after_summary
- reason
- request_id
- correlation_id
- source_ip
- created_at

### 3.13. Risk/abuse flag
Fields:
- risk_flag_id
- tenant_id
- target_type: user, order, service, payment, ip, domain
- target_id
- risk_type
- severity
- status
- evidence_summary
- created_by

## 4. API flow P0 và acceptance criteria
### 4.1. Register/Login
Acceptance:
- User được gắn đúng tenant từ domain context.
- Email trùng trong cùng tenant bị chặn theo policy.
- Login failed bị rate limit.
- Admin login không có 2FA bị chặn nếu policy bắt buộc.

### 4.2. Create top-up request
Acceptance:
- Request thuộc đúng tenant/user.
- Amount > 0.
- Payment method enabled.
- Sinh reference duy nhất.
- Không credit wallet trước khi approved.

### 4.3. Approve top-up
Acceptance:
- Actor có quyền approve.
- Một request chỉ approve một lần.
- Credit wallet và ledger xảy ra cùng một control flow.
- Audit có actor/reason/reference.

### 4.4. Checkout Admin client
Acceptance:
- Plan active.
- Stock available.
- Client wallet đủ tiền.
- Reservation atomic.
- Wallet debit + ledger + order paid + provisioning job được tạo an toàn.
- Double click không double debit.

### 4.5. Checkout Reseller client
Acceptance:
- Client wallet đủ selling price.
- Reseller wallet đủ reseller cost.
- Thiếu một trong hai ví thì không provision.
- Order lưu selling price snapshot và reseller cost snapshot.
- Report reseller profit tính đúng.

### 4.6. Renew service
Acceptance:
- Service active hoặc suspended trong grace.
- Term mới cộng từ old_term_end theo cycle.
- Nếu thuộc reseller, debit cả client wallet và reseller wallet.
- Terminated service không renew.

### 4.7. Cancel/refund
Acceptance:
- Chỉ plan có allow_mid_cycle_cancel mới hiện action.
- Refund dùng policy snapshot.
- Refund không vượt charge gốc.
- Service không còn active sau refund toàn phần/cancel hợp lệ.
- Audit có reason.

### 4.8. Service action
Các action: start, stop, reboot, reinstall, change password, console, change IP.

Acceptance:
- Action chỉ hiện nếu capability_snapshot cho phép.
- Tenant ownership được check backend.
- Paid action debit wallet trước khi gọi provider.
- Provider fail thì refund/rollback theo policy action.
- Action log/audit được ghi.

### 4.9. Reveal credential
Acceptance:
- Credential masked mặc định.
- Reveal yêu cầu quyền hợp lệ.
- Cross-tenant bị chặn.
- Audit `credential.revealed` được ghi.

### 4.10. Reseller pricing update
Acceptance:
- Chỉ Owner/staff có quyền mới update.
- Giá mới không ảnh hưởng order cũ.
- Negative margin bị cảnh báo hoặc block theo policy.
- Audit `pricing.tenant_plan.updated`.

### 4.11. Domain mapping
Acceptance:
- Domain phải verify trước active.
- Không active domain trùng tenant khác.
- Mapping change audit.
- Domain suspended không route vào tenant.

## 5. Error states P0
UI/API cần trả trạng thái rõ cho:
- insufficient_client_balance
- insufficient_reseller_balance
- out_of_stock
- reservation_expired
- plan_disabled
- source_disabled
- provider_unavailable
- permission_denied
- tenant_scope_mismatch
- provisioning_manual_review
- payment_verification_required
- refund_not_allowed

## 6. Acceptance test tối thiểu trước production
- Cross-tenant API tests.
- Double-click checkout tests.
- Concurrent reservation tests.
- Top-up approve idempotency tests.
- Reseller settlement tests.
- Provider timeout/partial success tests.
- Refund boundary tests.
- Credential redaction tests.
- Audit correlation tests.

---


# ===== 13_Abuse_Fraud_Operational_Policy_Phase1.md =====

# Tài liệu 13 - Abuse, Fraud & Operational Policy Phase 1

## 1. Mục tiêu tài liệu
Phase 1 chưa cần anti-fraud automation phức tạp, nhưng phải có manual abuse/fraud control tối thiểu. VPS/proxy là nhóm sản phẩm dễ bị spam, scan, brute-force, fraud payment và provider takedown.

Mục tiêu: giảm rủi ro provider khóa nguồn, giảm chargeback, có bằng chứng xử lý khi xảy ra abuse.

## 2. Nguyên tắc lõi
- Không provision tự động cho order rủi ro cao nếu policy manual review bật.
- Mọi suspend/terminate vì abuse phải có reason và evidence summary.
- Không xóa dấu vết abuse.
- Không để provider complaint trôi không owner.
- Terms/AUP acceptance phải có log.

## 3. AUP/Terms tối thiểu
Khi đăng ký hoặc checkout, user cần chấp nhận policy cấm:
- spam/phishing
- brute force/scanning trái phép
- malware/botnet/C2
- credential stuffing
- DDoS
- fraud payment/chargeback abuse
- vi phạm luật/provider AUP
- sử dụng proxy/VPS để che giấu hành vi gây hại

Không cần legal dài ở phase 1, nhưng phải có bản policy ngắn rõ ràng và log acceptance.

## 4. Risk flags
Risk flag có thể gắn vào:
- user
- order
- service
- payment/top-up
- IP/domain/email
- provider/source

Trạng thái:
- open
- under_review
- cleared
- action_taken
- closed

Severity:
- low
- medium
- high
- critical

## 5. Manual review triggers phase 1
Đưa order/payment vào manual review nếu có một hoặc nhiều dấu hiệu:
- Account mới tạo mua volume lớn.
- Nạp tiền lớn bất thường.
- Nhiều top-up nhỏ liên tục.
- Email/domain khả nghi theo blacklist nội bộ.
- Nhiều account cùng IP/device.
- Order nhiều proxy/VPS trong thời gian ngắn.
- Payment proof không khớp amount/reference.
- Reseller có lịch sử abuse cao.
- Provider/source đang có complaint spike.

## 6. Payment fraud control
### 6.1. Manual payment
- Mã tham chiếu bắt buộc.
- Amount/currency phải khớp.
- Nếu thiếu/chênh, chuyển pending_verification.
- Approve phải có actor và audit.

### 6.2. PayPal/USDT
Nếu chưa tích hợp chắc chắn:
- Dùng manual/semi-auto fallback.
- Không auto credit khi giao dịch chưa final/confirmed theo rule payment method.
- Chargeback/dispute phải tạo risk flag.

## 7. Abuse complaint workflow
Khi có complaint từ provider/third party:
1. Tạo abuse case/risk flag.
2. Gắn target service/user/tenant/provider.
3. Lưu evidence summary.
4. Xác định severity.
5. Nếu high/critical, suspend service trước để giảm thiệt hại.
6. Thông báo reseller/client theo policy.
7. Resolve: clear, warn, suspend, terminate, blacklist, refund/no refund.
8. Ghi audit.

## 8. Suspension/termination policy
Suspension reason:
- expiry
- manual_admin
- manual_reseller
- abuse
- payment_risk
- provider_complaint
- system_issue

Rule:
- Abuse suspension phải có reason và evidence summary.
- Terminate vì abuse cần quyền Admin hoặc Reseller Owner tùy policy.
- Refund abuse có thể bằng 0 nếu AUP quy định rõ.
- Nếu suspend nhầm, khôi phục phải audit.

## 9. Blacklist tối thiểu
Phase 1 nên có blacklist thủ công cho:
- email
- email domain
- IP đăng ký/login
- payment reference
- wallet/payment identity nếu có
- provider resource/IP bị complaint
- client/user
- reseller tenant nếu nghiêm trọng

Blacklist action:
- block_register
- block_topup
- block_checkout
- manual_review_only
- block_service_action

## 10. Account maturity limit
Để giảm abuse, có thể đặt limit theo tuổi account:
- Account mới: giới hạn số lượng order/service.
- Account mới: order lớn phải manual review.
- Sau một số ngày hoặc sau lịch sử thanh toán sạch, tăng limit.

Rule này nên cấu hình được theo tenant và provider/source.

## 11. Provider protection
Mỗi provider/source nên có:
- abuse complaint count
- recent complaint window
- source risk state: normal, watch, restricted, disabled
- manual review required flag

Nếu provider bị complaint tăng mạnh:
- tạm bật manual review cho source
- giảm limit order mới
- kiểm tra reseller/client gây rủi ro

## 12. Notifications
Thông báo tối thiểu:
- Service suspended vì abuse/payment risk.
- Order chuyển manual review.
- Top-up cần xác minh thêm.
- Provider complaint cần xử lý.
- Reseller có client bị abuse case.

## 13. Audit actions
- abuse.flag.created
- abuse.flag.updated
- abuse.service.suspended
- abuse.service.terminated
- abuse.user.blacklisted
- abuse.payment.marked_risky
- abuse.case.cleared
- risk.manual_review.required

## 14. KPI vận hành
- số abuse case theo ngày/tuần
- abuse case theo provider/source
- abuse case theo reseller
- manual review rate
- false positive rate nếu có
- suspend/terminate vì abuse
- chargeback/dispute count
- provider complaint response time

## 15. Acceptance criteria
- Có thể gắn abuse flag vào user/order/service.
- Suspend vì abuse bắt buộc có reason/evidence summary.
- Order rủi ro có thể chuyển manual review trước provisioning.
- Blacklist email/IP có hiệu lực ở register/checkout theo rule.
- Provider complaint có owner và trạng thái xử lý.
- Audit đầy đủ cho suspend/terminate/clear case.

---


# ===== 14_System_Architecture_Blueprint.md =====

# 14 - System Architecture Blueprint

## 1. Mục tiêu tài liệu

Tài liệu này mô tả kiến trúc tổng thể cho nền tảng VPS/Proxy hybrid multi-tenant, white-label reseller, wallet-first.

Mục tiêu là để team kỹ thuật thống nhất:
- Hệ thống gồm những khối nào.
- Dữ liệu đi qua các lớp nào.
- Tenant context được xác định và bảo vệ ở đâu.
- Billing, provisioning, inventory, credential và audit được đặt ở tầng nào.
- Module nào xử lý sync realtime, module nào xử lý bất đồng bộ qua worker/cron.
- Các nguyên tắc không được phá khi chọn framework hoặc ngôn ngữ triển khai.

Tài liệu này **chưa phải code**, không ràng buộc framework cụ thể.

---

## 2. Nguyên tắc kiến trúc nền

### 2.1 Wallet-first

Không tạo tài nguyên thật nếu chưa có trạng thái tiền hợp lệ.

Với client trực tiếp của platform:
```text
client_wallet đủ tiền -> debit wallet -> tạo order/provisioning job
```

Với client thuộc reseller:
```text
client_wallet đủ selling_price
và reseller_wallet đủ reseller_cost
-> debit cả 2 lớp theo rule settlement
-> tạo order/provisioning job
```

### 2.2 Tenant-first

Mọi request nghiệp vụ phải có tenant context hợp lệ trước khi đọc/ghi dữ liệu.

Không tin các field:
```text
tenant_id
seller_id
reseller_id
owner_id
```

nếu chúng đến từ request body/client. Các giá trị này phải lấy từ:
```text
domain mapping
authenticated session/token
server-side user membership
admin emergency access context
```

### 2.3 Ledger-first

Mọi thay đổi số dư phải đi qua ledger immutable.

Không có ledger entry thì xem như giao dịch không tồn tại.

### 2.4 Idempotency-first

Các thao tác có rủi ro lặp, đặc biệt là checkout/provisioning/payment webhook/retry, phải có idempotency key.

Cấm retry mù khi không biết provider đã tạo tài nguyên hay chưa.

### 2.5 Audit-by-default

Mọi thao tác liên quan tiền, credential, tenant, domain, quyền, provisioning, suspend/terminate phải ghi audit.

Audit không được chứa plaintext secret/credential.

### 2.6 Capability masking

UI và API không hiển thị/không cho gọi action nếu provider/source/service capability snapshot không hỗ trợ.

Ví dụ:
```text
Service từ source không hỗ trợ change_ip
-> UI không hiện nút change_ip
-> API vẫn phải chặn nếu user gọi trực tiếp endpoint change_ip
```

---

## 3. Kiến trúc tổng thể

```text
[ Admin Portal ]        [ Reseller Portal ]        [ Client Portal ]
       |                        |                         |
       +------------------------+-------------------------+
                                |
                         [ Web / API Layer ]
                                |
               [ Auth + Tenant Context + RBAC Guard ]
                                |
     +--------------------------+------------------------------+
     |                          |                              |
[ Core Domain Modules ]   [ Financial Core ]            [ Security Core ]
     |                          |                              |
     |                          |                              |
Catalog / Product         Wallet / Ledger                Audit Log
Order / Checkout          Top-up / Adjustment            Secret/Credential
Inventory / Reservation   Settlement                     Rate Limit
Service Lifecycle         Invoice Request                Abuse/Risk
                                |
                         [ Queue / Job Bus ]
                                |
                 +--------------+--------------+
                 |                             |
         [ Provisioning Worker ]       [ Cron/Scheduler ]
                 |                             |
          [ Provider Adapter Layer ]     Expiry/Renewal/Sync
                 |
       +---------+---------+----------+----------+
       |                   |                     |
   Proxmox              VPS API              Proxy Source
   OVH/Hetzner          Manual Source        Upstream Provider
```

---

## 4. Lớp portal

### 4.1 Admin Portal

Dành cho platform owner/team vận hành.

Phạm vi:
- Quản tenant/reseller.
- Quản master catalog.
- Quản provider/source.
- Quản top-up reseller và client trực tiếp.
- Xem toàn hệ thống, revenue, provider health, failed provisioning.
- Emergency access có reason và audit.

Không nên để Admin Portal dùng cùng routing logic với Client Portal mà thiếu guard. Admin có quyền cao nhưng phải bị audit nặng hơn.

### 4.2 Reseller Portal

Dành cho reseller owner/staff.

Phạm vi:
- Quản branding/storefront riêng.
- Quản client thuộc tenant của reseller.
- Clone sản phẩm từ master catalog.
- Set giá bán riêng.
- Duyệt top-up client nội bộ.
- Xem settlement cost/profit.
- Quản ticket/support/notification của tenant.

Reseller không được thấy:
```text
client của reseller khác
provider API key
master cost ngoài quyền được hiển thị
tenant data của platform
ledger platform-level không liên quan
```

### 4.3 Client Portal

Dành cho khách mua VPS/proxy.

Phạm vi:
- Đăng ký/đăng nhập.
- Nạp ví theo tenant.
- Mua/gia hạn dịch vụ.
- Xem dịch vụ, credential masked/reveal.
- Gửi support/request cancel.
- Xem invoice/transaction của chính mình.

Client không bao giờ được truyền tenant_id để tự chọn tenant. Tenant được xác định từ domain/session.

---

## 5. API layer

API layer chịu trách nhiệm:
- Xác thực user.
- Xác định tenant context.
- Kiểm tra RBAC/permission.
- Validate request.
- Gắn correlation_id/request_id.
- Ghi audit event cho action nhạy cảm.
- Không chứa secret provider trong response.

Mỗi request nghiệp vụ nên có context nội bộ:

```text
request_context:
- request_id
- actor_id
- actor_role
- actor_tenant_id
- actor_permissions
- domain_tenant_id
- effective_tenant_id
- is_emergency_access
- emergency_reason
- ip_address
- user_agent
```

Rule:
```text
effective_tenant_id phải do server xác định.
Nếu actor là platform admin, effective_tenant_id có thể là target tenant nhưng phải audit.
Nếu actor là reseller/client, effective_tenant_id phải bằng tenant của actor/domain.
```

---

## 6. Core domain modules

### 6.1 Catalog Module

Nhiệm vụ:
- Quản master product/plan/source.
- Clone catalog sang tenant.
- Lưu version/snapshot.
- Kiểm soát margin floor.
- Ẩn plan/source disabled/out-of-stock.

Không xử lý tiền trực tiếp. Catalog chỉ cung cấp giá/snapshot/capability để Checkout Module dùng.

### 6.2 Order & Checkout Module

Nhiệm vụ:
- Tạo order.
- Validate plan/source/capability.
- Reserve inventory.
- Debit wallet/settlement.
- Tạo provisioning job.
- Gắn order_status, payment_status, provisioning_status.

Checkout phải chạy theo transaction boundary rõ:
```text
validate -> reserve -> debit ledger -> create order/provisioning job
```

Nếu không thể hoàn thành một bước P0 thì phải rollback hoặc chuyển trạng thái manual_review có kiểm soát.

### 6.3 Inventory & Reservation Module

Nhiệm vụ:
- Quản capacity.
- Atomic reserve/release/allocate.
- Chống oversell.
- Quản reservation expiry.

Không cho provider worker tự trừ/tăng stock trực tiếp ngoài module này.

### 6.4 Service Lifecycle Module

Nhiệm vụ:
- Tạo service instance sau provisioning success.
- Quản active/suspended/expired/terminated.
- Quản renew/cancel/refund lifecycle.
- Gắn lifecycle event.

Mọi thay đổi service_status phải có transition hợp lệ và audit.

### 6.5 Provisioning Module

Nhiệm vụ:
- Tạo provisioning job.
- Gọi provider adapter qua worker.
- Xử lý retry/partial success/manual review.
- Gắn external_resource_id.
- Đồng bộ status từ provider.

API request của user không nên gọi provider trực tiếp trong cùng request lâu. Nên tạo job và worker xử lý.

---

## 7. Financial Core

### 7.1 Wallet Module

Nhiệm vụ:
- Quản wallet theo user/tenant/reseller.
- Tính available balance từ ledger.
- Không sửa balance bằng tay ngoài adjustment có audit.
- Phân biệt client wallet và reseller settlement wallet.

### 7.2 Ledger Module

Ledger entry là bất biến.

Thông tin tối thiểu:
```text
ledger_id
wallet_id
tenant_id
actor_id
entry_type
direction
amount
currency
balance_after
reference_type
reference_id
correlation_id
created_at
```

Không update/delete ledger entry sau khi posted.

### 7.3 Settlement Module

Nhiệm vụ:
- Khi client thuộc reseller mua hàng, ghi nhận selling_price, reseller_cost, margin.
- Debit đúng wallet.
- Tạo report profit theo snapshot.
- Xử lý refund/rollback theo policy.

---

## 8. Security Core

### 8.1 Credential Service

Nhiệm vụ:
- Encrypt credential at rest.
- Mask credential khi hiển thị.
- Ghi audit khi reveal.
- Không log plaintext.
- Quản rotation/reissue nếu provider hỗ trợ.

Các field nhạy cảm:
```text
provider_api_key
provider_secret
vps_root_password
proxy_username
proxy_password
ssh_private_key
console_url_token
```

### 8.2 Audit Service

Audit service nhận event từ API và worker.

Audit event phải có:
```text
actor
tenant
target
action
before/after redacted
request_id/correlation_id
ip/user_agent nếu từ request
worker_job_id nếu từ worker
```

### 8.3 Risk/Abuse Service

Phase 1 chưa cần scoring phức tạp, nhưng phải có:
- risk flag.
- abuse case.
- manual review.
- suspend reason.
- blacklist tối thiểu.

---

## 9. Queue và worker

Các thao tác nên đi qua queue:
- Provision VPS/proxy.
- Retry provisioning.
- Sync provider resource.
- Send notification.
- Generate report/export.
- Expiry reminder.
- Suspend/terminate due to expiry.
- Provider health check.

Queue message phải có:
```text
job_id
job_type
tenant_id
reference_type
reference_id
idempotency_key
attempt_count
correlation_id
created_at
```

Worker phải ghi:
```text
started_at
finished_at
status
provider_response_summary redacted
error_code
next_retry_at
manual_review_required
```

---

## 10. Cron/Scheduler

Các job định kỳ P0:
- reservation_expiry_job.
- service_expiry_job.
- suspension_job.
- termination_job.
- provider_health_check_job.
- provider_inventory_sync_job.
- renewal_reminder_job.
- pending_topup_review_reminder_job.
- audit_retention_job.

Cron job phải idempotent: chạy lại không được làm sai số dư, stock hoặc trạng thái.

---

## 11. Data boundary

### 11.1 Bảng thuộc tenant

Các bảng này bắt buộc có `tenant_id`:
```text
users
wallets
orders
order_items
reservations
services
service_credentials
tenant_products
tenant_plans
topup_requests
support_tickets
audit_logs
risk_flags
abuse_cases
notifications
```

### 11.2 Bảng platform-level

Các bảng platform-level:
```text
master_products
master_plans
provider_sources
provider_accounts
platform_settings
exchange_rates
system_jobs
```

Bảng platform-level vẫn cần audit khi sửa.

---

## 12. Observability

Mọi request/job quan trọng cần có trace:

```text
correlation_id:
topup -> ledger -> order -> reservation -> provisioning_job -> provider_request -> service -> notification
```

Log không được chứa:
```text
plaintext credential
provider secret
payment proof sensitive image raw URL nếu private
full token/session
```

Metric tối thiểu:
- checkout success rate.
- provisioning success/fail/manual_review rate.
- provider API latency/error rate.
- reservation expired count.
- wallet adjustment count.
- credential reveal count.
- abuse suspension count.
- top-up approval time.
- reseller low balance count.

---

## 13. Deployment logical units

Có thể triển khai mono-repo hoặc multi-service, nhưng logical unit nên giữ rõ:

```text
web_frontend
backend_api
worker
scheduler
database
cache/queue
object_storage
logging/monitoring
secret_manager
```

Không bắt buộc tách microservice phase 1. Một monolith có module boundary rõ + worker riêng thường thực dụng hơn.

Khuyến nghị phase 1:
```text
modular monolith + queue worker + scheduler
```

Lý do:
- Nhanh ra hàng.
- Ít overhead vận hành.
- Dễ giữ transaction billing/order/inventory.
- Sau này tách module khi tải thật sự lớn.

---

## 14. Nguyên tắc thất bại an toàn

Nếu hệ thống không chắc chắn, ưu tiên:
```text
không provision hơn provision sai
manual_review hơn retry mù
ẩn credential hơn lộ credential
không trừ tiền hơn trừ tiền không trace được
suspend có reason hơn terminate vội
```

Câu nền: **hệ thống hạ tầng không cần “ảo thuật nhanh”; nó cần không làm sai khi provider, payment, network và người dùng cùng gây nhiễu.**

---


# ===== 15_Database_Schema_And_ERD.md =====

# 15 - Database Schema And ERD

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa schema logic cho nền tảng VPS/Proxy. Đây là **data contract**, không phải migration code.

Mục tiêu:
- Khóa các bảng chính.
- Khóa quan hệ tenant/user/order/wallet/service/provider.
- Khóa các field bắt buộc, constraint, index, encryption, audit.
- Giảm rủi ro sai tiền, sai tenant, sai stock, cấp trùng tài nguyên.
- Làm nền cho API, worker, QA và reporting.

---

## 2. Nguyên tắc database bắt buộc

### 2.1 Tenant scope

Mọi bảng chứa dữ liệu thuộc khách/tenant phải có:

```text
tenant_id
created_at
updated_at
```

Các query của reseller/client/staff phải scope theo `tenant_id`.

Không dùng `tenant_id` từ request body để ghi DB. Tenant phải đến từ request context server-side.

### 2.2 Immutable financial ledger

Bảng ledger không được update/delete sau khi entry ở trạng thái `posted`.

Nếu sai, tạo entry điều chỉnh:
```text
adjustment_credit
adjustment_debit
refund
reversal
```

### 2.3 Snapshot cho giao dịch

Order/service phải lưu snapshot để không bị ảnh hưởng bởi thay đổi giá/policy tương lai.

Snapshot tối thiểu:
```text
product_snapshot
plan_snapshot
price_snapshot
billing_cycle_snapshot
capability_snapshot
provider_source_snapshot
reseller_cost_snapshot
fx_snapshot
policy_snapshot
```

### 2.4 Credential encryption

Mọi credential phải encrypt at rest.

Không lưu plaintext trong:
```text
audit_logs
worker_logs
provider_requests raw response
notifications
support tickets public notes
```

### 2.5 Idempotency uniqueness

Các thao tác có khả năng retry phải có unique key:
```text
checkout_idempotency_key
provisioning_idempotency_key
provider_request_idempotency_key
payment_webhook_idempotency_key
```

---

## 3. ERD logic tổng quan

```text
tenants
  ├── tenant_domains
  ├── users
  │     ├── user_roles
  │     ├── wallets
  │     │     └── wallet_ledger_entries
  │     ├── topup_requests
  │     └── orders
  │           ├── order_items
  │           ├── reservations
  │           └── services
  │                 ├── service_credentials
  │                 └── service_lifecycle_events
  ├── tenant_products
  │     └── tenant_plans
  ├── notifications
  ├── audit_logs
  ├── risk_flags
  └── abuse_cases

master_products
  └── master_plans
        └── plan_sources
              └── provider_sources
                    ├── provider_accounts
                    ├── provider_inventory
                    └── provider_resource_mappings

provisioning_jobs
  └── provider_requests
```

---

## 4. Core tenancy tables

### 4.1 `tenants`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| tenant_id | id | yes | Primary key |
| parent_tenant_id | id/null | no | Null với platform/root, set với reseller hierarchy nếu dùng |
| tenant_type | enum | yes | `platform`, `reseller`, `direct_store` |
| name | string | yes | Tên hiển thị |
| slug | string | yes | Unique |
| status | enum | yes | `active`, `suspended`, `disabled`, `pending_setup` |
| default_currency | string | yes | Ví dụ USD, VND |
| timezone | string | yes | Mặc định Asia/Ho_Chi_Minh nếu không có |
| owner_user_id | id/null | no | Owner chính |
| branding_settings | json | no | Logo, màu, footer, support links |
| billing_settings | json | no | Minimum balance, reseller threshold |
| risk_settings | json | no | Manual review thresholds |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(slug)
tenant_type in allowed enum
status in allowed enum
```

Index:
```text
tenant_type
status
owner_user_id
```

### 4.2 `tenant_domains`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| domain_id | id | yes | Primary key |
| tenant_id | id | yes | FK tenants |
| domain | string | yes | Lowercase normalized |
| domain_type | enum | yes | `system_subdomain`, `custom_domain` |
| verification_status | enum | yes | `pending`, `verified`, `failed`, `disabled` |
| verification_token_hash | string | no | Không lưu token plaintext nếu không cần |
| tls_status | enum | yes | `pending`, `active`, `failed`, `expired` |
| is_primary | bool | yes |  |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(domain)
only one primary domain per tenant
```

Security:
```text
domain -> tenant_id mapping là nguồn xác định tenant context cho storefront.
```

---

## 5. User, role, permission tables

### 5.1 `users`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| user_id | id | yes | Primary key |
| tenant_id | id | yes | FK tenants |
| email | string | yes | Normalize lowercase |
| email_verified_at | datetime/null | no |  |
| password_hash | string | yes | Không bao giờ log |
| full_name | string | no |  |
| user_type | enum | yes | `platform_staff`, `reseller_staff`, `client` |
| status | enum | yes | `active`, `suspended`, `disabled`, `pending_verification` |
| two_factor_status | enum | yes | `required`, `enabled`, `disabled` |
| last_login_at | datetime/null | no |  |
| failed_login_count | number | yes | rate limit |
| created_at / updated_at | datetime | yes |  |

Constraints:
```text
unique(tenant_id, email)
```

Important:
```text
Cùng một email có thể tồn tại ở nhiều tenant nếu business cho phép.
Nếu muốn global identity, cần bảng identities riêng. Phase 1 có thể giữ email unique theo tenant để đơn giản.
```

### 5.2 `roles`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| role_id | id | yes | Primary key |
| tenant_id | id/null | no | Null = system role, set = custom tenant role |
| role_key | string | yes | `platform_super_admin`, `reseller_owner`, `client` |
| name | string | yes |  |
| is_system | bool | yes |  |
| created_at / updated_at | datetime | yes |  |

### 5.3 `permissions`

| Field | Type logic | Required | Note |
|---|---|---:|---|
| permission_id | id | yes |  |
| permission_key | string | yes | Ví dụ `wallet.topup.approve` |
| module | string | yes | `wallet`, `catalog`, `service` |
| risk_level | enum | yes | `low`, `medium`, `high`, `critical` |

### 5.4 `role_permissions`

```text
role_id
permission_id
created_at
```

Constraint:
```text
unique(role_id, permission_id)
```

### 5.5 `user_roles`

```text
user_id
tenant_id
role_id
created_at
```

Constraint:
```text
unique(user_id, tenant_id, role_id)
```

---

## 6. Catalog tables

### 6.1 `master_products`

| Field | Required | Note |
|---|---:|---|
| product_id | yes | Primary key |
| product_type | yes | `vps`, `proxy`, `service_addon` |
| name | yes | Master name |
| description | no |  |
| status | yes | `draft`, `active`, `disabled`, `archived` |
| display_order | no |  |
| created_by | yes | Admin user |
| created_at / updated_at | yes |  |

### 6.2 `master_plans`

| Field | Required | Note |
|---|---:|---|
| plan_id | yes | Primary key |
| product_id | yes | FK master_products |
| plan_code | yes | Internal SKU |
| name | yes |  |
| specs | yes | json: CPU/RAM/SSD/location/bandwidth/IP |
| billing_cycle_type | yes | `day`, `month_30d`, `calendar_month`, `custom` |
| billing_cycle_value | yes | number |
| base_cost | yes | platform internal cost |
| suggested_price | yes |  |
| reseller_min_price | no | margin guard |
| status | yes | `active`, `disabled`, `archived` |
| version | yes | increment khi đổi giá/spec/policy |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(plan_code, version)
```

### 6.3 `provider_sources`

| Field | Required | Note |
|---|---:|---|
| source_id | yes | Primary key |
| source_type | yes | `proxmox`, `ovh`, `hetzner`, `manual`, `proxy_upstream`, `custom_api` |
| name | yes |  |
| provider_account_id | no | FK nếu có |
| location | no | country/region/datacenter |
| status | yes | `active`, `disabled`, `maintenance`, `out_of_stock` |
| capability_profile | yes | json |
| inventory_mode | yes | `finite_stock`, `provider_live`, `manual_unlimited`, `preloaded_list` |
| risk_level | yes | `low`, `medium`, `high` |
| created_at / updated_at | yes |  |

### 6.4 `plan_sources`

| Field | Required | Note |
|---|---:|---|
| plan_source_id | yes |  |
| plan_id | yes | FK master_plans |
| source_id | yes | FK provider_sources |
| priority | yes | Provider priority |
| cost_override | no | Nếu source có cost riêng |
| capacity_policy | yes | json |
| status | yes | `active`, `disabled` |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(plan_id, source_id)
```

### 6.5 `tenant_products`

| Field | Required | Note |
|---|---:|---|
| tenant_product_id | yes | Primary key |
| tenant_id | yes | FK tenants |
| master_product_id | yes | FK master_products |
| name_override | no |  |
| description_override | no |  |
| status | yes | `active`, `hidden`, `disabled` |
| clone_version | yes | version của master lúc clone/sync |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(tenant_id, master_product_id)
```

### 6.6 `tenant_plans`

| Field | Required | Note |
|---|---:|---|
| tenant_plan_id | yes |  |
| tenant_id | yes |  |
| tenant_product_id | yes |  |
| master_plan_id | yes |  |
| selling_price | yes | Giá bán cho client tenant |
| reseller_cost | yes | Giá platform thu reseller |
| currency | yes |  |
| margin_policy | yes | json |
| visibility | yes | `public`, `hidden`, `private` |
| status | yes | `active`, `disabled`, `margin_risk`, `archived` |
| source_policy_snapshot | yes | json |
| plan_snapshot | yes | json |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(tenant_id, master_plan_id)
selling_price >= 0
reseller_cost >= 0
```

Index:
```text
tenant_id, status, visibility
```

---

## 7. Wallet and billing tables

### 7.1 `wallets`

| Field | Required | Note |
|---|---:|---|
| wallet_id | yes |  |
| tenant_id | yes | Owner tenant |
| owner_type | yes | `tenant`, `user`, `reseller_settlement`, `platform` |
| owner_id | yes | user_id hoặc tenant_id |
| currency | yes |  |
| status | yes | `active`, `frozen`, `closed` |
| available_balance_cache | yes | Cache, không là source of truth |
| locked_balance_cache | yes | Nếu dùng lock/pending |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(owner_type, owner_id, currency)
```

Important:
```text
available_balance_cache chỉ để đọc nhanh.
Source of truth là wallet_ledger_entries.
```

### 7.2 `wallet_ledger_entries`

| Field | Required | Note |
|---|---:|---|
| ledger_id | yes | Primary key |
| wallet_id | yes | FK wallets |
| tenant_id | yes | Scope |
| direction | yes | `credit`, `debit` |
| amount | yes | positive decimal |
| currency | yes |  |
| entry_type | yes | `topup`, `purchase`, `reseller_cost`, `refund`, `adjustment`, `reversal`, `commission`, `lock`, `unlock` |
| status | yes | `posted`, `voided_by_reversal` |
| balance_after | yes | Balance sau entry |
| reference_type | yes | `order`, `topup_request`, `service`, `manual_adjustment`, `settlement` |
| reference_id | yes |  |
| idempotency_key | yes | Unique theo wallet/action |
| created_by | no | actor |
| reason | no | required với adjustment |
| correlation_id | yes | trace |
| created_at | yes | immutable |

Constraints:
```text
unique(wallet_id, idempotency_key)
amount > 0
posted entry cannot be updated/deleted
```

Index:
```text
wallet_id, created_at
tenant_id, reference_type, reference_id
correlation_id
```

### 7.3 `topup_requests`

| Field | Required | Note |
|---|---:|---|
| topup_request_id | yes |  |
| tenant_id | yes |  |
| wallet_id | yes |  |
| requested_by | yes | user |
| amount | yes |  |
| currency | yes |  |
| payment_method | yes | `bank_transfer`, `crypto`, `manual`, `other` |
| payment_reference | no |  |
| proof_attachment_id | no | private file |
| status | yes | `draft`, `submitted`, `under_review`, `approved`, `rejected`, `expired`, `cancelled` |
| reviewed_by | no |  |
| reviewed_at | no |  |
| review_note | no |  |
| ledger_id | no | set when approved |
| created_at / updated_at | yes |  |

Constraint:
```text
approved topup must have ledger_id
```

---

## 8. Order, reservation, service tables

### 8.1 `orders`

| Field | Required | Note |
|---|---:|---|
| order_id | yes | Primary key |
| tenant_id | yes | Tenant selling to client |
| client_user_id | yes | Buyer |
| seller_tenant_id | yes | Reseller tenant hoặc platform tenant |
| order_number | yes | Human readable unique |
| order_type | yes | `new_service`, `renewal`, `upgrade`, `addon` |
| order_status | yes | `draft`, `pending_payment`, `paid`, `provisioning`, `active`, `failed`, `manual_review`, `cancelled`, `refunded`, `expired` |
| payment_status | yes | `unpaid`, `paid`, `partially_refunded`, `refunded`, `failed` |
| provisioning_status | yes | `not_started`, `queued`, `running`, `success`, `failed`, `partial_success`, `manual_review` |
| subtotal | yes |  |
| discount_amount | yes |  |
| total_amount | yes | selling price total |
| currency | yes |  |
| client_wallet_debit_ledger_id | no |  |
| reseller_wallet_debit_ledger_id | no |  |
| idempotency_key | yes | Unique per buyer checkout |
| correlation_id | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(tenant_id, order_number)
unique(client_user_id, idempotency_key)
```

### 8.2 `order_items`

| Field | Required | Note |
|---|---:|---|
| order_item_id | yes |  |
| order_id | yes | FK orders |
| tenant_id | yes |  |
| tenant_plan_id | yes |  |
| master_plan_id | yes |  |
| quantity | yes | phase 1 thường = 1 |
| unit_price | yes | selling price snapshot |
| reseller_unit_cost | yes | reseller cost snapshot |
| billing_cycle_snapshot | yes | json |
| product_snapshot | yes | json |
| plan_snapshot | yes | json |
| capability_snapshot | yes | json |
| provider_source_snapshot | yes | json |
| policy_snapshot | yes | json |
| created_at | yes |  |

### 8.3 `reservations`

| Field | Required | Note |
|---|---:|---|
| reservation_id | yes |  |
| tenant_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| source_id | yes | provider source |
| quantity | yes |  |
| status | yes | `reserved`, `allocated`, `released`, `expired`, `failed` |
| reserved_at | yes |  |
| expires_at | yes | thường now + 5 phút |
| allocated_at | no |  |
| released_at | no |  |
| idempotency_key | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(source_id, idempotency_key)
allocated reservation cannot expire/release again
```

### 8.4 `services`

| Field | Required | Note |
|---|---:|---|
| service_id | yes |  |
| tenant_id | yes |  |
| client_user_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| reservation_id | yes |  |
| service_type | yes | `vps`, `proxy` |
| service_status | yes | `pending`, `active`, `suspended`, `expired`, `terminated`, `cancelled`, `failed` |
| billing_status | yes | `paid`, `due_soon`, `overdue`, `grace`, `cancelled` |
| suspension_reason | no | required if suspended |
| source_id | yes | provider source |
| external_resource_id | no | Provider resource id |
| service_identifier | no | IP/hostname/proxy endpoint |
| term_start_at | yes |  |
| term_end_at | yes |  |
| grace_until_at | no |  |
| terminate_after_at | no |  |
| capability_snapshot | yes | json |
| plan_snapshot | yes | json |
| price_snapshot | yes | json |
| created_at / updated_at | yes |  |

Index:
```text
tenant_id, client_user_id, service_status
term_end_at
source_id, external_resource_id
```

### 8.5 `service_credentials`

| Field | Required | Note |
|---|---:|---|
| credential_id | yes |  |
| service_id | yes | FK services |
| tenant_id | yes |  |
| credential_type | yes | `vps_root`, `proxy_auth`, `ssh_key`, `console_url`, `api_token` |
| encrypted_payload | yes | encrypted JSON |
| secret_version | yes | key rotation support |
| masked_hint | no | ví dụ `root / ********` |
| last_revealed_at | no |  |
| last_revealed_by | no | user |
| status | yes | `active`, `rotated`, `revoked` |
| created_at / updated_at | yes |  |

Security:
```text
Không lưu plaintext ở bất kỳ field nào.
Mỗi reveal phải ghi audit credential.revealed.
```

### 8.6 `service_lifecycle_events`

| Field | Required | Note |
|---|---:|---|
| event_id | yes |  |
| tenant_id | yes |  |
| service_id | yes |  |
| event_type | yes | `created`, `activated`, `renewed`, `expired`, `suspended`, `unsuspended`, `terminated`, `cancelled`, `failed` |
| from_status | no |  |
| to_status | no |  |
| reason | no | required with suspend/terminate |
| actor_id | no | user or system |
| job_id | no | if worker/cron |
| correlation_id | yes |  |
| created_at | yes |  |

---

## 9. Provider and provisioning tables

### 9.1 `provider_accounts`

| Field | Required | Note |
|---|---:|---|
| provider_account_id | yes |  |
| name | yes |  |
| provider_type | yes | `proxmox`, `ovh`, `hetzner`, `manual`, `proxy_upstream` |
| status | yes | `active`, `disabled`, `maintenance` |
| encrypted_credentials | yes | Provider API credentials |
| allowed_ips | no | API allowlist if used |
| health_status | yes | `unknown`, `healthy`, `degraded`, `down` |
| last_health_check_at | no |  |
| created_at / updated_at | yes |  |

Security:
```text
provider credentials encrypt at rest.
Không trả về qua API thông thường.
```

### 9.2 `provider_inventory`

| Field | Required | Note |
|---|---:|---|
| inventory_id | yes |  |
| source_id | yes |  |
| capacity_total | no | null nếu live/unlimited |
| reserved_count | yes |  |
| allocated_count | yes |  |
| available_count_cache | yes | derived/cache |
| last_synced_at | no |  |
| status | yes | `active`, `out_of_stock`, `sync_failed`, `disabled` |
| updated_at | yes |  |

Rule:
```text
available = capacity_total - reserved_count - allocated_count
Nếu capacity_total null và inventory_mode là provider_live thì phải check live trước reserve.
```

### 9.3 `provisioning_jobs`

| Field | Required | Note |
|---|---:|---|
| job_id | yes |  |
| tenant_id | yes |  |
| order_id | yes |  |
| order_item_id | yes |  |
| reservation_id | yes |  |
| source_id | yes |  |
| job_type | yes | `provision`, `renew`, `suspend`, `unsuspend`, `terminate`, `sync`, `reset_password`, `change_ip` |
| status | yes | `queued`, `running`, `success`, `failed`, `partial_success`, `manual_review`, `cancelled` |
| idempotency_key | yes | unique |
| attempt_count | yes |  |
| max_attempts | yes |  |
| next_retry_at | no |  |
| last_error_code | no |  |
| last_error_message_redacted | no |  |
| manual_review_reason | no |  |
| correlation_id | yes |  |
| created_at / updated_at | yes |  |

Constraints:
```text
unique(source_id, job_type, idempotency_key)
```

### 9.4 `provider_requests`

| Field | Required | Note |
|---|---:|---|
| provider_request_id | yes |  |
| job_id | yes | FK provisioning_jobs |
| source_id | yes |  |
| request_type | yes | same as job/action |
| external_request_id | no | nếu provider có |
| external_resource_id | no | nếu biết |
| status | yes | `sent`, `success`, `failed`, `timeout`, `unknown`, `manual_review` |
| retry_safety | yes | `safe_retry`, `unsafe_retry`, `do_not_retry`, `manual_review_required` |
| request_summary_redacted | no | không chứa secret |
| response_summary_redacted | no | không chứa secret |
| error_code | no |  |
| sent_at | yes |  |
| received_at | no |  |
| correlation_id | yes |  |

### 9.5 `provider_resource_mappings`

| Field | Required | Note |
|---|---:|---|
| mapping_id | yes |  |
| tenant_id | yes |  |
| service_id | yes |  |
| source_id | yes |  |
| external_resource_id | yes |  |
| external_status | no |  |
| last_synced_at | no |  |
| created_at / updated_at | yes |  |

Constraint:
```text
unique(source_id, external_resource_id)
```

---

## 10. Audit, risk, abuse, notification tables

### 10.1 `audit_logs`

| Field | Required | Note |
|---|---:|---|
| audit_id | yes |  |
| tenant_id | no | null nếu platform-level |
| actor_id | no | null nếu system |
| actor_type | yes | `user`, `system`, `worker`, `provider_webhook` |
| action | yes | e.g. `wallet.topup.approved` |
| target_type | yes |  |
| target_id | yes |  |
| before_snapshot_redacted | no | json |
| after_snapshot_redacted | no | json |
| metadata_redacted | no | json |
| ip_address | no |  |
| user_agent | no |  |
| correlation_id | yes |  |
| created_at | yes | immutable |

Index:
```text
tenant_id, created_at
actor_id, created_at
target_type, target_id
correlation_id
action
```

### 10.2 `risk_flags`

| Field | Required | Note |
|---|---:|---|
| risk_flag_id | yes |  |
| tenant_id | yes |  |
| user_id | no |  |
| service_id | no |  |
| order_id | no |  |
| flag_type | yes | `new_account_high_value`, `payment_mismatch`, `abuse_history`, `manual_blacklist`, `provider_risk` |
| severity | yes | `low`, `medium`, `high`, `critical` |
| status | yes | `open`, `reviewing`, `cleared`, `confirmed` |
| note | no |  |
| created_by | no |  |
| created_at / updated_at | yes |  |

### 10.3 `abuse_cases`

| Field | Required | Note |
|---|---:|---|
| abuse_case_id | yes |  |
| tenant_id | yes |  |
| service_id | no |  |
| user_id | no |  |
| provider_source_id | no |  |
| case_type | yes | `spam`, `scan`, `bruteforce`, `copyright`, `fraud`, `aup_violation`, `other` |
| severity | yes | `low`, `medium`, `high`, `critical` |
| status | yes | `open`, `investigating`, `warning_sent`, `suspended`, `resolved`, `closed` |
| evidence_private | no | json/file refs |
| action_taken | no |  |
| created_at / updated_at | yes |  |

### 10.4 `notifications`

| Field | Required | Note |
|---|---:|---|
| notification_id | yes |  |
| tenant_id | yes |  |
| recipient_user_id | no |  |
| channel | yes | `email`, `dashboard`, `telegram`, `webhook` |
| template_key | yes |  |
| status | yes | `queued`, `sent`, `failed`, `cancelled` |
| payload_redacted | yes | json |
| reference_type | no |  |
| reference_id | no |  |
| sent_at | no |  |
| created_at / updated_at | yes |  |

---

## 11. Recommended enums

### 11.1 Order status

```text
draft
pending_payment
paid
provisioning
active
failed
manual_review
cancelled
refunded
expired
```

### 11.2 Provisioning status/job status

```text
not_started
queued
running
success
failed
partial_success
manual_review
cancelled
```

### 11.3 Service status

```text
pending
active
suspended
expired
terminated
cancelled
failed
```

### 11.4 Reservation status

```text
reserved
allocated
released
expired
failed
```

### 11.5 Wallet ledger entry type

```text
topup
purchase
reseller_cost
refund
adjustment
reversal
commission
lock
unlock
```

---

## 12. Critical indexes

P0 indexes:
```text
users(tenant_id, email)
tenant_domains(domain)
tenant_plans(tenant_id, status, visibility)
wallets(owner_type, owner_id, currency)
wallet_ledger_entries(wallet_id, created_at)
wallet_ledger_entries(correlation_id)
orders(tenant_id, client_user_id, created_at)
orders(tenant_id, order_number)
reservations(source_id, status, expires_at)
services(tenant_id, client_user_id, service_status)
services(term_end_at, service_status)
provisioning_jobs(status, next_retry_at)
provisioning_jobs(correlation_id)
provider_resource_mappings(source_id, external_resource_id)
audit_logs(tenant_id, created_at)
audit_logs(correlation_id)
```

---

## 13. Data retention

Suggested phase 1:
```text
wallet_ledger_entries: keep forever
orders/services: keep forever or legal retention
audit_logs financial/security: keep long-term
provider_request detailed redacted body: keep 90-180 days
notification payload: keep 30-90 days
temporary payment proof files: follow finance retention policy
```

Không xóa cứng dữ liệu tài chính. Dùng status/archival.

---

## 14. Acceptance criteria database

Schema được xem là đạt khi:

- Mọi bảng tenant-owned có `tenant_id`.
- Ledger posted không thể sửa/xóa ở tầng application.
- Credential/provider secret không có plaintext field.
- Checkout có thể trace bằng `correlation_id`.
- Provisioning có unique idempotency key.
- Provider external resource có unique mapping.
- Reservation không thể allocate/release hai lần.
- Service/order lưu đủ snapshot để xử lý tranh chấp sau này.
- Có index cho các query P0: wallet, order, service, provisioning, audit, expiry.
- Admin emergency access để lại audit đầy đủ.

Câu nền: **database không chỉ lưu dữ liệu; nó là hàng rào chống con người, provider và bug làm sai tiền.**

---


# ===== 16_API_Contract_And_Permission_Spec.md =====

# 16 - API Contract And Permission Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa API contract logic cho backend/frontend. Đây là hợp đồng hành vi, không phải code.

Mỗi API cần rõ:
- Ai được gọi.
- Tenant scope lấy từ đâu.
- Request/response cần gì.
- Validate gì.
- Ghi audit gì.
- Error code nào trả về.
- Rate limit nào áp dụng.
- Có cần idempotency hay không.

---

## 2. API conventions

### 2.1 Base principles

```text
Không tin tenant_id từ body.
Không trả plaintext secret nếu không qua reveal action.
Không trả provider API key qua API thường.
Không cho action nếu capability snapshot không hỗ trợ.
Không tạo tài nguyên nếu chưa debit/settlement hợp lệ.
```

### 2.2 Request headers logic

Frontend nên gửi:
```text
Authorization
Idempotency-Key với action tạo giao dịch
X-Request-Id nếu có
```

Backend tự gắn:
```text
correlation_id
effective_tenant_id
actor_context
domain_context
```

### 2.3 Standard response

Thành công:
```text
{
  "success": true,
  "data": {},
  "request_id": "...",
  "correlation_id": "..."
}
```

Lỗi:
```text
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Human readable message",
    "details": {}
  },
  "request_id": "...",
  "correlation_id": "..."
}
```

### 2.4 Pagination

List API dùng:
```text
page
page_size
sort
filter
```

Response:
```text
items
pagination.total
pagination.page
pagination.page_size
```

### 2.5 Standard error codes

| Code | Meaning |
|---|---|
| `UNAUTHENTICATED` | Chưa đăng nhập |
| `FORBIDDEN` | Không có quyền |
| `TENANT_MISMATCH` | Resource không thuộc tenant |
| `VALIDATION_ERROR` | Sai input |
| `RESOURCE_NOT_FOUND` | Không thấy hoặc không được thấy resource |
| `RATE_LIMITED` | Vượt giới hạn |
| `IDEMPOTENCY_CONFLICT` | Idempotency key đã dùng cho payload khác |
| `INSUFFICIENT_CLIENT_BALANCE` | Ví client thiếu tiền |
| `INSUFFICIENT_RESELLER_BALANCE` | Ví reseller settlement thiếu tiền |
| `OUT_OF_STOCK` | Hết stock |
| `PLAN_DISABLED` | Plan bị tắt |
| `PROVIDER_UNAVAILABLE` | Provider/source không sẵn sàng |
| `CAPABILITY_NOT_SUPPORTED` | Action không được source hỗ trợ |
| `PROVISIONING_MANUAL_REVIEW` | Cần review thủ công |
| `CREDENTIAL_REVEAL_DENIED` | Không được xem credential |

---

## 3. Permission model per API

Mỗi endpoint phải khai báo:
```text
required_permission
allowed_roles
tenant_scope
risk_level
audit_action
```

Risk level:
```text
low: đọc list thông thường
medium: tạo order, cập nhật profile
high: duyệt top-up, suspend service, reveal credential
critical: chỉnh ledger, đổi provider, emergency access, terminate service
```

Critical action phải có:
```text
2FA nếu role yêu cầu
audit bắt buộc
reason nếu manual/admin action
correlation_id
```

---

## 4. Auth APIs

### 4.1 Register client

```text
POST /auth/register
```

Allowed:
```text
public on tenant storefront
```

Tenant scope:
```text
derived from domain
```

Request:
```text
email
password
full_name
accept_terms
```

Validation:
- domain phải map tenant active.
- email chưa tồn tại trong tenant.
- password đạt policy.
- accept_terms = true.

Response:
```text
user_id
email
status
email_verification_required
```

Audit:
```text
auth.registered
```

Errors:
```text
TENANT_NOT_ACTIVE
EMAIL_ALREADY_EXISTS
WEAK_PASSWORD
TERMS_NOT_ACCEPTED
```

### 4.2 Login

```text
POST /auth/login
```

Validation:
- tenant/domain hợp lệ.
- user active.
- password đúng.
- 2FA nếu required/enabled.

Audit:
```text
auth.login.success
auth.login.failed
auth.2fa.challenge
```

Rate limit:
```text
per IP + per email + per tenant
```

### 4.3 Logout

```text
POST /auth/logout
```

Audit:
```text
auth.logout
```

---

## 5. Tenant and domain APIs

### 5.1 Admin create reseller tenant

```text
POST /admin/tenants
```

Allowed:
```text
platform_super_admin
platform_staff with tenant.create
```

Request:
```text
name
slug
owner_email
default_currency
timezone
initial_status
```

Validation:
- slug unique.
- owner email valid.
- actor has permission.

Audit:
```text
tenant.created
```

### 5.2 Reseller update branding

```text
PATCH /reseller/tenant/branding
```

Allowed:
```text
reseller_owner
reseller_staff with tenant.branding.update
```

Request:
```text
logo_url
brand_color
support_email
support_telegram
footer_text
```

Validation:
- tenant active.
- assets allowed.
- no unsafe links.

Audit:
```text
tenant.branding.updated
```

### 5.3 Add custom domain

```text
POST /reseller/tenant/domains
```

Request:
```text
domain
```

Validation:
- domain chưa thuộc tenant khác.
- domain hợp lệ.
- tạo verification token.

Response:
```text
domain_id
verification_record_type
verification_record_value
verification_status
```

Audit:
```text
tenant.domain.created
```

---

## 6. Catalog APIs

### 6.1 Admin create master product

```text
POST /admin/catalog/products
```

Allowed:
```text
catalog.product.create
```

Request:
```text
product_type
name
description
status
display_order
```

Audit:
```text
catalog.product.created
```

### 6.2 Admin create master plan

```text
POST /admin/catalog/plans
```

Request:
```text
product_id
plan_code
name
specs
billing_cycle_type
billing_cycle_value
base_cost
suggested_price
reseller_min_price
status
```

Validation:
- product active/draft.
- billing cycle valid.
- prices >= 0.
- plan_code unique for version.

Audit:
```text
catalog.plan.created
```

### 6.3 Reseller list available master plans

```text
GET /reseller/catalog/master-plans
```

Allowed:
```text
reseller_owner
reseller_staff with catalog.view
```

Response:
```text
plans visible to reseller
suggested_price
reseller_cost
margin hints
capabilities
stock status
```

Do not expose:
```text
provider_api_key
sensitive provider account detail
```

### 6.4 Reseller clone/sync plan

```text
POST /reseller/catalog/plans/clone
```

Request:
```text
master_plan_id
selling_price
visibility
```

Validation:
- master plan active.
- selling_price >= reseller_cost or margin policy allows warning.
- tenant allowed to sell this product.
- source available.

Audit:
```text
catalog.tenant_plan.cloned
```

### 6.5 Client list catalog

```text
GET /client/catalog
```

Allowed:
```text
public or authenticated depending tenant setting
```

Tenant scope:
```text
domain/session tenant
```

Response:
```text
products
plans
selling_price
public specs
stock availability summary
```

Do not expose:
```text
reseller_cost
provider source internal ID unless safe alias
internal margin
```

---

## 7. Wallet APIs

### 7.1 Client get wallet

```text
GET /client/wallet
```

Allowed:
```text
client
```

Response:
```text
currency
available_balance
locked_balance
recent_ledger_entries
pending_topups
```

Audit:
```text
optional wallet.viewed
```

### 7.2 Submit top-up request

```text
POST /client/wallet/topups
```

Request:
```text
amount
currency
payment_method
payment_reference
proof_attachment_id
```

Validation:
- amount >= min topup.
- currency allowed.
- method enabled.
- proof required if method requires proof.

Status:
```text
submitted
```

Audit:
```text
wallet.topup.submitted
```

### 7.3 Approve client top-up

```text
POST /reseller/wallet/topups/{topup_request_id}/approve
```

Allowed:
```text
reseller_owner
reseller_staff with wallet.topup.approve
platform_admin for emergency/support
```

Validation:
- top-up belongs to actor tenant.
- status submitted/under_review.
- amount/currency match actual payment.
- reviewer has 2FA if required.

Result:
- create ledger credit to client wallet.
- set topup_request approved.
- notify client.

Audit:
```text
wallet.topup.approved
wallet.ledger.posted
```

Errors:
```text
TOPUP_ALREADY_REVIEWED
PERMISSION_DENIED
```

### 7.4 Admin approve reseller top-up

```text
POST /admin/resellers/{tenant_id}/wallet/topups/{topup_request_id}/approve
```

Allowed:
```text
platform_super_admin
finance_agent
```

Result:
- credit reseller settlement wallet.
- audit high risk.

---

## 8. Checkout and order APIs

### 8.1 Create checkout order

```text
POST /client/orders
```

Allowed:
```text
client
```

Idempotency:
```text
required
```

Request:
```text
tenant_plan_id
quantity
billing_cycle
coupon_code optional
```

Validation:
- user active.
- tenant active.
- tenant_plan belongs to tenant.
- plan visible and active.
- source active.
- stock available.
- client wallet >= selling_price.
- if tenant is reseller: reseller settlement wallet >= reseller_cost.
- abuse/risk check passes or routes manual_review.
- idempotency key not reused with different payload.

Success flow:
```text
create order
create order_item snapshots
reserve inventory
debit client wallet
debit reseller wallet if reseller tenant
create provisioning_job
order_status = provisioning
provisioning_status = queued
```

Response:
```text
order_id
order_number
order_status
provisioning_status
estimated_activation_message
```

Audit:
```text
order.created
reservation.created
wallet.client.debited
wallet.reseller_cost.debited
provisioning.job.created
```

Errors:
```text
INSUFFICIENT_CLIENT_BALANCE
INSUFFICIENT_RESELLER_BALANCE
OUT_OF_STOCK
PLAN_DISABLED
PROVIDER_UNAVAILABLE
RISK_MANUAL_REVIEW_REQUIRED
```

### 8.2 Get order detail

```text
GET /client/orders/{order_id}
```

Scope:
```text
client can view own order only
reseller can view orders inside tenant
admin can view all with permission
```

Tenant mismatch returns:
```text
404 or FORBIDDEN depending security policy
```

### 8.3 Cancel pending order

```text
POST /client/orders/{order_id}/cancel
```

Allowed only if:
```text
order_status in draft, pending_payment, manual_review_before_provision
provisioning not started
reservation still reserved
```

Result:
```text
release reservation
reverse debit if any
order_status = cancelled
```

Audit:
```text
order.cancelled
reservation.released
wallet.reversal.posted
```

---

## 9. Service APIs

### 9.1 Client list services

```text
GET /client/services
```

Filters:
```text
status
service_type
expiring_before
```

Response:
```text
service_id
name/spec summary
status
expiry
public endpoint masked
actions available by capability_snapshot
```

### 9.2 Client service detail

```text
GET /client/services/{service_id}
```

Response:
```text
service info
billing info
capability actions
credential masked_hint only
lifecycle events public subset
```

Never return plaintext credential here.

### 9.3 Reveal credential

```text
POST /client/services/{service_id}/credentials/{credential_id}/reveal
```

Allowed:
```text
client owns service
reseller owner/staff with credential.reveal
platform admin with credential.reveal
```

Validation:
- service belongs to tenant.
- credential belongs to service.
- actor has permission.
- 2FA required for staff/admin/reseller owner if policy requires.
- rate limit reveal.

Response:
```text
plaintext credential once
masked hint
reveal_expires_message
```

Audit:
```text
credential.revealed
```

Security:
- response not cached.
- no plaintext in audit/log.
- optional require re-auth for high-risk roles.

### 9.4 Renew service

```text
POST /client/services/{service_id}/renew
```

Idempotency:
```text
required
```

Validation:
- service active/expired/grace depending policy.
- plan still renewable.
- wallet sufficient.
- reseller wallet sufficient if reseller tenant.
- renew cycle valid.
- not terminated.

Result:
```text
debit wallet(s)
extend term_end_at based on policy
create lifecycle event service.renewed
if provider requires renew call -> create provisioning_job type renew
```

Audit:
```text
service.renewed
wallet.client.debited
wallet.reseller_cost.debited
```

### 9.5 Cancel service / request termination

```text
POST /client/services/{service_id}/cancel-request
```

Request:
```text
reason
cancel_at: immediate/end_of_term
```

Validation:
- service belongs to client.
- policy allows cancel.
- if immediate termination affects refund, route review if required.

Audit:
```text
service.cancel_requested
```

### 9.6 Admin/reseller suspend service

```text
POST /admin/services/{service_id}/suspend
POST /reseller/services/{service_id}/suspend
```

Request:
```text
reason
notify_client
```

Validation:
- reason required.
- permission required.
- source supports suspend or internal suspension only.
- if provider action required, create job.

Audit:
```text
service.suspended
provisioning.job.created optional
```

---

## 10. Provider and provisioning APIs

### 10.1 Admin list provider sources

```text
GET /admin/provider-sources
```

Allowed:
```text
provider.view
```

Response:
```text
source status
health status
inventory summary
capability profile
last sync
```

No provider secrets.

### 10.2 Admin create/update provider source

```text
POST /admin/provider-sources
PATCH /admin/provider-sources/{source_id}
```

Allowed:
```text
provider.manage
```

Validation:
- provider account credentials stored encrypted.
- capability profile valid.
- inventory policy valid.

Audit:
```text
provider.source.created
provider.source.updated
```

### 10.3 Manual retry provisioning job

```text
POST /admin/provisioning-jobs/{job_id}/retry
```

Allowed:
```text
provisioning.retry
```

Validation:
- job status failed/manual_review.
- retry_safety != do_not_retry.
- if unsafe_retry, require explicit override reason and high permission.
- max attempts not exceeded unless override.

Audit:
```text
provisioning.job.retry_requested
```

### 10.4 Mark manual review resolved

```text
POST /admin/provisioning-jobs/{job_id}/resolve
```

Request:
```text
resolution_type: success_failed_cancelled_link_existing_resource
external_resource_id optional
service_id optional
note required
```

Validation:
- note required.
- if linking resource, external_resource_id unique.
- ledger/reservation/service consistency checked.

Audit:
```text
provisioning.manual_review.resolved
```

---

## 11. Audit and report APIs

### 11.1 List audit logs

```text
GET /admin/audit-logs
GET /reseller/audit-logs
```

Scope:
```text
admin can view platform/global with permission
reseller can view own tenant
client generally cannot view raw audit logs
```

Filters:
```text
actor
action
target_type
target_id
date range
correlation_id
```

Redaction:
```text
always redacted
```

### 11.2 Reseller profit report

```text
GET /reseller/reports/profit
```

Calculation:
```text
gross_revenue = sum(client selling_price posted)
platform_cost = sum(reseller_cost posted)
gross_profit = gross_revenue - platform_cost - refunds/adjustments
```

Must use:
```text
ledger/order snapshots, not current plan price
```

### 11.3 Admin provider health report

```text
GET /admin/reports/provider-health
```

Includes:
```text
source status
success rate
failure rate
manual_review count
avg provisioning time
out_of_stock count
```

---

## 12. Abuse and risk APIs

### 12.1 Create manual risk flag

```text
POST /admin/risk-flags
POST /reseller/risk-flags
```

Request:
```text
user_id/service_id/order_id
flag_type
severity
note
```

Audit:
```text
risk.flag.created
```

### 12.2 Create abuse case

```text
POST /admin/abuse-cases
```

Request:
```text
service_id
case_type
severity
evidence_private
recommended_action
```

Actions:
```text
warning
suspend
terminate
provider_notice
```

Audit:
```text
abuse.case.created
```

---

## 13. Rate limit matrix

| Action | Limit logic |
|---|---|
| login | per IP + email + tenant |
| register | per IP + domain |
| top-up submit | per user + tenant |
| checkout | per user + tenant + plan |
| credential reveal | per user + service |
| reset password/change IP | per service + source |
| admin retry provisioning | per admin + job |
| domain verification | per tenant + domain |
| support message | per user + ticket |

---

## 14. API acceptance criteria

API spec được xem là đạt khi:
- Endpoint nào cũng có role/permission rõ.
- Endpoint nào cũng có tenant scope rõ.
- Checkout có idempotency bắt buộc.
- Reveal credential là action riêng, không nằm trong service detail.
- Reseller/client không thể thấy dữ liệu tenant khác.
- Admin action critical yêu cầu audit + reason.
- Error codes thống nhất để frontend hiển thị đúng.
- Provider/provisioning retry không cho unsafe retry vô thức.
- Report dùng snapshot/ledger, không dùng giá hiện tại.
- API response không chứa provider secret hoặc credential plaintext ngoài reveal endpoint.

Câu nền: **API không chỉ là đường truyền dữ liệu; nó là cổng kiểm soát quyền, tiền, tenant và rủi ro.**

---


# ===== 17_RBAC_Permission_Matrix.md =====

# 17 - RBAC Permission Matrix

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa quyền truy cập theo vai trò cho platform VPS/Proxy multi-tenant.

RBAC phải bảo vệ 5 thứ:
```text
tiền
credential
tenant data
provider/source
hành động phá hủy như suspend/terminate
```

Không được dựa vào frontend để ẩn nút là đủ. Backend bắt buộc kiểm tra permission.

---

## 2. Role chuẩn phase 1

### 2.1 Platform Super Admin

Quyền cao nhất trên platform.

Có thể:
- Quản mọi tenant.
- Quản master catalog.
- Quản provider/source.
- Duyệt top-up reseller.
- Xem report toàn hệ thống.
- Emergency access tenant.

Yêu cầu:
```text
2FA bắt buộc
mọi critical action có audit
không dùng cho tác vụ hằng ngày nếu có role thấp hơn
```

### 2.2 Platform Admin / Staff

Dành cho nhân sự vận hành.

Có thể được cấp theo permission:
- support.
- finance.
- catalog.
- provisioning.
- abuse.
- read-only audit.

Không mặc định có toàn quyền.

### 2.3 Finance Agent

Chuyên xử lý top-up, refund, adjustment.

Có thể:
- xem top-up request.
- approve/reject top-up theo scope.
- tạo adjustment nếu được cấp.
- xem ledger/report.

Không nên có quyền:
- xem credential.
- sửa provider.
- terminate service.

### 2.4 Support Agent

Có thể:
- xem user/order/service metadata.
- hỗ trợ ticket.
- tạo note.
- request suspend/unsuspend theo workflow.

Không mặc định có quyền:
- reveal credential.
- approve top-up.
- adjustment ledger.
- đổi provider source.

### 2.5 Provisioning Operator

Có thể:
- xem failed jobs.
- retry safe jobs.
- resolve manual review.
- link existing provider resource nếu có quyền cao.

Không nên có quyền:
- duyệt tiền.
- chỉnh catalog price.
- xem toàn bộ ledger nếu không cần.

### 2.6 Reseller Owner

Owner của reseller tenant.

Có thể:
- quản branding/domain.
- quản staff.
- quản client trong tenant.
- clone catalog/set giá.
- duyệt top-up client.
- xem report profit tenant.
- suspend service trong tenant theo policy.

Yêu cầu:
```text
2FA bật mặc định, khuyến nghị bắt buộc.
critical action cần audit.
```

### 2.7 Reseller Staff

Nhân sự của reseller.

Quyền theo role con:
- sales.
- support.
- finance.
- catalog manager.
- read-only.

Không được có quyền owner mặc định.

### 2.8 Client

Khách mua dịch vụ.

Có thể:
- xem/mua/gia hạn dịch vụ của mình.
- xem ledger/top-up của mình.
- reveal credential của service thuộc mình.
- gửi support/cancel request.

Không thể:
- xem client khác.
- xem reseller cost.
- chỉnh giá.
- duyệt top-up.
- xem provider/source nội bộ.

### 2.9 Read-only Auditor

Có thể xem dữ liệu phục vụ kiểm toán nhưng không sửa.

Không được:
- reveal credential.
- approve tiền.
- retry provisioning.
- sửa tenant/domain/provider.

---

## 3. Permission naming convention

Format:
```text
module.resource.action
```

Ví dụ:
```text
tenant.view
tenant.create
tenant.update
tenant.domain.manage

catalog.master.view
catalog.master.create
catalog.master.update
catalog.tenant.clone
catalog.tenant.price_update

wallet.view
wallet.topup.submit
wallet.topup.approve
wallet.adjustment.create
wallet.ledger.view

order.view
order.create
order.cancel
order.refund

service.view
service.credential.reveal
service.renew
service.suspend
service.unsuspend
service.terminate

provider.view
provider.manage
provisioning.job.view
provisioning.job.retry
provisioning.manual_review.resolve

audit.view
report.view
risk.flag.manage
abuse.case.manage
```

---

## 4. Permission risk level

| Risk level | Examples | Control |
|---|---|---|
| Low | view catalog, view own order | auth + tenant scope |
| Medium | create order, submit top-up | auth + validation + audit optional |
| High | approve top-up, reveal credential, suspend service | permission + audit + rate limit |
| Critical | ledger adjustment, provider manage, terminate, emergency access | 2FA + reason + audit + restricted role |

---

## 5. Matrix tổng quan

| Action | Super Admin | Platform Staff | Finance | Support | Provisioning | Reseller Owner | Reseller Staff | Client | Auditor |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Xem tất cả tenant | Yes | Permission | No | Permission-limited | No | No | No | No | Read-only |
| Tạo reseller tenant | Yes | Permission | No | No | No | No | No | No | No |
| Sửa branding tenant | Yes | Permission | No | No | No | Own tenant | Permission | No | No |
| Quản custom domain | Yes | Permission | No | No | No | Own tenant | Permission | No | Read-only |
| Tạo master product/plan | Yes | Permission | No | No | No | No | No | No | Read-only |
| Clone plan về reseller | Yes | Permission | No | No | No | Yes | Permission | No | Read-only |
| Sửa giá bán reseller | Yes | Permission | No | No | No | Yes | Permission | No | Read-only |
| Xem reseller cost | Yes | Permission | Finance view | No | No | Own tenant | Permission | No | Read-only |
| Xem client wallet | Yes | Permission | Permission | Permission | No | Own tenant | Permission | Own wallet | Read-only |
| Submit top-up | No | No | No | No | No | For own reseller wallet | No | Own wallet | No |
| Approve client top-up | Yes | Permission | Permission | No | No | Own tenant | Permission | No | Read-only |
| Approve reseller top-up | Yes | Permission | Yes | No | No | No | No | No | Read-only |
| Ledger adjustment | Yes | Permission critical | Permission critical | No | No | Optional own tenant only | No | No | Read-only |
| Create order | No | No | No | No | No | For client via admin action optional | No | Own account | No |
| View order | Yes | Permission | Finance-related | Support-related | Provisioning-related | Own tenant | Permission | Own order | Read-only |
| Cancel pending order | Yes | Permission | No | Permission | No | Own tenant | Permission | Own order if policy | No |
| View service | Yes | Permission | No/limited | Permission | Permission | Own tenant | Permission | Own service | Read-only |
| Reveal credential | Yes + audit | Permission + audit | No | Usually No | Optional + audit | Own tenant + audit | Permission + audit | Own service + audit | No |
| Suspend service | Yes | Permission | No | Request/Permission | Permission if technical | Own tenant | Permission | No | No |
| Terminate service | Yes | Critical permission | No | No | Critical permission | Own tenant if policy | Permission critical | No | No |
| Retry provisioning job | Yes | Permission | No | No | Yes | No | No | No | Read-only |
| Resolve manual provisioning | Yes | Permission | No | No | Yes | No | No | No | Read-only |
| Manage provider source | Yes | Permission critical | No | No | No | No | No | No | Read-only |
| View audit logs | Yes | Permission | Finance scope | Support scope | Provisioning scope | Own tenant | Permission | No | Read-only |
| Manage abuse case | Yes | Permission | No | Permission | No | Own tenant limited | Permission | No | Read-only |

---

## 6. Data visibility rules

### 6.1 Platform admin

Can see:
```text
all tenants
all users
all orders
all services
all provider health
all financial reports
```

But:
- Credential reveal still requires explicit action.
- Provider secrets hidden by default.
- Emergency access into tenant must have reason.

### 6.2 Reseller

Can see:
```text
tenant clients
tenant orders
tenant services
tenant wallet/client wallet
tenant catalog clones
tenant reports
```

Cannot see:
```text
other tenant data
platform provider secret
global profit except own tenant
master cost beyond reseller cost exposed to them
```

### 6.3 Client

Can see:
```text
own profile
own wallet
own order
own service
own notifications
own ticket
```

Cannot see:
```text
reseller margin
provider source details
other users
raw audit logs
```

---

## 7. Emergency access policy

Emergency access là khi platform admin/staff truy cập dữ liệu tenant để hỗ trợ/sửa sự cố.

Required:
```text
permission: tenant.emergency_access
2FA: required
reason: required
target_tenant_id: required
time-bound session: recommended
audit action: tenant.emergency_access.started / ended
```

Forbidden:
```text
không dùng emergency access để duyệt tiền không có chứng từ
không reveal credential nếu không có reason rõ
không chỉnh ledger nếu không có finance permission
```

Audit metadata:
```text
actor_id
target_tenant_id
reason
started_at
ended_at
actions performed
correlation_id
```

---

## 8. Credential reveal policy

Credential reveal là action high risk.

Allowed by default:
```text
client owns service
reseller owner for own tenant
platform super admin
staff with service.credential.reveal
```

Controls:
```text
2FA for admin/reseller/staff
rate limit
audit every reveal
mask by default
no plaintext in log
```

Optional policy:
```text
support staff không được reveal, chỉ được reset password nếu provider hỗ trợ.
```

---

## 9. Finance permission policy

Top-up approval và ledger adjustment phải tách quyền.

### 9.1 Approve top-up

Yêu cầu:
```text
wallet.topup.approve
scope: client tenant hoặc reseller wallet
audit
review note recommended
```

### 9.2 Ledger adjustment

Yêu cầu:
```text
wallet.adjustment.create
reason required
supporting reference required
2FA required
audit critical
```

Không cho cùng một staff vừa tạo adjustment vừa approve nếu cần four-eyes policy phase sau.

---

## 10. Provider permission policy

Provider/source là tài sản platform-level.

Actions critical:
```text
provider.manage
provider.credentials.update
provider.source.disable
provider.source.delete/archive
```

Required:
```text
2FA
audit
reason
change summary
```

Không cấp provider.manage cho reseller.

---

## 11. Service destructive action policy

### 11.1 Suspend

Required:
```text
reason
notify_client decision
audit
provider capability check
```

### 11.2 Unsuspend

Required:
```text
reason
verify billing/abuse resolved
audit
```

### 11.3 Terminate

Required:
```text
critical permission
reason
confirmation
audit
check refund/cancel policy
provider terminate capability
```

Terminated service không được renew.

---

## 12. Role templates phase 1

### 12.1 Platform Finance

```text
wallet.view
wallet.topup.approve
wallet.ledger.view
report.finance.view
audit.view.finance
```

### 12.2 Platform Support

```text
user.view
order.view
service.view
ticket.manage
audit.view.support
risk.flag.create
```

### 12.3 Platform Provisioning

```text
provider.view
provisioning.job.view
provisioning.job.retry
provisioning.manual_review.resolve
service.view
audit.view.provisioning
```

### 12.4 Reseller Finance

```text
wallet.view
wallet.topup.approve
wallet.ledger.view
report.tenant_profit.view
```

### 12.5 Reseller Support

```text
client.view
order.view
service.view
ticket.manage
risk.flag.create
```

### 12.6 Reseller Catalog Manager

```text
catalog.master.view
catalog.tenant.clone
catalog.tenant.price_update
catalog.tenant.visibility_update
```

---

## 13. RBAC acceptance criteria

RBAC được xem là đạt khi:
- Backend check permission ở mọi API nhạy cảm.
- Tenant scope luôn được kiểm tra cùng permission.
- Client không thể truy cập service/order/wallet của user khác.
- Reseller không thể truy cập tenant khác.
- Staff không có quyền mặc định ngoài role được cấp.
- Credential reveal luôn tạo audit.
- Ledger adjustment cần reason + permission critical.
- Admin emergency access có reason + audit.
- UI ẩn nút theo permission, nhưng backend vẫn chặn nếu gọi trực tiếp.
- Test case cross-tenant access phải pass.

Câu nền: **quyền không phải để làm khó người dùng; quyền là ranh giới giữa platform có thể mở rộng và platform tự đâm xuyên dữ liệu của chính mình.**

---


# ===== 18_Provider_Adapter_Technical_Spec.md =====

# 18 - Provider Adapter Technical Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa contract kỹ thuật cho provider adapter layer. Đây là tài liệu interface/behavior, **không phải code triển khai**.

Provider adapter có nhiệm vụ biến các provider khác nhau thành một lớp hành vi thống nhất cho hệ thống:
```text
Proxmox
OVH
Hetzner
manual VPS provider
proxy upstream
preloaded proxy pool
custom API provider
```

Điểm quan trọng: không ép mọi provider hỗ trợ cùng một action. Adapter phải khai báo capability rõ ràng, UI/API phải dựa vào capability để cho phép hoặc chặn action.

---

## 2. Vai trò của provider adapter

Adapter layer nằm giữa:
```text
provisioning_worker/service_lifecycle_module
và
provider/source thật
```

Adapter chịu trách nhiệm:
- Kiểm tra health provider.
- Kiểm tra stock/capacity.
- Provision tài nguyên.
- Lấy trạng thái tài nguyên.
- Suspend/unsuspend/terminate nếu provider hỗ trợ.
- Reset password/reinstall/change IP nếu provider hỗ trợ.
- Fetch/generate credential.
- Chuẩn hóa lỗi provider thành error code nội bộ.
- Đánh dấu retry safety.
- Không làm lộ provider secret/credential trong log.

Adapter không chịu trách nhiệm:
- Debit wallet.
- Quyết định giá.
- Quyết định tenant permission.
- Quyết định order có được mua hay không.
- Tự ý retry ngoài policy của job.

---

## 3. Adapter capability profile

Mỗi source phải khai báo `capability_profile`.

### 3.1 Capability chung

```text
supports_health_check
supports_live_stock_check
supports_auto_provision
supports_manual_provision
supports_status_sync
supports_suspend
supports_unsuspend
supports_terminate
supports_renew
supports_reset_password
supports_reinstall
supports_change_ip
supports_bandwidth_usage
supports_console
supports_reverse_dns
supports_snapshot
supports_backup
supports_credential_fetch
supports_credential_rotation
```

### 3.2 VPS-specific capability

```text
supports_os_template_selection
supports_custom_hostname
supports_ipv6
supports_private_network
supports_resize
supports_rescue_mode
supports_vnc_console
supports_ssh_key_injection
```

### 3.3 Proxy-specific capability

```text
supports_proxy_protocol_http
supports_proxy_protocol_socks5
supports_rotating_proxy
supports_static_proxy
supports_geo_selection
supports_ip_whitelist
supports_userpass_auth
supports_bandwidth_quota
supports_thread_limit
supports_change_exit_ip
```

### 3.4 Capability rule

Nếu capability = false:
```text
UI không hiện action
API trả CAPABILITY_NOT_SUPPORTED nếu user gọi trực tiếp
worker không tạo job action đó
```

Nếu capability thay đổi sau khi service đã bán:
```text
service dùng capability_snapshot tại thời điểm mua
nhưng admin có thể override nếu provider/source bị disable/maintenance
```

---

## 4. Adapter operation contract

Mỗi operation phải trả về kết quả chuẩn hóa:

```text
operation_result:
- status: success / failed / partial_success / unknown / manual_review_required
- external_request_id
- external_resource_id
- provider_status
- credentials_encrypted_payload optional
- public_service_identifier optional
- retry_safety: safe_retry / unsafe_retry / do_not_retry / manual_review_required
- error_code optional
- error_message_redacted optional
- raw_response_reference optional
```

Không trả:
```text
provider_api_secret
plaintext credential trong log
full raw response nếu chứa secret
```

---

## 5. Required operations

### 5.1 `checkHealth`

Mục tiêu:
```text
xác định provider/source còn gọi được không
```

Output:
```text
healthy
degraded
down
unknown
```

Nếu provider trả 401/403:
```text
source health = down
provider credential invalid
alert admin
không retry provisioning mới qua source này
```

### 5.2 `checkStock`

Mục tiêu:
```text
xác định source còn capacity cho plan/source không
```

Output:
```text
available
out_of_stock
unknown
capacity_count optional
```

Rule:
- Nếu source finite stock: dùng DB inventory là chính, có thể sync provider.
- Nếu provider live stock: phải check live trước reserve hoặc theo cache TTL.
- Nếu unknown: tùy policy, có thể cho manual_review hoặc chặn checkout.

### 5.3 `provision`

Input logic:
```text
tenant context
order item snapshot
plan specs
source config
idempotency_key
correlation_id
```

Output:
```text
external_resource_id
service_identifier
credentials
provider_status
```

Status:
- `success`: provider đã tạo tài nguyên, có thể tạo service active.
- `failed`: provider chắc chắn không tạo tài nguyên.
- `partial_success`: provider có thể đã tạo hoặc tạo thiếu dữ liệu.
- `unknown`: mất kết nối/timeout không biết kết quả.
- `manual_review_required`: cần người xử lý.

### 5.4 `getStatus`

Mục tiêu:
```text
sync trạng thái tài nguyên thật với service trong hệ thống
```

Output:
```text
external_status
is_running
is_suspended
is_terminated
usage metrics optional
```

Không tự đổi service_status nếu không qua lifecycle module.

### 5.5 `suspend` / `unsuspend`

Rule:
- Nếu provider không hỗ trợ: trả `CAPABILITY_NOT_SUPPORTED`.
- Nếu suspend do billing/abuse, reason phải được truyền vào job.
- Nếu suspend thành công provider nhưng hệ thống fail update, worker phải retry update nội bộ trước khi gọi provider lần nữa.

### 5.6 `terminate`

Terminate là destructive action.

Rule:
```text
không retry mù terminate nếu provider response unknown
manual review nếu không chắc tài nguyên đã bị xóa hay chưa
service chỉ chuyển terminated khi xác nhận hoặc admin resolve
```

### 5.7 `renew`

Một số provider cần gọi renew API, một số chỉ cần internal billing.

Rule:
```text
nếu provider renew required -> tạo job renew
nếu không -> chỉ extend term_end_at internal
```

### 5.8 `resetPassword`

Output:
```text
new encrypted credential
provider_status
```

Rule:
```text
credential cũ chuyển rotated/revoked nếu reset thành công
audit credential.rotated
```

### 5.9 `reinstall`

Reinstall có thể xóa dữ liệu.

Rule:
```text
client confirmation required
warning required
provider capability required
audit service.reinstall_requested
```

### 5.10 `changeIp`

Rule:
```text
check capability
check rate limit
possible extra cost
provider may return new endpoint/credential
update service_identifier and credential if needed
audit service.ip_changed
```

---

## 6. Error normalization

Provider error phải map về internal error.

| Provider situation | Internal code | Retry safety | Action |
|---|---|---|---|
| API key invalid / 401 | `PROVIDER_AUTH_FAILED` | do_not_retry | disable/alert |
| Rate limit / 429 | `PROVIDER_RATE_LIMITED` | safe_retry | backoff |
| Out of stock | `PROVIDER_OUT_OF_STOCK` | do_not_retry | release reservation |
| Timeout before response | `PROVIDER_TIMEOUT_UNKNOWN` | unsafe_retry | check resource/manual review |
| Timeout after provider request id known | `PROVIDER_TIMEOUT_REQUEST_KNOWN` | manual_review_required | query status first |
| 500 temporary | `PROVIDER_TEMPORARY_ERROR` | safe_retry | retry limited |
| Plan/template not found | `PROVIDER_PLAN_NOT_FOUND` | do_not_retry | disable source/plan mapping |
| Success but missing credential | `PROVIDER_CREDENTIAL_MISSING` | manual_review_required | fetch credential/manual |
| Resource already exists by idempotency | `PROVIDER_RESOURCE_ALREADY_EXISTS` | safe if mapped | link existing resource |
| Provider says terminated but system active | `PROVIDER_STATE_DRIFT` | manual_review_required | reconcile |

---

## 7. Retry policy

### 7.1 Safe retry

Có thể retry nếu:
```text
provider chắc chắn chưa tạo tài nguyên
hoặc operation idempotent
hoặc provider hỗ trợ idempotency key mạnh
```

Ví dụ:
```text
429 rate limit
temporary 500 trước khi gửi create
health check fail
status sync fail
```

### 7.2 Unsafe retry

Không retry tự động nếu:
```text
timeout sau khi đã gửi create request
provider không hỗ trợ idempotency
không biết external_resource_id
response mâu thuẫn
```

Action:
```text
manual_review_required
query provider by external_request_id nếu có
query by metadata/tag/idempotency nếu provider hỗ trợ
```

### 7.3 Do not retry

Không retry nếu:
```text
API key invalid
plan/template không tồn tại
source disabled
out_of_stock confirmed
permission denied
provider account suspended
```

---

## 8. Idempotency with provider

### 8.1 Internal idempotency

Mỗi job có:
```text
job_id
idempotency_key
correlation_id
```

Adapter phải nhận idempotency_key từ job, không tự tạo mới khi retry.

### 8.2 Provider idempotency

Nếu provider hỗ trợ:
```text
truyền idempotency_key vào provider request
lưu external_request_id
lưu external_resource_id
```

Nếu provider không hỗ trợ:
```text
phải gắn metadata/tag nếu được
hoặc dùng manual review khi unknown
```

### 8.3 Duplicate detection

Trước khi gọi create lần nữa, worker/adapter phải kiểm tra:
```text
provider_resource_mappings
provider_requests with same idempotency_key
provider external lookup if supported
```

---

## 9. Provider request logging

Mỗi provider call tạo `provider_requests`.

Lưu:
```text
request_type
external_request_id
external_resource_id
status
retry_safety
request_summary_redacted
response_summary_redacted
error_code
sent_at
received_at
correlation_id
```

Không lưu:
```text
API secret
root password plaintext
proxy password plaintext
full auth header
session token
```

Nếu raw response cần lưu để debug:
```text
lưu private object storage
redact trước
set retention
restrict permission
```

---

## 10. Credential handling

Adapter có thể nhận credential từ provider hoặc tự generate.

Rule:
```text
credential phải được encrypt trước khi lưu DB
plaintext chỉ tồn tại trong memory ngắn hạn
không đưa credential vào audit/provider_request log
service detail API chỉ trả masked_hint
reveal credential là API/action riêng
```

Credential payload logic:
```text
vps:
- hostname/ip
- username
- password or ssh key
- port
- console link if supported

proxy:
- host
- port
- protocol
- username
- password
- geo
- expiry/quota if any
```

---

## 11. Manual provider/source

Phase 1 có thể có provider manual.

Manual source flow:
```text
checkout -> reserve -> debit wallet -> provisioning_job manual_review
operator nhập external_resource_id/service_identifier/credential
system tạo service active
```

Manual source vẫn phải:
- có source_id.
- có reservation.
- có ledger.
- có audit.
- có credential encrypted.
- có lifecycle event.

Không được bypass hệ thống bằng cách tạo service tay không có order/ledger nếu là paid service.

---

## 12. Proxmox adapter notes

Proxmox thường có:
```text
node
storage
template
vmid
network bridge
ip allocation
cloud-init/user/password
```

Cần chú ý:
- VMID unique.
- clone template có thể mất thời gian.
- lock VM trong lúc clone/reinstall.
- password/cloud-init credential có thể không fetch lại dễ dàng.
- partial success có thể xảy ra khi clone thành công nhưng API timeout.

Capability thường:
```text
auto_provision yes
suspend yes via stop/disable policy
terminate yes
reset_password depends cloud-init/guest agent
reinstall yes but destructive
console maybe yes
```

---

## 13. Proxy upstream adapter notes

Proxy source có thể là:
```text
preloaded list
upstream API
rotating proxy account
dedicated static proxy
```

Cần chú ý:
- Một proxy credential có thể là tài nguyên thật.
- Oversell dễ xảy ra nếu preloaded list không lock atomic.
- Change IP có thể là rotate endpoint hoặc cấp proxy mới.
- Bandwidth/quota sync nếu provider hỗ trợ.

Preloaded proxy inventory:
```text
proxy_pool_items
status: available/reserved/allocated/disabled
reservation_id
service_id
```

---

## 14. Provider onboarding checklist

Trước khi bật source active:
- Provider account credential được encrypt.
- Health check pass.
- Capability profile được xác nhận.
- Test provision sandbox/manual nhỏ.
- Test timeout behavior.
- Test credential retrieval.
- Test terminate/suspend nếu hỗ trợ.
- Define retry safety map.
- Define stock mode.
- Define plan/source mapping.
- Define provider cost.
- Define abuse/takedown process.
- Add monitoring alert.

---

## 15. Adapter acceptance criteria

Adapter layer đạt khi:
- Mọi provider action trả result chuẩn hóa.
- Capability profile điều khiển UI/API/worker.
- Timeout create không retry mù.
- Credential không xuất hiện trong logs/audit.
- Idempotency key được truyền xuyên job/provider request.
- Out-of-stock release reservation đúng.
- Provider auth fail disable/alert đúng.
- Manual provider vẫn có order/ledger/reservation/service/audit.
- Có error mapping table cho từng provider.
- Có manual review flow cho partial_success/unknown.

Câu nền: **adapter tốt không làm provider giống nhau; adapter tốt làm hệ thống biết provider khác nhau ở đâu và phản ứng an toàn.**

---


# ===== 19_Worker_Queue_And_Cron_Jobs_Spec.md =====

# 19 - Worker Queue And Cron Jobs Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa các worker/job nền cần có cho platform VPS/Proxy. Đây là job behavior spec, không phải code.

Dự án này không nên xử lý tất cả bằng request realtime. Các việc như provisioning, provider sync, expiry, suspend/terminate, notification phải chạy qua queue/cron để:
- tránh request timeout.
- kiểm soát retry.
- giữ idempotency.
- ghi audit/correlation.
- xử lý provider chậm/lỗi.

---

## 2. Nguyên tắc worker/job

### 2.1 Idempotent

Job chạy lại không được:
```text
trừ tiền lần hai
reserve stock lần hai
allocate stock lần hai
tạo service trùng
gọi provider create trùng nếu không an toàn
gửi notification critical trùng vô hạn
```

### 2.2 Observable

Mỗi job phải có:
```text
job_id
job_type
status
attempt_count
max_attempts
correlation_id
started_at
finished_at
last_error_code
next_retry_at
manual_review_reason
```

### 2.3 Retry có phân loại

Không retry theo kiểu “cứ lỗi là chạy lại”.

Retry dựa vào:
```text
safe_retry
unsafe_retry
do_not_retry
manual_review_required
```

### 2.4 Tenant scoped

Job liên quan tenant phải lưu `tenant_id` và không xử lý resource ngoài tenant đó.

### 2.5 Audit important transitions

Job làm thay đổi:
```text
service_status
order_status
wallet/ledger
reservation
credential
provider/resource
```

thì phải audit.

---

## 3. Queue message contract

Message logic:

```text
job_id
job_type
tenant_id
reference_type
reference_id
source_id optional
idempotency_key
correlation_id
attempt_count
created_at
scheduled_at optional
```

Reference examples:
```text
order_id
reservation_id
service_id
topup_request_id
provider_source_id
abuse_case_id
```

---

## 4. Job status

```text
queued
running
success
failed
retry_scheduled
manual_review
cancelled
```

Transition:
```text
queued -> running
running -> success
running -> failed
running -> retry_scheduled
running -> manual_review
retry_scheduled -> queued
manual_review -> queued/success/failed/cancelled by operator
```

---

## 5. Core worker: `provisioning_worker`

### 5.1 Trigger

Created by checkout when:
```text
order paid
reservation created
wallet debit posted
risk check passed
```

### 5.2 Input

```text
provisioning_job_id
order_id
order_item_id
reservation_id
source_id
idempotency_key
correlation_id
```

### 5.3 Preconditions

- job status queued/retry_scheduled.
- order payment_status paid.
- reservation status reserved.
- source active.
- provider health not down unless manual override.
- no existing service for same order_item unless resolving.

### 5.4 Steps

```text
1. mark job running
2. load order/order_item/reservation/source snapshot
3. check idempotency/provider request history
4. call adapter.provision
5. if success:
   - create provider_resource_mapping
   - create encrypted service_credentials
   - create service status active
   - reservation reserved -> allocated
   - order provisioning_status success
   - order/service lifecycle event
   - send activation notification
6. if failed and retry_safety safe_retry:
   - schedule retry if attempts left
7. if failed do_not_retry:
   - release reservation
   - reverse wallet debit if policy requires
   - order failed or manual_review depending policy
8. if partial_success/unknown/unsafe_retry:
   - job manual_review
   - order provisioning_status manual_review
   - alert operator
```

### 5.5 Failure handling

| Situation | Handling |
|---|---|
| Provider timeout unknown | manual_review, no retry mù |
| Provider out of stock | release reservation, refund/reversal, notify |
| Provider auth fail | disable source, alert admin, do not retry |
| Credential missing | manual_review, attempt credential fetch if safe |
| DB update fails after provider success | retry internal reconciliation before new provider call |

### 5.6 Audit

```text
provisioning.job.started
provider.request.sent
provisioning.job.succeeded
service.activated
reservation.allocated
provisioning.job.failed
provisioning.job.manual_review
```

---

## 6. Worker: `provider_sync_worker`

### 6.1 Trigger

- scheduled by provider sync cron.
- manual admin action.
- after unknown/partial provisioning.

### 6.2 Purpose

Sync external provider status with internal service mapping.

### 6.3 Rules

- Không tự terminate service chỉ vì provider không thấy resource một lần.
- Nếu mismatch, tạo provider_state_drift risk/manual review.
- Nếu service internal active nhưng provider terminated, alert admin.
- Nếu provider suspended nhưng internal active, mark drift and review.

### 6.4 Audit

```text
provider.sync.started
provider.sync.completed
provider.state_drift.detected
```

---

## 7. Worker: `service_action_worker`

Handles:
```text
suspend
unsuspend
terminate
reset_password
reinstall
change_ip
renew provider-side
```

### 7.1 Preconditions

- action supported by capability_snapshot/current policy.
- actor/action permission already validated before job creation.
- service belongs to tenant.
- action not conflicting with another running action.

### 7.2 Conflict examples

```text
cannot reinstall while terminate running
cannot change_ip while service terminated
cannot unsuspend if abuse case still active unless override
cannot renew terminated service
```

### 7.3 Audit

```text
service.action.started
service.suspended
service.unsuspended
service.terminated
credential.rotated
service.ip_changed
service.action.failed
```

---

## 8. Worker: `notification_worker`

### 8.1 Trigger

Created by modules/jobs when notification should be sent.

### 8.2 Channels

```text
email
dashboard
telegram
webhook
```

### 8.3 Rules

- Notification payload redacted.
- Credential activation email should avoid plaintext password unless policy explicitly allows. Safer: send “login to reveal”.
- Critical admin alert should include correlation_id.
- Failed send retry with backoff.
- Do not spam duplicate event; use dedupe key.

### 8.4 Dedupe key examples

```text
service_expiry_reminder:{service_id}:{days_before}
reseller_low_balance:{tenant_id}:{date}
provider_down:{source_id}:{date_hour}
```

---

## 9. Cron: `reservation_expiry_job`

### 9.1 Frequency

```text
every 1 minute
```

### 9.2 Logic

Find:
```text
reservations status = reserved
expires_at < now
```

For each:
```text
mark reservation expired
decrement reserved_count
if order not paid/provisioning -> order expired/cancelled
if wallet already debited due to edge case -> reversal per policy
audit reservation.expired
```

### 9.3 Idempotency

- If reservation already allocated/released/expired, skip.
- Use lock on reservation row/source inventory.
- Running twice must not decrement stock twice.

---

## 10. Cron: `service_expiry_job`

### 10.1 Frequency

```text
every 5-15 minutes
```

### 10.2 Logic

Find active services:
```text
term_end_at < now
billing_status not overdue/grace
```

Then:
```text
service_status may remain active or become expired depending policy
billing_status -> overdue/grace
send expired/grace notification
create lifecycle event
```

If policy says auto suspend immediately:
```text
create service_action_job suspend
```

### 10.3 Notes

Do not terminate immediately on expiry unless product policy says no grace.

---

## 11. Cron: `suspension_job`

### 11.1 Frequency

```text
every 15 minutes
```

### 11.2 Logic

Find:
```text
services billing_status = grace/overdue
grace_until_at < now
service_status active/expired
```

Then:
```text
create suspend job
reason = billing_overdue
notify client/reseller
```

### 11.3 Guard

- Skip if service renewed.
- Skip if already suspended/terminated.
- If provider does not support suspend, mark internal suspended and alert if required.

---

## 12. Cron: `termination_job`

### 12.1 Frequency

```text
daily or every few hours depending product risk
```

### 12.2 Logic

Find:
```text
services suspended/expired
terminate_after_at < now
policy auto_terminate = true
```

Then:
```text
create terminate job
```

### 12.3 Guard

Terminate is destructive:
- require product policy.
- keep audit.
- allow admin hold flag to prevent auto terminate.
- if provider response unknown, manual review.

---

## 13. Cron: `renewal_reminder_job`

### 13.1 Frequency

```text
daily
```

### 13.2 Reminder windows

Suggested:
```text
7 days before
3 days before
1 day before
on expiry
after grace start
```

### 13.3 Dedupe

```text
service_id + reminder_type + date/window
```

---

## 14. Cron: `provider_health_check_job`

### 14.1 Frequency

```text
every 5-15 minutes for active sources
```

### 14.2 Logic

For each active source:
```text
adapter.checkHealth
update health_status
if down/degraded -> alert admin
if recovered -> alert/reopen source depending policy
```

### 14.3 Guard

- Do not disable source on one transient failure unless threshold reached.
- Disable auto-provision if repeated critical failure.

---

## 15. Cron: `provider_inventory_sync_job`

### 15.1 Frequency

```text
every 15-60 minutes
```

### 15.2 Logic

- Sync capacity for provider_live source.
- Update available_count_cache.
- Detect out_of_stock.
- Alert if stock below threshold.

### 15.3 Guard

Do not override reserved/allocated counts incorrectly. Provider live availability and internal reservations must reconcile.

---

## 16. Cron: `pending_topup_review_reminder_job`

### 16.1 Frequency

```text
every 30-60 minutes
```

### 16.2 Logic

Find:
```text
topup_requests submitted/under_review older than SLA threshold
```

Notify:
```text
finance/admin/reseller owner
```

---

## 17. Cron: `reseller_low_balance_job`

### 17.1 Frequency

```text
daily or every few hours
```

### 17.2 Trigger

```text
reseller settlement wallet < configured threshold
or reseller balance < estimated cost for next N orders/days
```

Notify reseller:
```text
new client orders may not provision if reseller wallet is insufficient
```

---

## 18. Cron: `audit_retention_job`

### 18.1 Frequency

```text
daily/weekly
```

### 18.2 Rule

- Financial/security audit retained long-term.
- Low-risk debug logs can be archived/deleted per retention.
- Never delete ledger.
- Never delete critical audit before retention policy allows.

---

## 19. Manual review queue

Manual review items should appear in Admin Portal:

Types:
```text
provisioning_unknown
provider_state_drift
credential_missing
payment_mismatch
abuse_case
high_risk_order
unsafe_retry
```

Each item must show:
```text
severity
age
tenant
order/service
correlation_id
recommended next action
```

Allowed resolutions:
```text
retry_safe
link_existing_resource
mark_failed_and_refund
mark_success
cancel_order
suspend_service
clear_flag
escalate
```

Every resolution requires note/audit.

---

## 20. Monitoring metrics

P0 metrics:
```text
jobs queued by type
jobs running by type
jobs failed by type
jobs manual_review by type
avg provisioning time
provider success/fail rate
reservation expiry count
service expiry/suspend/terminate count
notification failure count
topup pending age
```

Alert examples:
```text
provider down > 5 minutes
provisioning manual_review > threshold
queue backlog > threshold
reservation expiry spike
ledger adjustment spike
credential reveal spike
```

---

## 21. Worker acceptance criteria

Worker/cron spec đạt khi:
- Mọi job idempotent.
- Provisioning success tạo service/credential/reservation allocation đúng.
- Unknown provider response không retry mù.
- Reservation expired không double-release.
- Expiry/suspend/terminate job không đụng service đã renewed.
- Notification không spam trùng vô hạn.
- Provider health issue alert đúng.
- Manual review queue có resolution rõ.
- Có correlation_id xuyên job/provider/order/service.
- Job failure có trạng thái, error code, next step.

Câu nền: **worker là nơi hệ thống thể hiện bản lĩnh thật: không phải lúc mọi thứ chạy tốt, mà khi provider timeout, mạng rớt, queue chạy lại và dữ liệu vẫn không sai.**

---


# ===== 20_UI_Wireflow_And_Screen_Spec.md =====

# 20 - UI Wireflow And Screen Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa wireflow và screen spec cho 3 portal:
```text
Admin Portal
Reseller Portal
Client Portal
```

Đây không phải tài liệu thiết kế mỹ thuật/Figma. Mục tiêu là khóa:
- Màn nào cần có.
- Ai được xem.
- Dữ liệu hiển thị.
- Action nào có trên màn.
- Empty/error/loading state.
- Audit/notification liên quan.
- Security rule đặc biệt.

---

## 2. Nguyên tắc UI

### 2.1 UI phải theo permission và capability

Frontend phải ẩn action nếu:
```text
user không có permission
service/source capability không hỗ trợ
service status không cho action
tenant policy không cho action
```

Backend vẫn phải chặn nếu user gọi API trực tiếp.

### 2.2 Credential luôn masked mặc định

Service detail không hiển thị plaintext credential.

Credential flow:
```text
masked -> click Reveal -> confirm/2FA nếu cần -> API reveal -> hiển thị tạm thời -> audit
```

### 2.3 Money state phải rõ

Các màn ví/order phải hiển thị:
```text
available balance
pending top-up
ledger history
order payment status
refund/reversal nếu có
```

### 2.4 Manual review phải dễ thấy

Admin/reseller nên có dashboard cảnh báo:
```text
failed provisioning
manual review
pending top-up
provider down
abuse case
low reseller balance
```

---

## 3. Client Portal wireflow

### 3.1 Public storefront

Path logic:
```text
/
 /products
 /products/{product}
 /login
 /register
```

Data:
- Brand logo/theme theo tenant.
- Product categories: VPS, Proxy.
- Plan cards.
- Price.
- Specs.
- Stock status summary.
- Billing cycle.
- Terms/AUP link.

Actions:
- Register.
- Login.
- Buy now.

Empty state:
```text
No active products available.
```

Error state:
```text
Tenant disabled/domain not verified.
```

Security:
- Không hiển thị reseller_cost/provider internals.

### 3.2 Client dashboard

Purpose:
```text
tổng quan ví, order, service, expiry
```

Data:
- Wallet balance.
- Active services count.
- Expiring soon services.
- Pending orders.
- Recent notifications.

Actions:
- Top up wallet.
- Buy service.
- Renew expiring service.
- Open support.

### 3.3 Client wallet screen

Data:
```text
available balance
pending top-up requests
ledger entries
top-up instructions
```

Actions:
- Create top-up request.
- Upload proof.
- Cancel draft/submitted top-up if policy allows.

States:
- Pending review.
- Approved.
- Rejected with reason.
- Expired.

Acceptance:
- Ledger entries must show reference/order/topup.
- Balance must not be editable by client.

### 3.4 Client checkout flow

Steps:
```text
1. Select plan
2. Confirm specs and price
3. Check wallet balance
4. Accept terms/AUP
5. Submit order
6. Show provisioning status
```

UI validations:
- Insufficient wallet -> show top-up CTA.
- Plan out of stock -> disable buy.
- Provider unavailable -> disable/notify.
- Risk/manual review -> show pending review.

Result states:
```text
order provisioning queued
service active
manual review
failed with refund/reversal
```

### 3.5 Client orders list/detail

List columns:
```text
order_number
product/plan
amount
payment_status
provisioning_status
created_at
```

Detail:
```text
order timeline
ledger references
reservation/provisioning status
service link if active
refund/reversal if any
```

Actions:
- Cancel if allowed.
- Contact support.

### 3.6 Client services list

Columns/cards:
```text
service name/type
status
identifier masked/summary
expiry date
renew button
actions available
```

Filters:
```text
active
expired
suspended
terminated
vps/proxy
```

### 3.7 Client service detail

Data:
```text
service_id
product/plan snapshot
status
term_start_at
term_end_at
billing_status
service identifier
credential masked_hint
capability actions
lifecycle timeline
support link
```

Actions:
```text
Reveal credential
Renew
Request cancel
Reset password if supported
Reinstall if supported
Change IP if supported
Open ticket
```

Rules:
- Reveal credential is separate action and audit.
- Reinstall/change IP require confirmation.
- Terminated service has no renew/action except view history.

### 3.8 Client support/ticket screen

Data:
```text
tickets
messages
related service/order
status
```

Actions:
- Create ticket.
- Reply.
- Attach proof/evidence if allowed.

Security:
- Attachment private.
- Do not paste credential automatically.

---

## 4. Reseller Portal wireflow

### 4.1 Reseller dashboard

Data:
```text
reseller wallet balance
low balance warning
client count
active services
pending top-ups
profit summary
manual review items
expiring services
abuse warnings
```

Actions:
- Top up reseller wallet.
- Approve client top-up.
- Manage catalog.
- View failed/manual review orders.
- Manage clients.

### 4.2 Branding and storefront settings

Data:
```text
logo
brand color
support email
telegram/support link
footer
terms/AUP overrides if allowed
```

Actions:
- Update branding.
- Preview storefront.

Audit:
```text
tenant.branding.updated
```

### 4.3 Domain settings

Data:
```text
system subdomain
custom domains
verification status
TLS status
primary domain
```

Actions:
- Add domain.
- Verify.
- Set primary.
- Disable/remove domain.

Error states:
- Domain already used.
- Verification failed.
- TLS pending/failed.

Audit:
```text
tenant.domain.created
tenant.domain.verified
tenant.domain.primary_changed
```

### 4.4 Reseller catalog management

Data:
```text
available master plans
tenant cloned plans
selling price
reseller cost
margin
visibility
status
stock/source summary
```

Actions:
- Clone plan.
- Update selling price.
- Hide/show plan.
- Sync with master version.
- Disable plan.

Warnings:
```text
selling_price < reseller_cost
master plan has new version
source out of stock
provider degraded
```

Acceptance:
- Reseller never sees provider secret.
- Cost shown must be reseller cost, not necessarily platform internal base cost.

### 4.5 Client management

Data:
```text
client list
status
wallet balance
orders count
services count
risk flags
last login
```

Actions:
- Create/invite client if allowed.
- Suspend/disable client.
- View client detail.
- View wallet/order/service within tenant.
- Add risk note.

Security:
- Only clients under reseller tenant.

### 4.6 Client top-up review

Data:
```text
request amount
payment method
payment reference
proof attachment
client info
history
```

Actions:
- Approve.
- Reject with reason.
- Mark under review.
- Request more info.

Controls:
- Finance permission required.
- Approved creates ledger credit.
- Cannot approve twice.

Audit:
```text
wallet.topup.approved
wallet.topup.rejected
```

### 4.7 Reseller wallet/settlement

Data:
```text
reseller settlement wallet balance
ledger entries
platform cost debits
top-up history
low balance threshold
```

Important message:
```text
If reseller balance is insufficient, new client orders may not be provisioned even if client wallet has funds.
```

### 4.8 Reseller report

Reports:
```text
gross revenue from clients
platform reseller cost
gross profit
refunds/adjustments
top products
top clients
expiring services
```

Must use:
```text
ledger/order snapshots
```

Not current plan price.

---

## 5. Admin Portal wireflow

### 5.1 Admin dashboard

Cards:
```text
total tenants
active services
today orders
provisioning failures
manual review queue
provider health
pending reseller top-ups
abuse cases
revenue summary
```

Critical alerts:
- Provider down.
- Queue backlog.
- Failed provisioning spike.
- Manual review aging.
- Low stock.
- Credential reveal spike.

### 5.2 Tenant management

Data:
```text
tenant name/type/status
owner
domain
wallet balance
clients/services
risk status
created_at
```

Actions:
- Create tenant.
- Suspend/disable tenant.
- View tenant detail.
- Emergency access.
- Manage owner.

Emergency access:
- Requires reason.
- 2FA.
- Audit.

### 5.3 Master catalog

Screens:
```text
products list/detail
plans list/detail
plan-source mapping
version history
```

Actions:
- Create/update product.
- Create/update plan.
- Disable/archive.
- Create new version.
- Map source.
- Set reseller cost/min price.

Warnings:
- Plan used by active services.
- Changing cost may affect reseller margin.
- Source disabled/out-of-stock.

### 5.4 Provider/source management

Data:
```text
provider account
source type
health status
capability profile
inventory mode
capacity
last sync
failure rate
```

Actions:
- Add provider account.
- Update encrypted credentials.
- Add source.
- Test health.
- Disable source.
- Sync inventory.
- View provider request logs redacted.

Security:
- Provider secrets masked.
- Credential update requires critical permission/2FA.

### 5.5 Top-up and finance operations

Screens:
```text
reseller top-up queue
client direct top-up queue
ledger search
adjustment request/create
finance reports
```

Actions:
- Approve/reject top-up.
- Create adjustment.
- Export report.
- Search by correlation_id.

Controls:
- Adjustment requires reason/reference.
- Cannot edit ledger.
- Audit critical.

### 5.6 Provisioning operations

Data:
```text
jobs queued/running/failed/manual_review
source
tenant
order/service
attempt_count
retry_safety
error_code
correlation_id
```

Actions:
- Retry safe.
- Resolve manual review.
- Link existing external resource.
- Mark failed and refund.
- Disable problematic source.

Warnings:
- Unsafe retry requires high permission and reason.
- Provider unknown status should not retry blindly.

### 5.7 Service operations

Data:
```text
service status
tenant/client
provider mapping
billing status
lifecycle events
credential masked
```

Actions:
- Suspend.
- Unsuspend.
- Terminate.
- Reset password if supported.
- Sync provider.
- Reveal credential if permission.

Controls:
- Terminate critical.
- Reveal audit.
- Reason required for suspend/terminate.

### 5.8 Audit logs

Filters:
```text
tenant
actor
action
target
correlation_id
date range
risk level
```

Display:
```text
redacted before/after
metadata
ip/user agent
job id if worker
```

No plaintext secrets.

### 5.9 Abuse/risk center

Data:
```text
open abuse cases
risk flags
high-risk orders
provider notices
blacklist
```

Actions:
- Create case.
- Attach evidence.
- Warn client.
- Suspend service.
- Close case.
- Blacklist marker.

---

## 6. Common UI components

### 6.1 Status badges

Standard badges:
```text
active
pending
provisioning
manual review
failed
suspended
expired
terminated
out of stock
provider degraded
```

### 6.2 Timeline component

Used for:
```text
order timeline
service lifecycle
top-up review
provisioning job
abuse case
```

Timeline event:
```text
time
action
actor/system
status
note
correlation_id optional
```

### 6.3 Confirmation modal

Required for:
```text
reveal credential
reinstall
change IP
suspend
terminate
ledger adjustment
provider disable
domain primary change
```

### 6.4 Empty states

Examples:
```text
No services yet. Buy your first VPS/proxy.
No pending top-ups.
No manual review items.
No active provider sources.
```

### 6.5 Error states

Error should show:
```text
human message
error code
support correlation_id if needed
next action
```

Do not expose provider raw secret/error.

---

## 7. UI acceptance criteria

UI spec đạt khi:
- 3 portal có screen map rõ.
- Mỗi màn có role/data/action/error state.
- Credential luôn masked, reveal qua action riêng.
- Checkout hiển thị rõ insufficient client/reseller balance.
- Reseller thấy margin/cost của tenant mình.
- Admin thấy manual review/provider health/finance alerts.
- Action destructive có confirmation + reason.
- UI dùng capability để ẩn action không hỗ trợ.
- Cross-tenant data không thể xuất hiện trong UI.
- Mỗi lỗi quan trọng có correlation_id để support tra.

Câu nền: **UI tốt không chỉ đẹp; nó làm trạng thái hệ thống trở nên thật, để người vận hành biết phải làm gì trước khi lỗi thành tiền mất.**

---


# ===== 21_QA_Test_Cases_And_Acceptance_Plan.md =====

# 21 - QA Test Cases And Acceptance Plan

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa test plan và acceptance criteria cho dự án VPS/Proxy.

Mục tiêu:
- Test đúng nghiệp vụ tiền, tenant, provisioning, credential.
- Phát hiện lỗi trước khi lên production.
- Giảm kiểu test “click thấy chạy là được”.
- Làm căn cứ nghiệm thu dev.

---

## 2. Test principles

### 2.1 Test theo rủi ro

Ưu tiên test:
```text
money
tenant isolation
credential
provisioning idempotency
inventory oversell
refund/reversal
permission
abuse/suspend/terminate
```

### 2.2 Test negative case bắt buộc

Không chỉ test happy path. Dự án này chết ở edge case:
```text
provider timeout
reseller wallet thiếu
cross-tenant access
double click checkout
reservation expired
top-up approve twice
unsafe retry
credential reveal unauthorized
```

### 2.3 Test bằng correlation_id

Mỗi test nghiệp vụ phức tạp phải kiểm tra trace:
```text
order -> ledger -> reservation -> provisioning_job -> provider_request -> service -> audit
```

---

## 3. Test environment requirements

Minimum:
```text
staging environment
test tenant platform
test reseller A
test reseller B
test clients under each reseller
mock/manual provider source
at least one finite inventory source
test wallet balances
email/notification sandbox
```

Provider test:
```text
successful provider
out-of-stock provider
timeout/unknown provider
auth-failed provider
manual provider
```

---

## 4. Tenant isolation test cases

### TI-001 Client cannot access another client's service

Setup:
```text
Client A and Client B in same tenant.
Client B owns service S.
```

Action:
```text
Client A requests service S detail.
```

Expected:
```text
403 or 404.
No credential returned.
Audit/security event optional.
```

### TI-002 Reseller A cannot access Reseller B client

Setup:
```text
Reseller A tenant.
Reseller B tenant.
Client B belongs to Reseller B.
```

Action:
```text
Reseller A staff requests Client B profile/order/service.
```

Expected:
```text
403/404.
No data leak.
```

### TI-003 Body tenant_id injection ignored

Action:
```text
Client under tenant A submits request with tenant_id = tenant B in body.
```

Expected:
```text
Backend uses tenant A from context.
Request fails if resource not in tenant A.
No write under tenant B.
```

### TI-004 Domain maps correct tenant

Action:
```text
Open storefront from reseller custom domain.
```

Expected:
```text
Catalog/user registration belongs to mapped tenant.
```

---

## 5. RBAC test cases

### RBAC-001 Support cannot approve top-up

Setup:
```text
User role = support_agent without wallet.topup.approve.
```

Action:
```text
Approve top-up request.
```

Expected:
```text
FORBIDDEN.
No ledger entry.
Audit permission denied optional.
```

### RBAC-002 Finance cannot reveal credential

Action:
```text
Finance agent calls credential reveal.
```

Expected:
```text
FORBIDDEN.
No plaintext returned.
No reveal audit except denied event optional.
```

### RBAC-003 Reseller staff only with permission can update plan price

Expected:
```text
Staff without catalog.tenant.price_update cannot update.
Staff with permission can update own tenant plan only.
```

### RBAC-004 Admin emergency access requires reason

Action:
```text
Admin attempts tenant emergency access without reason.
```

Expected:
```text
VALIDATION_ERROR.
No emergency session.
```

---

## 6. Wallet and ledger test cases

### WL-001 Top-up approval creates exactly one ledger entry

Setup:
```text
Client submits top-up 100.
```

Action:
```text
Reseller approves top-up.
```

Expected:
```text
topup status approved.
one credit ledger entry amount 100.
wallet balance +100.
audit wallet.topup.approved.
```

### WL-002 Top-up cannot be approved twice

Action:
```text
Approve same top-up twice.
```

Expected:
```text
second attempt returns TOPUP_ALREADY_REVIEWED.
no second ledger entry.
balance unchanged.
```

### WL-003 Ledger posted cannot be edited

Action:
```text
Attempt update/delete posted ledger entry through API/admin.
```

Expected:
```text
forbidden/not supported.
adjustment required.
```

### WL-004 Manual adjustment requires reason

Action:
```text
Finance creates adjustment without reason.
```

Expected:
```text
VALIDATION_ERROR.
No ledger entry.
```

### WL-005 Client wallet sufficient but reseller wallet insufficient

Setup:
```text
Client wallet = 100.
Plan selling_price = 20.
Reseller settlement wallet = 5.
Reseller cost = 12.
```

Action:
```text
Client checkout.
```

Expected:
```text
Error INSUFFICIENT_RESELLER_BALANCE
No provisioning job.
No service.
No stock allocated.
Client wallet not debited, unless policy explicitly creates pending order without debit.
```

---

## 7. Checkout and reservation test cases

### CO-001 Successful checkout creates order/reservation/ledger/job

Expected:
```text
order created.
order_item snapshots stored.
reservation status reserved then allocated after provisioning.
client wallet debited.
reseller wallet debited if reseller tenant.
provisioning_job queued.
correlation_id present.
```

### CO-002 Double click checkout with same idempotency key

Action:
```text
Submit same checkout twice same idempotency key.
```

Expected:
```text
Only one order.
Only one wallet debit.
Only one reservation.
Only one provisioning job.
Second response returns same result.
```

### CO-003 Same idempotency key with different payload

Expected:
```text
IDEMPOTENCY_CONFLICT.
No new order.
```

### CO-004 Out of stock prevents checkout

Setup:
```text
source available = 0.
```

Expected:
```text
OUT_OF_STOCK.
No wallet debit.
No reservation.
No provisioning job.
```

### CO-005 Reservation expires

Setup:
```text
reservation created but payment/provisioning not completed.
expires_at passed.
```

Action:
```text
reservation_expiry_job runs.
```

Expected:
```text
reservation expired.
reserved_count decremented once.
order expired/cancelled.
audit reservation.expired.
```

### CO-006 Concurrent checkout for last stock

Setup:
```text
source capacity available = 1.
Two clients checkout same plan concurrently.
```

Expected:
```text
Only one succeeds.
Other gets OUT_OF_STOCK.
No oversell.
reserved/allocated counts correct.
```

---

## 8. Provisioning test cases

### PR-001 Provider success activates service

Expected:
```text
provider_request success.
external_resource_id stored.
service active.
credential encrypted.
reservation allocated.
order provisioning_status success.
activation notification queued.
```

### PR-002 Provider out of stock after reservation

Expected:
```text
job failed do_not_retry.
reservation released.
wallet reversal/refund per policy.
order failed.
notification sent.
```

### PR-003 Provider timeout unknown

Action:
```text
Adapter returns timeout unknown after create request sent.
```

Expected:
```text
job manual_review.
No automatic retry.
Order provisioning_status manual_review.
Operator alert.
No duplicate provider create.
```

### PR-004 Provider auth failure

Expected:
```text
source marked degraded/down or disabled per policy.
job failed/manual_review.
alert admin.
no retry.
```

### PR-005 Success but credential missing

Expected:
```text
job manual_review or credential fetch attempt if safe.
service not shown as fully active until credential ready, unless policy allows active_pending_credential.
```

### PR-006 Manual provider activation

Setup:
```text
manual source.
```

Expected:
```text
checkout creates paid order + reservation + manual_review job.
operator enters resource/credential.
service active.
credential encrypted.
audit manual resolution.
```

---

## 9. Service lifecycle test cases

### SL-001 Renew active service

Setup:
```text
service active, term_end_at future.
wallet sufficient.
```

Expected:
```text
wallet debit.
term_end_at extended from old term_end_at.
lifecycle event service.renewed.
```

### SL-002 Renew expired service in grace

Expected according to policy:
```text
term_end_at extends from old term_end_at if policy says no free days.
service unsuspended/active if required.
```

### SL-003 Cannot renew terminated service

Expected:
```text
VALIDATION_ERROR or SERVICE_NOT_RENEWABLE.
No wallet debit.
```

### SL-004 Auto expiry

Action:
```text
service_expiry_job runs after term_end_at.
```

Expected:
```text
billing_status overdue/grace.
notification queued.
no immediate terminate unless policy.
```

### SL-005 Auto suspend after grace

Expected:
```text
suspend job created.
service suspended.
suspension_reason = billing_overdue.
audit service.suspended.
```

### SL-006 Terminate after hold period

Expected:
```text
terminate job created only if policy auto_terminate true.
reason required.
provider terminate result handled.
```

---

## 10. Credential security test cases

### CR-001 Service detail returns masked credential only

Expected:
```text
no plaintext password/token in response.
masked_hint present.
```

### CR-002 Reveal credential creates audit

Action:
```text
Client reveals own credential.
```

Expected:
```text
plaintext returned only in reveal response.
audit credential.revealed.
last_revealed_at updated.
```

### CR-003 Unauthorized reveal denied

Action:
```text
Client tries reveal credential of another service.
```

Expected:
```text
403/404.
No plaintext.
No last_revealed update.
```

### CR-004 Logs/audit do not contain secret

Action:
```text
Search provider_request/audit/notification payload after activation.
```

Expected:
```text
no root password/proxy password/API key plaintext.
```

---

## 11. Catalog/pricing test cases

### CP-001 Reseller cannot sell disabled master plan

Expected:
```text
PLAN_DISABLED.
```

### CP-002 Price snapshot unaffected by future price update

Setup:
```text
Client buys plan at 20.
Admin later changes suggested price to 25.
```

Expected:
```text
existing order price_snapshot remains 20.
reports for order use 20.
```

### CP-003 Margin risk warning

Setup:
```text
reseller selling_price < reseller_cost.
```

Expected:
```text
plan status margin_risk or update blocked depending policy.
checkout blocked if policy requires non-negative margin.
```

### CP-004 Capability masking

Setup:
```text
source does not support change_ip.
```

Expected:
```text
UI hides change_ip.
API returns CAPABILITY_NOT_SUPPORTED if called.
```

---

## 12. Abuse/risk test cases

### AB-001 New high-risk order routes manual review

Setup:
```text
risk rule triggered.
```

Expected:
```text
order manual_review before provisioning.
no provider job until approved.
```

### AB-002 Abuse suspend requires reason/evidence

Action:
```text
Admin suspends service for abuse without reason.
```

Expected:
```text
VALIDATION_ERROR.
```

### AB-003 Abuse case workflow

Expected:
```text
case open -> investigating -> warning/suspended/resolved.
audit all transitions.
client notification if policy.
```

---

## 13. Notification test cases

### NT-001 Activation notification

Expected:
```text
sent/queued after service active.
does not contain plaintext credential unless policy explicitly allows.
contains link to service detail.
```

### NT-002 Expiry reminder dedupe

Action:
```text
renewal_reminder_job runs twice same day/window.
```

Expected:
```text
only one notification for same service/window.
```

### NT-003 Reseller low balance warning

Expected:
```text
notification sent when settlement wallet below threshold.
not spammed repeatedly beyond dedupe policy.
```

---

## 14. Report test cases

### RP-001 Reseller profit calculation

Setup:
```text
client purchase selling_price 20.
reseller_cost 12.
```

Expected:
```text
gross revenue 20.
platform cost 12.
gross profit 8 before refund/adjustment.
```

### RP-002 Refund affects report

Expected:
```text
refund/reversal reflected.
uses ledger entries and snapshots.
```

### RP-003 Admin provider health report

Expected:
```text
failed/manual_review/success count matches provisioning_jobs/provider_requests.
```

---

## 15. Deployment smoke tests

Before production:
- Login/register works per tenant domain.
- Client top-up submit works.
- Top-up approval posts ledger.
- Checkout success with manual/mock provider.
- Service active detail shows masked credential.
- Credential reveal works and audits.
- Cross-tenant access blocked.
- Reservation expiry job works.
- Service expiry reminder works.
- Provider health job reports.
- Backup job completed and restore tested in staging.

---

## 16. Acceptance gates

### Gate 1: Foundation

Pass if:
```text
tenant isolation tests pass
RBAC high-risk tests pass
wallet ledger tests pass
```

### Gate 2: Commerce

Pass if:
```text
checkout/reservation/idempotency tests pass
reseller settlement tests pass
refund/reversal tests pass
```

### Gate 3: Provisioning

Pass if:
```text
provider success/fail/timeout/manual tests pass
credential security tests pass
service lifecycle tests pass
```

### Gate 4: Production readiness

Pass if:
```text
notification/report tests pass
monitoring/backup/restore smoke tests pass
abuse/manual review tests pass
```

---

## 17. QA acceptance criteria

Project không nên production nếu fail bất kỳ P0:
- Cross-tenant access leaks data.
- Ledger double posts.
- Checkout double debits.
- Reservation oversells.
- Provider timeout creates duplicate resources.
- Credential plaintext appears in logs/audit.
- Reseller client order provisions while reseller wallet lacks cost.
- Staff without permission can approve money/reveal credential.
- Terminated service can be renewed accidentally.
- Backup/restore not tested.

Câu nền: **QA của hệ thống hạ tầng không phải tìm lỗi giao diện; QA là đóng những cánh cửa nơi tiền, tenant và credential có thể rò ra ngoài.**

---


# ===== 22_Deployment_DevOps_And_Environment_Runbook.md =====

# 22 - Deployment DevOps And Environment Runbook

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa runbook triển khai và vận hành môi trường cho platform VPS/Proxy.

Đây không phải tài liệu chọn framework/cloud cụ thể. Mục tiêu là khóa:
- Environment cần có.
- Secret quản lý thế nào.
- Database/queue/worker triển khai thế nào.
- Backup/restore cần gì.
- Monitoring/alert cần gì.
- Release/rollback ra sao.
- Checklist trước production.

---

## 2. Environment model

### 2.1 Local development

Dành cho dev build/test chức năng.

Cho phép:
```text
mock provider
manual provider
sandbox email
fake payment proof
seed tenant data
```

Không cho phép:
```text
provider production API key
production database copy chưa mask dữ liệu
real credential plaintext
```

### 2.2 Staging

Dành cho QA/UAT.

Yêu cầu:
```text
gần giống production
có queue/worker/scheduler thật
có test provider/sandbox provider
có test email/telegram sandbox
có backup/restore test
```

Staging phải test được:
- top-up flow.
- checkout.
- provisioning success/fail/timeout.
- tenant isolation.
- credential reveal.
- expiry/suspend cron.

### 2.3 Production

Dành cho khách thật.

Yêu cầu:
```text
secret management chuẩn
backup tự động
monitoring/alert
audit retention
rate limit
2FA admin
access control chặt
rollback plan
```

---

## 3. Logical deployment units

Không bắt buộc microservice phase 1. Khuyến nghị:
```text
modular monolith API
frontend portals
worker process
scheduler/cron process
database
queue/cache
object storage
log/monitoring stack
secret manager
```

### 3.1 Frontend

Có thể deploy một app hoặc tách:
```text
admin portal
reseller portal
client storefront/portal
```

Requirements:
- tenant branding/domain support.
- no secret in frontend.
- environment-specific API endpoint.
- CSP/security headers recommended.

### 3.2 Backend API

Responsibilities:
- auth.
- tenant context.
- RBAC.
- catalog/order/wallet/service APIs.
- audit.
- create queue jobs.

Backend không nên giữ long-running provider call trong request user.

### 3.3 Worker

Responsibilities:
- provisioning.
- provider sync.
- service actions.
- notification send.
- export/report long-running tasks.

Worker phải chạy cùng version với backend API hoặc có compatibility contract.

### 3.4 Scheduler/Cron

Responsibilities:
- reservation expiry.
- service expiry.
- suspension/termination.
- renewal reminder.
- provider health.
- inventory sync.
- retention.

Scheduler cần lock để tránh nhiều instance chạy cùng job một lúc.

---

## 4. Environment variables / config groups

Không ghi secret thật trong tài liệu/repo.

Config groups:
```text
APP_ENV
APP_URL
ADMIN_URL
DATABASE_URL
QUEUE_URL
CACHE_URL
OBJECT_STORAGE_CONFIG
EMAIL_PROVIDER_CONFIG
TELEGRAM_BOT_CONFIG
ENCRYPTION_KEY_REFERENCE
JWT/SESSION_SECRET_REFERENCE
PROVIDER_SECRET_REFERENCE
PAYMENT_METHOD_CONFIG
RATE_LIMIT_CONFIG
LOGGING_CONFIG
```

### 4.1 Secret classification

Critical:
```text
database password
session/JWT secret
encryption master key
provider API key/secret
email provider secret
telegram bot token
object storage secret
```

Sensitive:
```text
payment instructions internal refs
support webhook
backup storage credential
```

Public:
```text
public app URL
brand assets
feature flags non-sensitive
```

Rule:
```text
critical secrets never committed.
critical secrets rotate-able.
critical secrets not visible to normal developers in production.
```

---

## 5. Secret management

### 5.1 Required controls

- Provider credentials encrypted at rest.
- Service credentials encrypted at rest.
- Master encryption key stored outside database.
- Secret rotation plan.
- Access to secret manager audited.
- Staging/prod secrets separated.

### 5.2 Encryption key rotation logic

Minimum design:
```text
secret_version stored with encrypted payload
new credentials encrypted with current key version
old credentials decryptable until rotated
rotation job can re-encrypt if needed
```

### 5.3 Never log

```text
Authorization header
session token
provider API key
root password
proxy password
encryption key
payment proof private URL if sensitive
```

---

## 6. Database operations

### 6.1 Migration rule

Every schema change needs:
```text
migration description
backward compatibility note
rollback note
data migration risk
tested on staging
```

Financial/tenant tables require extra review:
```text
wallets
wallet_ledger_entries
orders
reservations
services
service_credentials
audit_logs
provider_resource_mappings
```

### 6.2 Migration safety

Before production migration:
- backup completed.
- migration tested on staging copy.
- long-running migration estimated.
- rollback or forward-fix plan.
- maintenance window if needed.

### 6.3 Data integrity checks

Periodic checks:
```text
wallet balance cache equals ledger sum
orders paid have ledger entries
services active have provider mapping
reservations allocated have service
provider mappings unique
credential rows encrypted
ledger entries immutable
```

---

## 7. Queue/worker deployment

### 7.1 Worker scaling

Workers can scale horizontally if:
```text
job locking prevents duplicate processing
idempotency enforced
provider rate limits respected
```

### 7.2 Job locking

Each job should be claimed by one worker:
```text
queued -> running with atomic lock
heartbeat/timeout for stuck jobs
stuck running job recovery policy
```

### 7.3 Stuck job policy

If job running too long:
```text
mark as stale
do not blindly rerun provider create
move to manual_review if provider action may have been sent
```

### 7.4 Provider rate limit

Provider-specific worker throttles:
```text
max concurrent jobs per source
min delay/backoff on 429
source maintenance mode
```

---

## 8. Object storage

Used for:
```text
payment proof attachments
abuse evidence files
private export files
possibly redacted provider raw response
brand assets
```

Rules:
- Private files require signed/authorized access.
- Payment proofs and abuse evidence not public.
- Brand assets can be public.
- Retention policy by file type.
- Virus/malware scan if attachments from users.

---

## 9. Logging and monitoring

### 9.1 Application logs

Logs should include:
```text
timestamp
level
request_id
correlation_id
tenant_id when safe
actor_id when safe
module
action
error_code
```

No plaintext secret.

### 9.2 Metrics

P0 metrics:
```text
api error rate
api latency
login failures
checkout success/fail rate
wallet ledger posting failures
provisioning success/fail/manual_review
queue backlog by job type
provider health
provider latency/error
reservation expiry count
service expiry/suspend/terminate count
notification failure rate
credential reveal count
ledger adjustment count
```

### 9.3 Alerts

Critical alerts:
```text
database down
queue down
worker not running
provider down repeated
provisioning manual_review spike
wallet ledger posting error
checkout failure spike
backup failed
credential reveal spike
admin login brute-force spike
```

Warning alerts:
```text
reseller low balance
source stock low
top-up pending too long
provider degraded
email notification failure
```

---

## 10. Backup and restore

### 10.1 Backup scope

Must back up:
```text
database
object storage private files
encryption key references/secret manager backup procedure
configuration snapshots
```

Database backup is useless if encryption keys are lost.

### 10.2 Backup frequency

Suggested:
```text
database: automated daily + point-in-time if possible
object storage: daily/incremental
config/secret references: whenever changed
```

### 10.3 Restore test

At least before production and periodically:
```text
restore staging from backup
verify users/tenants/orders/ledger/services
verify encrypted credential can decrypt with key reference
verify audit logs accessible
verify provider mappings intact
```

### 10.4 Restore acceptance

Restore is valid only if:
- Ledger balances reconcile.
- Services still map to provider resources.
- Credentials decrypt.
- Audit logs preserved.
- Tenant domains/config present.
- App can boot with restored config.

---

## 11. Release process

### 11.1 Release checklist

Before release:
- changelog written.
- migrations reviewed.
- staging tests passed.
- QA P0 passed.
- backup completed.
- rollback plan documented.
- monitoring ready.
- support/admin notified if needed.

### 11.2 Deployment order

Typical order:
```text
1. backup
2. run compatible DB migration
3. deploy backend API
4. deploy worker
5. deploy scheduler
6. deploy frontend
7. run smoke tests
8. monitor
```

If breaking change:
```text
use maintenance window or two-phase migration
```

### 11.3 Rollback

Rollback plan must define:
- app rollback.
- worker rollback.
- migration rollback/forward-fix.
- queue compatibility.
- jobs in progress handling.
- provider calls already sent.

Important:
```text
Không rollback bừa nếu migration đã thay đổi financial ledger semantics.
Trong hệ thống tiền, forward-fix đôi khi an toàn hơn rollback.
```

---

## 12. Incident response

### 12.1 Incident severity

| Severity | Example | Response |
|---|---|---|
| SEV1 | wallet/ledger wrong, credential leak, tenant data leak | immediate freeze affected functions |
| SEV2 | provider outage, provisioning failing widely | disable source, manual review |
| SEV3 | notification delayed, report wrong | fix within normal SLA |
| SEV4 | UI cosmetic | backlog |

### 12.2 Immediate actions

For wallet/ledger issue:
```text
pause checkout if needed
pause top-up approval if needed
snapshot affected ledger/orders
investigate correlation_id
create adjustment only after root cause known
```

For credential leak:
```text
disable reveal if needed
rotate affected credentials if possible
notify impacted users per policy
audit access logs
```

For provider duplicate resource:
```text
stop retry for source
list provider resources by time/correlation metadata
link or terminate duplicates carefully
adjust ledger/refund if needed
```

For tenant data leak:
```text
disable affected endpoint
audit access
notify affected tenant per policy/legal need
fix tenant guard
add regression test
```

---

## 13. Security hardening checklist

Minimum production:
- Admin 2FA required.
- Reseller owner 2FA enabled/default.
- Strong password policy.
- Rate limit login/register/checkout/reveal.
- CSP/security headers.
- HTTPS only.
- Secure cookies/session.
- Provider API IP allowlist if possible.
- Secrets not in repo/logs.
- DB access restricted.
- Audit critical actions.
- Backup encrypted.
- Least privilege staff roles.

---

## 14. Go-live checklist

Do not go production until:
- Tenant isolation QA passed.
- Wallet/ledger QA passed.
- Checkout/reservation QA passed.
- Provisioning success/fail/timeout QA passed.
- Credential security QA passed.
- Backup/restore tested.
- Monitoring/alert configured.
- Admin 2FA active.
- Provider source tested with small live resource.
- Abuse/manual suspend SOP ready.
- Support flow ready.
- Terms/AUP visible on storefront.

---

## 15. Runbook acceptance criteria

Deployment/runbook đạt khi:
- Có environment local/staging/prod rõ.
- Secret không commit, không log.
- Worker/scheduler có deployment riêng.
- Backup/restore test được.
- Monitoring/alert bao phủ money/provisioning/provider/credential.
- Release có checklist và rollback/forward-fix plan.
- Incident response có bước pause/freeze cho chức năng nguy hiểm.
- Production không chạy nếu chưa pass P0 QA.

Câu nền: **production không phải nơi để chứng minh dev chạy được; production là nơi chứng minh hệ thống vẫn đúng khi mọi thứ xung quanh bắt đầu sai.**

---


# ===== 23_Notification_Email_Telegram_Template_Spec.md =====

# 23 - Notification Email Telegram Template Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa notification events, channels và template logic cho nền tảng VPS/Proxy.

Notification không chỉ để “gửi email cho đẹp”. Nó giúp:
- giảm support.
- giảm tranh chấp tiền.
- cảnh báo trước khi mất dịch vụ.
- báo lỗi provisioning/provider sớm.
- nhắc reseller nạp ví settlement.
- lưu dấu vận hành.

---

## 2. Notification channels

### 2.1 Dashboard notification

Dùng cho:
```text
mọi user đã đăng nhập
client/reseller/admin inbox
```

Ưu điểm:
- không phụ thuộc email.
- giữ lịch sử trong hệ thống.
- có thể gắn link đến order/service/top-up.

### 2.2 Email

Dùng cho:
```text
registration/verification
top-up status
order/service status
expiry reminder
abuse warning
support update
```

Không nên gửi plaintext credential qua email mặc định.

### 2.3 Telegram/Admin alert

Dùng cho:
```text
provider down
provisioning failed/manual review
reseller low balance
abuse critical
queue backlog
backup failed
```

Telegram cho admin/reseller team nên chứa correlation_id và link admin/reseller portal.

### 2.4 Webhook optional

Phase sau hoặc reseller advanced:
```text
service activated
service expiring
top-up approved
order failed
```

---

## 3. Notification event naming

Format:
```text
module.event
```

Examples:
```text
auth.email_verification
wallet.topup.submitted
wallet.topup.approved
wallet.topup.rejected
wallet.reseller_low_balance
order.created
order.failed
order.manual_review
service.activated
service.expiring
service.expired
service.suspended
service.terminated
credential.revealed_alert_optional
provisioning.failed
provisioning.manual_review
provider.down
provider.recovered
abuse.warning
abuse.suspended
support.ticket_created
support.ticket_replied
```

---

## 4. Notification payload rules

Payload must include:
```text
tenant_id
recipient_user_id or recipient_group
template_key
channel
reference_type
reference_id
correlation_id
dedupe_key optional
priority
```

Payload must not include:
```text
plaintext password
provider API key
private token
full payment proof URL if public
sensitive abuse evidence
```

---

## 5. Priority levels

| Priority | Meaning | Examples |
|---|---|---|
| Low | Informational | catalog update, general support note |
| Normal | User flow status | top-up submitted, order created |
| High | Action needed | service expiring, top-up rejected |
| Critical | Operational risk | provider down, abuse critical, backup failed |

Critical notifications should go to dashboard + admin Telegram/email depending config.

---

## 6. Dedupe policy

Some events must not spam.

Dedupe examples:
```text
service_expiring:{service_id}:{window}
reseller_low_balance:{tenant_id}:{date}
provider_down:{source_id}:{hour}
queue_backlog:{job_type}:{hour}
topup_pending:{topup_request_id}:{day}
```

Critical provider alert can repeat after threshold:
```text
first alert immediately
repeat every N minutes/hours while unresolved
send recovered notification
```

---

## 7. Template variables

Common variables:
```text
tenant_name
brand_name
user_name
order_number
service_name
service_id_short
plan_name
amount
currency
wallet_balance
term_end_at
days_remaining
support_link
dashboard_link
correlation_id
```

Admin variables:
```text
source_name
provider_type
job_id
error_code
manual_review_reason
queue_backlog_count
tenant_name
reseller_name
```

---

## 8. Client templates

### 8.1 Email verification

Event:
```text
auth.email_verification
```

Subject:
```text
Verify your account for {brand_name}
```

Body:
```text
Hi {user_name},

Please verify your email address to activate your account on {brand_name}.

Verification link:
{verification_link}

If you did not create this account, you can ignore this message.
```

Channels:
```text
email
```

### 8.2 Top-up submitted

Event:
```text
wallet.topup.submitted
```

Subject:
```text
Top-up request received: {amount} {currency}
```

Body:
```text
Hi {user_name},

We received your wallet top-up request for {amount} {currency}.
Your request is now waiting for review.

Reference:
{topup_reference}

You will be notified once it is approved or rejected.
```

Channels:
```text
email optional
dashboard
```

### 8.3 Top-up approved

Event:
```text
wallet.topup.approved
```

Subject:
```text
Wallet top-up approved
```

Body:
```text
Hi {user_name},

Your top-up of {amount} {currency} has been approved.
Your updated wallet balance is {wallet_balance} {currency}.

You can now purchase or renew services from your dashboard.
```

Channels:
```text
email
dashboard
```

### 8.4 Top-up rejected

Event:
```text
wallet.topup.rejected
```

Subject:
```text
Wallet top-up could not be approved
```

Body:
```text
Hi {user_name},

Your top-up request for {amount} {currency} could not be approved.

Reason:
{review_note}

Please check your payment details or contact support if you believe this is a mistake.
```

Channels:
```text
email
dashboard
```

### 8.5 Order created / provisioning queued

Event:
```text
order.created
```

Subject:
```text
Order {order_number} has been created
```

Body:
```text
Hi {user_name},

Your order {order_number} has been created and is being processed.

Plan:
{plan_name}

Status:
{order_status}

You can track the order from your dashboard.
```

Channels:
```text
dashboard
email optional
```

### 8.6 Service activated

Event:
```text
service.activated
```

Subject:
```text
Your service is active: {service_name}
```

Body:
```text
Hi {user_name},

Your service is now active.

Service:
{service_name}

Plan:
{plan_name}

Expiry:
{term_end_at}

For security, login to your dashboard to reveal credentials.
```

Important:
```text
Do not include plaintext password by default.
```

Channels:
```text
email
dashboard
```

### 8.7 Provisioning failed

Event:
```text
order.failed
```

Subject:
```text
Order {order_number} could not be provisioned
```

Body:
```text
Hi {user_name},

We could not provision your order {order_number}.

Status:
{failure_status}

If payment was already captured from your wallet, a refund/reversal will be handled according to the platform policy.

Reference:
{correlation_id}
```

Channels:
```text
email
dashboard
```

### 8.8 Service expiring reminder

Event:
```text
service.expiring
```

Subject:
```text
Your service expires in {days_remaining} day(s)
```

Body:
```text
Hi {user_name},

Your service {service_name} will expire on {term_end_at}.

Please renew before expiry to avoid suspension or termination according to the service policy.

Renew here:
{service_link}
```

Channels:
```text
email
dashboard
```

Dedupe:
```text
service_expiring:{service_id}:{days_remaining_window}
```

### 8.9 Service expired

Event:
```text
service.expired
```

Subject:
```text
Your service has expired: {service_name}
```

Body:
```text
Hi {user_name},

Your service {service_name} expired on {term_end_at}.

Depending on the product policy, the service may enter a grace period before suspension or termination.
Please renew as soon as possible if you want to keep using it.
```

### 8.10 Service suspended

Event:
```text
service.suspended
```

Subject:
```text
Your service has been suspended
```

Body:
```text
Hi {user_name},

Your service {service_name} has been suspended.

Reason:
{suspension_reason}

Please contact support or resolve the related issue to request reactivation.
```

### 8.11 Service terminated

Event:
```text
service.terminated
```

Subject:
```text
Your service has been terminated
```

Body:
```text
Hi {user_name},

Your service {service_name} has been terminated.

Reason:
{termination_reason}

Terminated services may not be recoverable. Please contact support if you need clarification.
```

### 8.12 Abuse warning

Event:
```text
abuse.warning
```

Subject:
```text
Important notice about your service
```

Body:
```text
Hi {user_name},

We received an abuse or policy notice related to your service {service_name}.

Issue:
{abuse_case_type}

Please review and resolve this immediately. Continued violations may lead to suspension or termination.
```

---

## 9. Reseller templates

### 9.1 Reseller wallet top-up submitted

Event:
```text
wallet.reseller_topup.submitted
```

Recipient:
```text
reseller owner/staff
platform finance/admin alert optional
```

Body:
```text
Your reseller wallet top-up request for {amount} {currency} has been submitted and is waiting for platform review.
```

### 9.2 Reseller top-up approved

Body:
```text
Your reseller settlement wallet has been credited with {amount} {currency}.
New client orders can continue to provision as long as your balance covers reseller cost.
```

### 9.3 Reseller low balance

Event:
```text
wallet.reseller_low_balance
```

Subject:
```text
Your reseller wallet balance is low
```

Body:
```text
Hi {reseller_name},

Your reseller settlement wallet balance is currently {wallet_balance} {currency}.

New client orders may not be provisioned if your balance is lower than the platform reseller cost.

Please top up your reseller wallet to avoid checkout/provisioning interruption.
```

Channels:
```text
email
dashboard
telegram optional
```

Dedupe:
```text
reseller_low_balance:{tenant_id}:{date}
```

### 9.4 Client top-up pending review

Event:
```text
wallet.client_topup.pending_review
```

Body:
```text
A client top-up request is waiting for review.

Client:
{client_name}

Amount:
{amount} {currency}

Reference:
{topup_reference}
```

### 9.5 Reseller plan margin risk

Event:
```text
catalog.margin_risk
```

Body:
```text
One or more plans in your catalog may have low or negative margin.

Plan:
{plan_name}

Selling price:
{selling_price} {currency}

Reseller cost:
{reseller_cost} {currency}

Please update your selling price or disable the plan.
```

---

## 10. Admin templates

### 10.1 Provider down

Event:
```text
provider.down
```

Subject:
```text
Provider/source down: {source_name}
```

Telegram/dashboard body:
```text
Provider/source health check failed.

Source:
{source_name}

Provider:
{provider_type}

Error:
{error_code}

Last check:
{last_health_check_at}

Impact:
New provisioning through this source may fail or be paused.

Correlation:
{correlation_id}
```

Channels:
```text
telegram
dashboard
email optional
```

### 10.2 Provider recovered

Event:
```text
provider.recovered
```

Body:
```text
Provider/source has recovered.

Source:
{source_name}

Previous status:
{previous_status}

Current status:
healthy
```

### 10.3 Provisioning manual review

Event:
```text
provisioning.manual_review
```

Body:
```text
A provisioning job requires manual review.

Job:
{job_id}

Tenant:
{tenant_name}

Order:
{order_number}

Source:
{source_name}

Reason:
{manual_review_reason}

Retry safety:
{retry_safety}

Correlation:
{correlation_id}
```

### 10.4 Provisioning failed

Event:
```text
provisioning.failed
```

Body:
```text
Provisioning failed.

Job:
{job_id}

Order:
{order_number}

Source:
{source_name}

Error:
{error_code}

Attempts:
{attempt_count}/{max_attempts}

Next action:
{recommended_action}
```

### 10.5 Queue backlog

Event:
```text
system.queue_backlog
```

Body:
```text
Queue backlog threshold exceeded.

Job type:
{job_type}

Backlog:
{queue_backlog_count}

Oldest job age:
{oldest_job_age}

Please check worker health and provider status.
```

### 10.6 Backup failed

Event:
```text
system.backup_failed
```

Body:
```text
Production backup failed.

Backup type:
{backup_type}

Environment:
{environment}

Error:
{error_code}

Immediate action required.
```

### 10.7 Ledger adjustment created

Event:
```text
wallet.adjustment.created
```

Body:
```text
A wallet ledger adjustment was created.

Actor:
{actor_name}

Wallet:
{wallet_reference}

Amount:
{amount} {currency}

Direction:
{direction}

Reason:
{reason}

Correlation:
{correlation_id}
```

Critical audit/notification for finance/admin.

### 10.8 Credential reveal spike

Event:
```text
security.credential_reveal_spike
```

Body:
```text
Credential reveal activity is higher than normal.

Tenant:
{tenant_name}

Actor:
{actor_name}

Count:
{count}

Window:
{time_window}

Please review audit logs.
```

---

## 11. Notification timing matrix

| Event | Client | Reseller | Admin |
|---|---|---|---|
| top-up submitted | immediate dashboard/email | if client top-up: reseller alert | optional |
| top-up approved | immediate | optional | no |
| order created | immediate | optional tenant order alert | no |
| service activated | immediate | optional | no |
| provisioning failed | client if final failed | reseller if tenant client affected | immediate |
| manual review | maybe pending status only | optional | immediate |
| service expiring | 7/3/1 days | reseller summary optional | no |
| service expired | immediate | optional | no |
| service suspended | immediate | optional | if abuse/admin action |
| provider down | no | if affects reseller maybe optional | immediate |
| reseller low balance | no | immediate/daily dedupe | optional |
| abuse warning | immediate | reseller owner | admin abuse queue |

---

## 12. Template localization

Phase 1 can support one default language. Recommended structure:
```text
template_key
language
subject
body
channel
variables_schema
```

If platform targets international reseller/client:
```text
en
vi
```

Tenant can override safe text:
```text
brand greeting
footer
support link
```

Tenant should not override security/legal core wording unless approved.

---

## 13. Notification audit

Important notification events should create audit or notification record:
```text
notification.queued
notification.sent
notification.failed
```

For financial/security:
```text
top-up approved
ledger adjustment
credential reveal alert
provider down
abuse warning
```

Store:
```text
recipient
channel
template_key
reference_id
status
sent_at
error redacted
correlation_id
```

---

## 14. Security rules

Do not send:
```text
root password
proxy password
provider secret
private abuse evidence
payment proof file without authorized link
session/auth token
```

Use links:
```text
login to dashboard to reveal credentials
view top-up request
view service
view admin job
```

Links should expire or require auth.

---

## 15. Notification acceptance criteria

Notification spec đạt khi:
- Mỗi lifecycle tiền/order/service có notification hợp lý.
- Credential không gửi plaintext mặc định.
- Admin nhận provider/provisioning/queue/backup critical alerts.
- Reseller nhận low balance và pending client top-up.
- Expiry reminders có dedupe.
- Payload redacted.
- Notification records trace được bằng correlation_id.
- Templates có biến rõ, không phụ thuộc text hardcode rải rác.
- Critical failure có “next action” trong message.

Câu nền: **notification tốt là hệ thần kinh của platform: nó không làm thay việc vận hành, nhưng nó báo đau đúng chỗ trước khi cơ thể bị thương nặng.**

---


# ===== CHANGELOG_FIXES.md =====

# CHANGELOG_FIXES - Bản vá tài liệu v1.1

## Mục tiêu bản vá
Bản v1.1 vá các lỗ hổng P0 trong bộ tài liệu gốc để biến blueprint thành build-spec rõ hơn cho dev/ops, vẫn giữ phạm vi “chưa code”.

## Các lỗi/thiếu sót đã vá
1. Thêm mô hình reseller settlement: client wallet, reseller wallet, platform revenue, reseller profit.
2. Khóa rule không provision nếu reseller wallet thiếu reseller cost.
3. Bổ sung tenant enforcement: không tin tenant_id từ body, backend bắt buộc scope theo tenant context.
4. Bổ sung emergency access có reason và audit.
5. Chuyển 2FA Admin thành P0 phase 1; Reseller Owner bật mặc định/khuyến nghị bắt buộc.
6. Bổ sung catalog versioning, snapshot, propagation rule và margin guard.
7. Bổ sung top-up state machine và financial invariants.
8. Bổ sung idempotency key, partial success, retry safety và no retry mù.
9. Bổ sung atomic inventory reservation, reserved/allocated counters và expiry release.
10. Bổ sung credential security: encrypt, masked reveal, audit, redaction.
11. Bổ sung renew/cancel/refund guard và calendar month rule.
12. Bổ sung acceptance criteria cho flow/API/data model.
13. Bổ sung manual abuse/fraud controls phase 1.
14. Bổ sung audit naming, redaction và correlation_id xuyên flow.
15. Bổ sung report formula cho Admin và Reseller.

## File mới được thêm
- `09_Reseller_Settlement_Ledger_Model.md`
- `10_Tenant_Security_Access_Control_Spec.md`
- `11_Provisioning_Idempotency_And_Inventory_Locking.md`
- `12_API_Data_Model_Acceptance_Criteria.md`
- `13_Abuse_Fraud_Operational_Policy_Phase1.md`

## File cũ đã được vá
- `00_README.md`
- `01_Product_Scope_Business_Model.md`
- `02_Tenant_Model_Role_Architecture.md`
- `03_Product_Catalog_Pricing_Rules.md`
- `04_Billing_Wallet_Ledger_Spec.md`
- `05_Provisioning_Provider_Adapter_Spec.md`
- `06_Order_Service_Lifecycle_State_Machine.md`
- `07_Portal_Functional_Spec.md`
- `08_Audit_Reports_Operational_Control.md`

## Gợi ý dùng với dev
Đưa dev đọc theo thứ tự:
1. `00_README.md`
2. `CHANGELOG_FIXES.md`
3. `09`, `10`, `11` trước vì đây là ba tài liệu P0 nhất.
4. `12` để chuyển sang backlog/API/data model.
5. `13` để khóa vận hành rủi ro trước production.

## Bổ sung bản v1.2 - Technical Handoff
Bản v1.2 không thay đổi nguyên tắc vá lỗi P0 của v1.1, mà bổ sung bộ tài liệu `14–23` để dev/backend/frontend/QA/DevOps có thể chuyển blueprint thành kế hoạch build cụ thể hơn.

Các nhóm tài liệu mới:
- Architecture
- Database schema/ERD
- API contract
- RBAC matrix
- Provider adapter contract
- Worker/queue/cron
- UI wireflow
- QA acceptance plan
- DevOps runbook
- Notification templates

---


# ===== CHANGELOG_TECHNICAL_HANDOFF_v1_2.md =====

# CHANGELOG_TECHNICAL_HANDOFF_v1_2

## Mục tiêu bản v1.2
Bản v1.2 bổ sung bộ Technical Build Handoff Package để chuyển bộ blueprint/vá lỗi v1.1 thành tài liệu dev/backend/frontend/QA/DevOps có thể dùng để build.

Bản này vẫn giữ nguyên nguyên tắc:
- Chưa code.
- Không chọn framework cụ thể.
- Ưu tiên khóa behavior, data contract, permission, state, queue, QA và vận hành.

## File mới được thêm

### `14_System_Architecture_Blueprint.md`
Bổ sung kiến trúc tổng thể:
- portal layer.
- API/auth/tenant/RBAC.
- core modules.
- financial core.
- security core.
- queue/worker/cron.
- provider adapter layer.
- observability và fail-safe principle.

### `15_Database_Schema_And_ERD.md`
Bổ sung data contract:
- tenants/domains/users/roles/permissions.
- catalog/product/plan/source.
- wallet/ledger/top-up.
- order/item/reservation/service/credential.
- provider/provisioning/resource mapping.
- audit/risk/abuse/notification.
- index, constraint, enum và acceptance criteria.

### `16_API_Contract_And_Permission_Spec.md`
Bổ sung API contract:
- conventions.
- auth/tenant/catalog/wallet/order/service/provider/audit/risk APIs.
- role/permission.
- validation.
- error codes.
- idempotency.
- rate limit.
- audit action.

### `17_RBAC_Permission_Matrix.md`
Bổ sung quyền chi tiết:
- Platform Super Admin.
- Platform Staff.
- Finance Agent.
- Support Agent.
- Provisioning Operator.
- Reseller Owner.
- Reseller Staff.
- Client.
- Read-only Auditor.
- permission matrix và critical controls.

### `18_Provider_Adapter_Technical_Spec.md`
Bổ sung adapter contract:
- capability profile.
- operation result.
- provision/status/suspend/terminate/renew/reset/reinstall/change IP.
- retry safety.
- error normalization.
- idempotency.
- credential handling.
- manual provider.
- provider onboarding checklist.

### `19_Worker_Queue_And_Cron_Jobs_Spec.md`
Bổ sung job nền:
- provisioning_worker.
- provider_sync_worker.
- service_action_worker.
- notification_worker.
- reservation_expiry_job.
- service_expiry/suspension/termination jobs.
- provider health/inventory sync.
- manual review queue.
- monitoring metrics.

### `20_UI_Wireflow_And_Screen_Spec.md`
Bổ sung screen spec:
- Client Portal.
- Reseller Portal.
- Admin Portal.
- service detail/credential reveal.
- checkout flow.
- wallet/top-up.
- catalog/pricing.
- admin provisioning/manual review.
- common status/timeline/confirmation/error states.

### `21_QA_Test_Cases_And_Acceptance_Plan.md`
Bổ sung test plan:
- tenant isolation.
- RBAC.
- wallet/ledger.
- checkout/reservation.
- provisioning.
- service lifecycle.
- credential security.
- catalog/pricing.
- abuse/risk.
- notification/report.
- deployment smoke tests.

### `22_Deployment_DevOps_And_Environment_Runbook.md`
Bổ sung runbook:
- local/staging/production.
- secret management.
- DB migration.
- queue/worker.
- logging/monitoring/alert.
- backup/restore.
- release/rollback.
- incident response.
- go-live checklist.

### `23_Notification_Email_Telegram_Template_Spec.md`
Bổ sung notification:
- event naming.
- channels.
- priority/dedupe.
- client/reseller/admin templates.
- payload security.
- timing matrix.
- audit/notification record.

## File được cập nhật
- `00_README.md`: cập nhật mô tả v1.2, danh sách tài liệu và cách đọc.
- `MANIFEST.txt`: cập nhật danh sách file.
- `VPS_Proxy_Project_Master_Document.md`: tổng hợp lại toàn bộ tài liệu `00–23`.

## Điểm khóa thêm ở bản v1.2
1. Dev có schema logic đủ để thiết kế DB.
2. Backend có API behavior và error code rõ.
3. Frontend có wireflow/screen/action/state rõ.
4. QA có test case để nghiệm thu P0.
5. DevOps có checklist production, backup, monitoring, rollback.
6. Provider adapter có retry/idempotency/partial success behavior rõ.
7. RBAC không còn nằm mơ hồ trong tenant doc.
8. Notification được xem là hệ thống vận hành, không phải phần trang trí.

## Khuyến nghị bước tiếp theo
Sau v1.2, nếu chuẩn bị thuê dev hoặc triển khai sprint, nên làm tiếp:
- Backlog/Sprint plan theo milestone.
- User stories + acceptance criteria theo từng sprint.
- Data migration checklist khi chọn framework/database thật.
- Provider-specific adapter spec cho provider đầu tiên, ví dụ Proxmox hoặc một proxy upstream cụ thể.

---


---

# v1.3 Execution, Operations & Launch Layer

**Date:** 2026-04-22

The following docs extend the project from build specification into team execution and operating readiness:

```text
24 — Project Roadmap, Milestones & Sprint Plan
25 — Backlog, Epics, User Stories & Task Breakdown
26 — MVP Scope Lock & Non-Goals
27 — Developer Onboarding Guide
28 — Finance Reconciliation SOP
29 — Customer Support SOP & Macro Templates
30 — Provider Onboarding & Scoring Checklist
31 — Incident Response & Disaster Recovery Playbook
32 — Abuse, Compliance & Takedown SOP
33 — Launch Checklist & Go/No-Go Criteria
34 — Beta Pilot Program & Feedback Loop
35 — Reseller Acquisition & Enablement Playbook
36 — Go-To-Market Positioning & Offer Strategy
```

The practical purpose of this layer is to keep scope controlled, protect money and tenant boundaries, prepare support/finance/ops, and launch through pilot rather than guesswork.


---

# v1.4 Architecture Deep Dive Layer

**Date:** 2026-04-22

The following docs extend the project from handoff and execution planning into implementation-grade architecture alignment:

```text
37 — Go Backend Architecture & Module Boundaries
38 — PostgreSQL Data Consistency & Transaction Design
39 — Async Worker, Outbox & Job Architecture
40 — Provider Adapter Runtime & Error Taxonomy
41 — Tenant Isolation, RBAC & Security Architecture
42 — Secrets, Credential Encryption & Audit Architecture
43 — Observability, Logging, Metrics & Tracing Spec
44 — Scaling, Performance & Failure Mode Architecture
45 — Architecture Decision Records (ADR)
```

The practical purpose of this layer is to lock the backend implementation baseline before code starts: Go modular monolith, PostgreSQL as source of truth, outbox/job async processing, provider runtime safety, tenant/RBAC enforcement, encrypted credentials, observability, scaling/failure behavior, and documented ADR decisions.
