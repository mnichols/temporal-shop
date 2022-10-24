package workflows

import "github.com/temporalio/temporal-shop/services/go/pkg/messages/values"

type Ping struct {
	Value string
}
type Pong struct {
	Value string
}

type OrderProducts struct {
	OrderID          string                     `json:"order-id" validate:"required"`
	CustomerEmail    string                     `json:"customer_email" validate:"required"`
	CapturePayment   bool                       `json:"capture_payment"`
	AddPaymentMethod bool                       `json:"add_payment_method"`
	Products         map[string]*values.Product `json:"products" validate:"required"`
	ReminderInterval int64                      `json:"reminder_interval"`
}
