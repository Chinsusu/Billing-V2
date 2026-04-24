package order

import (
	"context"
	"encoding/json"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	ProvisioningJobType       jobs.Type          = "provider.provision"
	ProvisioningReferenceType jobs.ReferenceType = OrderAggregateType
	defaultProvisionPriority  int                = 50
)

type ProvisioningOrderStore interface {
	GetOrder(ctx context.Context, lookup OrderLookup) (Order, error)
}

type ProvisioningSourceResolver interface {
	ResolveOrderProvisioningSource(ctx context.Context, input ResolveOrderProvisioningSourceInput) (ProvisioningSource, error)
}

type ProvisioningQueueService struct {
	orders  ProvisioningOrderStore
	queue   jobs.QueueStore
	sources ProvisioningSourceResolver
}

type QueueProvisioningInput struct {
	OrderID          OrderID
	TenantID         tenant.ID
	ProviderSourceID catalog.ProviderSourceID
	ProviderType     provider.Type
}

type QueuePaidOrderProvisioningInput struct {
	OrderID  OrderID
	TenantID tenant.ID
}

type ResolveOrderProvisioningSourceInput struct {
	TenantID     tenant.ID
	TenantPlanID catalog.TenantPlanID
}

type ProvisioningSource struct {
	ProviderSourceID catalog.ProviderSourceID
	ProviderType     provider.Type
}

type ProvisioningQueuePayload struct {
	OrderID          OrderID                  `json:"order_id"`
	OrderDisplayID   int64                    `json:"order_display_id"`
	TenantID         tenant.ID                `json:"tenant_id"`
	BuyerUserID      identity.UserID          `json:"buyer_user_id"`
	TenantPlanID     catalog.TenantPlanID     `json:"tenant_plan_id"`
	ProviderSourceID catalog.ProviderSourceID `json:"provider_source_id"`
	ProviderType     provider.Type            `json:"provider_type"`
	Currency         string                   `json:"currency"`
	TotalMinor       int64                    `json:"total_minor"`
}

func NewProvisioningQueueService(orders ProvisioningOrderStore, queue jobs.QueueStore) *ProvisioningQueueService {
	return &ProvisioningQueueService{orders: orders, queue: queue}
}

func NewProvisioningQueueServiceWithSourceResolver(
	orders ProvisioningOrderStore,
	queue jobs.QueueStore,
	sources ProvisioningSourceResolver,
) *ProvisioningQueueService {
	return &ProvisioningQueueService{orders: orders, queue: queue, sources: sources}
}

func (service *ProvisioningQueueService) QueueOrderProvisioning(ctx context.Context, input QueueProvisioningInput) (jobs.Job, error) {
	if err := service.ready(); err != nil {
		return jobs.Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return jobs.Job{}, err
	}
	order, err := service.orders.GetOrder(ctx, OrderLookup{ID: input.OrderID, TenantID: input.TenantID})
	if err != nil {
		return jobs.Job{}, err
	}
	return service.queueProvisioningJob(ctx, order, input.ProviderSourceID, input.ProviderType)
}

func (service *ProvisioningQueueService) QueuePaidOrderProvisioning(ctx context.Context, input QueuePaidOrderProvisioningInput) (jobs.Job, error) {
	if err := service.ready(); err != nil {
		return jobs.Job{}, err
	}
	if service.sources == nil {
		return jobs.Job{}, ErrProvisioningSourceNotFound
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return jobs.Job{}, err
	}
	order, err := service.orders.GetOrder(ctx, OrderLookup{ID: input.OrderID, TenantID: input.TenantID})
	if err != nil {
		return jobs.Job{}, err
	}
	if err := ensureProvisionableOrder(order); err != nil {
		return jobs.Job{}, err
	}
	source, err := service.sources.ResolveOrderProvisioningSource(ctx, ResolveOrderProvisioningSourceInput{
		TenantID:     order.TenantID,
		TenantPlanID: order.TenantPlanID,
	})
	if err != nil {
		return jobs.Job{}, err
	}
	return service.queueProvisioningJob(ctx, order, source.ProviderSourceID, source.ProviderType)
}

func (service *ProvisioningQueueService) queueProvisioningJob(
	ctx context.Context,
	order Order,
	providerSourceID catalog.ProviderSourceID,
	providerType provider.Type,
) (jobs.Job, error) {
	if err := ensureProvisionableOrder(order); err != nil {
		return jobs.Job{}, err
	}
	if providerSourceID.Empty() {
		return jobs.Job{}, ErrProviderSourceIDMissing
	}
	if providerType == "" {
		return jobs.Job{}, ErrProviderTypeMissing
	}
	payload, err := provisioningQueuePayloadJSON(order, providerSourceID, providerType)
	if err != nil {
		return jobs.Job{}, err
	}
	return service.queue.CreateJob(ctx, jobs.CreateJobInput{
		TenantID:       order.TenantID,
		Type:           ProvisioningJobType,
		ReferenceType:  ProvisioningReferenceType,
		ReferenceID:    jobs.ReferenceID(order.ID),
		SourceID:       jobs.SourceID(providerSourceID),
		PayloadJSON:    payload,
		Priority:       defaultProvisionPriority,
		IdempotencyKey: provisioningQueueKey(order.TenantID, order.ID, providerSourceID),
		MaxAttempts:    5,
		CorrelationID:  jobs.CorrelationID(order.ID),
	})
}

func (service *ProvisioningQueueService) ready() error {
	if service == nil || service.orders == nil {
		return ErrServiceStoreMissing
	}
	if service.queue == nil {
		return ErrProvisioningQueueMissing
	}
	return nil
}

func (input QueueProvisioningInput) Normalize() QueueProvisioningInput {
	return QueueProvisioningInput{
		OrderID:          OrderID(trim(string(input.OrderID))),
		TenantID:         tenant.ID(trim(string(input.TenantID))),
		ProviderSourceID: catalog.ProviderSourceID(trim(string(input.ProviderSourceID))),
		ProviderType:     provider.Type(trim(string(input.ProviderType))),
	}
}

func (input QueueProvisioningInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ProviderSourceID.Empty() {
		return ErrProviderSourceIDMissing
	}
	if input.ProviderType == "" {
		return ErrProviderTypeMissing
	}
	return nil
}

func (input QueuePaidOrderProvisioningInput) Normalize() QueuePaidOrderProvisioningInput {
	return QueuePaidOrderProvisioningInput{
		OrderID:  OrderID(trim(string(input.OrderID))),
		TenantID: tenant.ID(trim(string(input.TenantID))),
	}
}

func (input QueuePaidOrderProvisioningInput) Validate() error {
	if input.OrderID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}

func (input ResolveOrderProvisioningSourceInput) Normalize() ResolveOrderProvisioningSourceInput {
	return ResolveOrderProvisioningSourceInput{
		TenantID:     tenant.ID(trim(string(input.TenantID))),
		TenantPlanID: catalog.TenantPlanID(trim(string(input.TenantPlanID))),
	}
}

func (input ResolveOrderProvisioningSourceInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.TenantPlanID.Empty() {
		return ErrTenantPlanIDMissing
	}
	return nil
}

func ensureProvisionableOrder(order Order) error {
	if order.OrderStatus != OrderStatusPaid || order.BillingStatus != BillingStatusPaid {
		return ErrProvisioningQueueNotPaid
	}
	return nil
}

func provisioningQueueKey(tenantID tenant.ID, orderID OrderID, providerSourceID catalog.ProviderSourceID) string {
	return "provisioning:" + string(tenantID) + ":" + string(orderID) + ":" + string(providerSourceID)
}

func provisioningQueuePayloadJSON(order Order, providerSourceID catalog.ProviderSourceID, providerType provider.Type) (json.RawMessage, error) {
	payload := ProvisioningQueuePayload{
		OrderID:          order.ID,
		OrderDisplayID:   order.DisplayID,
		TenantID:         order.TenantID,
		BuyerUserID:      order.BuyerUserID,
		TenantPlanID:     order.TenantPlanID,
		ProviderSourceID: providerSourceID,
		ProviderType:     providerType,
		Currency:         order.Currency,
		TotalMinor:       order.TotalMinor,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}
