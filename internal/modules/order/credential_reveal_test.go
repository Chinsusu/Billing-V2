package order

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

func TestServiceRevealServiceCredentialDecryptsMarksAndAudits(t *testing.T) {
	cipher, encrypted := testCredentialRevealCipher(t, `{"username":"root","password":"fixture-access"}`)
	now := time.Date(2026, 5, 13, 9, 30, 0, 0, time.UTC)
	store := &fakeCredentialRevealStore{
		service: ServiceInstance{ID: "service_1", DisplayID: 50001, TenantID: "tenant_1"},
		credential: ServiceCredential{
			ID:               "credential_1",
			TenantID:         "tenant_1",
			ServiceID:        "service_1",
			Type:             CredentialTypeVPSRoot,
			EncryptedPayload: encrypted,
			MaskedHint:       "root / ****",
			Status:           CredentialStatusActive,
		},
	}
	limiter := &fakeCredentialRevealLimiter{counter: CredentialRevealRateLimitCounter{AttemptCount: 1}}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithOptions(ServiceOptions{
		Store:                  store,
		Credentials:            store,
		Audit:                  auditLog,
		CredentialCipher:       cipher,
		CredentialRevealLimits: limiter,
		Now:                    func() time.Time { return now },
	})

	result, err := service.RevealServiceCredential(context.Background(), RevealServiceCredentialInput{
		TenantID:     tenant.ID("tenant_1"),
		ServiceID:    ServiceID("service_1"),
		CredentialID: CredentialID("credential_1"),
		ActorID:      identity.UserID("buyer_1"),
		BuyerUserID:  identity.UserID("buyer_1"),
		ClientIP:     "203.0.113.10",
		UserAgent:    "billing-test",
		Reason:       "customer requested access",
	})
	if err != nil {
		t.Fatalf("expected reveal result: %v", err)
	}
	if !strings.Contains(string(result.Payload), "fixture-access") || result.RevealedAt != now {
		t.Fatalf("unexpected reveal result: %+v", result)
	}
	if store.serviceLookup.BuyerUserID != identity.UserID("buyer_1") {
		t.Fatalf("expected buyer scope on service lookup, got %+v", store.serviceLookup)
	}
	if store.markInput.ActorID != identity.UserID("buyer_1") || store.markInput.RevealedAt != now {
		t.Fatalf("expected reveal mark, got %+v", store.markInput)
	}
	if limiter.input.WindowStart != now.Truncate(credentialRevealRateLimitWindow) {
		t.Fatalf("expected truncated rate window, got %+v", limiter.input)
	}
	if auditLog.calls != 1 ||
		auditLog.input.Action != credentialAuditActionRevealed ||
		auditLog.input.TargetID != audit.TargetID("credential_1") ||
		auditLog.input.IPAddress != "203.0.113.10" ||
		auditLog.input.UserAgent != "billing-test" {
		t.Fatalf("unexpected audit input: %+v", auditLog.input)
	}
	metadata := string(auditLog.input.MetadataRedacted)
	if strings.Contains(metadata, "fixture-access") || strings.Contains(metadata, encrypted) {
		t.Fatalf("audit metadata must not contain plaintext or encrypted payload: %s", metadata)
	}
	if !strings.Contains(metadata, "customer requested access") || !strings.Contains(metadata, "50001") {
		t.Fatalf("expected context in audit metadata, got %s", metadata)
	}
}

func TestServiceRevealServiceCredentialStopsWhenRateLimited(t *testing.T) {
	cipher, encrypted := testCredentialRevealCipher(t, `{"value":"fixture-access"}`)
	store := &fakeCredentialRevealStore{
		service: ServiceInstance{ID: "service_1", TenantID: "tenant_1"},
		credential: ServiceCredential{
			ID:               "credential_1",
			TenantID:         "tenant_1",
			ServiceID:        "service_1",
			Type:             CredentialTypeProxyAuth,
			EncryptedPayload: encrypted,
			MaskedHint:       "proxy / ****",
			Status:           CredentialStatusActive,
		},
	}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithOptions(ServiceOptions{
		Store:                  store,
		Credentials:            store,
		Audit:                  auditLog,
		CredentialCipher:       cipher,
		CredentialRevealLimits: &fakeCredentialRevealLimiter{counter: CredentialRevealRateLimitCounter{AttemptCount: credentialRevealRateLimitMaxAttempts + 1}},
		Now:                    func() time.Time { return time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC) },
	})

	_, err := service.RevealServiceCredential(context.Background(), RevealServiceCredentialInput{
		TenantID:     tenant.ID("tenant_1"),
		ServiceID:    ServiceID("service_1"),
		CredentialID: CredentialID("credential_1"),
		ActorID:      identity.UserID("buyer_1"),
	})
	if !errors.Is(err, ErrCredentialRevealRateLimited) {
		t.Fatalf("expected rate limit error, got %v", err)
	}
	if store.markCalls != 0 || auditLog.calls != 0 {
		t.Fatalf("rate limited reveal should not mark or audit success, mark=%d audit=%d", store.markCalls, auditLog.calls)
	}
}

func TestServiceRevealServiceCredentialPropagatesOwnerScopeDenial(t *testing.T) {
	cipher, _ := testCredentialRevealCipher(t, `{"value":"fixture-access"}`)
	store := &fakeCredentialRevealStore{getServiceErr: ErrServiceNotFound}
	service := NewServiceWithOptions(ServiceOptions{
		Store:                  store,
		Credentials:            store,
		CredentialCipher:       cipher,
		CredentialRevealLimits: &fakeCredentialRevealLimiter{counter: CredentialRevealRateLimitCounter{AttemptCount: 1}},
	})

	_, err := service.RevealServiceCredential(context.Background(), RevealServiceCredentialInput{
		TenantID:     tenant.ID("tenant_1"),
		ServiceID:    ServiceID("service_1"),
		CredentialID: CredentialID("credential_1"),
		ActorID:      identity.UserID("buyer_1"),
		BuyerUserID:  identity.UserID("buyer_1"),
	})
	if !errors.Is(err, ErrServiceNotFound) {
		t.Fatalf("expected service not found, got %v", err)
	}
	if store.credentialLookup.ID != "" {
		t.Fatalf("credential should not be loaded after owner denial, got %+v", store.credentialLookup)
	}
}

func testCredentialRevealCipher(t *testing.T, plaintext string) (secrets.Cipher, string) {
	t.Helper()
	cipher, err := secrets.NewAESGCMCipher("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("NewAESGCMCipher returned error: %v", err)
	}
	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}
	return cipher, encrypted
}

type fakeCredentialRevealStore struct {
	fakeOrderStore

	service          ServiceInstance
	serviceLookup    ServiceInstanceLookup
	getServiceErr    error
	credential       ServiceCredential
	credentialLookup ServiceCredentialLookup
	markInput        MarkServiceCredentialRevealedInput
	markCalls        int
}

func (store *fakeCredentialRevealStore) GetServiceInstance(_ context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	store.serviceLookup = lookup
	if store.getServiceErr != nil {
		return ServiceInstance{}, store.getServiceErr
	}
	return store.service, nil
}

func (store *fakeCredentialRevealStore) ListServiceCredentials(_ context.Context, filter ServiceCredentialFilter) ([]ServiceCredential, error) {
	if filter.TenantID != store.credential.TenantID || filter.ServiceID != store.credential.ServiceID {
		return nil, ErrCredentialNotFound
	}
	return []ServiceCredential{store.credential}, nil
}

func (store *fakeCredentialRevealStore) GetServiceCredential(_ context.Context, lookup ServiceCredentialLookup) (ServiceCredential, error) {
	store.credentialLookup = lookup
	if lookup.ID != store.credential.ID || lookup.TenantID != store.credential.TenantID || lookup.ServiceID != store.credential.ServiceID {
		return ServiceCredential{}, ErrCredentialNotFound
	}
	return store.credential, nil
}

func (store *fakeCredentialRevealStore) MarkServiceCredentialRevealed(_ context.Context, input MarkServiceCredentialRevealedInput) (ServiceCredential, error) {
	store.markCalls++
	store.markInput = input
	store.credential.LastRevealedAt = input.RevealedAt
	store.credential.LastRevealedBy = input.ActorID
	return store.credential, nil
}

type fakeCredentialRevealLimiter struct {
	input   CredentialRevealRateLimitInput
	counter CredentialRevealRateLimitCounter
	err     error
}

func (limiter *fakeCredentialRevealLimiter) IncrementCredentialRevealRateLimit(_ context.Context, input CredentialRevealRateLimitInput) (CredentialRevealRateLimitCounter, error) {
	limiter.input = input
	if limiter.err != nil {
		return CredentialRevealRateLimitCounter{}, limiter.err
	}
	return limiter.counter, nil
}
