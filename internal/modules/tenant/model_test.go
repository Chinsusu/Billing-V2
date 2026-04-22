package tenant

import (
	"errors"
	"testing"
)

func TestCreateTenantInputNormalize(t *testing.T) {
	input := CreateTenantInput{
		Type:            TypeReseller,
		Name:            " Demo Reseller ",
		Slug:            " Demo-Reseller ",
		DefaultCurrency: "usd",
	}

	normalized := input.Normalize()

	if normalized.Name != "Demo Reseller" {
		t.Fatalf("expected trimmed name, got %q", normalized.Name)
	}
	if normalized.Slug != "demo-reseller" {
		t.Fatalf("expected lower slug, got %q", normalized.Slug)
	}
	if normalized.DefaultCurrency != "USD" {
		t.Fatalf("expected upper currency, got %q", normalized.DefaultCurrency)
	}
	if normalized.Status != StatusPendingSetup {
		t.Fatalf("expected default pending status, got %q", normalized.Status)
	}
}

func TestCreateTenantInputValidateRejectsInvalidType(t *testing.T) {
	input := CreateTenantInput{Name: "Demo", Slug: "demo", DefaultCurrency: "USD", Type: TypeAdmin}

	if err := input.Validate(); !errors.Is(err, ErrTenantTypeInvalid) {
		t.Fatalf("expected invalid type, got %v", err)
	}
}

func TestCreateDomainInputNormalizeAndValidate(t *testing.T) {
	input := CreateDomainInput{TenantID: "tenant_1", Domain: " Store.Example.COM ", Type: DomainTypeCustomDomain}
	normalized := input.Normalize()

	if normalized.Domain != "store.example.com" {
		t.Fatalf("expected lower domain, got %q", normalized.Domain)
	}
	if normalized.VerificationStatus != DomainVerificationPending {
		t.Fatalf("expected default verification status, got %q", normalized.VerificationStatus)
	}
	if normalized.TLSStatus != TLSStatusPending {
		t.Fatalf("expected default tls status, got %q", normalized.TLSStatus)
	}
	if err := normalized.Validate(); err != nil {
		t.Fatalf("expected valid domain input, got %v", err)
	}
}
