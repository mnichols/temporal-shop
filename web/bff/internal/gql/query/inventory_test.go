package query

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/mocks"
	"testing"
)

type mockShopper struct {
	mock.Mock
}

func (m *mockShopper) Shopper(ctx context.Context, input *model.ShopperInput) (*model.Shopper, error) {
	args := m.Called(ctx, input)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Shopper), args.Error(1)
	}
	return nil, args.Error(1)
}
func Test_Inventory(t *testing.T) {
	id := cuid.New()
	inventoryID := cuid.New()
	email := "me@example.org"
	queryType := orchestrations.QueryName(&queries.GetInventoryRequest{})
	simpleInventory := &queries.GetInventoryResponse{Games: []*values.Game{
		{
			Id: cuid.New(),
		},
		{
			Id: cuid.New(),
		},
	}}

	/*
		QueryWorkflow(ctx context.Context, workflowID string, runID string, queryType string, args ...interface{}) (converter.EncodedValue, error)
	*/
	cases := []struct {
		name      string
		response  *queries.GetInventoryResponse
		id        string
		expect    *model.Inventory
		shopper   *model.Shopper
		queryErr  error
		expectErr error
	}{
		{
			name: "happy path",
			id:   id,
			shopper: &model.Shopper{
				ID:          id,
				Email:       email,
				InventoryID: inventoryID,
			},
			response: simpleInventory,
			expect:   inventoryFromProto(simpleInventory),
		},
		{
			name:      "inventory is missing",
			id:        id,
			response:  nil,
			expect:    nil,
			queryErr:  &serviceerror.NotFound{},
			expectErr: &serviceerror.NotFound{},
			shopper: &model.Shopper{
				ID:          id,
				Email:       email,
				InventoryID: inventoryID,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			A := assert.New(t)
			mc := &mocks.Client{}
			var ev *mockEncodedValue
			if tt.response != nil {
				ev = &mockEncodedValue{value: &queries.GetInventoryResponse{Games: tt.response.Games}}
			}
			ms := &mockShopper{}
			mc.On("QueryWorkflow", mock.Anything, inventoryID, "", queryType).Return(ev, tt.queryErr)
			sut := inventory{temporal: mc, shopper: ms.Shopper}
			ctx := context.Background()

			if tt.shopper != nil {
				var input *model.ShopperInput
				ms.On("Shopper", mock.Anything, input).Return(tt.shopper, nil)
			}
			actual, err := sut.Inventory(ctx, nil)
			if tt.expectErr != nil {
				A.Error(err, tt.expectErr)
			} else {
				A.NoError(err)
				A.Empty(cmp.Diff(tt.expect, actual))
			}

			mc.AssertExpectations(t)
			ms.AssertExpectations(t)
		})
	}
}
