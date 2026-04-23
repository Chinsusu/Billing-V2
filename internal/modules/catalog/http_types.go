package catalog

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type createProductRequest struct {
	Type         ProductType   `json:"product_type"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       ProductStatus `json:"status"`
	DisplayOrder int           `json:"display_order"`
}

func (request createProductRequest) toInput(createdBy UserID) CreateProductInput {
	return CreateProductInput{
		Type:         request.Type,
		Name:         request.Name,
		Description:  request.Description,
		Status:       request.Status,
		DisplayOrder: request.DisplayOrder,
		CreatedBy:    createdBy,
	}
}

type createPlanRequest struct {
	ProductID             ProductID        `json:"product_id"`
	Code                  string           `json:"plan_code"`
	Name                  string           `json:"name"`
	Specs                 json.RawMessage  `json:"specs"`
	BillingCycleType      BillingCycleType `json:"billing_cycle_type"`
	BillingCycleValue     int              `json:"billing_cycle_value"`
	BaseCostMinor         int64            `json:"base_cost_minor"`
	SuggestedPriceMinor   int64            `json:"suggested_price_minor"`
	ResellerMinPriceMinor int64            `json:"reseller_min_price_minor"`
	Currency              string           `json:"currency"`
	Status                PlanStatus       `json:"status"`
	Version               int              `json:"version"`
}

func (request createPlanRequest) toInput() CreatePlanInput {
	return CreatePlanInput{
		ProductID: request.ProductID,
		Code:      request.Code,
		Name:      request.Name,
		Specs:     request.Specs,
		BillingCycle: BillingCycle{
			Type:  request.BillingCycleType,
			Value: request.BillingCycleValue,
		},
		BaseCostMinor:         request.BaseCostMinor,
		SuggestedPriceMinor:   request.SuggestedPriceMinor,
		ResellerMinPriceMinor: request.ResellerMinPriceMinor,
		Currency:              request.Currency,
		Status:                request.Status,
		Version:               request.Version,
	}
}

type createProviderSourceRequest struct {
	Type              provider.Type              `json:"source_type"`
	Name              string                     `json:"name"`
	ProviderAccountID provider.AccountID         `json:"provider_account_id"`
	Location          string                     `json:"location"`
	Status            ProviderSourceStatus       `json:"status"`
	CapabilityProfile provider.CapabilityProfile `json:"capability_profile"`
	InventoryMode     InventoryMode              `json:"inventory_mode"`
	RiskLevel         RiskLevel                  `json:"risk_level"`
}

func (request createProviderSourceRequest) toInput() CreateProviderSourceInput {
	return CreateProviderSourceInput{
		Type:              request.Type,
		Name:              request.Name,
		ProviderAccountID: request.ProviderAccountID,
		Location:          request.Location,
		Status:            request.Status,
		CapabilityProfile: request.CapabilityProfile,
		InventoryMode:     request.InventoryMode,
		RiskLevel:         request.RiskLevel,
	}
}

type createPlanSourceRequest struct {
	PlanID             PlanID           `json:"plan_id"`
	SourceID           ProviderSourceID `json:"source_id"`
	Priority           int              `json:"priority"`
	CostOverrideMinor  int64            `json:"cost_override_minor"`
	CapacityPolicy     json.RawMessage  `json:"capacity_policy"`
	CapabilityOverride json.RawMessage  `json:"capability_override"`
	Status             PlanSourceStatus `json:"status"`
}

func (request createPlanSourceRequest) toInput() CreatePlanSourceInput {
	return CreatePlanSourceInput{
		PlanID:             request.PlanID,
		SourceID:           request.SourceID,
		Priority:           request.Priority,
		CostOverrideMinor:  request.CostOverrideMinor,
		CapacityPolicy:     request.CapacityPolicy,
		CapabilityOverride: request.CapabilityOverride,
		Status:             request.Status,
	}
}

type cloneTenantProductRequest struct {
	MasterProductID     ProductID           `json:"master_product_id"`
	NameOverride        string              `json:"name_override"`
	DescriptionOverride string              `json:"description_override"`
	Status              TenantProductStatus `json:"status"`
	CloneVersion        int                 `json:"clone_version"`
}

func (request cloneTenantProductRequest) toInput(tenantID tenant.ID) CreateTenantProductInput {
	return CreateTenantProductInput{
		TenantID:            tenantID,
		MasterProductID:     request.MasterProductID,
		NameOverride:        request.NameOverride,
		DescriptionOverride: request.DescriptionOverride,
		Status:              request.Status,
		CloneVersion:        request.CloneVersion,
	}
}

type cloneTenantPlanRequest struct {
	TenantProductID    TenantProductID      `json:"tenant_product_id"`
	MasterPlanID       PlanID               `json:"master_plan_id"`
	SellingPriceMinor  int64                `json:"selling_price_minor"`
	ResellerCostMinor  int64                `json:"reseller_cost_minor"`
	Currency           string               `json:"currency"`
	MarginPolicy       json.RawMessage      `json:"margin_policy"`
	Visibility         TenantPlanVisibility `json:"visibility"`
	Status             TenantPlanStatus     `json:"status"`
	CloneVersion       int                  `json:"clone_version"`
	ProductSnapshot    json.RawMessage      `json:"product_snapshot"`
	PlanSnapshot       json.RawMessage      `json:"plan_snapshot"`
	PriceSnapshot      json.RawMessage      `json:"price_snapshot"`
	CapabilitySnapshot json.RawMessage      `json:"capability_snapshot"`
}

func (request cloneTenantPlanRequest) toInput(tenantID tenant.ID) CreateTenantPlanInput {
	return CreateTenantPlanInput{
		TenantID:           tenantID,
		TenantProductID:    request.TenantProductID,
		MasterPlanID:       request.MasterPlanID,
		SellingPriceMinor:  request.SellingPriceMinor,
		ResellerCostMinor:  request.ResellerCostMinor,
		Currency:           request.Currency,
		MarginPolicy:       request.MarginPolicy,
		Visibility:         request.Visibility,
		Status:             request.Status,
		CloneVersion:       request.CloneVersion,
		ProductSnapshot:    request.ProductSnapshot,
		PlanSnapshot:       request.PlanSnapshot,
		PriceSnapshot:      request.PriceSnapshot,
		CapabilitySnapshot: request.CapabilitySnapshot,
	}
}

type productResponse struct {
	ID           ProductID     `json:"id"`
	DisplayID    int64         `json:"display_id"`
	Type         ProductType   `json:"product_type"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       ProductStatus `json:"status"`
	DisplayOrder int           `json:"display_order"`
	CreatedBy    UserID        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func newProductResponse(product Product) productResponse {
	return productResponse{
		ID:           product.ID,
		DisplayID:    product.DisplayID,
		Type:         product.Type,
		Name:         product.Name,
		Description:  product.Description,
		Status:       product.Status,
		DisplayOrder: product.DisplayOrder,
		CreatedBy:    product.CreatedBy,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

type billingCycleResponse struct {
	Type  BillingCycleType `json:"type"`
	Value int              `json:"value"`
}

type planResponse struct {
	ID                    PlanID               `json:"id"`
	DisplayID             int64                `json:"display_id"`
	ProductID             ProductID            `json:"product_id"`
	Code                  string               `json:"plan_code"`
	Name                  string               `json:"name"`
	Specs                 json.RawMessage      `json:"specs"`
	BillingCycle          billingCycleResponse `json:"billing_cycle"`
	BaseCostMinor         int64                `json:"base_cost_minor"`
	SuggestedPriceMinor   int64                `json:"suggested_price_minor"`
	ResellerMinPriceMinor int64                `json:"reseller_min_price_minor"`
	Currency              string               `json:"currency"`
	Status                PlanStatus           `json:"status"`
	Version               int                  `json:"version"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

func newPlanResponse(plan Plan) planResponse {
	return planResponse{
		ID:        plan.ID,
		DisplayID: plan.DisplayID,
		ProductID: plan.ProductID,
		Code:      plan.Code,
		Name:      plan.Name,
		Specs:     plan.Specs,
		BillingCycle: billingCycleResponse{
			Type:  plan.BillingCycle.Type,
			Value: plan.BillingCycle.Value,
		},
		BaseCostMinor:         plan.BaseCostMinor,
		SuggestedPriceMinor:   plan.SuggestedPriceMinor,
		ResellerMinPriceMinor: plan.ResellerMinPriceMinor,
		Currency:              plan.Currency,
		Status:                plan.Status,
		Version:               plan.Version,
		CreatedAt:             plan.CreatedAt,
		UpdatedAt:             plan.UpdatedAt,
	}
}

func newPlanResponses(plans []Plan) []planResponse {
	responses := make([]planResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, newPlanResponse(plan))
	}
	return responses
}

type providerSourceResponse struct {
	ID                ProviderSourceID           `json:"id"`
	DisplayID         int64                      `json:"display_id"`
	Type              provider.Type              `json:"source_type"`
	Name              string                     `json:"name"`
	ProviderAccountID provider.AccountID         `json:"provider_account_id"`
	Location          string                     `json:"location"`
	Status            ProviderSourceStatus       `json:"status"`
	CapabilityProfile provider.CapabilityProfile `json:"capability_profile"`
	InventoryMode     InventoryMode              `json:"inventory_mode"`
	RiskLevel         RiskLevel                  `json:"risk_level"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

func newProviderSourceResponse(source ProviderSource) providerSourceResponse {
	return providerSourceResponse{
		ID:                source.ID,
		DisplayID:         source.DisplayID,
		Type:              source.Type,
		Name:              source.Name,
		ProviderAccountID: source.ProviderAccountID,
		Location:          source.Location,
		Status:            source.Status,
		CapabilityProfile: source.CapabilityProfile,
		InventoryMode:     source.InventoryMode,
		RiskLevel:         source.RiskLevel,
		CreatedAt:         source.CreatedAt,
		UpdatedAt:         source.UpdatedAt,
	}
}

type planSourceResponse struct {
	ID                 PlanSourceID     `json:"id"`
	DisplayID          int64            `json:"display_id"`
	PlanID             PlanID           `json:"plan_id"`
	SourceID           ProviderSourceID `json:"source_id"`
	Priority           int              `json:"priority"`
	CostOverrideMinor  int64            `json:"cost_override_minor"`
	CapacityPolicy     json.RawMessage  `json:"capacity_policy"`
	CapabilityOverride json.RawMessage  `json:"capability_override"`
	Status             PlanSourceStatus `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

func newPlanSourceResponse(source PlanSource) planSourceResponse {
	return planSourceResponse{
		ID:                 source.ID,
		DisplayID:          source.DisplayID,
		PlanID:             source.PlanID,
		SourceID:           source.SourceID,
		Priority:           source.Priority,
		CostOverrideMinor:  source.CostOverrideMinor,
		CapacityPolicy:     source.CapacityPolicy,
		CapabilityOverride: source.CapabilityOverride,
		Status:             source.Status,
		CreatedAt:          source.CreatedAt,
		UpdatedAt:          source.UpdatedAt,
	}
}

type tenantProductResponse struct {
	ID                  TenantProductID     `json:"id"`
	DisplayID           int64               `json:"display_id"`
	TenantID            tenant.ID           `json:"tenant_id"`
	MasterProductID     ProductID           `json:"master_product_id"`
	NameOverride        string              `json:"name_override"`
	DescriptionOverride string              `json:"description_override"`
	Status              TenantProductStatus `json:"status"`
	CloneVersion        int                 `json:"clone_version"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

func newTenantProductResponse(product TenantProduct) tenantProductResponse {
	return tenantProductResponse{
		ID:                  product.ID,
		DisplayID:           product.DisplayID,
		TenantID:            product.TenantID,
		MasterProductID:     product.MasterProductID,
		NameOverride:        product.NameOverride,
		DescriptionOverride: product.DescriptionOverride,
		Status:              product.Status,
		CloneVersion:        product.CloneVersion,
		CreatedAt:           product.CreatedAt,
		UpdatedAt:           product.UpdatedAt,
	}
}

func newTenantProductResponses(products []TenantProduct) []tenantProductResponse {
	responses := make([]tenantProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, newTenantProductResponse(product))
	}
	return responses
}

type tenantPlanResponse struct {
	ID                 TenantPlanID         `json:"id"`
	DisplayID          int64                `json:"display_id"`
	TenantID           tenant.ID            `json:"tenant_id"`
	TenantProductID    TenantProductID      `json:"tenant_product_id"`
	MasterPlanID       PlanID               `json:"master_plan_id"`
	SellingPriceMinor  int64                `json:"selling_price_minor"`
	ResellerCostMinor  int64                `json:"reseller_cost_minor"`
	Currency           string               `json:"currency"`
	MarginPolicy       json.RawMessage      `json:"margin_policy"`
	Visibility         TenantPlanVisibility `json:"visibility"`
	Status             TenantPlanStatus     `json:"status"`
	CloneVersion       int                  `json:"clone_version"`
	ProductSnapshot    json.RawMessage      `json:"product_snapshot"`
	PlanSnapshot       json.RawMessage      `json:"plan_snapshot"`
	PriceSnapshot      json.RawMessage      `json:"price_snapshot"`
	CapabilitySnapshot json.RawMessage      `json:"capability_snapshot"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

func newTenantPlanResponse(plan TenantPlan) tenantPlanResponse {
	return tenantPlanResponse{
		ID:                 plan.ID,
		DisplayID:          plan.DisplayID,
		TenantID:           plan.TenantID,
		TenantProductID:    plan.TenantProductID,
		MasterPlanID:       plan.MasterPlanID,
		SellingPriceMinor:  plan.SellingPriceMinor,
		ResellerCostMinor:  plan.ResellerCostMinor,
		Currency:           plan.Currency,
		MarginPolicy:       plan.MarginPolicy,
		Visibility:         plan.Visibility,
		Status:             plan.Status,
		CloneVersion:       plan.CloneVersion,
		ProductSnapshot:    plan.ProductSnapshot,
		PlanSnapshot:       plan.PlanSnapshot,
		PriceSnapshot:      plan.PriceSnapshot,
		CapabilitySnapshot: plan.CapabilitySnapshot,
		CreatedAt:          plan.CreatedAt,
		UpdatedAt:          plan.UpdatedAt,
	}
}

func newTenantPlanResponses(plans []TenantPlan) []tenantPlanResponse {
	responses := make([]tenantPlanResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, newTenantPlanResponse(plan))
	}
	return responses
}

type tenantPlanPublicResponse struct {
	ID                TenantPlanID         `json:"id"`
	DisplayID         int64                `json:"display_id"`
	TenantProductID   TenantProductID      `json:"tenant_product_id"`
	MasterPlanID      PlanID               `json:"master_plan_id"`
	SellingPriceMinor int64                `json:"selling_price_minor"`
	Currency          string               `json:"currency"`
	Visibility        TenantPlanVisibility `json:"visibility"`
	Status            TenantPlanStatus     `json:"status"`
	CloneVersion      int                  `json:"clone_version"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

func newTenantPlanPublicResponses(plans []TenantPlan) []tenantPlanPublicResponse {
	responses := make([]tenantPlanPublicResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, tenantPlanPublicResponse{
			ID:                plan.ID,
			DisplayID:         plan.DisplayID,
			TenantProductID:   plan.TenantProductID,
			MasterPlanID:      plan.MasterPlanID,
			SellingPriceMinor: plan.SellingPriceMinor,
			Currency:          plan.Currency,
			Visibility:        plan.Visibility,
			Status:            plan.Status,
			CloneVersion:      plan.CloneVersion,
			CreatedAt:         plan.CreatedAt,
			UpdatedAt:         plan.UpdatedAt,
		})
	}
	return responses
}

type tenantCatalogResponse struct {
	Products []tenantProductResponse `json:"products"`
	Plans    []tenantPlanResponse    `json:"plans"`
}

func newTenantCatalogResponse(catalog TenantCatalog) tenantCatalogResponse {
	return tenantCatalogResponse{
		Products: newTenantProductResponses(catalog.Products),
		Plans:    newTenantPlanResponses(catalog.Plans),
	}
}

type tenantCatalogPublicResponse struct {
	Products []tenantProductResponse    `json:"products"`
	Plans    []tenantPlanPublicResponse `json:"plans"`
}

func newTenantCatalogPublicResponse(catalog TenantCatalog) tenantCatalogPublicResponse {
	return tenantCatalogPublicResponse{
		Products: newTenantProductResponses(catalog.Products),
		Plans:    newTenantPlanPublicResponses(catalog.Plans),
	}
}
