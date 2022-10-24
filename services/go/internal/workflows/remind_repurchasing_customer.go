package workflows

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/hashicorp/go-multierror"
//	"math"
//	"time"
//
//	serr "github.com/bdandy/go-errors"
//	"github.com/temporalio/temporal-shop/services/go/internal/repurchasing"
//	"github.com/temporalio/temporal-shop/services/go/internal/validation"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/commands"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/values"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/views"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
//	msgworkflows "github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
//	tlog "go.temporal.io/sdk/log"
//	"go.temporal.io/sdk/temporal"
//	"go.temporal.io/sdk/workflow"
//)
//
//const secondsPerWeek = 60 * 60 * 24 * 7
//const IntervalSecondsRetrySpecificationDiscovery int64 = secondsPerWeek
//const safeHoursLayout = "15:04:05+00:00"
//
//func RemindRepurchasingCustomerWorkflowID(storeID, customerID string) string {
//	return fmt.Sprintf("rep_rem_%s_%s", storeID, customerID)
//}
//func (w *Workflows) RemindRepurchasingCustomer(
//	ctx workflow.Context, params *msgworkflows.RemindRepurchasingCustomer) error {
//
//	v := validation.MustGetValidator(context.Background())
//	if err := v.Struct(params); err != nil {
//		return err
//	}
//	helper := remindRepurchasingCustomerHelper{}
//
//	ao := workflow.ActivityOptions{
//		StartToCloseTimeout: 10 * time.Second,
//		RetryPolicy: &temporal.RetryPolicy{
//			MaximumAttempts: 5,
//			InitialInterval: time.Second * 1,
//		},
//	}
//	ctx = workflow.WithActivityOptions(ctx, ao)
//
//	// GA: subscribe signals here
//
//	// flag to control whether we want to continue
//	var shouldRescheduleReminder = true
//	var productVariantIDs = helper.toProductVariantIDs(params.Purchases)
//
//	// contexts
//	notificationCtx, _ := workflow.WithCancel(ctx)
//
//	logger := tlog.With(workflow.GetLogger(ctx),
//		"store_id", params.StoreID,
//		"customer_id", params.CustomerID,
//		"product_variant_ids", productVariantIDs,
//		"reminder_intervals", params.Reminders,
//		"proposed_reminder", params.ProposedReminder,
//	)
//	logger.Debug("remind_repurchasing_customer")
//
//	// ensure we have reminder intervals, with the latest interval being the first item
//	if params.ProposedReminder == nil {
//		// state
//		reps, err := helper.getCurrentSpecifications(ctx, params, logger)
//
//		if err != nil {
//			logger.Error("err", err)
//			return err
//		}
//
//		var rerr error
//		if params.ProposedReminder, _, rerr = helper.proposeReminder(
//			ctx,
//			params,
//			logger,
//			reps); rerr != nil {
//
//			logger.Error("err", rerr)
//			return fmt.Errorf("failed to discover reminder intervals: %w", rerr)
//		}
//	}
//
//	// let's sleep until it is time to send a reminder
//	remindDuration := time.Second * time.Duration(params.ProposedReminder.IntervalSeconds)
//	err := workflow.Sleep(notificationCtx, remindDuration)
//	if err != nil {
//		if temporal.IsCanceledError(err) {
//			logger.Info("reminder was canceled")
//			return nil
//		}
//	}
//
//	// we have woken up and are now ready to send our reminder notification
//
//	reps, err := helper.getCurrentSpecifications(ctx, params, logger)
//
//	if err != nil {
//		logger.Error("err", err)
//		return err
//	}
//
//	nextProposedReminder, spec, err := helper.proposeReminder(ctx, params, logger, reps)
//	if err != nil {
//		logger.Error(err.Error(), "err", err, "stack", serr.Stack(err))
//		if errors.Is(err, RepurchaseSpecificationNotFoundError) {
//			retrySpecCtx, _ := workflow.WithCancel(ctx)
//			err2 := workflow.Sleep(retrySpecCtx, time.Duration(IntervalSecondsRetrySpecificationDiscovery)*time.Second)
//			if err2 != nil {
//				if temporal.IsCanceledError(err2) {
//					logger.Info("reminder was canceled")
//					return nil
//				}
//			}
//			return helper.replayWorkflow(ctx, params)
//		}
//		return err
//	}
//
//	if nextProposedReminder.ProductVariantID == params.ProposedReminder.ProductVariantID && nextProposedReminder.IntervalSeconds == params.ProposedReminder.IntervalSeconds {
//		secondsUntilNextSafeHourWindow, err := helper.calculateSecondsUntilNextSafeHourWindow(ctx, params, logger)
//		if err != nil {
//			logger.Error("err", err)
//			return fmt.Errorf("unable to determine safe hours: %w", err)
//		}
//		if secondsUntilNextSafeHourWindow > 0 {
//			params.ProposedReminder.IntervalSeconds = secondsUntilNextSafeHourWindow
//			return helper.replayWorkflow(ctx, params)
//		}
//		// we are within the safe hours window &
//		// there have been no changes while we were sleeping #bullock. go ahead and send the reminder
//		if err := helper.sendSMSReminder(ctx,
//			params,
//			logger,
//			productVariantIDs,
//			spec,
//			params.ProposedReminder,
//		); err != nil {
//			logger.Error("err", err)
//			return err
//		}
//		// update our history of reminders by prepending the proposed reminder
//		params.Reminders = append([]*values.RepurchaseReminder{params.ProposedReminder}, params.Reminders...)
//	} else {
//		// there were changes while we were sleeping
//		// let's replay with the updated proposed reminder
//		if nextProposedReminder.IntervalSeconds > params.ProposedReminder.IntervalSeconds {
//			nextProposedReminder.IntervalSeconds = nextProposedReminder.IntervalSeconds - params.ProposedReminder.IntervalSeconds
//		} else if nextProposedReminder.IntervalSeconds < params.ProposedReminder.IntervalSeconds {
//			// this effectively causes and immediate message to be sent
//			nextProposedReminder.IntervalSeconds = 0
//		}
//		params.ProposedReminder = nextProposedReminder
//	}
//
//	if !shouldRescheduleReminder {
//		logger.Info("abandoning reminders")
//		return nil
//	}
//	// this tells temporal to "restart" the workflow with the new params,
//	// reusing the same ID but without the history
//	return helper.replayWorkflow(ctx, params)
//}
//
//type remindRepurchasingCustomerHelper struct{}
//
//func (h *remindRepurchasingCustomerHelper) sendSMSReminder(ctx workflow.Context,
//	params *workflows.RemindRepurchasingCustomer,
//	logger tlog.Logger,
//	productVariantIDs []string,
//	spec *views.Repurchase,
//	reminder *values.RepurchaseReminder,
//) error {
//
//	var reminderMessage *rsms.RenderMessageResponse
//	renderMessageParams := rsms.RenderMessageParams{
//		StoreID:           params.StoreID,
//		CustomerID:        params.CustomerID,
//		Template:          spec.Templates.Reminder,
//		ProductVariantIDs: h.sortProductVariantIDs(reminder.ProductVariantID, productVariantIDs),
//		CustomVariables:   map[string]interface{}{},
//	}
//
//	if spec.DiscountPercentage > 0 {
//		renderMessageParams.CustomVariables["applied_discount"] = spec.DiscountPercentage
//	}
//	if err := workflow.ExecuteActivity(ctx, rsms.TypeClient.RenderMessage, renderMessageParams).
//		Get(ctx, &reminderMessage); err != nil {
//		logger.Error("err", err)
//		return fmt.Errorf("failed to render message: %w", err)
//	}
//
//	if reminderMessage.View == nil {
//		return fmt.Errorf("failed to create message body")
//	}
//
//	sms := rsms.SendCustomerSMSParams{
//		StoreID:     params.StoreID,
//		CustomerID:  params.CustomerID,
//		Body:        reminderMessage.View.Message,
//		MessageType: rsms.MessageTypeRepurchase,
//	}
//	logger.Debug("reminding repurchase customer", sms)
//	var sentSMS *rsms.SendCustomerSMSResponse
//	if err := workflow.ExecuteActivity(ctx, rsms.TypeClient.SendCustomerSMS, sms).
//		Get(ctx, &sentSMS); err != nil {
//		logger.Error("err", err)
//		return err
//	}
//	return nil
//}
//func (h *remindRepurchasingCustomerHelper) getCurrentSpecifications(ctx workflow.Context,
//	params *workflows.RemindRepurchasingCustomer,
//	logger tlog.Logger,
//) (*rsms.GetRepurchasesResponse, error) {
//	// state
//	var reps *rsms.GetRepurchasesResponse
//	var productVariantIDs = h.toProductVariantIDs(params.Purchases)
//
//	logger.Debug("discovering reminder interval")
//	if err := workflow.ExecuteActivity(ctx,
//		rsms.TypeClient.GetRepurchases,
//		rsms.GetRepurchasesParams{
//			StoreID:           params.StoreID,
//			ProductVariantIDs: productVariantIDs,
//		},
//	).Get(ctx, &reps); err != nil {
//		logger.Error("err", err)
//		return nil, fmt.Errorf("failed to get repurchases: %w", err)
//	}
//	if reps.View == nil {
//		return nil, fmt.Errorf("view is missing %v", reps)
//	}
//	return reps, nil
//}
//
//// proposeReminder tries to return the appropriate ordered reminders
//// while considering that a repurchase spec might have been removed
//// and while only considering products that have actually been purchased by the customer.
//func (h *remindRepurchasingCustomerHelper) proposeReminder(
//	ctx workflow.Context,
//	params *workflows.RemindRepurchasingCustomer,
//	logger tlog.Logger,
//	reps *rsms.GetRepurchasesResponse,
//) (*values.RepurchaseReminder, *views.Repurchase, error) {
//
//	candidates, err := h.reduceProposedReminders(params.Purchases, reps.View.Repurchases)
//	if err != nil {
//		logger.Error("err", err)
//		return nil, nil, fmt.Errorf("failed to reduce reminders %w", err)
//	}
//
//	cmd := &commands.SelectRepurchaseReminders{
//		History:    params.Reminders,
//		Candidates: candidates,
//	}
//
//	var selected *views.SelectedRepurchaseReminder
//	if err := workflow.ExecuteActivity(ctx,
//		repurchasing.TypeHandlers.SelectReminder,
//		cmd).Get(ctx, &selected); err != nil {
//		logger.Error("err", err)
//		return nil, nil, fmt.Errorf("failed to select reminder: %w", err)
//	}
//
//	spec, exists := h.trySeekRepurchase(selected.SelectedReminder, reps.View.Repurchases)
//	if !exists {
//		return nil, nil, RepurchaseSpecificationNotFoundError.New(selected.SelectedReminder.ProductVariantID)
//	}
//
//	return selected.SelectedReminder, spec, nil
//}
//
//func (h *remindRepurchasingCustomerHelper) toProductVariantIDs(purchases map[string]*values.Purchase) []string {
//	var productVariantIDs = make([]string, len(purchases))
//	i := 0
//	for pid := range purchases {
//		productVariantIDs[i] = pid
//		i++
//	}
//	return productVariantIDs
//}
//
//// sortProductVariantIDs simply plops `subject` to the first item while removing from the passed in ids
//func (h *remindRepurchasingCustomerHelper) sortProductVariantIDs(subject string, productVariantIDs []string) []string {
//	for i, v := range productVariantIDs {
//		if v == subject {
//			productVariantIDs = append(productVariantIDs[:i], productVariantIDs[i+1:]...)
//		}
//	}
//	return append([]string{subject}, productVariantIDs...)
//}
//
//// trySeekRepurchase finds the related Repurchase spec for the input reminder
//func (h *remindRepurchasingCustomerHelper) trySeekRepurchase(needle *values.RepurchaseReminder, haystack []*views.Repurchase) (*views.Repurchase, bool) {
//	for _, r := range haystack {
//		for _, pid := range r.ProductVariantIDs {
//			if pid == needle.ProductVariantID {
//				return r, true
//			}
//		}
//	}
//	return nil, false
//}
//
//// reduceProposedReminders ensure repurchase reminders are relevant based on purchases
//func (h *remindRepurchasingCustomerHelper) reduceProposedReminders(purchases map[string]*values.Purchase, repurchases []*views.Repurchase) ([]*values.RepurchaseReminder, error) {
//	var err *multierror.Error
//	out := []*values.RepurchaseReminder{}
//	pid2intervals := map[string]int64{}
//	for _, r := range repurchases {
//		for _, pid := range r.ProductVariantIDs {
//			pid2intervals[pid] = r.ReminderIntervalSeconds
//		}
//	}
//
//	for pid := range purchases {
//		ntrvl, exists := pid2intervals[pid]
//		if !exists {
//			err = multierror.Append(
//				err,
//				RepurchaseSpecificationNotFoundError.New(pid).WithStack())
//		} else {
//			out = append(out, &values.RepurchaseReminder{
//				ProductVariantID: pid,
//				IntervalSeconds:  ntrvl,
//			})
//		}
//	}
//
//	return out, err.ErrorOrNil()
//}
//
//func (h *remindRepurchasingCustomerHelper) replayWorkflow(
//	ctx workflow.Context,
//	params *msgworkflows.RemindRepurchasingCustomer,
//) error {
//	var limit = int(math.Min(float64(len(params.Reminders)), 5))
//	return workflow.NewContinueAsNewError(
//		ctx,
//		TypeWorkflows.RemindRepurchasingCustomer,
//		&msgworkflows.RemindRepurchasingCustomer{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//			Purchases:  params.Purchases,
//			Promotions: params.Promotions,
//			// only keep 5 reminders in history
//			Reminders:        params.Reminders[:limit],
//			ProposedReminder: params.ProposedReminder,
//		},
//	)
//}
//
//func (h *remindRepurchasingCustomerHelper) calculateSecondsUntilNextSafeHourWindow(ctx workflow.Context,
//	params *workflows.RemindRepurchasingCustomer,
//	logger tlog.Logger) (int64, error) {
//	getCustomerParams := &rsms.GetCustomerParams{CustomerID: params.CustomerID, StoreID: params.StoreID}
//	var getCustomerResponse *rsms.GetCustomerResponse
//	if err := workflow.ExecuteActivity(ctx, rsms.TypeClient.GetCustomer, getCustomerParams).
//		Get(ctx, &getCustomerResponse); err != nil {
//		logger.Error("err", err)
//		return 0, fmt.Errorf("unable to get customer: %w", err)
//	}
//	start, end, err := extractSafeHours(ctx, getCustomerResponse)
//	if err != nil {
//		return 0, fmt.Errorf("unable to extract safe hours: %w", err)
//	}
//	now := workflow.Now(ctx).UTC()
//	if now.Before(start) {
//		// safe hours window is upcoming: wait for it
//		logger.Info("safe hours have not started")
//		return int64(start.Sub(now).Seconds()), nil
//	}
//	if now.After(end) {
//		// safe hours window already passed: try again tomorrow
//		logger.Info("safe hours already passed")
//		start = start.Add(24 * time.Hour)
//		return int64(start.Sub(now).Seconds()), nil
//	}
//	return 0, nil
//}
//
//func extractSafeHours(ctx workflow.Context, r *rsms.GetCustomerResponse) (time.Time, time.Time, error) {
//	var safeHourStart, safeHourEnd time.Time
//	if r.View == nil {
//		return safeHourStart, safeHourEnd, fmt.Errorf("view is missing and required")
//	}
//	if len(r.View.SafeHours) < 2 {
//		return safeHourStart, safeHourEnd, fmt.Errorf("safe hours missing and required")
//	}
//	now := workflow.Now(ctx).UTC()
//	var err error
//	safeHourStart, err = time.Parse(safeHoursLayout, r.View.SafeHours[0])
//	if err != nil {
//		return safeHourStart, safeHourEnd, fmt.Errorf("unable to parse safe hours start")
//	}
//	safeHourStart = time.Date(now.Year(), now.Month(), now.Day(), safeHourStart.Hour(), safeHourStart.Minute(), 0, 0, time.UTC)
//	safeHourEnd, err = time.Parse(safeHoursLayout, r.View.SafeHours[1])
//	if err != nil {
//		return safeHourStart, safeHourEnd, fmt.Errorf("unable to parse safe hours end")
//	}
//	safeHourEnd = time.Date(now.Year(), now.Month(), now.Day(), safeHourEnd.Hour(), safeHourEnd.Minute(), 0, 0, time.UTC)
//	return safeHourStart, safeHourEnd, nil
//}
