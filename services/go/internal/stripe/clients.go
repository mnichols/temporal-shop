package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type PaymentIntentClienter interface {
	New(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error)
	Get(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error)
	Update(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error)
	Cancel(id string, params *stripe.PaymentIntentCancelParams) (*stripe.PaymentIntent, error)
	Capture(id string, params *stripe.PaymentIntentCaptureParams) (*stripe.PaymentIntent, error)
	Confirm(id string, params *stripe.PaymentIntentConfirmParams) (*stripe.PaymentIntent, error)
	List(listParams *stripe.PaymentIntentListParams) *paymentintent.Iter
	Search(params *stripe.PaymentIntentSearchParams) *paymentintent.SearchIter
}

type SearchClienter interface {
}
