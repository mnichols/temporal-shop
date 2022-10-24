package stripe

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

const ErrTypeOrderPaid = "order_paid"
const MetadataOrderID = "order_id"

var queryOrderHasBeenPaid = func(orderID string) string {
	return fmt.Sprintf("status:'succeeded' AND metadata['%s']:'%s'", MetadataOrderID, orderID)
}

type Handlers struct {
	PaymentIntentClient PaymentIntentClienter
}

// CreatePaymentIntent creates a payment intent for an Order if it has not successfully been paid
func (h *Handlers) CreatePaymentIntent(ctx context.Context, cmd *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	var orderID = cmd.Metadata[MetadataOrderID]

	// some input validation
	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	if cmd.IdempotencyKey == nil {
		return nil, fmt.Errorf("idempotency key is required")
	}

	var intent *stripe.PaymentIntent

	// first check that order has not been paid
	searchParams := &stripe.PaymentIntentSearchParams{}
	searchParams.Context = ctx
	searchParams.Query = queryOrderHasBeenPaid(orderID)
	
	activity.RecordHeartbeat(ctx)
	i := h.PaymentIntentClient.Search(searchParams)
	for i.Next() {
		intent = i.PaymentIntent()
	}
	if i.Err() != nil {
		return nil, fmt.Errorf("failed to check for payment intent %w", i.Err())
	}

	if intent != nil {
		err := fmt.Errorf("order id %s has already been paid", orderID)

		// explicitly short fuse the retry from the inside of activity with this API
		// alternatively, this can also be handled from the outside by passing in the error type (here, "order_paid") inside the RetryPolicy
		//     p := &temporal.RetryPolicy{
		//     		NonRetryableErrorTypes: []string{ ErrTypeOrderPaid }
		//     }
		return nil, temporal.NewNonRetryableApplicationError(err.Error(), ErrTypeOrderPaid, err, cmd)
	}

	activity.RecordHeartbeat(ctx)
	intent, err := h.PaymentIntentClient.New(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent %w", err)
	}
	return intent, nil
}
func (h *Handlers) CapturePayment(ctx context.Context, cmd *stripe.CaptureParams) error {
	return fmt.Errorf("not implemented")
}
func (h *Handlers) AddPaymentMethod(ctx context.Context, cmd *stripe.PaymentMethodParams) error {
	return fmt.Errorf("not implemented")
}
