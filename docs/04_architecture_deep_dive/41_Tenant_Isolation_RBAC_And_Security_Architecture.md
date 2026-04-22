# 41 - Tenant Isolation RBAC And Security Architecture

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Tenant context, domain mapping, identity boundary, RBAC, emergency access, tenant-safe repository rules, security controls  
Related docs: 02, 10, 12, 16, 17, 21, 22, 37, 38, 42, 43

---

## 1. Mục tiêu tài liệu

Tài liệu này khóa kiến trúc bảo vệ tenant và quyền truy cập.

Dự án có nhiều reseller, nhiều client, nhiều ví, nhiều service và credential trên cùng codebase. Lỗi tenant isolation là P0 vì có thể gây:

```text
client xem service của client khác
reseller thấy khách của reseller khác
staff duyệt tiền ngoài phạm vi
credential lộ sai tenant
platform admin thao tác cross-tenant không audit
```

Kết luận:

```text
Tenant context phải được resolve server-side.
Client/reseller API không tin tenant_id từ body/query.
Mọi tenant-owned query phải scope tenant.
RBAC check bắt buộc ở backend.
Emergency access phải có reason, 2FA, audit.
```

---

## 2. Tenant context sources

Tenant context được xác định từ:

```text
domain mapping
authenticated session/token
server-side user membership
platform admin explicit target context
emergency access grant
```

Không tin:

```text
tenant_id trong request body
seller_id trong request body
reseller_id trong request body
owner_id trong request body
client-sent role/permission
```

Request context chuẩn:

```text
request_id
correlation_id
domain
domain_tenant_id
actor_id
actor_tenant_id
actor_type
role_ids
permissions
effective_tenant_id
is_platform_admin
is_emergency_access
emergency_reason
ip_address
user_agent
```

---

## 3. Domain-to-tenant resolution

Storefront/client request:

```text
Host header -> normalized domain -> tenant_domains -> tenant_id
```

Rules:

```text
domain must be verified and active
tenant must be active
primary domain preference is display concern, not auth bypass
custom domain cannot map to multiple tenants
```

If domain not found:

```text
return TENANT_NOT_FOUND
do not fall back to random/default tenant
```

If tenant disabled:

```text
return TENANT_INACTIVE or storefront disabled page
do not allow login/checkout
```

Admin portal may use platform domain and explicit admin routes. Admin target tenant must come from route/query only on admin endpoints and must be audited for high-risk access.

---

## 4. Identity model

Phase 1:

```text
user belongs to one tenant
email unique per tenant
same email may exist in multiple tenants
roles assigned per tenant
```

Login:

```text
domain resolves tenant
email/password checked inside tenant
session stores actor_id and actor_tenant_id
permissions loaded server-side
```

Admin login:

```text
platform admin domain
platform staff tenant
2FA required for privileged roles
```

Do not allow:

```text
login on reseller domain and silently authenticate platform admin
client choosing tenant from dropdown
session switching tenant without explicit admin flow
```

---

## 5. Effective tenant rule

For client/reseller/staff:

```text
effective_tenant_id = actor_tenant_id = domain_tenant_id
```

If mismatch:

```text
deny request
force re-auth if session domain changed
audit suspicious event if relevant
```

For platform admin:

```text
effective_tenant_id may be target tenant for admin route
target_tenant_id must be explicit
critical action needs permission + reason + 2FA policy
audit includes actor_tenant_id and target_tenant_id
```

For background job:

```text
effective_tenant_id = job.tenant_id
actor = system actor
job must load only resources under job.tenant_id
```

---

## 6. RBAC architecture

RBAC entities:

```text
roles
permissions
role_permissions
user_roles
permission risk level
```

Permission naming:

```text
module.resource.action
```

Examples:

```text
wallet.topup.approve
service.credential.reveal
provisioning.manual_review.resolve
provider.manage
tenant.domain.manage
```

Backend check pattern:

```text
1. Resolve actor context.
2. Load resource metadata with tenant scope.
3. Check permission.
4. Check resource policy/status.
5. Check risk controls: 2FA/reason/rate limit.
6. Execute mutation inside transaction.
7. Audit.
```

RBAC alone is not enough. Resource ownership and tenant scope must also be checked.

---

## 7. Policy checks by risk level

| Risk | Examples | Required controls |
|---|---|---|
| Low | view public catalog, view own order | auth where needed + tenant scope |
| Medium | checkout, submit top-up, update profile | permission + validation + rate limit if needed |
| High | approve top-up, suspend service, reveal credential | permission + audit + rate limit + reason for staff action |
| Critical | ledger adjustment, provider manage, terminate, emergency access | restricted permission + 2FA + reason + audit + optional approval |

Critical actions should include:

```text
actor_id
target_tenant_id
target_resource_id
reason
ip/user_agent
correlation_id
before/after redacted
```

---

## 8. Tenant-safe repository rules

Repository method names should make tenant scope obvious.

Good:

```text
GetServiceForTenant(ctx, tenantID, serviceID)
ListOrdersForClient(ctx, tenantID, clientID, filter)
GetWalletForOwner(ctx, tenantID, ownerType, ownerID)
CreateLedgerEntriesTx(ctx, tx, tenantID, entries)
```

Bad:

```text
GetService(id)
GetWallet(walletID)
FindOrder(orderID)
UpdateBalance(walletID, amount)
```

Rule:

```text
All tenant-owned table queries include tenant_id.
Admin cross-tenant reads use explicit admin repository methods.
No handler calls repository directly for business mutation.
```

---

## 9. Database-level tenant integrity

Application guard is mandatory, but DB should help catch mistakes.

Patterns:

```sql
ALTER TABLE services
ADD CONSTRAINT services_id_tenant_unique UNIQUE (service_id, tenant_id);

ALTER TABLE service_credentials
ADD CONSTRAINT service_credentials_service_tenant_fk
FOREIGN KEY (service_id, tenant_id)
REFERENCES services (service_id, tenant_id);
```

Use composite FK for high-risk child rows:

```text
wallet_ledger_entries -> wallets
orders -> users/tenant
order_items -> orders
reservations -> orders/order_items
services -> orders/reservations/users
service_credentials -> services
audit_logs -> tenant nullable by policy
```

Indexes for large tenant tables should start with:

```text
tenant_id
```

PostgreSQL RLS can be considered later, but MVP still needs app-level tenant guard because:

```text
admin routes need explicit cross-tenant policy
worker jobs need system context
tests must verify repository behavior
```

---

## 10. API route classes

### 10.1 Public tenant routes

```text
/auth/register
/auth/login
/catalog
```

Tenant from domain. No tenant_id body.

### 10.2 Client routes

```text
/wallet
/orders
/services
/services/{id}/credential/reveal
```

Scope:

```text
tenant_id = actor tenant
client_user_id = actor id
```

### 10.3 Reseller routes

```text
/reseller/clients
/reseller/catalog
/reseller/topups
/reseller/reports
```

Scope:

```text
tenant_id = reseller tenant
resource belongs to tenant
```

### 10.4 Platform admin routes

```text
/admin/tenants
/admin/provider-sources
/admin/reseller-topups
/admin/provisioning/jobs
```

Admin may target tenant, but target must be explicit and audited for sensitive access.

---

## 11. Emergency access

Emergency access is for rare support/security incidents.

Required:

```text
platform permission
2FA fresh check
reason
target_tenant_id
time-limited access
audit start/end/action
```

Not allowed:

```text
silent impersonation
credential reveal without separate audit
changing ledger without adjustment workflow
using emergency access for normal ops convenience
```

Emergency access session should display internally:

```text
who is accessing
which tenant
reason
expiry
correlation_id
```

---

## 12. Session and 2FA policy

Phase 1 minimum:

```text
platform admin 2FA required
reseller owner 2FA available and recommended
fresh 2FA for critical actions if role requires
session expiration configured
login rate limit
password reset rate limit
logout invalidates session
```

Session should store minimal data:

```text
session_id
actor_id
actor_tenant_id
issued_at
expires_at
2fa_verified_at
```

Permissions can be cached short-term but must be invalidated on role changes.

---

## 13. Rate limit security controls

Rate limit by:

```text
IP
tenant_id
actor_id
action
resource_id where relevant
```

Actions requiring rate limit:

```text
login
password reset
top-up submit
checkout
credential reveal
provider action request
admin critical action
```

Credential reveal should have stricter rate limits than normal service view.

---

## 14. Audit integration

Audit required for:

```text
login failure/security event
role/permission change
tenant/domain change
top-up approve/reject
ledger adjustment
checkout/refund
service suspend/terminate
credential reveal
provider config change
emergency access
cross-tenant admin access
```

Audit must include:

```text
actor_id
actor_tenant_id
target_tenant_id
action
resource_type
resource_id
reason if required
correlation_id
ip/user_agent
redacted before/after
```

Never include plaintext credential or provider secret.

---

## 15. Testing requirements

P0 tests:

```text
Client cannot access another client's service by guessed ID.
Reseller A cannot access Reseller B client/order/service.
tenant_id injection in body is ignored.
Admin cross-tenant access requires admin permission and audit.
Staff without permission cannot approve top-up.
Finance cannot reveal credential unless explicitly granted.
Credential reveal logs audit and remains tenant-scoped.
Worker job cannot process resource outside job tenant.
```

Repository tests should detect:

```text
query by id without tenant scope
composite FK violation for cross-tenant child row
role change invalidates permission cache
```

---

## 16. Acceptance criteria

Tenant/RBAC architecture đạt khi:

```text
Tenant context is resolved server-side for every request.
Client/reseller APIs never trust tenant_id from body/query.
Every tenant-owned repository query includes tenant_id or ActorContext.
RBAC middleware/policy protects all write and high-risk read actions.
Emergency access has 2FA, reason, expiry, and audit.
Credential reveal requires ownership/permission and audit.
Cross-tenant tests are part of CI before launch.
Audit can reconstruct actor, tenant, resource, action, and reason.
```

---

## 17. Tóm tắt quyết định

```text
Tenant is the data border.
RBAC is necessary but not sufficient without ownership checks.
Platform admin power must be explicit and audited.
Repository APIs must make tenant scope hard to forget.
Emergency access is controlled access, not impersonation.
Frontend permission hiding is UX only; backend is security.
```
