package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserIDMissing          = errors.New("user id missing")
	ErrEmailMissing           = errors.New("email missing")
	ErrPasswordHashMissing    = errors.New("password hash missing")
	ErrUserTypeInvalid        = errors.New("user type invalid")
	ErrUserStatusInvalid      = errors.New("user status invalid")
	ErrTwoFactorStatusInvalid = errors.New("two factor status invalid")
)

type UserType string

const (
	UserTypePlatformStaff UserType = "platform_staff"
	UserTypeResellerStaff UserType = "reseller_staff"
	UserTypeClient        UserType = "client"
)

func (userType UserType) Valid() bool {
	switch userType {
	case UserTypePlatformStaff, UserTypeResellerStaff, UserTypeClient:
		return true
	default:
		return false
	}
}

type UserStatus string

const (
	UserStatusActive              UserStatus = "active"
	UserStatusSuspended           UserStatus = "suspended"
	UserStatusDisabled            UserStatus = "disabled"
	UserStatusPendingVerification UserStatus = "pending_verification"
)

func (status UserStatus) Valid() bool {
	switch status {
	case UserStatusActive, UserStatusSuspended, UserStatusDisabled, UserStatusPendingVerification:
		return true
	default:
		return false
	}
}

type TwoFactorStatus string

const (
	TwoFactorStatusRequired TwoFactorStatus = "required"
	TwoFactorStatusEnabled  TwoFactorStatus = "enabled"
	TwoFactorStatusDisabled TwoFactorStatus = "disabled"
)

func (status TwoFactorStatus) Valid() bool {
	switch status {
	case TwoFactorStatusRequired, TwoFactorStatusEnabled, TwoFactorStatusDisabled:
		return true
	default:
		return false
	}
}

type User struct {
	ID               UserID
	DisplayID        int64
	TenantID         tenant.ID
	Email            string
	EmailVerifiedAt  time.Time
	PasswordHash     string
	FullName         string
	Type             UserType
	Status           UserStatus
	TwoFactorStatus  TwoFactorStatus
	LastLoginAt      time.Time
	FailedLoginCount int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserListFilter struct {
	TenantID  tenant.ID
	Type      UserType
	Status    UserStatus
	DisplayID int64
	Email     string
	Limit     int
}

type UserSummary struct {
	User       User
	TenantName string
	TenantSlug string
}

type CreateUserInput struct {
	TenantID        tenant.ID
	Email           string
	EmailVerifiedAt time.Time
	PasswordHash    string
	FullName        string
	Type            UserType
	Status          UserStatus
	TwoFactorStatus TwoFactorStatus
}

func (input CreateUserInput) Normalize() CreateUserInput {
	output := input
	output.Email = strings.ToLower(strings.TrimSpace(output.Email))
	output.FullName = strings.TrimSpace(output.FullName)
	if output.Status == "" {
		output.Status = UserStatusPendingVerification
	}
	if output.TwoFactorStatus == "" {
		output.TwoFactorStatus = TwoFactorStatusDisabled
	}
	return output
}

func (input CreateUserInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.Email == "" {
		return ErrEmailMissing
	}
	if input.PasswordHash == "" {
		return ErrPasswordHashMissing
	}
	if !input.Type.Valid() {
		return ErrUserTypeInvalid
	}
	if input.Status != "" && !input.Status.Valid() {
		return ErrUserStatusInvalid
	}
	if input.TwoFactorStatus != "" && !input.TwoFactorStatus.Valid() {
		return ErrTwoFactorStatusInvalid
	}
	return nil
}

type UserStore interface {
	CreateUser(ctx context.Context, input CreateUserInput) (User, error)
	GetUserByID(ctx context.Context, tenantID tenant.ID, userID UserID) (User, error)
	FindUserByEmail(ctx context.Context, tenantID tenant.ID, email string) (User, error)
	ListUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error)
}
