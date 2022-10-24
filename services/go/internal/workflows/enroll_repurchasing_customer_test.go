package workflows

//
//import (
//	"errors"
//	"fmt"
//	"github.com/google/go-cmp/cmp"
//	"github.com/lucsky/cuid"
//	"github.com/stretchr/testify/mock"
//	"github.com/temporalio/temporal-shop/services/go/internal/clients/recharge"
//	"github.com/temporalio/temporal-shop/services/go/internal/clients/rsms"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/values"
//	"github.com/temporalio/temporal-shop/services/go/pkg/messages/views"
//	"go.temporal.io/sdk/converter"
//	"go.temporal.io/sdk/workflow"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/suite"
//	msgworkflows "github.com/temporalio/temporal-shop/services/go/pkg/messages/workflows"
//	"go.temporal.io/sdk/testsuite"
//)
//
//type EnrollRepurchasingCustomerTest struct {
//	suite.Suite
//	testsuite.WorkflowTestSuite
//	env *testsuite.TestWorkflowEnvironment
//}
//
//// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
//func (s *EnrollRepurchasingCustomerTest) SetupSuite() {
//}
//
//// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
//func (s *EnrollRepurchasingCustomerTest) SetupTest() {
//	s.env = s.NewTestWorkflowEnvironment()
//}
//
//// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
//func (s *EnrollRepurchasingCustomerTest) BeforeTest(suiteName, testName string) {
//
//}
//
//// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
//func (s *EnrollRepurchasingCustomerTest) AfterTest(suiteName, testName string) {
//	s.env.AssertExpectations(s.T())
//}
//
//func TestEnrollRepurchasingCustomer(t *testing.T) {
//	suite.Run(t, new(EnrollRepurchasingCustomerTest))
//
//}
//
//func (s *EnrollRepurchasingCustomerTest) TestUnenrolledCustomer_HasSMSOptedIn() {
//
//	productVariantID := cuid.New()
//	params := &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Purchases: map[string]*values.Purchase{
//			productVariantID: {Quantity: 3, PurchasedAt: s.env.Now().Add(-3 * time.Second), ProductVariantID: productVariantID},
//		},
//	}
//
//	repurchaseOptedIn := &views.CustomerNotificationPreferenceOptIn{
//		Status: recharge.NotificationPreferenceStatusAccepted.String(),
//	}
//	repurchaseSpec := &views.Repurchase{
//
//		ProductVariantIDs:       []string{productVariantID},
//		WorkflowID:              cuid.New(),
//		DiscountPercentage:      0,
//		ReminderIntervalSeconds: 0,
//		Templates: &views.RepurchaseTemplates{
//			Welcome:  "aüf wiédersehen",
//			Reminder: "",
//		},
//		Links: nil,
//	}
//	renderedWelcomeBody := cuid.New()
//	s.env.OnActivity(recharge.TypeClient.GetCustomerNotificationPreferences, mock.Anything, recharge.GetCustomerNotificationPreferencesParams{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//	}).Return(&recharge.GetCustomerNotificationPreferencesResponse{
//		View: &views.CustomerNotificationPreferences{
//			NotificationPreferences: struct {
//				SMS   *views.CustomerNotificationPreferenceChannel `json:"sms"`
//				Email *views.CustomerNotificationPreferenceChannel `json:"email"`
//			}{
//				SMS: &views.CustomerNotificationPreferenceChannel{
//					Transactional: nil,
//					Promotional:   nil,
//					Repurchase:    repurchaseOptedIn,
//				},
//			},
//		},
//	}, nil)
//
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, rsms.GetRepurchasesParams{
//		StoreID:           params.StoreID,
//		ProductVariantIDs: []string{productVariantID},
//	}).Return(&rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{
//			StoreID:     params.StoreID,
//			Repurchases: []*views.Repurchase{repurchaseSpec},
//		},
//	}, nil)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, rsms.RenderMessageParams{
//		Template:          repurchaseSpec.Templates.Welcome,
//		StoreID:           params.StoreID,
//		CustomerID:        params.CustomerID,
//		ProductVariantIDs: []string{productVariantID},
//		CustomVariables:   make(map[string]interface{}),
//	}).Return(&rsms.RenderMessageResponse{
//		View:         &views.RenderedMessage{Message: renderedWelcomeBody},
//		HTTPResponse: nil,
//	}, nil)
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, rsms.SendCustomerSMSParams{
//		StoreID:     params.StoreID,
//		CustomerID:  params.CustomerID,
//		Body:        renderedWelcomeBody,
//		MessageType: rsms.MessageTypeRepurchase,
//	}).Return(&rsms.SendCustomerSMSResponse{}, nil)
//
//	w := &Workflows{}
//	s.env.RegisterWorkflow(w.RemindRepurchasingCustomer)
//	s.env.OnWorkflow(TypeWorkflows.RemindRepurchasingCustomer, mock.Anything, mock.Anything).
//		Return(func(ctx workflow.Context, wparams *msgworkflows.RemindRepurchasingCustomer) error {
//
//			diff := cmp.Diff(&msgworkflows.RemindRepurchasingCustomer{
//				StoreID:    params.StoreID,
//				CustomerID: params.CustomerID,
//				Purchases:  params.Purchases,
//			}, wparams)
//			if diff != "" {
//				s.Empty(diff, "failed to execute reminder correctly")
//			}
//
//			return nil
//		})
//
//	var state *msgworkflows.QueryEnrollRepurchasingCustomerResult
//
//	s.env.RegisterDelayedCallback(func() {
//		var result converter.EncodedValue
//		var queryErr error
//		result, queryErr = s.env.QueryWorkflow(msgworkflows.QueryEnrollRepurchasingCustomer)
//		s.NoError(queryErr)
//		queryErr = result.Get(&state)
//		s.NoError(queryErr)
//	}, time.Minute*1)
//
//	s.env.ExecuteWorkflow(TypeWorkflows.EnrollRepurchasingCustomer, params)
//
//	s.True(s.env.IsWorkflowCompleted())
//	s.Empty(cmp.Diff(&msgworkflows.QueryEnrollRepurchasingCustomerResult{
//		Welcomed:                true,
//		RepurchaseSpecification: repurchaseSpec,
//	}, state))
//
//	err := s.env.GetWorkflowError()
//
//	s.isContinuedAsNewWith(err, &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Welcomed:   true,
//		Attempts:   0,
//	})
//}
//func (s *EnrollRepurchasingCustomerTest) TestPreviouslyEnrolledCustomer() {
//
//	productVariantID := cuid.New()
//	params := &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    cuid.New(),
//		CustomerID: cuid.New(),
//		Purchases: map[string]*values.Purchase{
//			productVariantID: {Quantity: 3, PurchasedAt: s.env.Now().Add(-3 * time.Second), ProductVariantID: productVariantID},
//		},
//		Welcomed: true,
//	}
//
//	repurchaseSpec := &views.Repurchase{
//
//		ProductVariantIDs:       []string{productVariantID},
//		WorkflowID:              cuid.New(),
//		DiscountPercentage:      0,
//		ReminderIntervalSeconds: 0,
//		Templates: &views.RepurchaseTemplates{
//			Welcome:  "aüf wiédersehen",
//			Reminder: "",
//		},
//		Links: nil,
//	}
//	renderedWelcomeBody := cuid.New()
//	s.env.OnActivity(recharge.TypeClient.GetCustomerNotificationPreferences, mock.Anything, recharge.GetCustomerNotificationPreferencesParams{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//	}).Never()
//
//	s.env.OnActivity(rsms.TypeClient.GetRepurchases, mock.Anything, rsms.GetRepurchasesParams{
//		StoreID:           params.StoreID,
//		ProductVariantIDs: []string{productVariantID},
//	}).Return(&rsms.GetRepurchasesResponse{
//		View: &views.StoreRepurchases{
//			StoreID:     params.StoreID,
//			Repurchases: []*views.Repurchase{repurchaseSpec},
//		},
//	}, nil)
//
//	s.env.OnActivity(rsms.TypeClient.RenderMessage, mock.Anything, rsms.RenderMessageParams{
//		Template:          repurchaseSpec.Templates.Welcome,
//		StoreID:           params.StoreID,
//		CustomerID:        params.CustomerID,
//		ProductVariantIDs: []string{productVariantID},
//		CustomVariables:   make(map[string]interface{}),
//	}).Never()
//
//	s.env.OnActivity(rsms.TypeClient.SendCustomerSMS, mock.Anything, rsms.SendCustomerSMSParams{
//		StoreID:     params.StoreID,
//		CustomerID:  params.CustomerID,
//		Body:        renderedWelcomeBody,
//		MessageType: rsms.MessageTypeRepurchase,
//	}).Never()
//	w := &Workflows{}
//	s.env.RegisterWorkflow(w.RemindRepurchasingCustomer)
//	s.env.OnWorkflow(TypeWorkflows.RemindRepurchasingCustomer, mock.Anything, mock.Anything).Never()
//
//	var state *msgworkflows.QueryEnrollRepurchasingCustomerResult
//
//	s.env.RegisterDelayedCallback(func() {
//		var result converter.EncodedValue
//		var queryErr error
//		result, queryErr = s.env.QueryWorkflow(msgworkflows.QueryEnrollRepurchasingCustomer)
//		s.NoError(queryErr)
//		queryErr = result.Get(&state)
//		s.NoError(queryErr)
//	}, time.Minute*1)
//
//	s.env.ExecuteWorkflow(TypeWorkflows.EnrollRepurchasingCustomer, params)
//
//	s.True(s.env.IsWorkflowCompleted())
//	s.Empty(cmp.Diff(&msgworkflows.QueryEnrollRepurchasingCustomerResult{
//		Welcomed:                true,
//		RepurchaseSpecification: repurchaseSpec,
//	}, state))
//
//	err := s.env.GetWorkflowError()
//
//	s.isContinuedAsNewWith(err, &msgworkflows.EnrollRepurchasingCustomer{
//		StoreID:    params.StoreID,
//		CustomerID: params.CustomerID,
//		Purchases:  params.Purchases,
//		Welcomed:   true,
//		Attempts:   0,
//	})
//}
//
//func (s *EnrollRepurchasingCustomerTest) isContinuedAsNewWith(err error, expectParams interface{}) {
//	s.T().Helper()
//	s.Error(err)
//	s.True(workflow.IsContinueAsNewError(err), fmt.Sprintf("expected continue as new but got %s", err.Error()))
//	cer := &workflow.ContinueAsNewError{}
//	s.True(errors.As(err, &cer))
//	dc := converter.GetDefaultDataConverter()
//	var actual *msgworkflows.EnrollRepurchasingCustomer
//	err = dc.FromPayloads(cer.Input, &actual)
//	s.NoError(err)
//
//	s.Empty(cmp.Diff(expectParams, actual))
//}
