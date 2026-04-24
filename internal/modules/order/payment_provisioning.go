package order

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

type PaymentFinalizer interface {
	FinalizePayment(ctx context.Context, input FinalizePaymentInput) (Order, error)
}

type PaidOrderProvisioningQueue interface {
	QueuePaidOrderProvisioning(ctx context.Context, input QueuePaidOrderProvisioningInput) (jobs.Job, error)
}

type PaymentFinalizationService struct {
	finalizer    PaymentFinalizer
	provisioning PaidOrderProvisioningQueue
}

func NewPaymentFinalizationService(finalizer PaymentFinalizer, provisioning PaidOrderProvisioningQueue) *PaymentFinalizationService {
	return &PaymentFinalizationService{finalizer: finalizer, provisioning: provisioning}
}

func (service *PaymentFinalizationService) FinalizePayment(ctx context.Context, input FinalizePaymentInput) (Order, error) {
	if service == nil || service.finalizer == nil {
		return Order{}, ErrServiceStoreMissing
	}
	record, err := service.finalizer.FinalizePayment(ctx, input)
	if err != nil {
		return Order{}, err
	}
	if service.provisioning == nil {
		return record, nil
	}
	_, err = service.provisioning.QueuePaidOrderProvisioning(ctx, QueuePaidOrderProvisioningInput{
		OrderID:  record.ID,
		TenantID: record.TenantID,
	})
	if err != nil {
		return Order{}, err
	}
	return record, nil
}
