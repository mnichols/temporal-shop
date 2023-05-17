package orchestrations

import (
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	orchestrations2 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/shopping"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
	"testing"
	"time"
)

// ShopperTestSuite
// https://docs.temporal.io/docs/go/testing/
type ShopperTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
func (s *ShopperTestSuite) SetupSuite() {

}

// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
// CAREFUL not to put this `env` inside the SetupSuite or else you will
// get interleaved test times between parallel tests (testify runs suite tests in parallel)
func (s *ShopperTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
func (s *ShopperTestSuite) BeforeTest(suiteName, testName string) {

}

// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
func (s *ShopperTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *ShopperTestSuite) Test_StartShopperSession_SpawnsInventoryAndCart() {
	s.env.RegisterWorkflow(TypeOrchestrations.Inventory)
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)
	params := &orchestrations2.StartShopperRequest{
		ShopperId:       cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		CartId:          cuid.New(),
		DurationSeconds: int64(1 * time.Second),
	}
	s.env.OnWorkflow(
		TypeOrchestrations.Inventory,
		mock.Anything,
		&orchestrations2.AllocateInventoryRequest{
			InventoryId: params.InventoryId,
			Email:       params.Email,
		},
	).Return(nil)
	s.env.OnWorkflow(
		TypeOrchestrations.Cart,
		mock.Anything,
		&orchestrations2.SetShoppingCartRequest{
			CartId:    params.CartId,
			Email:     params.Email,
			ShopperId: params.ShopperId,
		},
	).Return(nil)
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Return(nil)
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.CartId, "").Return(nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*1)
	s.env.ExecuteWorkflow(TypeOrchestrations.Shopper, params)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
func (s *ShopperTestSuite) Test_StartShopperSession_IsRefreshable_And_Cancelable() {
	s.env.RegisterWorkflow(TypeOrchestrations.Inventory)
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)
	params := &orchestrations2.StartShopperRequest{
		ShopperId:       cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		CartId:          cuid.New(),
		DurationSeconds: int64(4 * time.Second),
	}
	//var canceledInventory *workflow.Info
	//var canceledCart *workflow.Info
	var completedInventory *workflow.Info
	var completedCart *workflow.Info

	s.env.OnWorkflow(
		TypeOrchestrations.Inventory,
		mock.Anything,
		&orchestrations2.AllocateInventoryRequest{
			InventoryId: params.InventoryId,
			Email:       params.Email,
		},
	).Return(nil)
	s.env.OnWorkflow(
		TypeOrchestrations.Cart,
		mock.Anything,
		&orchestrations2.SetShoppingCartRequest{
			CartId:    params.CartId,
			ShopperId: params.ShopperId,
			Email:     params.Email,
		},
	).Return(nil)
	// this block will not be called because we are mocking out the child workflow execution
	//s.env.SetOnChildWorkflowCanceledListener(func(info *workflow.Info) {
	//	s.T().Log("SetOnChildWorkflowCanceledListener called")
	//	if info.WorkflowExecution.ID == params.InventoryId {
	//		canceledInventory = info
	//	} else if info.WorkflowExecution.ID == params.CartId {
	//		canceledCart = info
	//	}
	//})

	// this will be called though...seems inconsistent
	s.env.SetOnChildWorkflowCompletedListener(func(workflowInfo *workflow.Info, result converter.EncodedValue, err error) {
		s.T().Log("SetOnChildWorkflowCompletedListener called")
		if workflowInfo.WorkflowExecution.ID == params.InventoryId {
			completedInventory = workflowInfo
		} else if workflowInfo.WorkflowExecution.ID == params.CartId {
			completedCart = workflowInfo
		}
	})

	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Return(nil)
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.CartId, "").Return(nil)

	refreshes := []*commands.RefreshShopperRequest{
		{DurationSeconds: int64(time.Second * 1)},
		{DurationSeconds: 0},
	}
	for _, r := range refreshes {
		s.env.RegisterDelayedCallback(func() {
			s.env.SignalWorkflow(SignalName(r), r)
		}, time.Duration(r.DurationSeconds))
	}

	s.env.ExecuteWorkflow(TypeOrchestrations.Shopper, params)
	s.True(s.env.IsWorkflowCompleted())

	s.NoError(s.env.GetWorkflowError())
	// see note above...these will not be set by the canceled listener
	//s.NotNil(canceledInventory)
	//s.NotNil(canceledCart)
	s.NotNil(completedInventory)
	s.NotNil(completedCart)
}
func (s *ShopperTestSuite) Test_StartShopperSession_ContinuesAsNewAfterThresholdMet() {
	s.env.SetTestTimeout(time.Second * 5)
	s.env.RegisterWorkflow(TypeOrchestrations.Inventory)
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)
	params := &orchestrations2.StartShopperRequest{
		ShopperId:       cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		DurationSeconds: int64(4 * time.Second),
	}

	s.env.OnWorkflow(
		TypeOrchestrations.Inventory,
		mock.Anything,
		&orchestrations2.AllocateInventoryRequest{
			InventoryId: params.InventoryId,
			Email:       params.Email,
		},
	).Return(nil)
	s.env.OnWorkflow(
		TypeOrchestrations.Cart,
		mock.Anything,
		&orchestrations2.SetShoppingCartRequest{
			CartId:    shopping.CartID(params.ShopperId),
			ShopperId: params.ShopperId,
			Email:     params.Email,
		},
	).Return(nil)

	// never cancel the children
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Never()
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.CartId, "").Never()

	s.env.RegisterDelayedCallback(func() {
		// send up to the threshold
		// doing more than 500 or so hangs tests...you've been warned
		for i := 0; i <= ShopperRefreshCountThreshold; i++ {
			refresh := &commands.RefreshShopperRequest{DurationSeconds: int64(time.Second * 1)}
			s.env.SignalWorkflow(SignalName(refresh), refresh)
		}
	}, time.Second*1)

	s.env.ExecuteWorkflow(TypeOrchestrations.Shopper, params)
	s.True(s.env.IsWorkflowCompleted())
	s.True(workflow.IsContinueAsNewError(s.env.GetWorkflowError()))
}

func TestShopper(t *testing.T) {
	suite.Run(t, &ShopperTestSuite{})
}
