package provider

import (
	"context"
	"time"
)

type OperationName string

const (
	OperationCheckHealth   OperationName = "check_health"
	OperationCheckStock    OperationName = "check_stock"
	OperationProvision     OperationName = "provision"
	OperationGetStatus     OperationName = "get_status"
	OperationSuspend       OperationName = "suspend"
	OperationUnsuspend     OperationName = "unsuspend"
	OperationTerminate     OperationName = "terminate"
	OperationRenew         OperationName = "renew"
	OperationResetPassword OperationName = "reset_password"
	OperationChangeIP      OperationName = "change_ip"
)

type FakeAdapter struct {
	Provider     Type
	Capabilities CapabilityProfile
	Health       HealthResult
	Stock        StockResult
	Results      map[OperationName]OperationResult
	Errors       map[OperationName]error
	Calls        []OperationName
	ObservedAt   time.Time
}

func NewFakeAdapter(providerType Type) *FakeAdapter {
	now := time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC)
	return &FakeAdapter{
		Provider:     providerType,
		Capabilities: DefaultCapabilityProfile(providerType),
		Health: HealthResult{
			HealthStatus: HealthStatusHealthy,
			Result:       SuccessResult(now),
		},
		Stock: StockResult{
			StockStatus:   StockStatusAvailable,
			CapacityCount: 1,
			Result:        SuccessResult(now),
		},
		Results:    make(map[OperationName]OperationResult),
		Errors:     make(map[OperationName]error),
		ObservedAt: now,
	}
}

func (adapter *FakeAdapter) ProviderType() Type {
	return adapter.Provider
}

func (adapter *FakeAdapter) CapabilityProfile() CapabilityProfile {
	return adapter.Capabilities
}

func (adapter *FakeAdapter) SetResult(operation OperationName, result OperationResult) {
	adapter.Results[operation] = result
}

func (adapter *FakeAdapter) SetError(operation OperationName, err error) {
	adapter.Errors[operation] = err
}

func (adapter *FakeAdapter) CheckHealth(ctx context.Context, operation OperationContext, request HealthRequest) (HealthResult, error) {
	adapter.record(OperationCheckHealth)
	if err := adapter.err(OperationCheckHealth); err != nil {
		if adapterErr, ok := err.(AdapterError); ok {
			return HealthResult{HealthStatus: HealthStatusDown, Result: ResultFromError(adapterErr, adapter.ObservedAt)}, err
		}
		return HealthResult{HealthStatus: HealthStatusUnknown, Result: ResultFromError(NewError(ErrorTemporary, "health check failed"), adapter.ObservedAt)}, err
	}
	return adapter.Health, nil
}

func (adapter *FakeAdapter) CheckStock(ctx context.Context, operation OperationContext, request StockRequest) (StockResult, error) {
	adapter.record(OperationCheckStock)
	if err := adapter.err(OperationCheckStock); err != nil {
		if adapterErr, ok := err.(AdapterError); ok {
			return StockResult{StockStatus: StockStatusUnknown, Result: ResultFromError(adapterErr, adapter.ObservedAt)}, err
		}
		return StockResult{StockStatus: StockStatusUnknown, Result: ResultFromError(NewError(ErrorTemporary, "stock check failed"), adapter.ObservedAt)}, err
	}
	return adapter.Stock, nil
}

func (adapter *FakeAdapter) Provision(ctx context.Context, operation OperationContext, request ProvisionRequest) (OperationResult, error) {
	return adapter.operation(OperationProvision)
}

func (adapter *FakeAdapter) GetStatus(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationGetStatus)
}

func (adapter *FakeAdapter) Suspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationSuspend)
}

func (adapter *FakeAdapter) Unsuspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationUnsuspend)
}

func (adapter *FakeAdapter) Terminate(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationTerminate)
}

func (adapter *FakeAdapter) Renew(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationRenew)
}

func (adapter *FakeAdapter) ResetPassword(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationResetPassword)
}

func (adapter *FakeAdapter) ChangeIP(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.operation(OperationChangeIP)
}

func (adapter *FakeAdapter) operation(name OperationName) (OperationResult, error) {
	adapter.record(name)
	if err := adapter.err(name); err != nil {
		if adapterErr, ok := err.(AdapterError); ok {
			return ResultFromError(adapterErr, adapter.ObservedAt), err
		}
		return ResultFromError(NewError(ErrorTemporary, "provider operation failed"), adapter.ObservedAt), err
	}
	if result, ok := adapter.Results[name]; ok {
		return result, nil
	}
	return SuccessResult(adapter.ObservedAt), nil
}

func (adapter *FakeAdapter) record(name OperationName) {
	adapter.Calls = append(adapter.Calls, name)
}

func (adapter *FakeAdapter) err(name OperationName) error {
	if adapter.Errors == nil {
		return nil
	}
	return adapter.Errors[name]
}
