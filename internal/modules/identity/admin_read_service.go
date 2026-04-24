package identity

import (
	"context"
	"errors"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var ErrAdminReadStoreMissing = errors.New("admin account read store missing")

type TenantReader interface {
	ListTenants(ctx context.Context, filter tenant.ListTenantsFilter) ([]tenant.TenantSummary, error)
}

type UserReader interface {
	ListUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error)
}

type AdminReadService struct {
	tenants TenantReader
	users   UserReader
}

func NewAdminReadService(tenants TenantReader, users UserReader) *AdminReadService {
	return &AdminReadService{tenants: tenants, users: users}
}

func (service *AdminReadService) ListAdminTenants(ctx context.Context, filter tenant.ListTenantsFilter) ([]tenant.TenantSummary, error) {
	if service == nil || service.tenants == nil {
		return nil, ErrAdminReadStoreMissing
	}
	return service.tenants.ListTenants(ctx, filter)
}

func (service *AdminReadService) ListAdminUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error) {
	if service == nil || service.users == nil {
		return nil, ErrAdminReadStoreMissing
	}
	return service.users.ListUsers(ctx, filter)
}
