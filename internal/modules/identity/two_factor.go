package identity

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrTwoFactorStoreMissing   = errors.New("two factor store missing")
	ErrTwoFactorMethodNotFound = errors.New("two factor method not found")
	ErrTwoFactorNotAllowed     = errors.New("two factor not allowed")
	ErrTwoFactorSetupRequired  = errors.New("two factor setup required")
	ErrTwoFactorAlreadyEnabled = errors.New("two factor already enabled")
)

const TwoFactorMethodTOTP = "totp"

type TwoFactorMethod struct {
	TenantID         tenant.ID
	UserID           UserID
	Method           string
	SecretCiphertext string
	EnabledAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SetupTwoFactorResult struct {
	Method       string
	Secret       string
	ProvisionURI string
}

type VerifyTwoFactorResult struct {
	Session Session
	User    User
}

type TwoFactorStore interface {
	UpsertTOTPSecret(ctx context.Context, input UpsertTOTPSecretInput) (TwoFactorMethod, error)
	GetTOTPMethod(ctx context.Context, tenantID tenant.ID, userID UserID) (TwoFactorMethod, error)
	MarkTOTPEnabled(ctx context.Context, tenantID tenant.ID, userID UserID, now time.Time) error
	SetUserTwoFactorStatus(ctx context.Context, tenantID tenant.ID, userID UserID, status TwoFactorStatus) error
}

type UpsertTOTPSecretInput struct {
	TenantID         tenant.ID
	UserID           UserID
	SecretCiphertext string
}
