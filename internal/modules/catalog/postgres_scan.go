package catalog

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type catalogScanner interface {
	Scan(dest ...interface{}) error
}

func scanProduct(row catalogScanner) (Product, error) {
	var record Product
	var id, productType, status, createdBy string
	var description sql.NullString
	if err := row.Scan(&id, &record.DisplayID, &productType, &record.Name, &description, &status, &record.DisplayOrder, &createdBy, &record.CreatedAt, &record.UpdatedAt); err != nil {
		return Product{}, mapNotFound(err, ErrProductNotFound, "scan catalog product")
	}
	record.ID = ProductID(id)
	record.Type = ProductType(productType)
	record.Description = description.String
	record.Status = ProductStatus(status)
	record.CreatedBy = UserID(createdBy)
	return record, nil
}

func scanPlan(row catalogScanner) (Plan, error) {
	var record Plan
	var id, productID, cycleType, status string
	var specs []byte
	if err := row.Scan(
		&id, &record.DisplayID, &productID, &record.Code, &record.Name, &specs, &cycleType, &record.BillingCycle.Value,
		&record.BaseCostMinor, &record.SuggestedPriceMinor, &record.ResellerMinPriceMinor, &record.Currency,
		&status, &record.Version, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return Plan{}, mapNotFound(err, ErrPlanNotFound, "scan catalog plan")
	}
	record.ID = PlanID(id)
	record.ProductID = ProductID(productID)
	record.Specs = append(record.Specs, specs...)
	record.BillingCycle.Type = BillingCycleType(cycleType)
	record.Status = PlanStatus(status)
	return record, nil
}

func scanProviderSource(row catalogScanner) (ProviderSource, error) {
	var record ProviderSource
	var id, sourceType, status, inventoryMode, riskLevel string
	var providerAccountID, location sql.NullString
	var capabilityProfile []byte
	if err := row.Scan(
		&id, &record.DisplayID, &sourceType, &record.Name, &providerAccountID, &location, &status,
		&capabilityProfile, &inventoryMode, &riskLevel, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return ProviderSource{}, mapNotFound(err, ErrProviderSourceNotFound, "scan catalog provider source")
	}
	if err := json.Unmarshal(capabilityProfile, &record.CapabilityProfile); err != nil {
		return ProviderSource{}, fmt.Errorf("decode capability profile: %w", err)
	}
	record.ID = ProviderSourceID(id)
	record.Type = provider.Type(sourceType)
	record.ProviderAccountID = provider.AccountID(providerAccountID.String)
	record.Location = location.String
	record.Status = ProviderSourceStatus(status)
	record.InventoryMode = InventoryMode(inventoryMode)
	record.RiskLevel = RiskLevel(riskLevel)
	return record, nil
}

func scanPlanSource(row catalogScanner) (PlanSource, error) {
	var record PlanSource
	var id, planID, sourceID, status string
	var capacityPolicy, capabilityOverride []byte
	if err := row.Scan(
		&id, &record.DisplayID, &planID, &sourceID, &record.Priority, &record.CostOverrideMinor,
		&capacityPolicy, &capabilityOverride, &status, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return PlanSource{}, mapNotFound(err, ErrPlanSourceNotFound, "scan catalog plan source")
	}
	record.ID = PlanSourceID(id)
	record.PlanID = PlanID(planID)
	record.SourceID = ProviderSourceID(sourceID)
	record.CapacityPolicy = append(record.CapacityPolicy, capacityPolicy...)
	record.CapabilityOverride = append(record.CapabilityOverride, capabilityOverride...)
	record.Status = PlanSourceStatus(status)
	return record, nil
}

func scanTenantProduct(row catalogScanner) (TenantProduct, error) {
	var record TenantProduct
	var id, tenantID, productID, status string
	var nameOverride, descriptionOverride sql.NullString
	if err := row.Scan(
		&id, &record.DisplayID, &tenantID, &productID, &nameOverride, &descriptionOverride,
		&status, &record.CloneVersion, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return TenantProduct{}, mapNotFound(err, ErrTenantProductNotFound, "scan catalog tenant product")
	}
	record.ID = TenantProductID(id)
	record.TenantID = tenant.ID(tenantID)
	record.MasterProductID = ProductID(productID)
	record.NameOverride = nameOverride.String
	record.DescriptionOverride = descriptionOverride.String
	record.Status = TenantProductStatus(status)
	return record, nil
}

func scanTenantPlan(row catalogScanner) (TenantPlan, error) {
	var record TenantPlan
	var id, tenantID, tenantProductID, masterPlanID, visibility, status string
	var marginPolicy, productSnapshot, planSnapshot, priceSnapshot, capabilitySnapshot []byte
	if err := row.Scan(
		&id, &record.DisplayID, &tenantID, &tenantProductID, &masterPlanID, &record.SellingPriceMinor,
		&record.ResellerCostMinor, &record.Currency, &marginPolicy, &visibility, &status, &record.CloneVersion,
		&productSnapshot, &planSnapshot, &priceSnapshot, &capabilitySnapshot, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return TenantPlan{}, mapNotFound(err, ErrTenantPlanNotFound, "scan catalog tenant plan")
	}
	record.ID = TenantPlanID(id)
	record.TenantID = tenant.ID(tenantID)
	record.TenantProductID = TenantProductID(tenantProductID)
	record.MasterPlanID = PlanID(masterPlanID)
	record.MarginPolicy = append(record.MarginPolicy, marginPolicy...)
	record.Visibility = TenantPlanVisibility(visibility)
	record.Status = TenantPlanStatus(status)
	record.ProductSnapshot = append(record.ProductSnapshot, productSnapshot...)
	record.PlanSnapshot = append(record.PlanSnapshot, planSnapshot...)
	record.PriceSnapshot = append(record.PriceSnapshot, priceSnapshot...)
	record.CapabilitySnapshot = append(record.CapabilitySnapshot, capabilitySnapshot...)
	return record, nil
}

func mapNotFound(err error, notFound error, label string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFound
	}
	return fmt.Errorf("%s: %w", label, err)
}
