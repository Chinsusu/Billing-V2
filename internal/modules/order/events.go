package order

const (
	OrderAggregateType = "order"

	// OrderEventCreated is for order history and audit consumers.
	OrderEventCreated = "order.created"

	// OrderEventStatusChanged is the lifecycle event later provisioning tasks should consume.
	OrderEventStatusChanged = "order.status_changed"

	OrderProvisioningTrigger = OrderEventStatusChanged
)
