package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

const (
	targetCredentialRevealServiceID     = "00000000-0000-0000-0000-000000000909"
	targetCredentialRevealType          = "recovery_code"
	targetCredentialRevealMaskedHint    = "Recovery code / ****"
	targetCredentialRevealSecretVersion = "target-credential-reveal-smoke"
	targetCredentialRevealReason        = "target credential reveal smoke"
	targetCredentialRevealPayload       = `{"username":"target-smoke-user","password":"target-credential-smoke-secret","host":"target-smoke.invalid"}`
)

type targetCredentialRevealFixture struct {
	CredentialID     string
	EncryptedPayload string
	ServiceDisplayID int64
	ExpectedPayload  string
}

type targetCredentialRevealResponse struct {
	ID                   string          `json:"id"`
	Type                 string          `json:"credential_type"`
	MaskedHint           string          `json:"masked_hint"`
	Status               string          `json:"status"`
	Payload              json.RawMessage `json:"payload"`
	RevealedAt           time.Time       `json:"revealed_at"`
	RevealExpiresMessage string          `json:"reveal_expires_message"`
}

type targetCredentialRevealDBEvidence struct {
	LastRevealedByClient bool
	RateLimitAttempts    int
	AuditCount           int
	AuditDisplayID       int64
	AuditHasDisplayID    bool
	AuditLeakedSecret    bool
	AuditLeakedCipher    bool
}

func runDevTargetCredentialRevealSmoke(dsn string, baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev-target-credential-reveal smoke")
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-target-credential-reveal smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	encryptionKey := strings.TrimSpace(os.Getenv("ENCRYPTION_KEY"))
	if encryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required for dev-target-credential-reveal smoke")
	}
	cipher, err := secrets.NewAESGCMCipher(encryptionKey)
	if err != nil {
		return fmt.Errorf("build credential fixture cipher: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := platformdb.Open(ctx, platformdb.Config{DriverName: platformdb.DefaultDriverName, DSN: dsn})
	if err != nil {
		return err
	}
	defer conn.Close()

	fixture, err := prepareTargetCredentialRevealFixture(ctx, conn, cipher)
	if err != nil {
		return err
	}
	startedAt := time.Now().UTC().Add(-2 * time.Second)

	client := &http.Client{Timeout: timeout}
	cookieName := targetAuthSessionCookieName()
	credentials := targetAuthSmokeCredentialsFromEnv()
	login, cookie, err := loginForTargetAuthSmoke(ctx, client, baseURL, cookieName, demoTenantID, credentials.ClientEmail, credentials.ClientPassword)
	if err != nil {
		return err
	}
	defer func() { _ = logoutTargetAuthSmoke(context.Background(), client, baseURL, cookie) }()
	if login.ActorType != "client" || login.TenantID != demoTenantID {
		return fmt.Errorf("target credential reveal expected seeded client login")
	}

	reveal, err := revealTargetCredential(ctx, client, baseURL, cookie, fixture)
	if err != nil {
		return err
	}
	if reveal.ID != fixture.CredentialID {
		return fmt.Errorf("target credential reveal response returned unexpected credential id")
	}

	evidence, err := collectTargetCredentialRevealDBEvidence(ctx, conn, fixture, startedAt)
	if err != nil {
		return err
	}
	if err := validateTargetCredentialRevealDBEvidence(evidence); err != nil {
		return err
	}

	fmt.Printf("target credential reveal smoke passed: service_display_id=%d credential_type=%s client_session_cookie_only=pass no_store=pass audit_display_id=%d last_revealed_by=client_actor rate_limit_attempts=%d provider_mutation_routes_called=no money_mutation_routes_called=no\n",
		fixture.ServiceDisplayID,
		targetCredentialRevealType,
		evidence.AuditDisplayID,
		evidence.RateLimitAttempts,
	)
	fmt.Println("Target credential reveal smoke output intentionally excludes plaintext credentials, encrypted payloads, raw credential IDs, session tokens, cookies, DSNs, provider payloads, and provider credentials.")
	return nil
}

func prepareTargetCredentialRevealFixture(ctx context.Context, conn *sql.DB, cipher secrets.Cipher) (targetCredentialRevealFixture, error) {
	serviceDisplayID, err := seededServiceDisplayID(ctx, conn)
	if err != nil {
		return targetCredentialRevealFixture{}, err
	}
	encryptedPayload, err := cipher.Encrypt(targetCredentialRevealPayload)
	if err != nil {
		return targetCredentialRevealFixture{}, fmt.Errorf("encrypt target credential reveal fixture: %w", err)
	}

	var credentialID string
	err = conn.QueryRowContext(ctx, `
INSERT INTO service_credentials (
    tenant_id,
    service_instance_id,
    credential_type,
    encrypted_payload,
    encryption_key_version,
    encryption_algorithm,
    secret_version,
    masked_hint,
    status,
    last_revealed_at,
    last_revealed_by
)
VALUES ($1::uuid, $2::uuid, $3::service_credential_type, $4, 'v1', 'aes-256-gcm', $5, $6, 'active', NULL, NULL)
ON CONFLICT (tenant_id, service_instance_id, credential_type) WHERE status = 'active'
DO UPDATE SET
    encrypted_payload = EXCLUDED.encrypted_payload,
    encryption_key_version = EXCLUDED.encryption_key_version,
    encryption_algorithm = EXCLUDED.encryption_algorithm,
    secret_version = EXCLUDED.secret_version,
    masked_hint = EXCLUDED.masked_hint,
    last_revealed_at = NULL,
    last_revealed_by = NULL,
    updated_at = NOW()
RETURNING credential_id`,
		demoTenantID,
		targetCredentialRevealServiceID,
		targetCredentialRevealType,
		encryptedPayload,
		targetCredentialRevealSecretVersion,
		targetCredentialRevealMaskedHint,
	).Scan(&credentialID)
	if err != nil {
		return targetCredentialRevealFixture{}, fmt.Errorf("prepare target credential reveal fixture: %w", err)
	}

	if _, err := conn.ExecContext(ctx, `
DELETE FROM service_credential_reveal_rate_limits
WHERE tenant_id = $1::uuid
  AND actor_id = $2::uuid
  AND service_instance_id = $3::uuid`,
		demoTenantID,
		demoCustomerID,
		targetCredentialRevealServiceID,
	); err != nil {
		return targetCredentialRevealFixture{}, fmt.Errorf("reset target credential reveal rate limit fixture: %w", err)
	}

	return targetCredentialRevealFixture{
		CredentialID:     credentialID,
		EncryptedPayload: encryptedPayload,
		ServiceDisplayID: serviceDisplayID,
		ExpectedPayload:  targetCredentialRevealPayload,
	}, nil
}

func seededServiceDisplayID(ctx context.Context, conn *sql.DB) (int64, error) {
	var displayID int64
	var buyerUserID string
	if err := conn.QueryRowContext(ctx, `
SELECT service_instances.display_id, orders.buyer_user_id::text
FROM service_instances
JOIN orders ON orders.order_id = service_instances.order_id
  AND orders.tenant_id = service_instances.tenant_id
WHERE service_instances.service_instance_id = $1::uuid
  AND service_instances.tenant_id = $2::uuid
  AND service_instances.status = 'active'`,
		targetCredentialRevealServiceID,
		demoTenantID,
	).Scan(&displayID, &buyerUserID); err != nil {
		return 0, fmt.Errorf("load seeded target credential reveal service: %w", err)
	}
	if buyerUserID != demoCustomerID {
		return 0, fmt.Errorf("seeded target credential reveal service buyer mismatch")
	}
	return displayID, nil
}

func revealTargetCredential(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	cookie *http.Cookie,
	fixture targetCredentialRevealFixture,
) (targetCredentialRevealResponse, error) {
	var zero targetCredentialRevealResponse
	if cookie == nil || strings.TrimSpace(cookie.Name) == "" || strings.TrimSpace(cookie.Value) == "" {
		return zero, fmt.Errorf("target credential reveal session cookie missing")
	}
	fullURL, err := normalizedAPIURL(baseURL, "/client/services/"+targetCredentialRevealServiceID+"/credentials/"+fixture.CredentialID+"/reveal")
	if err != nil {
		return zero, err
	}
	payload, err := json.Marshal(map[string]string{"reason": targetCredentialRevealReason})
	if err != nil {
		return zero, fmt.Errorf("marshal target credential reveal request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(payload))
	if err != nil {
		return zero, fmt.Errorf("build target credential reveal request")
	}
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})

	response, err := client.Do(request)
	if err != nil {
		return zero, fmt.Errorf("request target credential reveal failed")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return zero, fmt.Errorf("read target credential reveal response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return zero, targetCredentialRevealStatusError(response.StatusCode, body)
	}
	return validateTargetCredentialRevealResponse(response.Header, body, fixture)
}

func validateTargetCredentialRevealResponse(header http.Header, body []byte, fixture targetCredentialRevealFixture) (targetCredentialRevealResponse, error) {
	var zero targetCredentialRevealResponse
	if !strings.Contains(strings.ToLower(header.Get("Cache-Control")), "no-store") {
		return zero, fmt.Errorf("target credential reveal response missing no-store cache header")
	}
	if strings.ToLower(header.Get("Pragma")) != "no-cache" {
		return zero, fmt.Errorf("target credential reveal response missing no-cache pragma")
	}
	if err := assertTargetCredentialRevealBodyRedaction(body, fixture.EncryptedPayload); err != nil {
		return zero, err
	}

	var envelope struct {
		Data      targetCredentialRevealResponse `json:"data"`
		RequestID string                         `json:"request_id"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return zero, fmt.Errorf("decode target credential reveal response")
	}
	if strings.TrimSpace(envelope.RequestID) == "" {
		return zero, fmt.Errorf("target credential reveal response missing request_id")
	}
	record := envelope.Data
	if record.ID == "" || record.ID != fixture.CredentialID {
		return zero, fmt.Errorf("target credential reveal response credential id mismatch")
	}
	if record.Type != targetCredentialRevealType || record.Status != "active" || record.MaskedHint != targetCredentialRevealMaskedHint {
		return zero, fmt.Errorf("target credential reveal response metadata mismatch")
	}
	if record.RevealedAt.IsZero() || strings.TrimSpace(record.RevealExpiresMessage) == "" {
		return zero, fmt.Errorf("target credential reveal response missing reveal timing metadata")
	}
	if err := assertTargetCredentialRevealPayload(record.Payload, fixture.ExpectedPayload); err != nil {
		return zero, err
	}
	return record, nil
}

func assertTargetCredentialRevealBodyRedaction(body []byte, encryptedPayload string) error {
	bodyText := string(body)
	bodyLower := strings.ToLower(bodyText)
	blockedFields := []string{
		`"encrypted_payload"`,
		`"encryption_key_version"`,
		`"encryption_algorithm"`,
		`"secret_version"`,
		`"token_hash"`,
	}
	for _, blocked := range blockedFields {
		if strings.Contains(bodyLower, blocked) {
			return fmt.Errorf("target credential reveal response exposed blocked credential metadata")
		}
	}
	if encryptedPayload != "" && strings.Contains(bodyText, encryptedPayload) {
		return fmt.Errorf("target credential reveal response exposed encrypted credential payload")
	}
	if strings.Contains(bodyText, targetCredentialRevealSecretVersion) {
		return fmt.Errorf("target credential reveal response exposed internal secret version")
	}
	return nil
}

func assertTargetCredentialRevealPayload(payload json.RawMessage, expected string) error {
	var got map[string]string
	if err := json.Unmarshal(payload, &got); err != nil {
		return fmt.Errorf("target credential reveal payload is not expected JSON")
	}
	var want map[string]string
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		return fmt.Errorf("target credential reveal expected payload fixture is invalid")
	}
	if len(got) != len(want) {
		return fmt.Errorf("target credential reveal payload mismatch")
	}
	for key, value := range want {
		if got[key] != value {
			return fmt.Errorf("target credential reveal payload mismatch")
		}
	}
	return nil
}

func targetCredentialRevealStatusError(gotStatus int, body []byte) error {
	var apiError errorEnvelope
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Code != "" {
		return fmt.Errorf("target credential reveal expected HTTP 200, got %d (%s)", gotStatus, apiError.Error.Code)
	}
	return fmt.Errorf("target credential reveal expected HTTP 200, got %d", gotStatus)
}

func collectTargetCredentialRevealDBEvidence(
	ctx context.Context,
	conn *sql.DB,
	fixture targetCredentialRevealFixture,
	startedAt time.Time,
) (targetCredentialRevealDBEvidence, error) {
	var evidence targetCredentialRevealDBEvidence
	if err := conn.QueryRowContext(ctx, `
SELECT COALESCE(last_revealed_at IS NOT NULL AND last_revealed_by = $4::uuid, false)
FROM service_credentials
WHERE credential_id = $1::uuid
  AND tenant_id = $2::uuid
  AND service_instance_id = $3::uuid
  AND status = 'active'`,
		fixture.CredentialID,
		demoTenantID,
		targetCredentialRevealServiceID,
		demoCustomerID,
	).Scan(&evidence.LastRevealedByClient); err != nil {
		return evidence, fmt.Errorf("verify target credential reveal metadata: %w", err)
	}

	if err := conn.QueryRowContext(ctx, `
SELECT COALESCE(SUM(attempt_count), 0)
FROM service_credential_reveal_rate_limits
WHERE tenant_id = $1::uuid
  AND actor_id = $2::uuid
  AND service_instance_id = $3::uuid`,
		demoTenantID,
		demoCustomerID,
		targetCredentialRevealServiceID,
	).Scan(&evidence.RateLimitAttempts); err != nil {
		return evidence, fmt.Errorf("verify target credential reveal rate limit: %w", err)
	}

	displayIDText := fmt.Sprintf("%d", fixture.ServiceDisplayID)
	if err := conn.QueryRowContext(ctx, `
SELECT
  COUNT(*),
  COALESCE(MAX(display_id), 0),
  COALESCE(BOOL_OR(metadata_redacted ->> 'service_display_id' = $4), false),
  COALESCE(BOOL_OR(POSITION($5 IN concat_ws(' ', metadata_redacted::text, before_snapshot_redacted::text, after_snapshot_redacted::text)) > 0), false),
  COALESCE(BOOL_OR(POSITION($6 IN concat_ws(' ', metadata_redacted::text, before_snapshot_redacted::text, after_snapshot_redacted::text)) > 0), false)
FROM audit_logs
WHERE tenant_id = $1::uuid
  AND action = 'credential.revealed'
  AND target_type = 'service_credential'
  AND target_id = $2::uuid
  AND created_at >= $3`,
		demoTenantID,
		fixture.CredentialID,
		startedAt,
		displayIDText,
		"target-credential-smoke-secret",
		fixture.EncryptedPayload,
	).Scan(&evidence.AuditCount, &evidence.AuditDisplayID, &evidence.AuditHasDisplayID, &evidence.AuditLeakedSecret, &evidence.AuditLeakedCipher); err != nil {
		return evidence, fmt.Errorf("verify target credential reveal audit: %w", err)
	}
	return evidence, nil
}

func validateTargetCredentialRevealDBEvidence(evidence targetCredentialRevealDBEvidence) error {
	if !evidence.LastRevealedByClient {
		return fmt.Errorf("target credential reveal metadata did not record client actor")
	}
	if evidence.RateLimitAttempts < 1 {
		return fmt.Errorf("target credential reveal rate limit was not recorded")
	}
	if evidence.AuditCount < 1 || evidence.AuditDisplayID == 0 {
		return fmt.Errorf("target credential reveal audit was not recorded")
	}
	if !evidence.AuditHasDisplayID {
		return fmt.Errorf("target credential reveal audit missing service display id")
	}
	if evidence.AuditLeakedSecret || evidence.AuditLeakedCipher {
		return fmt.Errorf("target credential reveal audit exposed credential material")
	}
	return nil
}
