package workflows

import (
	"fmt"
)

func EnrollRepurchasingCustomerWorkflowID(storeID, customerID string) string {
	return fmt.Sprintf("rep_enr_%s_%s", storeID, customerID)
}

//
//// EnrollRepurchasingCustomer is workflow for controlling the `welcomed` state of a customer for repurchasing as well
//// as exposing the specification and welcomed state for queries. This workflow is recalled (ContinueAsNew) periodically
//// to limit the state accumulation caused by queries.
//// Keeping `welcomed` state in the params avoids doing a lookup for welcomed state.
//// if we want to look up this status in a more durable storage though that welcomed check should
//// come first.
//func (w *Workflows) EnrollRepurchasingCustomer(ctx workflow.Context,
//	params *msgworkflows.EnrollRepurchasingCustomer) error {
//
//	const oneWeek = 24 * 7 * time.Hour
//	// TODO : how many attempts before we abandon trying to do this thing?
//
//	v := validation.MustGetValidator(context.Background())
//	if err := v.Struct(params); err != nil {
//		return err
//	}
//
//	logger := tlog.With(workflow.GetLogger(ctx),
//		"store_id", params.StoreID,
//		"customer_id", params.CustomerID)
//
//	helper := &enrollRepurchasingCustomerHelper{}
//	welcomePurchase, per := helper.selectWelcomePurchase(ctx, params.Purchases)
//	if per != nil {
//		return fmt.Errorf("failed to select welcome purchase: %w", per)
//	}
//
//	state := &msgworkflows.QueryEnrollRepurchasingCustomerResult{
//		// keeping state in the params avoids doing a lookup for welcomed state.
//		// if we want to look up this status in a more durable storage though it should
//		// come first.
//		Welcomed:                params.Welcomed,
//		RepurchaseSpecification: nil,
//	}
//
//	if err := workflow.SetQueryHandler(
//		ctx,
//		msgworkflows.QueryEnrollRepurchasingCustomer,
//		func() (*msgworkflows.QueryEnrollRepurchasingCustomerResult, error) {
//			logger.Debug(msgworkflows.QueryEnrollRepurchasingCustomer, "state", state)
//			return state, nil
//		},
//	); err != nil {
//		err = fmt.Errorf("failed to set query handler for '%s': %w", msgworkflows.QueryEnrollRepurchasingCustomer, err)
//		logger.Error(err.Error(), "query", msgworkflows.QueryEnrollRepurchasingCustomer, "err", err)
//		return err
//	}
//
//	// _NOTE_ that we are using retry policy for backoff up to one week
//	// for any error conditions we encounter
//	ao := workflow.ActivityOptions{
//		StartToCloseTimeout: 120 * time.Second,
//		RetryPolicy: &temporal.RetryPolicy{
//			InitialInterval:        12 * time.Hour,
//			BackoffCoefficient:     2,
//			MaximumInterval:        oneWeek, // one week
//			MaximumAttempts:        5,
//			NonRetryableErrorTypes: nil,
//		},
//	}
//	activityCtx := workflow.WithActivityOptions(ctx, ao)
//
//	logger = tlog.With(logger,
//		"welcome_purchase", welcomePurchase,
//	)
//
//	logger.Debug("enrolling repurchasing customer")
//
//	if !state.Welcomed {
//		// first let's grab notification preferences
//		var notificationPreferencesRes *recharge.GetCustomerNotificationPreferencesResponse
//		if err := workflow.ExecuteActivity(
//			activityCtx,
//			recharge.TypeClient.GetCustomerNotificationPreferences,
//			recharge.GetCustomerNotificationPreferencesParams{
//				StoreID:    params.StoreID,
//				CustomerID: params.CustomerID,
//			},
//		).Get(activityCtx, &notificationPreferencesRes); err != nil || notificationPreferencesRes.View == nil {
//			logger.Error("failed to get customer notification preferences", "err", err)
//			return helper.replayWorkflow(ctx, params)
//		}
//
//		selectedChannel, shouldSendWelcome, err := helper.determineNotificationType(activityCtx, params, notificationPreferencesRes.View, logger)
//		if err != nil {
//			logger.Error("failed to determine customer notification preferences", "err", err)
//			return helper.replayWorkflow(ctx, params)
//
//		}
//		if !shouldSendWelcome {
//			// customer doesnt want repurchase notifications so let's just bail
//			logger.Info("customer should not be welcomed to repurchasing", "preferences", notificationPreferencesRes.View)
//			return nil
//		}
//		if selectedChannel == recharge.NotificationChannelEmail {
//			return fmt.Errorf("email is not support for repurchase notifications yet")
//		}
//	}
//
//	// let's try to get the spec for the purchase
//	// why don't we do this first? because we want to opt in the customer as soon as we can without
//	// disruption from the state of the spec
//	if state.RepurchaseSpecification == nil {
//		var rer error
//		if state.RepurchaseSpecification, rer = helper.getRepurchaseSpec(activityCtx, params, logger, welcomePurchase); rer != nil {
//			// we couldnt find our spec...let's keep trying
//			logger.Error(rer.Error(), "err", rer)
//			return helper.replayWorkflow(ctx, params)
//		}
//	}
//
//	if !state.Welcomed {
//		logger.Info("sending message")
//		if err := helper.sendSMSWelcome(activityCtx, params, logger, []string{welcomePurchase.ProductVariantID}, state.RepurchaseSpecification); err != nil {
//			return helper.replayWorkflow(ctx, params)
//		}
//		childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
//			WorkflowID: RemindRepurchasingCustomerWorkflowID(params.StoreID, params.CustomerID),
//		})
//		logger.Info("starting reminders for customer")
//		workflow.ExecuteChildWorkflow(childCtx, TypeWorkflows.RemindRepurchasingCustomer, &msgworkflows.RemindRepurchasingCustomer{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//			Purchases:  params.Purchases,
//		})
//	}
//
//	state.Welcomed = true
//
//	// TODO the `cancel` func will be used when we support the UnenrollRepurchasingCustomer signal
//	resetCtx, _ := workflow.WithCancel(ctx)
//	resetTimer := workflow.NewTimer(resetCtx, oneWeek)
//	if err := resetTimer.Get(resetCtx, nil); err != nil {
//
//	}
//	logger.Debug("resetting orchestration state")
//	return workflow.NewContinueAsNewError(ctx, TypeWorkflows.EnrollRepurchasingCustomer, &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Welcomed:   state.Welcomed,
//		Attempts:   0, // we are welcomed so reset
//	})
//}
//
//type enrollRepurchasingCustomerHelper struct {
//}
//
//func (h *enrollRepurchasingCustomerHelper) replayWorkflow(ctx workflow.Context, params *msgworkflows.EnrollRepurchasingCustomer) error {
//	return workflow.NewContinueAsNewError(ctx, TypeWorkflows.EnrollRepurchasingCustomer, &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Welcomed:   params.Welcomed,
//		Attempts:   params.Attempts + 1,
//	})
//}
//func (h *enrollRepurchasingCustomerHelper) getRepurchaseSpec(
//	ctx workflow.Context,
//	params *msgworkflows.EnrollRepurchasingCustomer,
//	logger tlog.Logger,
//	purchase *values.Purchase,
//) (*views.Repurchase, error) {
//	var reps *rsms.GetRepurchasesResponse
//
//	if err := workflow.ExecuteActivity(ctx, rsms.TypeClient.GetRepurchases, rsms.GetRepurchasesParams{
//		StoreID:           params.StoreID,
//		ProductVariantIDs: []string{purchase.ProductVariantID},
//	}).Get(ctx, &reps); err != nil || reps.View == nil {
//		logger.Error("failed to get repurchase specs", "err", err)
//		return nil, err
//	}
//
//	if len(reps.View.Repurchases) != 1 {
//		err := fmt.Errorf("expected exactly one repurchase spec but received %d", len(reps.View.Repurchases))
//		return nil, err
//	}
//	return reps.View.Repurchases[0], nil
//}
//func (h *enrollRepurchasingCustomerHelper) selectWelcomePurchase(ctx workflow.Context, purchases map[string]*values.Purchase) (*values.Purchase, error) {
//
//	if len(purchases) == 0 {
//		return nil, fmt.Errorf("at least one purchase must be provided")
//	}
//	if len(purchases) == 1 {
//		for pid := range purchases {
//			return purchases[pid], nil
//		}
//	}
//	var latest time.Time
//	var ref string
//	for pid, p := range purchases {
//		if p.PurchasedAt.After(latest) {
//			latest = p.PurchasedAt
//			ref = pid
//		}
//	}
//	return purchases[ref], nil
//}
//func (h *enrollRepurchasingCustomerHelper) determineNotificationType(
//	ctx workflow.Context,
//	params *msgworkflows.EnrollRepurchasingCustomer,
//	prefs *views.CustomerNotificationPreferences,
//	logger tlog.Logger,
//) (recharge.NotificationChannel, bool, error) {
//	// we only support SMS right now
//	smsStatus := h.getSMSStatus(ctx, prefs)
//	emailStatus := h.getEmailStatus(ctx, prefs)
//
//	// early return
//	if smsStatus == recharge.NotificationPreferenceStatusDeclined && emailStatus == recharge.NotificationPreferenceStatusDeclined {
//		return recharge.NotificationChannelUnknown, false, nil
//	}
//	var selectedChannel = recharge.NotificationChannelUnknown
//	var channelsToAccept []recharge.NotificationChannel
//	if emailStatus == recharge.NotificationPreferenceStatusUnspecified {
//		channelsToAccept = append(channelsToAccept, recharge.NotificationChannelEmail)
//		selectedChannel = recharge.NotificationChannelEmail
//	}
//	if smsStatus == recharge.NotificationPreferenceStatusUnspecified {
//		channelsToAccept = append(channelsToAccept, recharge.NotificationChannelSMS)
//		selectedChannel = recharge.NotificationChannelSMS
//	}
//
//	for _, c := range channelsToAccept {
//		// let's make that accepted at recharge
//		if err := workflow.ExecuteActivity(ctx, recharge.TypeClient.SetCustomerNotificationPreferences, recharge.SetCustomerNotificationPreferencesParams{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//			Channel:    c,
//			Status:     recharge.NotificationPreferenceStatusAccepted,
//			Type:       recharge.NotificationPreferenceTypeRepurchase,
//		}).Get(ctx, nil); err != nil {
//			return recharge.NotificationChannelUnknown, false, err
//		}
//	}
//	return selectedChannel, true, nil
//}
//func (h *enrollRepurchasingCustomerHelper) getSMSStatus(ctx workflow.Context, prefs *views.CustomerNotificationPreferences) recharge.NotificationPreferenceStatus {
//	if prefs == nil ||
//		prefs.NotificationPreferences.SMS == nil ||
//		prefs.NotificationPreferences.SMS.Repurchase == nil {
//		return recharge.NotificationPreferenceStatusUnknown
//	}
//	return recharge.NotificationPreferenceStatus(prefs.NotificationPreferences.SMS.Repurchase.Status)
//}
//func (h *enrollRepurchasingCustomerHelper) getEmailStatus(ctx workflow.Context, prefs *views.CustomerNotificationPreferences) recharge.NotificationPreferenceStatus {
//	if prefs == nil ||
//		prefs.NotificationPreferences.Email == nil ||
//		prefs.NotificationPreferences.Email.Repurchase == nil {
//		return recharge.NotificationPreferenceStatusUnknown
//	}
//	return recharge.NotificationPreferenceStatus(prefs.NotificationPreferences.Email.Repurchase.Status)
//}
//func (h *enrollRepurchasingCustomerHelper) sendSMSWelcome(ctx workflow.Context, params *msgworkflows.EnrollRepurchasingCustomer, logger tlog.Logger, productVariantIDs []string, spec *views.Repurchase) error {
//
//	var reminderMessage *rsms.RenderMessageResponse
//	renderMessageParams := rsms.RenderMessageParams{
//		StoreID:           params.StoreID,
//		CustomerID:        params.CustomerID,
//		Template:          spec.Templates.Welcome,
//		ProductVariantIDs: productVariantIDs,
//		CustomVariables:   map[string]interface{}{},
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
