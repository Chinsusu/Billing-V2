package wallet

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateWalletInputNormalizeValidate(t *testing.T) {
	input := CreateWalletInput{
		TenantID:              tenant.ID("tenant-1"),
		OwnerType:             OwnerTypeUser,
		OwnerID:               OwnerID(" user-1 "),
		Currency:              " usd ",
		AvailableBalanceMinor: 1000,
		LockedBalanceMinor:    200,
		Metadata:              json.RawMessage(`{"source":"seed"}`),
	}.Normalize()

	if input.OwnerID != OwnerID("user-1") {
		t.Fatalf("expected trimmed owner id, got %q", input.OwnerID)
	}
	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.Status != StatusActive {
		t.Fatalf("expected active status, got %q", input.Status)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid wallet input: %v", err)
	}

	input.AvailableBalanceMinor = -1
	if err := input.Validate(); !errors.Is(err, ErrBalanceInvalid) {
		t.Fatalf("expected balance error, got %v", err)
	}
}

func TestCreateWalletInputRequiresOwner(t *testing.T) {
	err := CreateWalletInput{
		TenantID: tenant.ID("tenant-1"),
		Currency: "USD",
	}.Normalize().Validate()
	if !errors.Is(err, ErrOwnerTypeInvalid) {
		t.Fatalf("expected owner type error, got %v", err)
	}

	err = CreateWalletInput{
		TenantID:  tenant.ID("tenant-1"),
		OwnerType: OwnerTypeTenant,
		Currency:  "USD",
	}.Normalize().Validate()
	if !errors.Is(err, ErrOwnerIDMissing) {
		t.Fatalf("expected owner id error, got %v", err)
	}
}

func TestWalletEnums(t *testing.T) {
	for _, ownerType := range []OwnerType{OwnerTypeTenant, OwnerTypeUser, OwnerTypeResellerSettlement, OwnerTypePlatform} {
		if !ownerType.Valid() {
			t.Fatalf("expected valid owner type %q", ownerType)
		}
	}
	if OwnerType("bad").Valid() {
		t.Fatal("unexpected valid owner type")
	}
	if !StatusActive.Valid() || !StatusFrozen.Valid() || !StatusClosed.Valid() {
		t.Fatal("expected core statuses to be valid")
	}
	if Status("bad").Valid() {
		t.Fatal("unexpected valid status")
	}
}
