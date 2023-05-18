package orchestrations

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	inventory "github.com/temporalio/temporal-shop/services/go/api/generated/inventory/v1"
	commands "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/commands/v1"
	orchestrations2 "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/orchestrations/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	values "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
	inventory2 "github.com/temporalio/temporal-shop/services/go/internal/inventory"
	"github.com/temporalio/temporal-shop/services/go/internal/shopping"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
	"time"
)

// CartTestSuite
// https://docs.temporal.io/docs/go/testing/
type CartTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupSuite https://pkg.go.dev/github.com/stretchr/testify/suite#SetupAllSuite
func (s *CartTestSuite) SetupSuite() {

}

// SetupTest https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
// CAREFUL not to put this `env` inside the SetupSuite or else you will
// get interleaved test times between parallel tests (testify runs suite tests in parallel)
func (s *CartTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// BeforeTest https://pkg.go.dev/github.com/stretchr/testify/suite#BeforeTest
func (s *CartTestSuite) BeforeTest(suiteName, testName string) {

}

// AfterTest https://pkg.go.dev/github.com/stretchr/testify/suite#AfterTest
func (s *CartTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}
func (s *CartTestSuite) Test_SetItemsContinuesAsNew() {
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)

	params := &orchestrations2.SetShoppingCartRequest{
		ShopperId: cuid.New(),
		Email:     cuid.New(),
		CartId:    cuid.New(),
	}

	var gamePriceCents int64 = 50
	productIds := []string{"product_1", "product_2"}
	var games []*values.Game
	setItemsCommand := &commands.SetCartItemsRequest{
		CartId:               params.CartId,
		ProductIdsToQuantity: map[string]int64{},
	}
	for i, id := range productIds {
		games = append(games, &values.Game{Id: id, PriceCents: gamePriceCents})
		setItemsCommand.ProductIdsToQuantity[id] = int64(i + 1)
	}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalName(setItemsCommand), setItemsCommand)
	}, time.Second*1)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*3)

	s.env.ExecuteWorkflow(TypeOrchestrations.Cart, params)
	s.True(s.env.IsWorkflowCompleted())
	werr := s.env.GetWorkflowError()
	s.NotNil(werr)
	s.True(workflow.IsContinueAsNewError(werr))
	can := &workflow.ContinueAsNewError{}
	s.True(errors.As(werr, &can))
	s.Equal("Cart", can.WorkflowType.Name)
	nextParams := &orchestrations2.SetShoppingCartRequest{}
	converter := converter.GetDefaultDataConverter()
	s.NoError(converter.FromPayloads(can.Input, nextParams))
	s.Empty(cmp.Diff(nextParams.ProductIdsToQuantity, setItemsCommand.ProductIdsToQuantity))
}

func (s *CartTestSuite) Test_ClearingCartWithSetItemsContinuesAsNew() {
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)

	params := &orchestrations2.SetShoppingCartRequest{
		ShopperId: cuid.New(),
		Email:     cuid.New(),
		CartId:    cuid.New(),
	}

	clearCartCommand := &commands.SetCartItemsRequest{
		CartId: params.CartId,
		// no products specified
	}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalName(clearCartCommand), clearCartCommand)
	}, time.Second*1)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*3)

	s.env.ExecuteWorkflow(TypeOrchestrations.Cart, params)
	s.True(s.env.IsWorkflowCompleted())
	werr := s.env.GetWorkflowError()
	s.NotNil(werr)
	s.True(workflow.IsContinueAsNewError(werr))
	can := &workflow.ContinueAsNewError{}
	s.True(errors.As(werr, &can))
	s.Equal("Cart", can.WorkflowType.Name)
	nextParams := &orchestrations2.SetShoppingCartRequest{}
	converter := converter.GetDefaultDataConverter()
	s.NoError(converter.FromPayloads(can.Input, nextParams))
	s.Empty(nextParams.ProductIdsToQuantity)
	s.Equal(clearCartCommand.CartId, nextParams.CartId)
}

func (s *CartTestSuite) Test_SetShoppingCartHydratesState() {
	s.env.RegisterWorkflow(TypeOrchestrations.Cart)
	params := &orchestrations2.SetShoppingCartRequest{
		ShopperId:            cuid.New(),
		Email:                cuid.New(),
		CartId:               cuid.New(),
		ProductIdsToQuantity: map[string]int64{},
	}

	var gamePriceCents int64 = 50
	var totalQuantity int64 = 0
	productIds := []string{"product_1", "product_2"}
	var games []*values.Game

	for i, id := range productIds {
		games = append(games, &values.Game{Id: id, PriceCents: gamePriceCents})
		params.ProductIdsToQuantity[id] = int64(i + 1)
		totalQuantity = totalQuantity + params.ProductIdsToQuantity[id]
	}
	var queryResult converter.EncodedValue
	var queryErr error

	s.env.OnActivity(inventory2.TypeHandlers.GetGames, mock.Anything, mock.MatchedBy(func(cmd *inventory.GetGamesRequest) bool {
		if len(productIds) != len(cmd.IncludeProductIds) {
			return false
		}
		return cmp.Diff(productIds, cmd.IncludeProductIds, cmpopts.SortSlices(func(x, y string) bool {
			return x < y
		})) == ""
	})).Return(&inventory.GetGamesResponse{Games: games}, nil)

	s.env.RegisterDelayedCallback(func() {
		q := &queries.GetCartRequest{CartId: params.CartId}
		queryResult, queryErr = s.env.QueryWorkflow(QueryName(q), q)
	}, time.Second*2)
	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Second*3)

	s.env.ExecuteWorkflow(TypeOrchestrations.Cart, params)
	werr := s.env.GetWorkflowError()
	s.True(s.env.IsWorkflowCompleted())
	// this next assertion fails
	s.True(temporal.IsCanceledError(werr))

	s.NoError(queryErr)
	s.NotNil(queryResult)
	s.True(queryResult.HasValue())
	actual := &queries.GetCartResponse{}
	s.NoError(queryResult.Get(actual))
	expect := &queries.GetCartResponse{CartId: params.CartId, ShopperId: params.ShopperId,
		SubtotalCents:       totalQuantity * gamePriceCents,
		TaxRateBps:          shopping.DefaultTaxRateBPS,
		ProductIdToGame:     map[string]*values.Game{},
		ProductIdToQuantity: map[string]int64{},
	}
	expect.TotalCents = shopping.CalculateTotalCents(expect.SubtotalCents, expect.TaxRateBps)
	expect.TaxCents = shopping.CalculateTaxCents(expect.SubtotalCents, expect.TaxRateBps)
	for _, g := range games {
		expect.ProductIdToQuantity[g.Id] = params.ProductIdsToQuantity[g.Id]
		expect.ProductIdToGame[g.Id] = g
	}

	s.Empty(cmp.Diff(expect, actual, protocmp.Transform()))
}

func TestCart(t *testing.T) {
	suite.Run(t, &CartTestSuite{})
}
