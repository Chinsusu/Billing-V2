package identity

import (
	"context"
	"errors"
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

func newTestAuthService(passwordHash string, sessions *fakeSessionStore, now time.Time) *AuthService {
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
		Roles:      fakeRoleReader{roleIDs: []RoleID{"role_admin"}},
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
	created      CreateSessionInput
	identity     SessionIdentity
	resolvedHash string
	revokedHash  string
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

type fakeRoleReader struct {
	roleIDs []RoleID
}

func (reader fakeRoleReader) ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID UserID) ([]RoleID, error) {
	return append([]RoleID(nil), reader.roleIDs...), nil
}
