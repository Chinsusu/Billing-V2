package tenant

import "context"

type Store interface {
	Create(ctx context.Context, input CreateTenantInput) (Tenant, error)
	GetByID(ctx context.Context, tenantID ID) (Tenant, error)
	FindBySlug(ctx context.Context, slug string) (Tenant, error)
	ListTenants(ctx context.Context, filter ListTenantsFilter) ([]TenantSummary, error)
	CreateDomain(ctx context.Context, input CreateDomainInput) (Domain, error)
	FindActiveDomain(ctx context.Context, domain string) (Domain, error)
}
