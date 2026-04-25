# Git Workflow, Build, Test, PR, and Merge Guide

**Version:** v1.7
**Date:** 2026-04-22  
**Scope:** Branching, local development, build, test, commit, pull request, review, merge, release, and hotfix workflow.

## Mục tiêu

Tài liệu này định nghĩa quy trình Git chuẩn để team làm việc đều tay, dễ review, dễ rollback và ít lỗi khi nhiều người cùng sửa code.

Quy trình khuyến nghị: `main` được bảo vệ, mọi thay đổi đi qua branch ngắn hạn và pull request. Không dùng branch `develop` dài hạn trong giai đoạn đầu, trừ khi team cần quản lý nhiều release song song.

## Luật bắt buộc

1. Không commit trực tiếp vào `main`.
2. Mọi thay đổi code phải đi qua pull request.
3. Một branch chỉ xử lý một mục tiêu rõ.
4. Create task branches from latest `origin/main`; do not create them from another agent feature/task branch.
5. Một PR nên nhỏ, dễ review và có phạm vi rõ.
6. Build và test liên quan phải pass trước khi merge.
7. Không merge khi còn conflict, CI fail hoặc review chưa xử lý.
8. Không commit secret, password, private key, token, file `.env` thật hoặc dữ liệu khách hàng.
9. Không sửa lịch sử Git của branch người khác nếu chưa thống nhất.
10. Không dùng commit message mơ hồ như `fix`, `update`, `wip`, `change`, `done`.
11. Không bypass quy trình với flow tiền, tenant, quyền, credential, provider, provisioning hoặc audit.

## Workflow tổng thể

```text
origin/main mới nhất
  -> tạo branch ngắn hạn
  -> code và test local
  -> commit rõ nghĩa
  -> push branch
  -> mở PR
  -> review và CI
  -> merge vào main
  -> xóa branch
```

## Branch chính

`main` là branch ổn định nhất.

Rule cho `main`:

- Luôn build được.
- Luôn chạy được test nền.
- Chỉ nhận code qua PR.
- Bật branch protection khi đưa lên GitHub/GitLab.
- Cần ít nhất 1 review approve trước merge.
- Không force push.
- Không rewrite history.

## Loại branch

Dùng format:

```text
feat/<scope>-<short-name>
fix/<scope>-<short-name>
docs/<scope>-<short-name>
test/<scope>-<short-name>
chore/<scope>-<short-name>
refactor/<scope>-<short-name>
hotfix/<scope>-<short-name>
```

Scope nên là module hoặc vùng thay đổi: `wallet`, `ledger`, `order`, `provider`, `tenant`, `rbac`, `docs`, `infra`.

Tên branch viết thường, dùng dấu `-`, không dùng dấu cách, không dùng tên người.

## Tạo branch

Always start from latest `origin/main`, not from the current branch:

```bash
git fetch origin --prune
git switch -c feat/wallet-topup-api origin/main
```

For parallel agents, prefer an isolated worktree so the branch cannot inherit another task branch by accident:

```bash
git fetch origin --prune
git worktree add -b feat/wallet-topup-api /tmp/Billing-T011 origin/main
```

Do not create a new branch while currently on another agent feature/task branch. Stacked task branches make PRs include unrelated commits/files, create conflicts, and can merge an unreviewed task by accident.

If the branch base is wrong:

1. Stop coding on that branch.
2. Create a clean branch from `origin/main`.
3. Cherry-pick only commits that belong to your task.
4. Do not merge or rebase the old task branch into the clean branch.

Example:

```bash
git fetch origin --prune
git switch -c feat/wallet-topup-api-clean origin/main
git cherry-pick <your-task-commit>
```

## Dev local

Trước khi code:

- Đọc tài liệu liên quan trong `docs/`.
- Xác định module owner.
- Xác định file sẽ sửa.
- Xác định test cần thêm hoặc cần chạy.
- Kiểm tra rule 500 dòng mỗi file.

Trong khi code:

- Commit nhỏ theo từng ý nghĩa.
- Không để branch lệch `main` quá lâu.
- Không copy logic dùng chung nếu đã có owner.
- Không bỏ qua lỗi build/test vì "sẽ sửa sau".
- TODO phải có lý do rõ, không viết TODO mơ hồ.

Trước khi mở PR:

```bash
git status
git diff
git diff --stat
```

Nếu diff khó đọc, tách commit hoặc tách PR.

## Build

Use `docs/05_development_standards/63_Validation_Command_Matrix.md` as the source of truth for the exact validation set by change type.

Khi có Go code, build tối thiểu:

```bash
make test
go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker
```

Nếu có `Makefile`, dùng lệnh chuẩn:

```bash
make fmt
make test
make build
```

Build rule:

- Build fail thì không mở PR ready.
- CI fail thì không merge.
- Nếu thay đổi entrypoint trong `cmd/*`, phải build entrypoint đó.
- Nếu thay đổi shared package trong `internal/platform`, phải chạy test toàn repo.

## Format và lint

Go code phải chạy:

```bash
gofmt -w <files>
make test
```

Khi có lint tool:

```bash
golangci-lint run
```

Rule:

- Không tranh luận style bằng tay nếu formatter xử lý được.
- Không disable lint toàn repo vì một lỗi nhỏ.
- Nếu bỏ qua lint ở một dòng, phải ghi rõ lý do.
- Không để import thừa, code chết hoặc biến không dùng.

## Test

Use `docs/05_development_standards/63_Validation_Command_Matrix.md` for the command matrix by change type.

Mức test tối thiểu:

```text
docs only              không bắt buộc test code
handler/API            test route chính hoặc handler
service rule           test case chính và case lỗi
store/SQL              test query hoặc integration test
ledger/wallet          test debit, credit, reversal, idempotency
tenant/RBAC            test đúng tenant, sai tenant, thiếu quyền, đủ quyền
provider/provisioning  test success, fail, timeout, partial, retry/manual review
platform shared        chạy test toàn repo
migration              test trên database dev sạch nếu có thể
```

Test cho tiền, tenant, quyền, credential, provider và provisioning là bắt buộc trước khi merge tính năng thật.

Không dùng test chỉ để tăng coverage nếu không kiểm tra hành vi thật.

## Commit

Dùng format:

```text
<type>(<scope>): <short summary>
```

Type hợp lệ:

```text
feat      thêm tính năng
fix       sửa lỗi
docs      sửa tài liệu
test      thêm/sửa test
refactor  đổi cấu trúc nhưng không đổi hành vi
chore     việc phụ trợ như config, tooling
build     build system hoặc dependency
ci        CI workflow
perf      cải thiện hiệu năng
revert    revert commit
```

Ví dụ:

```text
feat(wallet): add topup service
fix(provider): classify timeout as retryable
docs(workflow): add git PR merge rules
test(ledger): cover reversal cases
chore(repo): add go module
```

Rule commit:

- Summary nên ngắn và rõ.
- Một commit nên có một ý nghĩa chính.
- Không commit file sinh ra không cần thiết.
- Không commit secret.
- Nếu cần giải thích thêm, viết body cho commit.

## Push

Push branch lần đầu:

```bash
git push -u origin feat/wallet-topup-api
```

Push các lần sau:

```bash
git push
```

Rule:

- Không push thẳng `main`.
- Không force push branch đang được review nếu chưa báo.
- Nếu cần sửa lịch sử branch cá nhân, dùng `git push --force-with-lease`.
- Không dùng `git push --force` thường.

## Pull request

PR phải có:

- Mục tiêu thay đổi.
- Module/file chính đã sửa.
- Cách test đã chạy.
- Rủi ro còn lại.
- Screenshot hoặc log nếu thay đổi UI/API behavior.
- Link issue/task nếu có.

PR title dùng style commit:

```text
feat(wallet): add topup API
fix(provider): retry timeout safely
docs(workflow): add git rules
```

PR nên nhỏ:

- Tốt: dưới 400 dòng diff nếu có thể.
- Cần cân nhắc tách: trên 800 dòng diff.
- Tránh PR vừa đổi kiến trúc, vừa đổi format, vừa thêm feature.

## PR checklist

Người mở PR kiểm tra:

- Branch bắt đầu từ `main` mới nhất chưa?
- Diff có đúng phạm vi không?
- Có file nào vượt 500 dòng không?
- Có logic copy đáng ra phải shared không?
- Có secret hoặc dữ liệu nhạy cảm không?
- Có migration cần rollback không?
- Có test cho rule quan trọng không?
- Build/test local đã chạy chưa?
- README hoặc docs cần cập nhật không?

Reviewer kiểm tra:

- Hành vi có đúng tài liệu không?
- Tên module, file, function có rõ không?
- Boundary giữa handler/service/store có đúng không?
- Transaction và ledger có an toàn không?
- Tenant/RBAC có bị bypass không?
- Provider/provisioning retry có an toàn không?
- Error message có dễ hiểu và không lộ secret không?
- Test có bắt được lỗi thật không?

## Review

Comment review phải nêu rõ:

- Vấn đề nằm ở đâu.
- Rủi ro là gì.
- Gợi ý sửa nếu có.

Không comment mơ hồ như `wrong`, `bad`, `clean this`, `why?`.

Tác giả PR phải phản hồi mọi comment actionable bằng cách sửa code, giải thích lý do không sửa, hoặc tạo follow-up task nếu đã thống nhất.

Không resolve comment khi chưa xử lý xong hoặc chưa có thống nhất.

## Cập nhật branch với main

Ưu tiên rebase cho branch cá nhân:

```bash
git fetch origin
git rebase origin/main
git push --force-with-lease
```

Nếu team chưa quen rebase, có thể dùng merge:

```bash
git fetch origin
git merge origin/main
```

Rule:

- Không rebase branch nhiều người cùng dùng nếu chưa thống nhất.
- Không rebase `main`.
- Sau rebase, chỉ dùng `--force-with-lease`.

## Merge

Chỉ merge khi:

- Có approve bắt buộc.
- CI pass.
- Không còn conflict.
- Comment blocking đã xử lý.
- Branch cập nhật với `main`.
- PR scope đúng và không lẫn thay đổi ngoài ý định.
- PR does not contain commits or files from another task branch.

Chiến lược khuyến nghị:

```text
Squash merge  mặc định cho feature/fix nhỏ và vừa
Merge commit  dùng khi cần giữ nhiều commit có ý nghĩa riêng
Rebase merge  chỉ dùng khi team thống nhất
```

Không merge bằng cách kéo code về local rồi push thẳng lên `main`.

## Sau khi merge

Dọn branch:

```bash
git switch main
git pull --ff-only origin main
git branch -d feat/wallet-topup-api
git push origin --delete feat/wallet-topup-api
```

Sau merge phải theo dõi CI/deploy. Nếu có lỗi nghiêm trọng, revert nhanh thay vì sửa nóng trên `main`.

## Revert và hotfix

Revert thay đổi đã merge:

```bash
git revert <commit-sha>
```

Rule revert:

- Không dùng `git reset --hard` trên branch shared.
- Không rewrite history của `main`.
- Revert cũng đi qua PR nếu không phải sự cố khẩn cấp.
- Mô tả rõ lý do revert.

Hotfix dùng khi production có lỗi nghiêm trọng:

```bash
git switch main
git pull --ff-only origin main
git switch -c hotfix/<scope>-<short-name>
```

Hotfix chỉ sửa phạm vi nhỏ nhất, có test đúng lỗi, có PR nhanh và không kèm refactor hoặc feature mới.

## Release và tag

Khi bắt đầu release chính thức:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

Không tag nếu build/test release chưa pass.

## CI khuyến nghị

Pipeline tối thiểu:

```text
format check
lint
unit test
build all cmd entrypoints
migration check
secret scan
```

Pipeline mở rộng khi hệ thống lớn hơn: integration test với PostgreSQL, race test cho wallet/ledger, API contract test, Docker image build và dependency vulnerability scan.

CI fail thì không merge. Nếu test flaky, sửa test hoặc cách chạy test, không ignore lâu dài.

## Secret và file cấm commit

Không commit:

```text
.env
.env.*
*.pem
*.key
*.p12
*.pfx
id_rsa
id_ed25519
database dump
production config
customer export
provider credential
```

Chỉ commit file mẫu như `.env.example` hoặc `config.example.yaml`.

Nếu lỡ commit secret:

1. Rotate secret ngay.
2. Xóa secret khỏi branch.
3. Báo cho owner.
4. Nếu đã push, xử lý lịch sử Git theo quy trình bảo mật riêng.

Không chỉ xóa ở commit sau rồi coi là xong.

## Docs và migration

Docs-only PR cần:

- Đúng manifest nếu thêm file mới.
- Link từ README hoặc `docs/00_README.md` nếu là tài liệu quan trọng.
- Không làm hỏng encoding tiếng Việt.
- Không tạo file quá 500 dòng nếu có thể tách.

Migration PR phải ghi rõ:

- Bảng/cột/index thay đổi.
- Có ảnh hưởng dữ liệu cũ không.
- Có cần backfill không.
- Có lock bảng lâu không.
- Cách rollback.

Không sửa migration đã merge nếu đã chạy ở môi trường shared. Tạo migration mới để sửa.

## Definition of done cho PR

Một PR được coi là xong khi:

- Scope rõ và không lẫn việc ngoài ý định.
- Build/test cần thiết đã pass.
- Review đã xử lý.
- Docs đã cập nhật nếu hành vi đổi.
- Không có secret.
- Không có file vượt 500 dòng.
- Không có tên mơ hồ theo coding standard.
- Có kế hoạch rollback nếu thay đổi rủi ro cao.
- Merge xong branch được xóa.

## Lệnh nhanh

```bash
git switch main
git pull --ff-only origin main
git switch -c feat/<scope>-<short-name>
git status
git diff
make test
git add <files>
git commit -m "feat(scope): short summary"
git push -u origin feat/<scope>-<short-name>
```
