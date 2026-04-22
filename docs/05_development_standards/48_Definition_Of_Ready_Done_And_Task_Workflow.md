# Definition of Ready, Definition of Done, and Task Workflow

**Version:** v1.7  
**Date:** 2026-04-22  
**Scope:** Task intake, readiness, development, review, completion, and follow-up rules.

## Mục tiêu

Tài liệu này khóa cách biến yêu cầu thành task có thể dev, review và nghiệm thu. Mục tiêu là không bắt đầu code khi yêu cầu còn mơ hồ, và không gọi là xong khi mới chỉ chạy được trên máy dev.

## Luật nền

1. Task phải có mục tiêu rõ trước khi dev.
2. Task phải có module owner.
3. Task phải có acceptance criteria.
4. Task chạm tiền, tenant, quyền, credential, provider, provisioning hoặc audit phải được đánh dấu high-risk.
5. Task high-risk phải có test và kế hoạch rollback.
6. Không trộn feature, refactor lớn và format cleanup trong cùng một task.
7. Không đóng task nếu docs, migration, test hoặc config liên quan chưa cập nhật.

## Trạng thái task

```text
Backlog
  -> Ready
  -> In Progress
  -> Review
  -> Changes Requested
  -> Approved
  -> Merged
  -> Verified
  -> Done
```

Ý nghĩa:

- `Backlog`: ý tưởng hoặc yêu cầu chưa đủ rõ.
- `Ready`: đủ điều kiện để dev bắt đầu.
- `In Progress`: đang code.
- `Review`: đã mở PR và chờ review/CI.
- `Changes Requested`: cần sửa theo review hoặc CI.
- `Approved`: đã được approve, chờ merge.
- `Merged`: đã merge vào `main`.
- `Verified`: đã kiểm tra sau merge hoặc trên môi trường phù hợp.
- `Done`: đã xong đủ tiêu chí.

## Definition of Ready

Một task chỉ được chuyển sang `Ready` khi có:

- Mục tiêu ngắn gọn.
- Phạm vi làm và không làm.
- Module owner hoặc khu vực code chính.
- Acceptance criteria có thể kiểm tra.
- API/data/UI impact nếu có.
- Test cần viết hoặc cần chạy.
- Rủi ro nếu đụng flow quan trọng.
- Link tài liệu liên quan trong `docs/`.

Task chưa ready nếu:

- Chưa rõ user/role nào dùng.
- Chưa rõ tenant hoặc quyền truy cập.
- Chưa rõ tiền bị debit/credit/lock ở đâu.
- Chưa rõ provider có thể timeout, fail hoặc partial success thế nào.
- Chưa rõ dữ liệu mới nằm ở bảng nào.
- Chưa rõ response API hoặc error code.

## Acceptance criteria

Acceptance criteria nên viết theo dạng có thể test:

```text
Given <điều kiện ban đầu>
When <hành động xảy ra>
Then <kết quả phải đúng>
```

Ví dụ:

```text
Given reseller wallet không đủ số dư
When client checkout service
Then hệ thống không tạo provisioning job và trả lỗi wallet balance is not enough
```

Acceptance criteria không nên viết mơ hồ:

```text
làm checkout tốt hơn
xử lý lỗi provider
tối ưu API
```

## Definition of Done chung

Một task được coi là done khi:

- Code đã merge vào `main`.
- CI/build/test cần thiết đã pass.
- Acceptance criteria đã được verify.
- PR review đã xử lý.
- Docs đã cập nhật nếu hành vi đổi.
- Không có file vượt 500 dòng.
- Không có secret hoặc dữ liệu nhạy cảm trong repo.
- Không còn TODO mơ hồ.
- Có rollback note nếu thay đổi rủi ro cao.

## Done cho backend feature

Backend feature done khi có:

- Handler chỉ nhận request và map response.
- Service chứa business rule.
- Store chứa SQL nếu có database.
- Error code rõ và ổn định.
- Log có request id, tenant id nếu có, nhưng không có secret.
- Unit test cho service rule chính.
- Integration test nếu có database/migration quan trọng.
- API contract hoặc docs cập nhật nếu endpoint đổi.

## Done cho flow tiền

Flow tiền gồm wallet, ledger, order payment, refund, settlement và reversal.

Done khi có thêm:

- Ledger entry được tạo đúng lúc.
- Không sửa transaction cũ, chỉ tạo adjustment/reversal.
- Idempotency được test.
- Retry không debit hai lần.
- Số dư ví không âm nếu policy không cho phép.
- Race condition quan trọng đã được test hoặc có lock rõ.
- Audit event được ghi.

## Done cho tenant và RBAC

Done khi có thêm:

- Test đúng tenant.
- Test sai tenant.
- Test thiếu quyền.
- Test đủ quyền.
- Không query dữ liệu tenant nếu thiếu tenant context.
- UI/API không lộ action khi role không được phép.

## Done cho provider và provisioning

Done khi có thêm:

- Test success.
- Test fail.
- Test timeout.
- Test partial success nếu provider có khả năng đó.
- Retry policy rõ.
- Manual review path rõ.
- Capability snapshot được dùng khi quyết định action.
- Không log credential hoặc raw response có secret.

## Done cho migration

Migration done khi có:

- Tên migration rõ mục đích.
- Up migration chạy được trên database sạch.
- Down migration hoặc rollback plan rõ.
- Không sửa migration đã merge nếu đã chạy ở shared environment.
- Backfill hoặc data migration có batch plan nếu dữ liệu lớn.
- Index mới được đánh giá lock/risk.

## Done cho docs-only task

Docs-only task done khi:

- Link trong `docs/MANIFEST.txt` đúng.
- `docs/00_README.md` cập nhật nếu tài liệu quan trọng.
- README top-level cập nhật nếu tài liệu là entrypoint.
- Encoding tiếng Việt không lỗi.
- File mới không vượt 500 dòng.

## Chia nhỏ task

Nên tách task khi:

- PR dự kiến trên 800 dòng diff.
- Đụng nhiều module không liên quan trực tiếp.
- Vừa đổi schema vừa đổi nhiều API.
- Vừa refactor vừa thêm feature.
- Một phần có thể merge độc lập và giảm rủi ro.

Ví dụ tách tốt:

```text
1. Add wallet schema and migration
2. Add wallet service debit/credit
3. Add wallet API handlers
4. Add checkout integration with wallet
```

## Blocker

Khi bị blocker, ghi rõ:

- Đang bị chặn bởi gì.
- Cần ai quyết định.
- Có lựa chọn nào.
- Rủi ro của từng lựa chọn.
- Task nào vẫn có thể làm song song.

Không để blocker dạng:

```text
blocked
need check
waiting
```

## Review handoff

Khi mở PR, tác giả phải ghi:

- Reviewer nên đọc file nào trước.
- Flow nào quan trọng nhất.
- Test đã chạy.
- Rủi ro còn lại.
- Chỗ nào cần reviewer chú ý.

Reviewer không nên phải tự đoán intent của PR.

## No-go rule

Không bắt đầu dev nếu task liên quan các vùng sau mà chưa có acceptance criteria:

- Wallet debit/credit/lock.
- Ledger entry.
- Refund/reversal.
- Tenant isolation.
- RBAC.
- Credential encryption.
- Provider provisioning.
- Manual review.
- Audit/compliance.

## Checklist nhanh

Trước dev:

- Task ready chưa?
- Module owner rõ chưa?
- Acceptance criteria test được chưa?
- Có rủi ro tiền/tenant/quyền/provider không?

Trước PR:

- Build/test đã chạy chưa?
- Docs cần update chưa?
- Diff có dễ review không?
- Có file vượt 500 dòng không?

Trước done:

- PR đã merge chưa?
- Acceptance criteria đã verify chưa?
- Có follow-up nào cần ghi lại không?
