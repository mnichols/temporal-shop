package orchestrations

import (
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	inventory2 "github.com/temporalio/temporal-shop/services/go/api/generated/inventory/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/orchestrations/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	values "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

// InventoryTestSuite
// https://docs.temporal.io/docs/go/testing/
type InventoryTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
func (s *InventoryTestSuite) SetupSuite() {

}

// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
// CAREFUL not to put this `env` inside the SetupSuite or else you will
// get interleaved test times between parallel tests (testify runs suite tests in parallel)
func (s *InventoryTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
func (s *InventoryTestSuite) BeforeTest(suiteName, testName string) {

}

// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
func (s *InventoryTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *InventoryTestSuite) Test_LoadsGamesAndExposesThem() {
	games := []*values.Game{
		{Id: cuid.New()},
		{Id: cuid.New()},
	}

	params := &orchestrations2.AllocateInventoryRequest{
		InventoryId: cuid.New(),
		Email:       cuid.New(),
	}

	s.env.RegisterActivity(inventory.TypeHandlers)
	s.env.OnActivity(
		inventory.TypeHandlers.GetGames,
		mock.Anything,
		&inventory2.GetGamesRequest{Version: "1"},
	).Return(
		&inventory2.GetGamesResponse{Games: games},
		nil,
	)
	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*1)
	s.env.ExecuteWorkflow(TypeOrchestrations.Inventory, params)
	s.True(s.env.IsWorkflowCompleted())
	s.Nil(s.env.GetWorkflowError())
	var getInventory *queries.GetInventoryResponse
	val, err := s.env.QueryWorkflow(QueryName(&queries.GetInventoryRequest{}))
	s.NoError(err)
	s.True(val.HasValue())
	s.NoError(val.Get(&getInventory))
	s.Equal(len(games), len(getInventory.Games))
}

func TestCreateInventory(t *testing.T) {
	suite.Run(t, &InventoryTestSuite{})
}
