# T041 - Frontend UI Phase 2: Interactive Components & Real Workflows

Status: IN_PROGRESS
Owner: Claude
Branch: feat/frontend-ui-phase2
PR: -
Risk: frontend/UI
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Phase 2 của frontend UI. Phase 1 (T009, T010) đã hoàn thành app shell với 3 portals, layout, navigation, và mock data.
Phase 2 bổ sung các thành phần tương tác còn thiếu để UI sẵn sàng kết nối API thật: login screen, modal/dialog, toast notification, form workflows cho các action quan trọng, filter/search trên table, và URL-based routing.

## Background & Analysis

Phân tích UI hiện tại (2026-04-23) cho thấy:

### Có sẵn (Phase 1)
- 3 portals: Admin (15 screens), Reseller (5 screens), Client (4 screens)
- Layout: AppShell, Sidebar (collapsible), Topbar
- UI components: KpiCard, StatusBadge, TablePagination
- Mock data: billingData.ts (đầy đủ data shape)
- Stack: Next.js 15, React 19, TypeScript, Tailwind CSS v4

### Thiếu hoàn toàn
1. **Auth** — Không có login screen, không có session/token handling
2. **Modal/Dialog** — Không có confirm dialog hay form modal nào
3. **Toast/Notification** — Không có feedback sau action
4. **Empty/Loading/Error states** — Table không có skeleton hoặc empty UI
5. **Filter/Search** — Không có filter trên bất kỳ table nào
6. **Action forms** — Các nút tồn tại nhưng không làm gì (Approve/Reject topup, tạo tenant, order, top-up)
7. **URL routing** — Navigation dùng state thay vì URL
8. **Detail pages** — Click vào row không đi đâu

## Scope

### In scope

**Phase 2A — Foundation components**
- `Modal` component dùng chung (controlled, Esc to close, backdrop click)
- `Toast` / notification system (success/error/warning/info, auto-dismiss 4s)
- `ConfirmDialog` — wrapper trên Modal cho confirm destructive actions
- `EmptyState` — component cho table rỗng
- `LoadingSkeleton` — skeleton rows cho table
- `SearchFilter` — search input + filter bar reusable

**Phase 2B — URL-based routing**
- Chuyển navigation từ `useState` sang Next.js file-system router
- Mỗi screen có URL riêng: `/admin/overview`, `/admin/tenants`, `/reseller/dashboard`, v.v.
- Protected route layout (redirect `/login` nếu chưa auth)
- Browser Back/Forward hoạt động

**Phase 2C — Login screen**
- `/login`: form email + password, validation, loading state
- Mock auth: hardcode credentials, lưu session localStorage
- Redirect sau login về portal đúng theo role

**Phase 2D — Admin action workflows**
- Topup approve/reject: ConfirmDialog (reject có reason field)
- Provisioning manual review: ConfirmDialog re-trigger / reject
- Tạo Tenant: Modal form (name, type, domain)
- Tạo Product: Modal form (SKU, name, unit, price)
- AdminSettings Save hoạt động (mock + toast)

**Phase 2E — Reseller action workflows**
- Tạo Client: Modal form (name, email, wallet)
- Request Top-up: Modal form (amount, method, reference)
- ResellerSettings Save hoạt động
- ResellerCatalog: edit selling price inline

**Phase 2F — Client action workflows**
- Top-up form: Modal (amount, method)
- Order flow: "Order now" → modal confirm → mock deduct wallet
- Renew service: confirm modal

### Out of scope
- Kết nối API thật (task riêng sau T041)
- Mobile responsive
- Dark mode
- WebSocket / real-time updates
- Export CSV/Excel

## Acceptance Criteria

- [ ] Modal đóng khi click backdrop hoặc nhấn Esc
- [ ] Toast hiển thị góc dưới phải, auto dismiss 4s, có nút close
- [ ] ConfirmDialog nhận `title`, `description`, `onConfirm`, `onCancel`, `danger` prop
- [ ] Mỗi screen có URL unique, browser Back/Forward hoạt động
- [ ] `/login` redirect về portal khi đã auth; unauthenticated redirect về `/login`
- [ ] Login form validate empty fields, hiển thị error
- [ ] Mỗi action button quan trọng có modal/confirm flow
- [ ] Sau action hiển thị toast success/error
- [ ] Form validate required fields trước submit
- [ ] `npm run build` pass không error trong `frontend/`
- [ ] Không có TypeScript error

## Implementation Order

```
1. Phase 2A — shared components (Modal, Toast, ConfirmDialog, EmptyState, Skeleton)
2. Phase 2B — URL routing
3. Phase 2C — Login screen + mock auth
4. Phase 2D — Admin actions
5. Phase 2E — Reseller actions
6. Phase 2F — Client actions
```

## Notes

- Giữ nguyên mock data layer — không đổi `billingData.ts`, `sampleData.ts`
- Tất cả form actions vẫn dùng mock state, chưa gọi API
- Toast và Modal dùng React Portal
- Đọc `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md` trước khi code

## Agent Log

- 2026-04-23: Task T041 created. Branch feat/frontend-ui-phase2 từ origin/main (HEAD: 4daacb0).
