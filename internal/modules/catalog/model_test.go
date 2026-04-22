package catalog

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateProductInputNormalizeValidate(t *testing.T) {
	input := CreateProductInput{
		Type:      ProductTypeVPS,
		Name:      "  VPS Singapore  ",
		CreatedBy: " admin-1 ",
	}.Normalize()

	if input.Name != "VPS Singapore" {
		t.Fatalf("expected trimmed name, got %q", input.Name)
	}
	if input.CreatedBy != "admin-1" {
		t.Fatalf("expected trimmed creator, got %q", input.CreatedBy)
	}
	if input.Status != ProductStatusDraft {
		t.Fatalf("expected draft status, got %q", input.Status)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid product input: %v", err)
	}

	input.Type = ProductType("bad")
	if err := input.Validate(); !errors.Is(err, ErrProductTypeInvalid) {
		t.Fatalf("expected product type error, got %v", err)
	}
}

func TestCreatePlanInputNormalizeValidate(t *testing.T) {
	input := CreatePlanInput{
		ProductID:             ProductID("product-1"),
		Code:                  " vps-2c4g-monthly ",
		Name:                  " VPS 2C4G monthly ",
		BillingCycle:          BillingCycle{Type: BillingCycleMonth30Days, Value: 1},
		BaseCostMinor:         500,
		SuggestedPriceMinor:   900,
		ResellerMinPriceMinor: 700,
		Currency:              " usd ",
	}.Normalize()

	if input.Code != "vps-2c4g-monthly" {
		t.Fatalf("expected trimmed code, got %q", input.Code)
	}
	if input.Currency != "USD" {
		t.Fatalf("expected uppercase currency, got %q", input.Currency)
	}
	if string(input.Specs) != "{}" {
		t.Fatalf("expected default specs JSON, got %s", input.Specs)
	}
	if input.Version != 1 {
		t.Fatalf("expected default version 1, got %d", input.Version)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid plan input: %v", err)
	}

	input.SuggestedPriceMinor = -1
	if err := input.Validate(); !errors.Is(err, ErrMoneyAmountInvalid) {
		t.Fatalf("expected money amount error, got %v", err)
	}

	input.SuggestedPriceMinor = 900
	input.Currency = "USDT"
	if err := input.Validate(); !errors.Is(err, ErrCurrencyInvalid) {
		t.Fatalf("expected currency error, got %v", err)
	}
}

func TestCreateProviderSourceInputDefaults(t *testing.T) {
	input := CreateProviderSourceInput{
		Type:          provider.TypeProxmox,
		Name:          "  Proxmox SG  ",
		InventoryMode: InventoryModeProviderLive,
	}.Normalize()

	if input.Name != "Proxmox SG" {
		t.Fatalf("expected trimmed name, got %q", input.Name)
	}
	if input.Status != ProviderSourceStatusDisabled {
		t.Fatalf("expected disabled status, got %q", input.Status)
	}
	if input.RiskLevel != RiskLevelMedium {
		t.Fatalf("expected medium risk, got %q", input.RiskLevel)
	}
	if !input.CapabilityProfile.SupportsAutoProvision {
		t.Fatal("expected default capability profile for proxmox")
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid provider source input: %v", err)
	}

	input.InventoryMode = InventoryMode("bad")
	if err := input.Validate(); !errors.Is(err, ErrInventoryModeInvalid) {
		t.Fatalf("expected inventory mode error, got %v", err)
	}

	input.InventoryMode = InventoryModeProviderLive
	input.Type = provider.Type("bad")
	if err := input.Validate(); !errors.Is(err, ErrSourceTypeInvalid) {
		t.Fatalf("expected provider source type error, got %v", err)
	}
}

func TestCreateTenantPlanInputNormalizeValidate(t *testing.T) {
	input := CreateTenantPlanInput{
		TenantID:          tenant.ID("tenant-1"),
		TenantProductID:   TenantProductID("tenant-product-1"),
		MasterPlanID:      PlanID("plan-1"),
		SellingPriceMinor: 1200,
		ResellerCostMinor: 800,
		Currency:          "vnd",
		ProductSnapshot:   json.RawMessage(`{"name":"VPS"}`),
	}.Normalize()

	if input.Currency != "VND" {
		t.Fatalf("expected uppercase currency, got %q", input.Currency)
	}
	if input.Visibility != TenantPlanVisibilityHidden {
		t.Fatalf("expected hidden visibility, got %q", input.Visibility)
	}
	if input.Status != TenantPlanStatusDisabled {
		t.Fatalf("expected disabled status, got %q", input.Status)
	}
	if string(input.MarginPolicy) != "{}" {
		t.Fatalf("expected default margin policy, got %s", input.MarginPolicy)
	}
	if string(input.PlanSnapshot) != "{}" {
		t.Fatalf("expected default plan snapshot, got %s", input.PlanSnapshot)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid tenant plan input: %v", err)
	}

	input.TenantID = ""
	if err := input.Validate(); !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant id error, got %v", err)
	}
}
