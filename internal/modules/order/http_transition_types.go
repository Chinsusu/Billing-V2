package order

type transitionOrderStatusRequest struct {
	FromStatus    OrderStatus   `json:"from_status"`
	ToStatus      OrderStatus   `json:"to_status"`
	BillingStatus BillingStatus `json:"billing_status"`
}
