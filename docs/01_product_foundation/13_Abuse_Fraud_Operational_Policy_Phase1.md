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
