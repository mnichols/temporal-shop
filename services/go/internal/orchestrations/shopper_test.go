package orchestrations

import (
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/suite"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"go.temporal.io/sdk/testsuite"
	"testing"
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

func (s *ShopperTestSuite) Test_StartShopperSession() {
	params := &orchestrations2.StartSessionRequest{
		Id:    cuid.New(),
		Email: cuid.New(),
	}
	s.env.OnWorkflow(TypeOrchestrations.CreateInventory, &orchestrations2.CreateInventoryRequest{Email: params.Email}).Return(nil)
	s.env.ExecuteWorkflow(TypeOrchestrations.Shopper, params)
	s.NoError(s.env.GetWorkflowError())
}

func TestShopper(t *testing.T) {
	suite.Run(t, &ShopperTestSuite{})
}
