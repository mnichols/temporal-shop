package workflows

//
//import (
//	"errors"
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/google/go-cmp/cmp"
//	"github.com/lucsky/cuid"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/suite"
//	"github.com/temporalio/temporal-shop/services/go/internal/clients/rsms"
//	"github.com/temporalio/temporal-shop/services/go/internal/repurchasing"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/commands"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/values"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/views"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
//	msgworkflows "github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
//	"go.temporal.io/sdk/converter"
//	"go.temporal.io/sdk/testsuite"
//	"go.temporal.io/sdk/workflow"
//)
//
//type RemindRepurchasingCustomerTest struct {
//	suite.Suite
//	testsuite.WorkflowTestSuite
//	env *testsuite.TestWorkflowEnvironment
//}
//
//// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
//func (s *RemindRepurchasingCustomerTest) SetupSuite() {
//
//}
//
//// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
//func (s *RemindRepurchasingCustomerTest) SetupTest() {
//	s.env = s.NewTestWorkflowEnvironment()
//}
//
//// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
//func (s *RemindRepurchasingCustomerTest) BeforeTest(suiteName, testName string) {
//
//}
//
//// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
//func (s *RemindRepurchasingCustomerTest) AfterTest(suiteName, testName string) {
//	s.env.AssertExpectations(s.T())
//}
//
//func TestRemindRepurchasingCustomer(t *testing.T) {
//	suite.Run(t, new(RemindRepurchasingCustomerTest))
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_DifferentNotificationIntervals() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: s.env.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 1,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//		GetCustomerParams: &rsms.GetCustomerParams{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//		},
//		GetCustomerResponse: &rsms.GetCustomerResponse{
//			View: &views.Customer{
//				SafeHours: []string{"00:00:00+00:00", "23:59:59+00:00"},
//			},
//		},
//		GetCustomerError: nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(2)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(2)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, mock.Anything).
//		Return(mockRSMS.RenderMessage).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, mock.Anything).
//		Return(mockRSMS.SendCustomerSMS).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.GetCustomer, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetCustomer).Times(1)
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.True(s.env.IsWorkflowCompleted())
//	err := s.env.GetWorkflowError()
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Promotions: params.Promotions,
//		Reminders: []*values.RepurchaseReminder{
//			{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//		ProposedReminder: &values.RepurchaseReminder{
//			ProductVariantID: productVariantIDs[0],
//			IntervalSeconds:  1,
//		},
//	})
//
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_DifferentNotificationIntervals_SpecNotFound() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: time.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 1,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	proposedReminder := &values.RepurchaseReminder{
//		ProductVariantID: productVariantIDs[0],
//		IntervalSeconds:  1,
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	missingRepurchaseSpecResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//		var isMissing = false
//		for _, pid := range r.ProductVariantIDs {
//			if pid == productVariantIDs[0] {
//				isMissing = true
//				break
//			}
//		}
//		if !isMissing {
//			missingRepurchaseSpecResponse.View.Repurchases = append(missingRepurchaseSpecResponse.View.Repurchases, r)
//		}
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//	}
//	mockRSMSEmpty := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: missingRepurchaseSpecResponse,
//		GetRepurchasesError:    nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: proposedReminder,
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMSEmpty.GetRepurchases).Times(1)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, mock.Anything).Never()
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, mock.Anything).Never()
//	var timerDurations []time.Duration
//	s.env.SetOnTimerScheduledListener(func(timerID string, duration time.Duration) {
//		timerDurations = append(timerDurations, duration)
//	})
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.Len(timerDurations, 2)
//	s.Equal(time.Second*1, timerDurations[0])
//	s.Equal(IntervalSecondsRetrySpecificationDiscovery, int64(timerDurations[1].Seconds()))
//	s.True(s.env.IsWorkflowCompleted())
//
//	err := s.env.GetWorkflowError()
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:          params.StoreID,
//		CustomerID:       params.CustomerID,
//		Promotions:       params.Promotions,
//		Purchases:        params.Purchases,
//		Reminders:        params.Reminders,
//		ProposedReminder: proposedReminder,
//	})
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_DifferentNotificationIntervals_NotificationIntervalIncreased() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//	intervalIncrease := int64(4)
//
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: time.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 1,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	increasedIntervalSpecResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//		var didIncrease = false
//		for _, pid := range r.ProductVariantIDs {
//			if pid == productVariantIDs[0] {
//				didIncrease = true
//				break
//			}
//		}
//		if didIncrease {
//			r = &views.Repurchase{
//				ProductVariantIDs:       r.ProductVariantIDs,
//				WorkflowID:              r.WorkflowID,
//				DiscountPercentage:      r.DiscountPercentage,
//				ReminderIntervalSeconds: r.ReminderIntervalSeconds + intervalIncrease, // the increase in interval
//				Templates:               r.Templates,
//				Links:                   r.Links,
//			}
//		}
//		increasedIntervalSpecResponse.View.Repurchases = append(increasedIntervalSpecResponse.View.Repurchases, r)
//
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//	}
//	mockRSMSIncreasedInterval := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: increasedIntervalSpecResponse,
//		GetRepurchasesError:    nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//	}
//	mockIncreaseRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  5,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  5,
//			},
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMSIncreasedInterval.GetRepurchases).Times(1)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(1)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockIncreaseRepurchasing.SelectReminder).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, mock.Anything).Never()
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, mock.Anything).Never()
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.True(s.env.IsWorkflowCompleted())
//	err := s.env.GetWorkflowError()
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Promotions: params.Promotions,
//		Reminders:  []*values.RepurchaseReminder{},
//		ProposedReminder: &values.RepurchaseReminder{
//			ProductVariantID: productVariantIDs[0],
//			IntervalSeconds:  intervalIncrease,
//		},
//	})
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_DifferentNotificationIntervals_NotificationIntervalDecreased() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: time.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 2,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	decreasedIntervalSpecResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//		var didIncrease = false
//		for _, pid := range r.ProductVariantIDs {
//			if pid == productVariantIDs[0] {
//				didIncrease = true
//				break
//			}
//		}
//		if didIncrease {
//			r = &views.Repurchase{
//				ProductVariantIDs:       r.ProductVariantIDs,
//				WorkflowID:              r.WorkflowID,
//				DiscountPercentage:      r.DiscountPercentage,
//				ReminderIntervalSeconds: 1, // the decrease in interval
//				Templates:               r.Templates,
//				Links:                   r.Links,
//			}
//		}
//		decreasedIntervalSpecResponse.View.Repurchases = append(decreasedIntervalSpecResponse.View.Repurchases, r)
//
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//	}
//	mockRSMSIncreasedInterval := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: decreasedIntervalSpecResponse,
//		GetRepurchasesError:    nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  4,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  2,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  2,
//			},
//		},
//	}
//	mockDecreaseRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  4,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1, // decrease
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMSIncreasedInterval.GetRepurchases).Times(1)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(1)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockDecreaseRepurchasing.SelectReminder).Times(1)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, mock.Anything).Never()
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, mock.Anything).Never()
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.True(s.env.IsWorkflowCompleted())
//	err := s.env.GetWorkflowError()
//
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Promotions: params.Promotions,
//		Reminders:  []*values.RepurchaseReminder{},
//		ProposedReminder: &values.RepurchaseReminder{
//			ProductVariantID: productVariantIDs[0],
//			IntervalSeconds:  0,
//		},
//	})
//
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_SafeHoursUpcoming() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//	startTime := time.Date(2022, time.March, 16, 16, 16, 0, 0, time.UTC)
//	s.env.SetStartTime(startTime)
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: s.env.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 1,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//		GetCustomerParams: &rsms.GetCustomerParams{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//		},
//		GetCustomerResponse: &rsms.GetCustomerResponse{
//			View: &views.Customer{
//				SafeHours: []string{
//					// by setting the startTime for the worfklow we can control
//					// for the customer's safe hours being a specific amount of
//					// time in the future
//					startTime.Add(1 * time.Hour).Format(safeHoursLayout),
//					startTime.Add(2 * time.Hour).Format(safeHoursLayout),
//				},
//			},
//		},
//		GetCustomerError: nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(2)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(2)
//
//	s.env.OnActivity(rsms.TypeClient.GetCustomer, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetCustomer).Times(1)
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.True(s.env.IsWorkflowCompleted())
//	err := s.env.GetWorkflowError()
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Promotions: params.Promotions,
//		Reminders:  []*values.RepurchaseReminder{},
//		ProposedReminder: &values.RepurchaseReminder{
//			ProductVariantID: productVariantIDs[0],
//			IntervalSeconds:  3599,
//		},
//	})
//
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_SafeHoursPassed() {
//	type sku struct {
//		productVariantID string
//		purchase         *values.Purchase
//		promotion        *values.Promotion
//	}
//
//	skus := make(map[string]sku)
//	repurchaseSpecs := make(map[string]*views.Repurchase)
//	var productVariantIDs []string
//	renderedMessageBody := cuid.New()
//	startTime := time.Date(2022, time.March, 16, 16, 16, 0, 0, time.UTC)
//	s.env.SetStartTime(startTime)
//	for i := 0; i < 3; i++ {
//		pid := cuid.New()
//		productVariantIDs = append(productVariantIDs, pid)
//
//		s := sku{
//			productVariantID: pid,
//			purchase: &values.Purchase{
//				Quantity:    i + 1,
//				PurchasedAt: s.env.Now(),
//			},
//		}
//		skus[s.productVariantID] = s
//		discount := 0
//		if i%2 != 0 {
//			discount = 0 + 15
//		}
//		repurchaseSpecs[s.productVariantID] = &views.Repurchase{
//			ProductVariantIDs:       []string{s.productVariantID},
//			WorkflowID:              "",
//			DiscountPercentage:      discount,
//			ReminderIntervalSeconds: int64(i) + 1,
//			Templates:               &views.RepurchaseTemplates{Welcome: fmt.Sprintf("welcome earthling %d", i), Reminder: fmt.Sprintf("reminding you now %d", i)},
//			Links:                   &views.RepurchaseLinks{Purchase: "https://buy.me.here"},
//		}
//	}
//
//	// first time doesn't have promotion or notification interval details
//	params := &msgworkflows.RemindRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Reminders:  []*values.RepurchaseReminder{},
//		Promotions: make(map[string]*values.Promotion),
//		Purchases:  make(map[string]*values.Purchase),
//	}
//	for k, s := range skus {
//		if s.promotion != nil {
//			params.Promotions[k] = s.promotion
//		}
//		if s.purchase != nil {
//			params.Purchases[k] = s.purchase
//		}
//
//	}
//	repurchaseResponse := &rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{},
//	}
//	for _, r := range repurchaseSpecs {
//		repurchaseResponse.View.Repurchases = append(repurchaseResponse.View.Repurchases, r)
//	}
//
//	mockRSMS := &rsms.MockClient{
//		GetRepurchasesParams: &rsms.GetRepurchasesParams{
//			ProductVariantIDs: productVariantIDs,
//			StoreID:           params.StoreID,
//		},
//		GetRepurchasesResponse: repurchaseResponse,
//		GetRepurchasesError:    nil,
//		RenderMessageParams: &rsms.RenderMessageParams{
//			StoreID:           params.StoreID,
//			CustomerID:        params.CustomerID,
//			Template:          "reminding you now 0",
//			ProductVariantIDs: productVariantIDs,
//			CustomVariables:   map[string]interface{}{},
//		},
//		RenderMessageResponse: &rsms.RenderMessageResponse{
//			View: &views.RenderedMessage{Message: renderedMessageBody},
//		},
//		RenderMessageError: nil,
//		SendCustomerSMSParams: &rsms.SendCustomerSMSParams{
//			StoreID:     params.StoreID,
//			CustomerID:  params.CustomerID,
//			MessageType: rsms.MessageTypeRepurchase,
//			Body:        renderedMessageBody,
//		},
//		SendCustomerMessageResponse: &rsms.SendCustomerSMSResponse{},
//		SendCustomerSMSError:        nil,
//		GetCustomerParams: &rsms.GetCustomerParams{
//			StoreID:    params.StoreID,
//			CustomerID: params.CustomerID,
//		},
//		GetCustomerResponse: &rsms.GetCustomerResponse{
//			View: &views.Customer{
//				SafeHours: []string{
//					startTime.Add(-8 * time.Hour).Format(safeHoursLayout),
//					startTime.Add(-1 * time.Hour).Format(safeHoursLayout),
//				},
//			},
//		},
//		GetCustomerError: nil,
//	}
//
//	mockRepurchasing := &repurchasing.MockHandlers{
//		SelectReminderParams: &commands.SelectRepurchaseReminders{
//			History: params.Reminders,
//			Candidates: []*values.RepurchaseReminder{
//				{
//					ProductVariantID: productVariantIDs[2],
//					IntervalSeconds:  3,
//				},
//				{
//					ProductVariantID: productVariantIDs[1],
//					IntervalSeconds:  2,
//				},
//				{
//					ProductVariantID: productVariantIDs[0],
//					IntervalSeconds:  1,
//				},
//			},
//		},
//		SelectReminderResponse: &views.SelectedRepurchaseReminder{
//			SelectedReminder: &values.RepurchaseReminder{
//				ProductVariantID: productVariantIDs[0],
//				IntervalSeconds:  1,
//			},
//		},
//	}
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetRepurchases).Times(2)
//
//	s.env.OnActivity(repurchasing.TypeHandlers.SelectReminder, mock.Anything, mock.Anything).
//		Return(mockRepurchasing.SelectReminder).Times(2)
//
//	s.env.OnActivity(rsms.TypeClient.GetCustomer, mock.Anything, mock.Anything).
//		Return(mockRSMS.GetCustomer).Times(1)
//
//	s.env.ExecuteWorkflow(TypeWorkflows.RemindRepurchasingCustomer, params)
//	s.True(s.env.IsWorkflowCompleted())
//	err := s.env.GetWorkflowError()
//	s.isContinuedAsNewWith(err, &workflows.RemindRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Promotions: params.Promotions,
//		Reminders:  []*values.RepurchaseReminder{},
//		ProposedReminder: &values.RepurchaseReminder{
//			ProductVariantID: productVariantIDs[0],
//			IntervalSeconds:  57599,
//		},
//	})
//
//}
//
///*
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_EarlyPurchaseReset() {
//}
//
//func (s *RemindRepurchasingCustomerTest) Test_MultiSku_CustomerPurchasedNewProductWithWorkflow() {
//}
//*/
//func (s *RemindRepurchasingCustomerTest) isContinuedAsNewWith(err error, expectParams interface{}) {
//	s.T().Helper()
//	s.Error(err)
//	s.True(workflow.IsContinueAsNewError(err), fmt.Sprintf("expected continue as new but got %s", err.Error()))
//	cer := &workflow.ContinueAsNewError{}
//	s.True(errors.As(err, &cer))
//	dc := converter.GetDefaultDataConverter()
//	var actual *workflows.RemindRepurchasingCustomer
//	err = dc.FromPayloads(cer.Input, &actual)
//	s.NoError(err)
//
//	s.Empty(cmp.Diff(expectParams, actual))
//}
