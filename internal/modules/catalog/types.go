package catalog

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrProductIDMissing            = errors.New("catalog product id missing")
	ErrProductTypeInvalid          = errors.New("catalog product type invalid")
	ErrProductNameMissing          = errors.New("catalog product name missing")
	ErrProductStatusInvalid        = errors.New("catalog product status invalid")
	ErrPlanIDMissing               = errors.New("catalog plan id missing")
	ErrPlanCodeMissing             = errors.New("catalog plan code missing")
	ErrPlanNameMissing             = errors.New("catalog plan name missing")
	ErrPlanStatusInvalid           = errors.New("catalog plan status invalid")
	ErrBillingCycleInvalid         = errors.New("catalog billing cycle invalid")
	ErrBillingCycleValue           = errors.New("catalog billing cycle value invalid")
	ErrCurrencyMissing             = errors.New("catalog currency missing")
	ErrCurrencyInvalid             = errors.New("catalog currency invalid")
	ErrMoneyAmountInvalid          = errors.New("catalog money amount invalid")
	ErrVersionInvalid              = errors.New("catalog version invalid")
	ErrSourceIDMissing             = errors.New("catalog provider source id missing")
	ErrSourceTypeInvalid           = errors.New("catalog provider source type invalid")
	ErrSourceNameMissing           = errors.New("catalog provider source name missing")
	ErrSourceStatusInvalid         = errors.New("catalog provider source status invalid")
	ErrInventoryModeInvalid        = errors.New("catalog inventory mode invalid")
	ErrRiskLevelInvalid            = errors.New("catalog risk level invalid")
	ErrPlanSourceIDMissing         = errors.New("catalog plan source id missing")
	ErrPlanSourceStatus            = errors.New("catalog plan source status invalid")
	ErrPlanSourcePriority          = errors.New("catalog plan source priority invalid")
	ErrTenantProductIDMissing      = errors.New("catalog tenant product id missing")
	ErrTenantProductStatus         = errors.New("catalog tenant product status invalid")
	ErrTenantPlanIDMissing         = errors.New("catalog tenant plan id missing")
	ErrTenantPlanStatus            = errors.New("catalog tenant plan status invalid")
	ErrTenantPlanVisibility        = errors.New("catalog tenant plan visibility invalid")
	ErrCreatedByMissing            = errors.New("catalog created by missing")
	ErrCatalogStoreExecutorMissing = errors.New("catalog store executor missing")
	ErrProductNotFound             = errors.New("catalog product not found")
	ErrPlanNotFound                = errors.New("catalog plan not found")
	ErrProviderSourceNotFound      = errors.New("catalog provider source not found")
	ErrPlanSourceNotFound          = errors.New("catalog plan source not found")
	ErrTenantProductNotFound       = errors.New("catalog tenant product not found")
	ErrTenantPlanNotFound          = errors.New("catalog tenant plan not found")
)

type ProductID string
type PlanID string
type ProviderSourceID string
type PlanSourceID string
type TenantProductID string
type TenantPlanID string
type UserID string

func (id ProductID) Empty() bool        { return strings.TrimSpace(string(id)) == "" }
func (id PlanID) Empty() bool           { return strings.TrimSpace(string(id)) == "" }
func (id ProviderSourceID) Empty() bool { return strings.TrimSpace(string(id)) == "" }
func (id PlanSourceID) Empty() bool     { return strings.TrimSpace(string(id)) == "" }
func (id TenantProductID) Empty() bool  { return strings.TrimSpace(string(id)) == "" }
func (id TenantPlanID) Empty() bool     { return strings.TrimSpace(string(id)) == "" }

func defaultJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return append(json.RawMessage(nil), value...)
}

func trim(value string) string {
	return strings.TrimSpace(value)
}

func upperTrim(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func validateCurrency(value string) error {
	if value == "" {
		return ErrCurrencyMissing
	}
	if len(value) != 3 {
		return ErrCurrencyInvalid
	}
	for _, letter := range value {
		if letter < 'A' || letter > 'Z' {
			return ErrCurrencyInvalid
		}
	}
	return nil
}

func validateMinorAmount(value int64) error {
	if value < 0 {
		return ErrMoneyAmountInvalid
	}
	return nil
}
