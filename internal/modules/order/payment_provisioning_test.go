package order

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestPaymentFinalizationServiceQueuesProvisioningAfterPaidOrder(t *testing.T) {
	finalizer := &fakePaymentFinalizer{record: paidProvisioningOrder()}
	queue := &fakePaidOrderProvisioningQueue{}
	service := NewPaymentFinalizationService(finalizer, queue)

	record, err := service.FinalizePayment(context.Background(), FinalizePaymentInput{
		ID:          "order-1",
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: "buyer-1",
	})
	if err != nil {
		t.Fatalf("expected finalization result: %v", err)
	}
	if record.ID != OrderID("order-1") {
		t.Fatalf("unexpected order result: %+v", record)
	}
	if queue.calls != 1 ||
		queue.input.OrderID != OrderID("order-1") ||
		queue.input.TenantID != tenant.ID("tenant-1") {
		t.Fatalf("unexpected queue input: calls=%d input=%+v", queue.calls, queue.input)
	}
}

func TestPaymentFinalizationServiceReturnsQueueError(t *testing.T) {
	service := NewPaymentFinalizationService(
		&fakePaymentFinalizer{record: paidProvisioningOrder()},
		&fakePaidOrderProvisioningQueue{err: ErrProvisioningSourceNotFound},
	)

	_, err := service.FinalizePayment(context.Background(), FinalizePaymentInput{
		ID:          "order-1",
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: "buyer-1",
	})
	if !errors.Is(err, ErrProvisioningSourceNotFound) {
		t.Fatalf("expected source error, got %v", err)
	}
}

type fakePaymentFinalizer struct {
	record Order
	input  FinalizePaymentInput
	err    error
}

func (finalizer *fakePaymentFinalizer) FinalizePayment(_ context.Context, input FinalizePaymentInput) (Order, error) {
	finalizer.input = input.Normalize()
	if finalizer.err != nil {
		return Order{}, finalizer.err
	}
	return finalizer.record, nil
}

type fakePaidOrderProvisioningQueue struct {
	input QueuePaidOrderProvisioningInput
	err   error
	calls int
}

func (queue *fakePaidOrderProvisioningQueue) QueuePaidOrderProvisioning(_ context.Context, input QueuePaidOrderProvisioningInput) (jobs.Job, error) {
	queue.calls++
	queue.input = input.Normalize()
	if queue.err != nil {
		return jobs.Job{}, queue.err
	}
	return jobs.Job{ID: jobs.ID("job-1")}, nil
}
