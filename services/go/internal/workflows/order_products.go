package workflows

import (
	"context"
	"errors"
	"fmt"
	stripesdk "github.com/stripe/stripe-go/v72"
	"github.com/temporalio/temporal-shop/services/go/internal/stripe"
	"github.com/temporalio/temporal-shop/services/go/internal/validation"
	"github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

var stripeActs *stripe.Handlers

func (w *Workflows) OrderProducts(ctx workflow.Context, params *workflows.OrderProducts) error {
	validator := validation.MustGetValidator(context.Background())
	if err := validator.Struct(params); err != nil {
		return fmt.Errorf("invalid products order: %w", err)
	}
	logger := log.With(
		workflow.GetLogger(ctx),
		"order_id",
		params.OrderID,
		"email",
		params.CustomerEmail,
	)
	var orderTotal = 0
	for _, p := range params.Products {
		orderTotal = orderTotal + p.Price
	}
	var orderTimedOut bool
	var paymentIntent *stripesdk.PaymentIntent
	idempotencyKey := fmt.Sprintf("%s-%d", params.OrderID, orderTotal)
	metadata := map[string]string{
		stripe.MetadataOrderID: params.OrderID,
	}
	// cleanup routine
	defer func() {
		// don't schedule for reminder if this was canceled or the order didn't time out
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) || !orderTimedOut {
			return
		}

		// When the Workflow is canceled, it has to get a new disconnected context to execute any Activities
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		// TODO : ExecuteActivity for scheduling a reminder to finish their order
		err := workflow.ExecuteActivity(newCtx, func() {}).Get(ctx, nil)
		if err != nil {
			logger.Error("reminder notification failed", "err", err)
		}
	}()

	// set up query handler ASAP
	workflow.SetQueryHandler(ctx, workflows.QueryOrder, func() (*workflows.QueryOrderResult, error) {
		if paymentIntent == nil {
			return nil, nil
		}
		return &workflows.QueryOrderResult{
			OrderID:         params.OrderID,
			PaymentIntentID: paymentIntent.ID,
		}, nil
	})

	// questions to ask
	/*
		1. does the activity accept an identifier that can be used for idempotence?
			a. I probably use stripe IdempotencyKey, but I need to consider what to do when customer revisits the order after allowing this one to expire
			b. also consider that order totals could change so would warrant a new payment intent
		2. also, we want to report this as quickly as possible but have two HTTP calls inside this activity
			a. therefore, we could extract the "order_paid" check and call it as a separate activity and lower the StartToClose on each
	*/
	paymentActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 1,
		HeartbeatTimeout:    time.Second * 3,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second * 3,
			MaximumAttempts: 3,
			// See CreatePaymentIntent activity for an alternative to this
			NonRetryableErrorTypes: []string{stripe.ErrTypeOrderPaid},
		},
	})

	if err := workflow.ExecuteActivity(paymentActivityCtx, stripeActs.CreatePaymentIntent, &stripesdk.PaymentIntentParams{
		Amount: stripesdk.Int64(int64(orderTotal)),
		Params: stripesdk.Params{
			IdempotencyKey: stripesdk.String(idempotencyKey),
			Metadata:       metadata,
		},
	}).Get(ctx, &paymentIntent); err != nil {
		// handle error for creating payment
	}

	//timerCtx, timerCancel := workflow.WithCancel(ctx)
	//cancelTimer := workflow.NewTimer(timerCtx, time.Second*3600)
	//// we want to give the order an hour to be paid
	//// after that, just cancel this order session
	//workflow.GoNamed(ctx, "order_stopwatch", func(ctx workflow.Context) {
	//	workflow.NewSelector(ctx).AddFuture(cancelTimer, func(f workflow.Future) {
	//		if err := f.Get(timerCtx, nil); err != nil {
	//			if paymentIntent == nil || paymentIntent.Status != stripesdk.PaymentIntentStatusSucceeded {
	//				orderTimedOut = true
	//			}
	//		}
	//	}).AddReceive(cancelChannel, false).Select(timerCtx)
	//})
	//
	//for !orderTimedOut && ctx.Err() != nil {
	//
	//	// schedule a cancellation if we don't get a signal that the payment intent has succeeded
	//	// maybe we schedule a reminder for the customer after _n_ days to try again
	//
	//	// pause and wait for payment details from signal
	//
	//	//
	//	return nil
	//}
	return nil
}
