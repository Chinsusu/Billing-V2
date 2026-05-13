package order

import "context"

func (service *Service) TransitionServiceLifecycle(ctx context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	if err := service.ready(); err != nil {
		return ServiceInstance{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ServiceInstance{}, err
	}
	record, err := service.store.TransitionServiceLifecycle(ctx, input)
	if err != nil {
		return ServiceInstance{}, err
	}
	if err := service.appendServiceLifecycleAudit(ctx, input, record); err != nil {
		return ServiceInstance{}, err
	}
	return record, nil
}

// RenewServiceTerm records lifecycle-side renewal effects after billing has succeeded.
func (service *Service) RenewServiceTerm(ctx context.Context, input RenewServiceTermInput) (ServiceInstance, error) {
	if err := service.ready(); err != nil {
		return ServiceInstance{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ServiceInstance{}, err
	}
	before, err := service.store.GetServiceInstance(ctx, ServiceInstanceLookup{
		ID:          input.ID,
		TenantID:    input.TenantID,
		BuyerUserID: input.BuyerUserID,
	})
	if err != nil {
		return ServiceInstance{}, err
	}
	if before.Status != input.FromStatus {
		return ServiceInstance{}, ErrServiceStatusConflict
	}
	newTermEnd, err := CalculateRenewedTermEnd(before, input.Cycle)
	if err != nil {
		return ServiceInstance{}, err
	}
	return service.TransitionServiceLifecycle(ctx, TransitionServiceLifecycleInput{
		ID:            input.ID,
		TenantID:      input.TenantID,
		BuyerUserID:   input.BuyerUserID,
		ActorID:       input.ActorID,
		ActorType:     input.ActorType,
		Action:        ServiceLifecycleActionRenew,
		FromStatus:    input.FromStatus,
		ToStatus:      ServiceStatusActive,
		BillingStatus: BillingStatusPaid,
		Reason:        input.Reason,
		TermEnd:       newTermEnd,
	})
}
