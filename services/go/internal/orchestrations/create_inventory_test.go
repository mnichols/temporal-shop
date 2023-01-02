package orchestrations

import (
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	inventory2 "github.com/temporalio/temporal-shop/api/inventory/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

// CreateInventoryTestSuite
// https://docs.temporal.io/docs/go/testing/
type CreateInventoryTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
func (s *CreateInventoryTestSuite) SetupSuite() {

}

// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
// CAREFUL not to put this `env` inside the SetupSuite or else you will
// get interleaved test times between parallel tests (testify runs suite tests in parallel)
func (s *CreateInventoryTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
func (s *CreateInventoryTestSuite) BeforeTest(suiteName, testName string) {

}

// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
func (s *CreateInventoryTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *CreateInventoryTestSuite) Test_LoadsGamesAndExposesThem() {
	games := []*inventory2.Game{
		{Id: cuid.New()},
		{Id: cuid.New()},
	}

	params := &orchestrations.CreateInventoryRequest{
		Id:    cuid.New(),
		Email: cuid.New(),
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
	s.env.ExecuteWorkflow(TypeOrchestrations.CreateInventory, params)
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
	suite.Run(t, &CreateInventoryTestSuite{})
}
