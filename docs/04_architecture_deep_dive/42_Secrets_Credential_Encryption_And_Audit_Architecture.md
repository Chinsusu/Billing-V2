# 42 - Secrets Credential Encryption And Audit Architecture

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Secret classification, provider secrets, service credential encryption, reveal flow, redaction, audit, rotation, backup safety  
Related docs: 05, 08, 10, 15, 16, 17, 18, 21, 22, 23, 37, 38, 41, 43

---

## 1. Mục tiêu tài liệu

Tài liệu này khóa cách hệ thống lưu, dùng, reveal, audit và rotate secret/credential.

Dữ liệu nhạy cảm trong dự án:

```text
provider API key
provider token/session secret
service root password
proxy username/password
SSH private key
console URL token
2FA secret
session/JWT secret
webhook secret
email/telegram/object storage secret
database password
encryption key
```

Kết luận:

```text
No plaintext secret in database except transient memory.
No plaintext credential in logs/audit/notification/support notes.
Credential reveal is an explicit audited action.
Provider secrets and service credentials use envelope encryption or secret manager references.
Rotation must be planned from MVP schema.
```

---

## 2. Secret classification

| Class | Examples | Storage |
|---|---|---|
| Platform critical secret | DB password, encryption master key, session secret | secret manager/env reference |
| Provider secret | API key/token for OVH/Hetzner/Proxmox/proxy upstream | encrypted at rest or secret manager reference |
| Service credential | VPS root password, proxy auth, SSH key | encrypted payload in DB |
| User auth secret | password hash, 2FA secret, reset token hash | password hash/secret encrypted/token hash |
| Operational secret | SMTP token, Telegram bot token, object storage key | secret manager/env reference |

Rule:

```text
password hash is not plaintext, but still sensitive.
2FA secret must be encrypted.
reset token should be hashed, not stored plaintext.
```

---

## 3. Encryption model

Recommended MVP approach:

```text
application-level envelope encryption
master key stored outside DB
per-record data key or derived key version
encrypted_payload JSONB/BYTEA in DB
secret_version/key_version stored with row
```

Fields:

```text
encrypted_payload
encryption_key_version
encryption_algorithm
masked_hint
created_at
rotated_at
status
```

Do not store:

```text
plaintext_payload
raw_password
raw_api_key
raw_proxy_auth
```

Acceptable encrypted payload examples:

```json
{
  "username": "root",
  "password": "...",
  "ssh_private_key": "...",
  "console_url": "..."
}
```

The JSON is encrypted as a whole or each secret value is encrypted. Metadata safe for display goes into `masked_hint`.

---

## 4. Provider secret storage

Provider account table should store either:

```text
encrypted_credentials
```

or:

```text
secret_reference
```

If using encrypted DB payload:

```text
provider_account_id
provider_type
encrypted_credentials
encryption_key_version
status
last_rotated_at
created_at
updated_at
```

If using external secret manager:

```text
provider_account_id
secret_reference
secret_version
status
```

API response must return:

```text
provider_account_id
provider_type
status
health_status
last_rotated_at
```

Never return:

```text
api_key
api_secret
token
password
full encrypted blob unless internal restricted path
```

---

## 5. Service credential storage

`service_credentials` should contain:

```text
credential_id
tenant_id
service_id
credential_type
encrypted_payload
encryption_key_version
masked_hint
last_revealed_at
last_revealed_by
status
created_at
updated_at
```

Credential types:

```text
vps_root
proxy_auth
ssh_key
console_url
api_token
recovery_code
```

Status:

```text
active
rotated
revoked
expired
```

Service detail response returns:

```text
credential_id
credential_type
masked_hint
status
last_revealed_at
```

Only reveal endpoint returns plaintext temporarily.

---

## 6. Credential reveal flow

Flow:

```text
1. User requests reveal.
2. Resolve tenant and actor.
3. Load service with tenant scope.
4. Verify ownership/permission.
5. Check service status and credential status.
6. Rate limit reveal action.
7. Require fresh 2FA if policy says.
8. Write audit intent or audit after success.
9. Decrypt credential.
10. Return plaintext only in response body.
11. Update last_revealed_at/by.
12. Write audit credential.revealed.
```

Response should include:

```text
credential payload
expires_in_seconds for UI display policy
request_id
correlation_id
```

UI should:

```text
show credential temporarily
provide copy button
hide on navigation/timeout
avoid persistent local storage
```

Do not send reveal credential through:

```text
email
telegram
notification center
support ticket public note
analytics event
```

---

## 7. Audit for secrets

Audit actions:

```text
credential.reveal_requested
credential.revealed
credential.reveal_denied
credential.rotated
credential.revoked
provider_secret.created
provider_secret.updated
provider_secret.rotated
provider_secret.accessed_by_system
encryption_key.rotated
```

Audit must include:

```text
actor_id
actor_tenant_id
target_tenant_id
service_id/provider_account_id
credential_type
reason if staff/admin
correlation_id
ip/user_agent
result
```

Audit must not include:

```text
plaintext credential
provider API key
encrypted payload
full token
password
```

If diff logging exists, redaction must happen before audit persistence.

---

## 8. Redaction policy

Fields to redact by key name:

```text
password
pass
secret
token
api_key
api_secret
credential
private_key
authorization
cookie
set_cookie
otp
totp
recovery_code
```

Replacement:

```text
[REDACTED]
```

Partial masking allowed only for safe hints:

```text
root / ********
proxy user abc*** / ********
token ending ****1234
```

Never rely only on exact key name. Provider raw payload can place secret in arbitrary fields, so provider adapter must explicitly sanitize known payload structure.

---

## 9. Logging rules

Application logs may contain:

```text
credential_id
service_id
provider_account_id
source_id
operation
error_code
correlation_id
```

Application logs must not contain:

```text
plaintext secret
encrypted credential blob
raw provider response with secret
authorization header
cookie/session token
```

In dev:

```text
do not add debug logs that print structs containing credential fields.
tests should scan logs/audit for secret sample values.
```

---

## 10. Key rotation

Key rotation types:

```text
master key rotation
data key rotation
provider API key rotation
service credential rotation/reset
2FA secret reset
session secret rotation
```

MVP schema must support:

```text
encryption_key_version
secret_version
status active/rotated/revoked
last_rotated_at
rotated_by
```

Rotation flow for encrypted DB payload:

```text
1. Load encrypted row.
2. Decrypt with old key version.
3. Encrypt with new key version.
4. Update key_version and audit.
5. Verify decrypt with new key.
```

Rotation job should be resumable and idempotent.

---

## 11. Credential lifecycle

### 11.1 Provisioning success

```text
provider returns credential
adapter sanitizes logs
credential service encrypts payload
service_credentials row created
service activation transaction commits
audit service.activated
```

### 11.2 Reset password

```text
check capability
create service_action job
provider reset returns new credential
old credential status = rotated
new credential status = active
audit credential.rotated
notify user without plaintext secret
```

### 11.3 Termination

On service terminated:

```text
credential may remain encrypted for retention/support policy
status can become revoked/expired
reveal may be blocked by policy
retention cleanup later deletes or tombstones encrypted payload
```

---

## 12. Backup and restore considerations

Database backup contains encrypted credentials.

Restore requires:

```text
database backup
matching encryption key versions
secret manager restore path
access control for restore operators
restore test that verifies decrypt works
```

If backup is restored to staging:

```text
do not restore production secrets unless approved
mask/anonymize tenant/user data
rotate provider credentials or disable provider accounts
disable outgoing notification
```

Losing encryption keys means encrypted credentials cannot be recovered. This is an incident.

---

## 13. Support and operations policy

Support should not ask users to paste root passwords/proxy passwords into tickets.

Support may:

```text
guide user to reveal credential in portal
verify masked_hint
request service_id/order_number
escalate credential reset
```

Support should not:

```text
send credential via chat/email
copy credential into internal note
reveal credential without permission/audit
```

Admin/operator reveal of client credential should require:

```text
explicit permission
reason
audit
possibly emergency access policy
```

---

## 14. Incident handling

Credential exposure incident triggers:

```text
freeze affected reveal/provider operations if needed
identify affected tenants/services
rotate provider secrets if exposed
rotate service credentials where possible
audit log review
notify affected parties per policy
postmortem
add regression test/redaction rule
```

Examples:

```text
plaintext password found in logs
provider API key committed to repo
audit diff stored credential
support ticket contains credential
backup restored without secret controls
```

---

## 15. Acceptance criteria

Secrets architecture đạt khi:

```text
Provider secrets are encrypted or secret-manager referenced.
Service credentials are encrypted at rest.
Service detail returns masked credential only.
Reveal endpoint enforces tenant, ownership/permission, rate limit, and audit.
No plaintext secret appears in logs/audit/provider_requests/notifications.
Encryption key version is stored with encrypted rows.
Credential rotation/reset can create new active credential and retire old one.
Backup/restore procedure includes encryption key dependency.
Tests scan logs/audit for known secret sample values.
```

P0 tests:

```text
unauthorized credential reveal denied.
authorized reveal creates audit.
service list/detail never returns plaintext credential.
provider request log redacts API key and returned password.
rotated credential no longer appears as active.
```

---

## 16. Tóm tắt quyết định

```text
Credentials are data assets, not normal fields.
Reveal is a privileged event, not a page load.
Audit proves access, but never stores the secret itself.
Logs are hostile to secrets by default.
Key versioning is required from MVP.
Support workflows must avoid moving credentials into human text channels.
```
