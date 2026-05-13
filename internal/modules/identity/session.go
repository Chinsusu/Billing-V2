package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrAuthStoreMissing    = errors.New("auth store missing")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserInactive        = errors.New("user inactive")
	ErrTwoFactorRequired   = errors.New("two factor required")
	ErrSessionTokenMissing = errors.New("session token missing")
	ErrSessionInvalid      = errors.New("session invalid")
	ErrSessionExpired      = errors.New("session expired")
	ErrLoginTenantInvalid  = errors.New("login tenant invalid")
	ErrSessionTTLMissing   = errors.New("session ttl missing")
	ErrSessionMissing      = errors.New("session missing")
)

type LoginInput struct {
	Email                  string
	Password               string
	Domain                 string
	LocalTenantID          tenant.ID
	AllowLocalTenantHeader bool
	UserAgent              string
	ClientIP               string
}

type LoginResult struct {
	Token                  string
	Session                Session
	User                   User
	RoleIDs                []RoleID
	ExpiresAt              time.Time
	TwoFactorRequired      bool
	TwoFactorSatisfied     bool
	TwoFactorSetupRequired bool
}

type Session struct {
	ID                   string
	TenantID             tenant.ID
	UserID               UserID
	TokenHash            string
	UserAgentHash        string
	ExpiresAt            time.Time
	RevokedAt            time.Time
	LastSeenAt           time.Time
	TwoFactorSatisfiedAt time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type SessionIdentity struct {
	Session Session
	User    User
	RoleIDs []RoleID
}

func (identity SessionIdentity) Actor() Actor {
	return NewActor(identity.User.ID, identity.User.TenantID, actorTypeForUser(identity.User), identity.RoleIDs...)
}

type AuthTenantStore interface {
	GetByID(ctx context.Context, tenantID tenant.ID) (tenant.Tenant, error)
	FindActiveDomain(ctx context.Context, domain string) (tenant.Domain, error)
}

type SessionStore interface {
	CreateSession(ctx context.Context, input CreateSessionInput) (Session, error)
	FindSessionIdentityByTokenHash(ctx context.Context, tokenHash string, now time.Time) (SessionIdentity, error)
	MarkSessionTwoFactorSatisfied(ctx context.Context, tokenHash string, now time.Time) (Session, error)
	RevokeSessionByTokenHash(ctx context.Context, tokenHash string, now time.Time) error
	RevokeUserSessions(ctx context.Context, tenantID tenant.ID, userID UserID, now time.Time) error
}

type UserRoleReader interface {
	ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID UserID) ([]RoleID, error)
}

type CreateSessionInput struct {
	TenantID      tenant.ID
	UserID        UserID
	TokenHash     string
	UserAgentHash string
	ExpiresAt     time.Time
}

type AuthService struct {
	tenants          AuthTenantStore
	users            UserStore
	sessions         SessionStore
	twoFactor        TwoFactorStore
	rateLimits       AuthRateLimitStore
	passwordResets   PasswordResetStore
	resetDelivery    PasswordResetDelivery
	roles            UserRoleReader
	cipher           SecretCipher
	audit            AuthAuditAppender
	sessionTTL       time.Duration
	passwordResetTTL time.Duration
	now              func() time.Time
}

type AuthServiceOptions struct {
	Tenants          AuthTenantStore
	Users            UserStore
	Sessions         SessionStore
	TwoFactor        TwoFactorStore
	RateLimits       AuthRateLimitStore
	PasswordResets   PasswordResetStore
	ResetDelivery    PasswordResetDelivery
	Roles            UserRoleReader
	Cipher           SecretCipher
	Audit            AuthAuditAppender
	SessionTTL       time.Duration
	PasswordResetTTL time.Duration
	Now              func() time.Time
}

func NewAuthService(options AuthServiceOptions) *AuthService {
	now := options.Now
	if now == nil {
		now = time.Now
	}
	return &AuthService{
		tenants:          options.Tenants,
		users:            options.Users,
		sessions:         options.Sessions,
		twoFactor:        options.TwoFactor,
		rateLimits:       options.RateLimits,
		passwordResets:   options.PasswordResets,
		resetDelivery:    options.ResetDelivery,
		roles:            options.Roles,
		cipher:           options.Cipher,
		audit:            options.Audit,
		sessionTTL:       options.SessionTTL,
		passwordResetTTL: options.PasswordResetTTL,
		now:              now,
	}
}

func (service *AuthService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	if service == nil || service.tenants == nil || service.users == nil || service.sessions == nil || service.roles == nil {
		return LoginResult{}, ErrAuthStoreMissing
	}
	if service.sessionTTL <= 0 {
		return LoginResult{}, ErrSessionTTLMissing
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Email == "" || input.Password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}
	tenantID, err := service.resolveLoginTenant(ctx, input)
	if err != nil {
		return LoginResult{}, err
	}
	if err := service.enforceAuthRateLimit(ctx, AuthActionLogin, tenantID, input.Email, input.ClientIP); err != nil {
		return LoginResult{}, err
	}
	user, err := service.users.FindUserByEmail(ctx, tenantID, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}
	if user.Status != UserStatusActive {
		return LoginResult{}, ErrUserInactive
	}
	ok, err := VerifyPasswordArgon2id(input.Password, user.PasswordHash)
	if err != nil || !ok {
		return LoginResult{}, ErrInvalidCredentials
	}
	roleIDs, err := service.roles.ListRoleIDsForUser(ctx, tenantID, user.ID)
	if err != nil {
		return LoginResult{}, err
	}
	twoFactorRequired := userRequiresTwoFactor(user)
	twoFactorSetupRequired := false
	if twoFactorRequired && service.twoFactor != nil {
		if _, err := service.twoFactor.GetTOTPMethod(ctx, tenantID, user.ID); errors.Is(err, ErrTwoFactorMethodNotFound) {
			twoFactorSetupRequired = true
		} else if err != nil {
			return LoginResult{}, err
		}
	}
	token, tokenHash, err := newSessionToken()
	if err != nil {
		return LoginResult{}, err
	}
	now := service.now().UTC()
	session, err := service.sessions.CreateSession(ctx, CreateSessionInput{
		TenantID:      tenantID,
		UserID:        user.ID,
		TokenHash:     tokenHash,
		UserAgentHash: HashSessionUserAgent(input.UserAgent),
		ExpiresAt:     now.Add(service.sessionTTL),
	})
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token:                  token,
		Session:                session,
		User:                   user,
		RoleIDs:                roleIDs,
		ExpiresAt:              session.ExpiresAt,
		TwoFactorRequired:      twoFactorRequired,
		TwoFactorSatisfied:     session.TwoFactorSatisfied(),
		TwoFactorSetupRequired: twoFactorSetupRequired,
	}, nil
}

func (service *AuthService) ResolveSession(ctx context.Context, token string) (SessionIdentity, error) {
	if service == nil || service.sessions == nil {
		return SessionIdentity{}, ErrAuthStoreMissing
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return SessionIdentity{}, ErrSessionTokenMissing
	}
	identity, err := service.sessions.FindSessionIdentityByTokenHash(ctx, HashSessionToken(token), service.now().UTC())
	if err != nil {
		return SessionIdentity{}, err
	}
	return identity, nil
}

func (session Session) TwoFactorSatisfied() bool {
	return !session.TwoFactorSatisfiedAt.IsZero()
}

func (service *AuthService) SetupTwoFactor(ctx context.Context, token string) (SetupTwoFactorResult, error) {
	if service == nil || service.sessions == nil || service.twoFactor == nil || service.cipher == nil {
		return SetupTwoFactorResult{}, ErrAuthStoreMissing
	}
	identity, err := service.ResolveSession(ctx, token)
	if err != nil {
		return SetupTwoFactorResult{}, err
	}
	if !userRequiresTwoFactor(identity.User) {
		return SetupTwoFactorResult{}, ErrTwoFactorNotAllowed
	}
	method, err := service.twoFactor.GetTOTPMethod(ctx, identity.User.TenantID, identity.User.ID)
	if err == nil && !method.EnabledAt.IsZero() {
		return SetupTwoFactorResult{}, ErrTwoFactorAlreadyEnabled
	}
	if err != nil && !errors.Is(err, ErrTwoFactorMethodNotFound) {
		return SetupTwoFactorResult{}, err
	}
	secret, err := GenerateTOTPSecret()
	if err != nil {
		return SetupTwoFactorResult{}, err
	}
	ciphertext, err := service.cipher.Encrypt(secret)
	if err != nil {
		return SetupTwoFactorResult{}, err
	}
	if _, err := service.twoFactor.UpsertTOTPSecret(ctx, UpsertTOTPSecretInput{
		TenantID:         identity.User.TenantID,
		UserID:           identity.User.ID,
		SecretCiphertext: ciphertext,
	}); err != nil {
		return SetupTwoFactorResult{}, err
	}
	if err := service.twoFactor.SetUserTwoFactorStatus(ctx, identity.User.TenantID, identity.User.ID, TwoFactorStatusRequired); err != nil {
		return SetupTwoFactorResult{}, err
	}
	if err := service.appendTwoFactorAudit(ctx, identity, AuditActionTwoFactorSetup); err != nil {
		return SetupTwoFactorResult{}, err
	}
	return SetupTwoFactorResult{
		Method:       TwoFactorMethodTOTP,
		Secret:       secret,
		ProvisionURI: totpProvisionURI(identity.User.Email, secret),
	}, nil
}

func (service *AuthService) VerifyTwoFactor(ctx context.Context, token string, code string) (VerifyTwoFactorResult, error) {
	if service == nil || service.sessions == nil || service.twoFactor == nil || service.cipher == nil {
		return VerifyTwoFactorResult{}, ErrAuthStoreMissing
	}
	identity, err := service.ResolveSession(ctx, token)
	if err != nil {
		return VerifyTwoFactorResult{}, err
	}
	method, err := service.twoFactor.GetTOTPMethod(ctx, identity.User.TenantID, identity.User.ID)
	if err != nil {
		if errors.Is(err, ErrTwoFactorMethodNotFound) {
			return VerifyTwoFactorResult{}, ErrTwoFactorSetupRequired
		}
		return VerifyTwoFactorResult{}, err
	}
	secret, err := service.cipher.Decrypt(method.SecretCiphertext)
	if err != nil {
		return VerifyTwoFactorResult{}, err
	}
	ok, err := VerifyTOTPCode(secret, code, service.now().UTC())
	if err != nil {
		_ = service.appendTwoFactorAudit(ctx, identity, AuditActionTwoFactorFailure)
		return VerifyTwoFactorResult{}, err
	}
	if !ok {
		_ = service.appendTwoFactorAudit(ctx, identity, AuditActionTwoFactorFailure)
		return VerifyTwoFactorResult{}, ErrTwoFactorCodeInvalid
	}
	now := service.now().UTC()
	if err := service.twoFactor.MarkTOTPEnabled(ctx, identity.User.TenantID, identity.User.ID, now); err != nil {
		return VerifyTwoFactorResult{}, err
	}
	if err := service.twoFactor.SetUserTwoFactorStatus(ctx, identity.User.TenantID, identity.User.ID, TwoFactorStatusEnabled); err != nil {
		return VerifyTwoFactorResult{}, err
	}
	session, err := service.sessions.MarkSessionTwoFactorSatisfied(ctx, HashSessionToken(token), now)
	if err != nil {
		return VerifyTwoFactorResult{}, err
	}
	if err := service.appendTwoFactorAudit(ctx, identity, AuditActionTwoFactorSuccess); err != nil {
		return VerifyTwoFactorResult{}, err
	}
	return VerifyTwoFactorResult{Session: session, User: identity.User}, nil
}

func (service *AuthService) Logout(ctx context.Context, token string) error {
	if service == nil || service.sessions == nil {
		return ErrAuthStoreMissing
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	return service.sessions.RevokeSessionByTokenHash(ctx, HashSessionToken(token), service.now().UTC())
}

func userRequiresTwoFactor(user User) bool {
	return user.Type == UserTypePlatformStaff ||
		user.TwoFactorStatus == TwoFactorStatusRequired ||
		user.TwoFactorStatus == TwoFactorStatusEnabled
}

func (service *AuthService) resolveLoginTenant(ctx context.Context, input LoginInput) (tenant.ID, error) {
	domain := normalizeDomain(input.Domain)
	if domain != "" {
		resolved, err := service.tenants.FindActiveDomain(ctx, domain)
		if err == nil {
			return resolved.TenantID, nil
		}
		if !errors.Is(err, tenant.ErrDomainNotFound) {
			return "", err
		}
	}
	if input.AllowLocalTenantHeader && !input.LocalTenantID.Empty() {
		record, err := service.tenants.GetByID(ctx, input.LocalTenantID)
		if err != nil {
			return "", err
		}
		if record.Status != tenant.StatusActive {
			return "", tenant.ErrTenantStatusInvalid
		}
		return record.ID, nil
	}
	return "", ErrLoginTenantInvalid
}

func newSessionToken() (string, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	return token, HashSessionToken(token), nil
}

func HashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func HashSessionUserAgent(userAgent string) string {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(userAgent))
	return hex.EncodeToString(sum[:])
}

func actorTypeForUser(user User) ActorType {
	switch user.Type {
	case UserTypePlatformStaff:
		return ActorTypePlatformStaff
	case UserTypeResellerStaff:
		return ActorTypeResellerStaff
	case UserTypeClient:
		return ActorTypeClient
	default:
		return ""
	}
}

func normalizeDomain(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if host, _, ok := strings.Cut(value, ":"); ok {
		value = host
	}
	return strings.Trim(value, ".")
}
