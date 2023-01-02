package orchestrations

import (
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
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

func (s *ShopperTestSuite) Test_StartShopperSession_CreatesInventory() {
	s.env.RegisterWorkflow(TypeOrchestrations.CreateInventory)
	params := &orchestrations2.StartShopperRequest{
		Id:              cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		DurationSeconds: int64(1 * time.Second),
	}
	s.env.OnWorkflow(
		TypeOrchestrations.CreateInventory,
		mock.Anything,
		&orchestrations2.CreateInventoryRequest{
			Id:    params.InventoryId,
			Email: params.Email,
		},
	).Return(nil)
	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Return(nil)
	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*1)
	s.env.ExecuteWorkflow(TypeOrchestrations.Shopper, params)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
func (s *ShopperTestSuite) Test_StartShopperSession_IsRefreshable_And_Cancelable() {
	s.env.RegisterWorkflow(TypeOrchestrations.CreateInventory)
	params := &orchestrations2.StartShopperRequest{
		Id:              cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		DurationSeconds: int64(4 * time.Second),
	}
	var canceledInventory *workflow.Info
	var completedInventory *workflow.Info

	s.env.OnWorkflow(
		TypeOrchestrations.CreateInventory,
		mock.Anything,
		&orchestrations2.CreateInventoryRequest{
			Id:    params.InventoryId,
			Email: params.Email,
		},
	).Return(nil)
	// this will not be called because we are mocking out the child workflow execution
	s.env.SetOnChildWorkflowCanceledListener(func(info *workflow.Info) {
		s.T().Log("SetOnChildWorkflowCanceledListener called")
		canceledInventory = info
	})
	// this will be called though...seems inconsistent
	s.env.SetOnChildWorkflowCompletedListener(func(workflowInfo *workflow.Info, result converter.EncodedValue, err error) {
		s.T().Log("SetOnChildWorkflowCompletedListener called")
		completedInventory = workflowInfo
	})

	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Return(nil)

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
	s.NotNil(canceledInventory)
	s.NotNil(completedInventory)
}
func (s *ShopperTestSuite) Test_StartShopperSession_ContinuesAsNewAfterThresholdMet() {
	s.env.SetTestTimeout(time.Second * 5)
	s.env.RegisterWorkflow(TypeOrchestrations.CreateInventory)
	params := &orchestrations2.StartShopperRequest{
		Id:              cuid.New(),
		Email:           cuid.New(),
		InventoryId:     cuid.New(),
		DurationSeconds: int64(4 * time.Second),
	}

	s.env.OnWorkflow(
		TypeOrchestrations.CreateInventory,
		mock.Anything,
		&orchestrations2.CreateInventoryRequest{
			Id:    params.InventoryId,
			Email: params.Email,
		},
	).Return(nil)

	s.env.OnRequestCancelExternalWorkflow(mock.Anything, params.InventoryId, "").Never() //Return(nil)

	s.env.RegisterDelayedCallback(func() {
		// send up to the threshold
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
