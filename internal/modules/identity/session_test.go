package identity

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestAuthServiceLoginCreatesSessionWithArgon2idPassword(t *testing.T) {
	passwordHash, err := HashPasswordArgon2idWithConfig("admin123", testArgon2idConfig())
	if err != nil {
		t.Fatalf("HashPasswordArgon2id returned error: %v", err)
	}
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{}
	service := newTestAuthService(passwordHash, sessions, now)

	result, err := service.Login(context.Background(), LoginInput{
		Email:                  " ADMIN@LOCAL.BILLING ",
		Password:               "admin123",
		LocalTenantID:          "tenant_1",
		AllowLocalTenantHeader: true,
		UserAgent:              "test-agent",
	})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected plaintext token for cookie")
	}
	if sessions.created.TokenHash == "" || sessions.created.TokenHash == result.Token {
		t.Fatalf("expected hashed token in session store, got %q", sessions.created.TokenHash)
	}
	if sessions.created.UserAgentHash == "" || sessions.created.UserAgentHash == "test-agent" {
		t.Fatalf("expected hashed user agent, got %q", sessions.created.UserAgentHash)
	}
	if sessions.created.ExpiresAt != now.Add(time.Hour) {
		t.Fatalf("expected session expiry, got %s", sessions.created.ExpiresAt)
	}
	if len(result.RoleIDs) != 1 || result.RoleIDs[0] != "role_admin" {
		t.Fatalf("expected role ids, got %+v", result.RoleIDs)
	}
	if !result.TwoFactorRequired || !result.TwoFactorSetupRequired || result.TwoFactorSatisfied {
		t.Fatalf("expected platform staff login to require 2FA setup, got %+v", result)
	}
}

func TestAuthServiceLoginRejectsInvalidPassword(t *testing.T) {
	passwordHash, err := HashPasswordArgon2idWithConfig("admin123", testArgon2idConfig())
	if err != nil {
		t.Fatalf("HashPasswordArgon2id returned error: %v", err)
	}
	service := newTestAuthService(passwordHash, &fakeSessionStore{}, time.Now())

	_, err = service.Login(context.Background(), LoginInput{
		Email:                  "admin@local.billing",
		Password:               "wrong",
		LocalTenantID:          "tenant_1",
		AllowLocalTenantHeader: true,
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthServiceLoginRequiresTenantContext(t *testing.T) {
	service := newTestAuthService("hash", &fakeSessionStore{}, time.Now())

	_, err := service.Login(context.Background(), LoginInput{
		Email:    "admin@local.billing",
		Password: "admin123",
	})
	if !errors.Is(err, ErrLoginTenantInvalid) {
		t.Fatalf("expected login tenant error, got %v", err)
	}
}

func TestAuthServiceResolveSessionReturnsIdentity(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{
		identity: SessionIdentity{
			Session: Session{ID: "session_1", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Type: UserTypePlatformStaff, Status: UserStatusActive},
			RoleIDs: []RoleID{"role_admin"},
		},
	}
	service := newTestAuthService("hash", sessions, now)

	identity, err := service.ResolveSession(context.Background(), "token")
	if err != nil {
		t.Fatalf("ResolveSession returned error: %v", err)
	}
	if sessions.resolvedHash != HashSessionToken("token") {
		t.Fatalf("expected token hash lookup, got %q", sessions.resolvedHash)
	}
	actor := identity.Actor()
	if actor.ID != "user_1" || !actor.IsPlatformAdmin || !actor.HasRole("role_admin") {
		t.Fatalf("unexpected actor from session: %+v", actor)
	}
}

func TestAuthServiceSetupTwoFactorEncryptsSecretAndAudits(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{
		identity: SessionIdentity{
			Session: Session{ID: "11111111-1111-1111-1111-111111111111", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Email: "admin@example.com", Type: UserTypePlatformStaff, Status: UserStatusActive},
		},
	}
	twoFactor := &fakeTwoFactorStore{}
	audit := &fakeAuthAuditAppender{}
	service := newTestAuthServiceWithStores("hash", sessions, twoFactor, audit, now)

	result, err := service.SetupTwoFactor(context.Background(), "token")
	if err != nil {
		t.Fatalf("SetupTwoFactor returned error: %v", err)
	}
	if result.Method != TwoFactorMethodTOTP || result.Secret == "" || result.ProvisionURI == "" {
		t.Fatalf("unexpected setup result: %+v", result)
	}
	if twoFactor.upsert.SecretCiphertext == "" || strings.Contains(twoFactor.upsert.SecretCiphertext, result.Secret) {
		t.Fatalf("expected encrypted secret, got %q", twoFactor.upsert.SecretCiphertext)
	}
	if twoFactor.status != TwoFactorStatusRequired {
		t.Fatalf("expected required status, got %q", twoFactor.status)
	}
	if len(audit.inputs) != 1 || audit.inputs[0].Action != AuditActionTwoFactorSetup {
		t.Fatalf("expected setup audit, got %+v", audit.inputs)
	}
}

func TestAuthServiceSetupTwoFactorRejectsEnabledMethod(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{
		identity: SessionIdentity{
			Session: Session{ID: "11111111-1111-1111-1111-111111111111", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Email: "admin@example.com", Type: UserTypePlatformStaff, Status: UserStatusActive},
		},
	}
	twoFactor := &fakeTwoFactorStore{
		method: TwoFactorMethod{
			TenantID:         "tenant_1",
			UserID:           "user_1",
			Method:           TwoFactorMethodTOTP,
			SecretCiphertext: "encrypted-secret",
			EnabledAt:        now.Add(-time.Hour),
		},
	}
	audit := &fakeAuthAuditAppender{}
	service := newTestAuthServiceWithStores("hash", sessions, twoFactor, audit, now)

	_, err := service.SetupTwoFactor(context.Background(), "token")
	if !errors.Is(err, ErrTwoFactorAlreadyEnabled) {
		t.Fatalf("expected already-enabled 2FA error, got %v", err)
	}
	if twoFactor.upsert.SecretCiphertext != "" {
		t.Fatalf("setup should not rotate enabled TOTP secret, got %+v", twoFactor.upsert)
	}
	if len(audit.inputs) != 0 {
		t.Fatalf("setup should not audit a rejected secret rotation, got %+v", audit.inputs)
	}
}

func TestAuthServiceVerifyTwoFactorMarksSessionSatisfied(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	cipher, err := NewAESGCMSecretCipher(testCipherKey())
	if err != nil {
		t.Fatalf("NewAESGCMSecretCipher returned error: %v", err)
	}
	encrypted, err := cipher.Encrypt("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}
	code, err := TOTPCodeAt("JBSWY3DPEHPK3PXP", now)
	if err != nil {
		t.Fatalf("TOTPCodeAt returned error: %v", err)
	}
	sessions := &fakeSessionStore{
		identity: SessionIdentity{
			Session: Session{ID: "11111111-1111-1111-1111-111111111111", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Email: "admin@example.com", Type: UserTypePlatformStaff, Status: UserStatusActive},
		},
	}
	twoFactor := &fakeTwoFactorStore{
		method: TwoFactorMethod{TenantID: "tenant_1", UserID: "user_1", Method: TwoFactorMethodTOTP, SecretCiphertext: encrypted},
	}
	audit := &fakeAuthAuditAppender{}
	service := newTestAuthServiceWithStores("hash", sessions, twoFactor, audit, now)

	result, err := service.VerifyTwoFactor(context.Background(), "token", code)
	if err != nil {
		t.Fatalf("VerifyTwoFactor returned error: %v", err)
	}
	if !result.Session.TwoFactorSatisfied() {
		t.Fatal("expected session to be marked 2FA satisfied")
	}
	if sessions.satisfiedHash != HashSessionToken("token") {
		t.Fatalf("expected session token hash, got %q", sessions.satisfiedHash)
	}
	if twoFactor.status != TwoFactorStatusEnabled || twoFactor.enabledAt.IsZero() {
		t.Fatalf("expected enabled two factor status, got status=%q enabled=%s", twoFactor.status, twoFactor.enabledAt)
	}
	if len(audit.inputs) != 1 || audit.inputs[0].Action != AuditActionTwoFactorSuccess {
		t.Fatalf("expected success audit, got %+v", audit.inputs)
	}
}

func TestAuthServiceVerifyTwoFactorAuditsFailure(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	cipher, err := NewAESGCMSecretCipher(testCipherKey())
	if err != nil {
		t.Fatalf("NewAESGCMSecretCipher returned error: %v", err)
	}
	encrypted, err := cipher.Encrypt("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}
	sessions := &fakeSessionStore{
		identity: SessionIdentity{
			Session: Session{ID: "11111111-1111-1111-1111-111111111111", TenantID: "tenant_1", UserID: "user_1"},
			User:    User{ID: "user_1", TenantID: "tenant_1", Email: "admin@example.com", Type: UserTypePlatformStaff, Status: UserStatusActive},
		},
	}
	twoFactor := &fakeTwoFactorStore{
		method: TwoFactorMethod{TenantID: "tenant_1", UserID: "user_1", Method: TwoFactorMethodTOTP, SecretCiphertext: encrypted},
	}
	audit := &fakeAuthAuditAppender{}
	service := newTestAuthServiceWithStores("hash", sessions, twoFactor, audit, now)

	_, err = service.VerifyTwoFactor(context.Background(), "token", "000000")
	if !errors.Is(err, ErrTwoFactorCodeInvalid) {
		t.Fatalf("expected invalid code error, got %v", err)
	}
	if len(audit.inputs) != 1 || audit.inputs[0].Action != AuditActionTwoFactorFailure {
		t.Fatalf("expected failure audit, got %+v", audit.inputs)
	}
}

func newTestAuthService(passwordHash string, sessions *fakeSessionStore, now time.Time) *AuthService {
	return newTestAuthServiceWithStores(passwordHash, sessions, &fakeTwoFactorStore{}, &fakeAuthAuditAppender{}, now)
}

func newTestAuthServiceWithStores(passwordHash string, sessions *fakeSessionStore, twoFactor *fakeTwoFactorStore, audit *fakeAuthAuditAppender, now time.Time) *AuthService {
	cipher, err := NewAESGCMSecretCipher(testCipherKey())
	if err != nil {
		panic(err)
	}
	return NewAuthService(AuthServiceOptions{
		Tenants: &fakeAuthTenantStore{
			tenants: map[tenant.ID]tenant.Tenant{
				"tenant_1": {ID: "tenant_1", Status: tenant.StatusActive},
			},
		},
		Users: &fakeAuthUserStore{
			user: User{
				ID:              "user_1",
				TenantID:        "tenant_1",
				Email:           "admin@local.billing",
				PasswordHash:    passwordHash,
				Type:            UserTypePlatformStaff,
				Status:          UserStatusActive,
				TwoFactorStatus: TwoFactorStatusDisabled,
			},
		},
		Sessions:   sessions,
		TwoFactor:  twoFactor,
		Roles:      fakeRoleReader{roleIDs: []RoleID{"role_admin"}},
		Cipher:     cipher,
		Audit:      audit,
		SessionTTL: time.Hour,
		Now:        func() time.Time { return now },
	})
}

type fakeAuthTenantStore struct {
	tenants map[tenant.ID]tenant.Tenant
	domains map[string]tenant.Domain
}

func (store *fakeAuthTenantStore) GetByID(ctx context.Context, tenantID tenant.ID) (tenant.Tenant, error) {
	record, ok := store.tenants[tenantID]
	if !ok {
		return tenant.Tenant{}, tenant.ErrTenantNotFound
	}
	return record, nil
}

func (store *fakeAuthTenantStore) FindActiveDomain(ctx context.Context, domain string) (tenant.Domain, error) {
	record, ok := store.domains[domain]
	if !ok {
		return tenant.Domain{}, tenant.ErrDomainNotFound
	}
	return record, nil
}

type fakeAuthUserStore struct {
	user User
}

func (store *fakeAuthUserStore) CreateUser(ctx context.Context, input CreateUserInput) (User, error) {
	return User{}, nil
}

func (store *fakeAuthUserStore) GetUserByID(ctx context.Context, tenantID tenant.ID, userID UserID) (User, error) {
	return User{}, nil
}

func (store *fakeAuthUserStore) FindUserByEmail(ctx context.Context, tenantID tenant.ID, email string) (User, error) {
	if tenantID != store.user.TenantID || email != store.user.Email {
		return User{}, ErrUserNotFound
	}
	return store.user, nil
}

func (store *fakeAuthUserStore) ListUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error) {
	return nil, nil
}

type fakeSessionStore struct {
	created       CreateSessionInput
	identity      SessionIdentity
	resolvedHash  string
	revokedHash   string
	satisfiedHash string
}

func (store *fakeSessionStore) CreateSession(ctx context.Context, input CreateSessionInput) (Session, error) {
	store.created = input
	return Session{
		ID:        "session_1",
		TenantID:  input.TenantID,
		UserID:    input.UserID,
		TokenHash: input.TokenHash,
		ExpiresAt: input.ExpiresAt,
	}, nil
}

func (store *fakeSessionStore) FindSessionIdentityByTokenHash(ctx context.Context, tokenHash string, now time.Time) (SessionIdentity, error) {
	store.resolvedHash = tokenHash
	if store.identity.Session.ID == "" {
		return SessionIdentity{}, ErrSessionInvalid
	}
	return store.identity, nil
}

func (store *fakeSessionStore) RevokeSessionByTokenHash(ctx context.Context, tokenHash string, now time.Time) error {
	store.revokedHash = tokenHash
	return nil
}

func (store *fakeSessionStore) MarkSessionTwoFactorSatisfied(ctx context.Context, tokenHash string, now time.Time) (Session, error) {
	store.satisfiedHash = tokenHash
	return Session{
		ID:                   "session_1",
		TenantID:             "tenant_1",
		UserID:               "user_1",
		TokenHash:            tokenHash,
		TwoFactorSatisfiedAt: now,
	}, nil
}

type fakeTwoFactorStore struct {
	upsert    UpsertTOTPSecretInput
	method    TwoFactorMethod
	status    TwoFactorStatus
	enabledAt time.Time
}

func (store *fakeTwoFactorStore) UpsertTOTPSecret(ctx context.Context, input UpsertTOTPSecretInput) (TwoFactorMethod, error) {
	store.upsert = input
	store.method = TwoFactorMethod{
		TenantID:         input.TenantID,
		UserID:           input.UserID,
		Method:           TwoFactorMethodTOTP,
		SecretCiphertext: input.SecretCiphertext,
	}
	return store.method, nil
}

func (store *fakeTwoFactorStore) GetTOTPMethod(ctx context.Context, tenantID tenant.ID, userID UserID) (TwoFactorMethod, error) {
	if store.method.SecretCiphertext == "" {
		return TwoFactorMethod{}, ErrTwoFactorMethodNotFound
	}
	return store.method, nil
}

func (store *fakeTwoFactorStore) MarkTOTPEnabled(ctx context.Context, tenantID tenant.ID, userID UserID, now time.Time) error {
	store.enabledAt = now
	return nil
}

func (store *fakeTwoFactorStore) SetUserTwoFactorStatus(ctx context.Context, tenantID tenant.ID, userID UserID, status TwoFactorStatus) error {
	store.status = status
	return nil
}

type fakeAuthAuditAppender struct {
	inputs []AuthAuditInput
}

func (appender *fakeAuthAuditAppender) AppendAuthAudit(ctx context.Context, input AuthAuditInput) error {
	appender.inputs = append(appender.inputs, input)
	return nil
}

type fakeRoleReader struct {
	roleIDs []RoleID
}

func (reader fakeRoleReader) ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID UserID) ([]RoleID, error) {
	return append([]RoleID(nil), reader.roleIDs...), nil
}
