package order

const (
	OrderAggregateType = "order"

	// OrderEventCreated is for order history and audit consumers.
	OrderEventCreated = "order.created"

	// OrderEventStatusChanged is the lifecycle event later provisioning tasks should consume.
	OrderEventStatusChanged = "order.status_changed"

	OrderProvisioningTrigger = OrderEventStatusChanged

	ServiceAggregateType = "service"

	ServiceEventRenewed     = "service.renewed"
	ServiceEventExpired     = "service.expired"
	ServiceEventGrace       = "service.grace_started"
	ServiceEventSuspended   = "service.suspended"
	ServiceEventUnsuspended = "service.unsuspended"
	ServiceEventTerminated  = "service.terminated"
)

func serviceLifecycleEventType(action ServiceLifecycleAction) string {
	switch action {
	case ServiceLifecycleActionRenew:
		return ServiceEventRenewed
	case ServiceLifecycleActionExpire:
		return ServiceEventExpired
	case ServiceLifecycleActionGrace:
		return ServiceEventGrace
	case ServiceLifecycleActionSuspend:
		return ServiceEventSuspended
	case ServiceLifecycleActionUnsuspend:
		return ServiceEventUnsuspended
	case ServiceLifecycleActionTerminate:
		return ServiceEventTerminated
	default:
		return "service.lifecycle_changed"
	}
}
