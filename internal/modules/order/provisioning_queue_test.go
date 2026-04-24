package order

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestProvisioningQueueServiceQueuesPaidOrder(t *testing.T) {
	orderStore := &fakeProvisioningOrderStore{order: paidProvisioningOrder()}
	queue := &fakeProvisioningJobQueue{}
	service := NewProvisioningQueueService(orderStore, queue)

	job, err := service.QueueOrderProvisioning(context.Background(), QueueProvisioningInput{
		OrderID:          " order-1 ",
		TenantID:         tenant.ID(" tenant-1 "),
		ProviderSourceID: catalog.ProviderSourceID(" source-1 "),
		ProviderType:     provider.TypeManual,
	})
	if err != nil {
		t.Fatalf("expected provisioning queue: %v", err)
	}
	if job.Type != ProvisioningJobType || job.ReferenceType != ProvisioningReferenceType || job.ReferenceID != jobs.ReferenceID("order-1") {
		t.Fatalf("unexpected queued job: %+v", job)
	}
	if queue.createCalls != 1 {
		t.Fatalf("expected one queue create call, got %d", queue.createCalls)
	}
	if orderStore.lookup.ID != OrderID("order-1") || orderStore.lookup.TenantID != tenant.ID("tenant-1") {
		t.Fatalf("expected tenant-scoped order lookup, got %+v", orderStore.lookup)
	}
	input := queue.inputs[0]
	if input.TenantID != tenant.ID("tenant-1") || input.SourceID != jobs.SourceID("source-1") ||
		input.Priority != defaultProvisionPriority || input.IdempotencyKey == "" {
		t.Fatalf("unexpected queue input: %+v", input)
	}
	var payload ProvisioningQueuePayload
	if err := json.Unmarshal(input.PayloadJSON, &payload); err != nil {
		t.Fatalf("expected payload json: %v", err)
	}
	if payload.OrderID != "order-1" || payload.OrderDisplayID != 30001 ||
		payload.BuyerUserID != identity.UserID("buyer-1") ||
		payload.ProviderSourceID != catalog.ProviderSourceID("source-1") ||
		payload.ProviderType != provider.TypeManual {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestProvisioningQueueServiceQueuesPaidOrderFromResolvedSource(t *testing.T) {
	orderStore := &fakeProvisioningOrderStore{order: paidProvisioningOrder()}
	sourceResolver := &fakeProvisioningSourceResolver{
		source: ProvisioningSource{
			ProviderSourceID: catalog.ProviderSourceID("source-1"),
			ProviderType:     provider.TypeManual,
		},
	}
	queue := &fakeProvisioningJobQueue{}
	service := NewProvisioningQueueServiceWithSourceResolver(orderStore, queue, sourceResolver)

	job, err := service.QueuePaidOrderProvisioning(context.Background(), QueuePaidOrderProvisioningInput{
		OrderID:  " order-1 ",
		TenantID: tenant.ID(" tenant-1 "),
	})
	if err != nil {
		t.Fatalf("expected paid order queue: %v", err)
	}
	if job.ReferenceID != jobs.ReferenceID("order-1") || job.SourceID != jobs.SourceID("source-1") {
		t.Fatalf("unexpected queued job: %+v", job)
	}
	if sourceResolver.input.TenantID != tenant.ID("tenant-1") ||
		sourceResolver.input.TenantPlanID != catalog.TenantPlanID("tenant-plan-1") {
		t.Fatalf("unexpected source lookup: %+v", sourceResolver.input)
	}
	if len(queue.jobsByKey) != 1 {
		t.Fatalf("expected one unique queued job, got %d", len(queue.jobsByKey))
	}
}

func TestProvisioningQueueServiceRejectsUnpaidStatuses(t *testing.T) {
	cases := []struct {
		name          string
		orderStatus   OrderStatus
		billingStatus BillingStatus
	}{
		{name: "pending", orderStatus: OrderStatusPendingPayment, billingStatus: BillingStatusUnpaid},
		{name: "paid_unpaid", orderStatus: OrderStatusPaid, billingStatus: BillingStatusUnpaid},
		{name: "cancelled", orderStatus: OrderStatusCancelled, billingStatus: BillingStatusUnpaid},
		{name: "failed", orderStatus: OrderStatusFailed, billingStatus: BillingStatusPaid},
		{name: "refunded", orderStatus: OrderStatusRefunded, billingStatus: BillingStatusRefunded},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			order := paidProvisioningOrder()
			order.OrderStatus = tc.orderStatus
			order.BillingStatus = tc.billingStatus
			queue := &fakeProvisioningJobQueue{}
			service := NewProvisioningQueueService(&fakeProvisioningOrderStore{order: order}, queue)

			_, err := service.QueueOrderProvisioning(context.Background(), validQueueProvisioningInput())
			if !errors.Is(err, ErrProvisioningQueueNotPaid) {
				t.Fatalf("expected not-paid error, got %v", err)
			}
			if queue.createCalls != 0 {
				t.Fatalf("expected no queue call, got %d", queue.createCalls)
			}
		})
	}
}

func TestProvisioningQueueServiceRejectsUnpaidPaidOrderQueue(t *testing.T) {
	record := paidProvisioningOrder()
	record.OrderStatus = OrderStatusPendingPayment
	record.BillingStatus = BillingStatusUnpaid
	service := NewProvisioningQueueServiceWithSourceResolver(
		&fakeProvisioningOrderStore{order: record},
		&fakeProvisioningJobQueue{},
		&fakeProvisioningSourceResolver{},
	)

	_, err := service.QueuePaidOrderProvisioning(context.Background(), QueuePaidOrderProvisioningInput{
		OrderID:  "order-1",
		TenantID: tenant.ID("tenant-1"),
	})
	if !errors.Is(err, ErrProvisioningQueueNotPaid) {
		t.Fatalf("expected not-paid error, got %v", err)
	}
}

func TestProvisioningQueueServiceIsIdempotentForDuplicateRequests(t *testing.T) {
	queue := &fakeProvisioningJobQueue{}
	service := NewProvisioningQueueService(&fakeProvisioningOrderStore{order: paidProvisioningOrder()}, queue)

	first, err := service.QueueOrderProvisioning(context.Background(), validQueueProvisioningInput())
	if err != nil {
		t.Fatalf("expected first queue: %v", err)
	}
	second, err := service.QueueOrderProvisioning(context.Background(), validQueueProvisioningInput())
	if err != nil {
		t.Fatalf("expected second queue: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected duplicate queue to return same job, got %q and %q", first.ID, second.ID)
	}
	if len(queue.jobsByKey) != 1 {
		t.Fatalf("expected one unique queued job, got %d", len(queue.jobsByKey))
	}
}

func TestProvisioningQueueServiceIsIdempotentForDuplicatePaidOrderQueue(t *testing.T) {
	queue := &fakeProvisioningJobQueue{}
	service := NewProvisioningQueueServiceWithSourceResolver(
		&fakeProvisioningOrderStore{order: paidProvisioningOrder()},
		queue,
		&fakeProvisioningSourceResolver{
			source: ProvisioningSource{ProviderSourceID: "source-1", ProviderType: provider.TypeManual},
		},
	)

	first, err := service.QueuePaidOrderProvisioning(context.Background(), validQueuePaidOrderProvisioningInput())
	if err != nil {
		t.Fatalf("expected first queue: %v", err)
	}
	second, err := service.QueuePaidOrderProvisioning(context.Background(), validQueuePaidOrderProvisioningInput())
	if err != nil {
		t.Fatalf("expected second queue: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected duplicate queue to return same job, got %q and %q", first.ID, second.ID)
	}
	if len(queue.jobsByKey) != 1 {
		t.Fatalf("expected one unique queued job, got %d", len(queue.jobsByKey))
	}
}

func TestProvisioningQueueServiceRequiresProviderSource(t *testing.T) {
	service := NewProvisioningQueueService(&fakeProvisioningOrderStore{order: paidProvisioningOrder()}, &fakeProvisioningJobQueue{})

	_, err := service.QueueOrderProvisioning(context.Background(), QueueProvisioningInput{
		OrderID:  "order-1",
		TenantID: tenant.ID("tenant-1"),
	})
	if !errors.Is(err, ErrProviderSourceIDMissing) {
		t.Fatalf("expected provider source error, got %v", err)
	}
}

func TestProvisioningQueueServiceRequiresSourceResolverForPaidOrderQueue(t *testing.T) {
	service := NewProvisioningQueueService(&fakeProvisioningOrderStore{order: paidProvisioningOrder()}, &fakeProvisioningJobQueue{})

	_, err := service.QueuePaidOrderProvisioning(context.Background(), validQueuePaidOrderProvisioningInput())
	if !errors.Is(err, ErrProvisioningSourceNotFound) {
		t.Fatalf("expected source resolver error, got %v", err)
	}
}

func TestProvisioningQueueServiceRequiresProviderType(t *testing.T) {
	service := NewProvisioningQueueService(&fakeProvisioningOrderStore{order: paidProvisioningOrder()}, &fakeProvisioningJobQueue{})

	_, err := service.QueueOrderProvisioning(context.Background(), QueueProvisioningInput{
		OrderID:          "order-1",
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
	})
	if !errors.Is(err, ErrProviderTypeMissing) {
		t.Fatalf("expected provider type error, got %v", err)
	}
}

func TestResolveOrderProvisioningSourceInputValidate(t *testing.T) {
	input := ResolveOrderProvisioningSourceInput{
		TenantID:     tenant.ID(" tenant-1 "),
		TenantPlanID: catalog.TenantPlanID(" tenant-plan-1 "),
	}.Normalize()
	if input.TenantID != tenant.ID("tenant-1") || input.TenantPlanID != catalog.TenantPlanID("tenant-plan-1") {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid source lookup: %v", err)
	}
}

func validQueueProvisioningInput() QueueProvisioningInput {
	return QueueProvisioningInput{
		OrderID:          "order-1",
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		ProviderType:     provider.TypeManual,
	}
}

func validQueuePaidOrderProvisioningInput() QueuePaidOrderProvisioningInput {
	return QueuePaidOrderProvisioningInput{
		OrderID:  "order-1",
		TenantID: tenant.ID("tenant-1"),
	}
}

func paidProvisioningOrder() Order {
	return Order{
		ID:            "order-1",
		DisplayID:     30001,
		TenantID:      tenant.ID("tenant-1"),
		BuyerUserID:   identity.UserID("buyer-1"),
		TenantPlanID:  catalog.TenantPlanID("tenant-plan-1"),
		Currency:      "USD",
		TotalMinor:    2500,
		OrderStatus:   OrderStatusPaid,
		BillingStatus: BillingStatusPaid,
	}
}

type fakeProvisioningOrderStore struct {
	order  Order
	lookup OrderLookup
	err    error
}

func (store *fakeProvisioningOrderStore) GetOrder(_ context.Context, lookup OrderLookup) (Order, error) {
	store.lookup = lookup
	if store.err != nil {
		return Order{}, store.err
	}
	return store.order, nil
}

type fakeProvisioningJobQueue struct {
	createCalls int
	inputs      []jobs.CreateJobInput
	jobsByKey   map[string]jobs.Job
}

type fakeProvisioningSourceResolver struct {
	source ProvisioningSource
	input  ResolveOrderProvisioningSourceInput
	err    error
}

func (resolver *fakeProvisioningSourceResolver) ResolveOrderProvisioningSource(_ context.Context, input ResolveOrderProvisioningSourceInput) (ProvisioningSource, error) {
	resolver.input = input.Normalize()
	if resolver.err != nil {
		return ProvisioningSource{}, resolver.err
	}
	return resolver.source, nil
}

func (queue *fakeProvisioningJobQueue) CreateJob(_ context.Context, input jobs.CreateJobInput) (jobs.Job, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return jobs.Job{}, err
	}
	queue.createCalls++
	queue.inputs = append(queue.inputs, input)
	if queue.jobsByKey == nil {
		queue.jobsByKey = map[string]jobs.Job{}
	}
	if job, ok := queue.jobsByKey[input.IdempotencyKey]; ok {
		return job, nil
	}
	job := jobs.Job{
		ID:             jobs.ID("job-" + input.IdempotencyKey),
		TenantID:       input.TenantID,
		Type:           input.Type,
		ReferenceType:  input.ReferenceType,
		ReferenceID:    input.ReferenceID,
		SourceID:       input.SourceID,
		PayloadJSON:    input.PayloadJSON,
		Priority:       input.Priority,
		IdempotencyKey: input.IdempotencyKey,
		MaxAttempts:    input.MaxAttempts,
		CorrelationID:  input.CorrelationID,
	}
	queue.jobsByKey[input.IdempotencyKey] = job
	return job, nil
}
