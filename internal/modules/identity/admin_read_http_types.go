package identity

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type adminTenantResponse struct {
	ID              tenant.ID     `json:"id"`
	DisplayID       int64         `json:"display_id"`
	ParentID        tenant.ID     `json:"parent_tenant_id,omitempty"`
	Type            tenant.Type   `json:"tenant_type"`
	Name            string        `json:"name"`
	Slug            string        `json:"slug"`
	Status          tenant.Status `json:"status"`
	DefaultCurrency string        `json:"default_currency"`
	Timezone        string        `json:"timezone"`
	OwnerUserID     string        `json:"owner_user_id,omitempty"`
	PrimaryDomain   string        `json:"primary_domain,omitempty"`
	UserCount       int64         `json:"user_count"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

func newAdminTenantResponse(record tenant.TenantSummary) adminTenantResponse {
	value := record.Tenant
	return adminTenantResponse{
		ID:              value.ID,
		DisplayID:       value.DisplayID,
		ParentID:        value.ParentID,
		Type:            value.Type,
		Name:            value.Name,
		Slug:            value.Slug,
		Status:          value.Status,
		DefaultCurrency: value.DefaultCurrency,
		Timezone:        value.Timezone,
		OwnerUserID:     value.OwnerUserID,
		PrimaryDomain:   record.PrimaryDomain,
		UserCount:       record.UserCount,
		CreatedAt:       value.CreatedAt,
		UpdatedAt:       value.UpdatedAt,
	}
}

func newAdminTenantResponses(records []tenant.TenantSummary) []adminTenantResponse {
	responses := make([]adminTenantResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, newAdminTenantResponse(record))
	}
	return responses
}

type adminAccountResponse struct {
	ID              UserID          `json:"id"`
	DisplayID       int64           `json:"display_id"`
	TenantID        tenant.ID       `json:"tenant_id"`
	TenantName      string          `json:"tenant_name"`
	TenantSlug      string          `json:"tenant_slug"`
	Email           string          `json:"email"`
	EmailVerifiedAt *time.Time      `json:"email_verified_at,omitempty"`
	FullName        string          `json:"full_name"`
	Type            UserType        `json:"user_type"`
	Status          UserStatus      `json:"status"`
	TwoFactorStatus TwoFactorStatus `json:"two_factor_status"`
	LastLoginAt     *time.Time      `json:"last_login_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func newAdminAccountResponse(record UserSummary) adminAccountResponse {
	value := record.User
	return adminAccountResponse{
		ID:              value.ID,
		DisplayID:       value.DisplayID,
		TenantID:        value.TenantID,
		TenantName:      record.TenantName,
		TenantSlug:      record.TenantSlug,
		Email:           value.Email,
		EmailVerifiedAt: timePtr(value.EmailVerifiedAt),
		FullName:        value.FullName,
		Type:            value.Type,
		Status:          value.Status,
		TwoFactorStatus: value.TwoFactorStatus,
		LastLoginAt:     timePtr(value.LastLoginAt),
		CreatedAt:       value.CreatedAt,
		UpdatedAt:       value.UpdatedAt,
	}
}

func newAdminAccountResponses(records []UserSummary) []adminAccountResponse {
	responses := make([]adminAccountResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, newAdminAccountResponse(record))
	}
	return responses
}

func timePtr(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
