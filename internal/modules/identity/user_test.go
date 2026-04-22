package identity

import (
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateUserInputNormalize(t *testing.T) {
	input := CreateUserInput{
		TenantID:     "tenant_a",
		Email:        " Client@Example.COM ",
		PasswordHash: "hash",
		FullName:     " Client User ",
		Type:         UserTypeClient,
	}

	normalized := input.Normalize()

	if normalized.Email != "client@example.com" {
		t.Fatalf("expected normalized email, got %q", normalized.Email)
	}
	if normalized.FullName != "Client User" {
		t.Fatalf("expected trimmed full name, got %q", normalized.FullName)
	}
	if normalized.Status != UserStatusPendingVerification {
		t.Fatalf("expected default status, got %q", normalized.Status)
	}
	if normalized.TwoFactorStatus != TwoFactorStatusDisabled {
		t.Fatalf("expected default two factor status, got %q", normalized.TwoFactorStatus)
	}
}

func TestCreateUserInputValidateRequiresTenant(t *testing.T) {
	input := CreateUserInput{Email: "client@example.com", PasswordHash: "hash", Type: UserTypeClient}

	if err := input.Validate(); !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected missing tenant, got %v", err)
	}
}

func TestCreateUserInputValidateRejectsInvalidType(t *testing.T) {
	input := CreateUserInput{TenantID: "tenant_a", Email: "client@example.com", PasswordHash: "hash", Type: "owner"}

	if err := input.Validate(); !errors.Is(err, ErrUserTypeInvalid) {
		t.Fatalf("expected invalid user type, got %v", err)
	}
}
