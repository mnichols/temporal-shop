package query

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/temporal-shop/api/temporal_shop/queries/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/middleware"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/mocks"
	"testing"
)

func Test_Shopper(t *testing.T) {
	id := cuid.New()
	email := "me@example.org"
	queryType := orchestrations.QueryName(&queries.GetShopperRequest{})

	/*
		QueryWorkflow(ctx context.Context, workflowID string, runID string, queryType string, args ...interface{}) (converter.EncodedValue, error)
	*/
	cases := []struct {
		name      string
		response  *queries.GetShopperResponse
		id        string
		expect    *model.Shopper
		auth      *mockAuth
		queryErr  error
		expectErr error
	}{
		{
			name:     "happy path",
			id:       id,
			auth:     &mockAuth{},
			response: &queries.GetShopperResponse{ShopperId: id},
			expect:   &model.Shopper{ID: id, Email: email},
		},
		{
			name:      "shopper has no session",
			id:        id,
			auth:      &mockAuth{},
			response:  nil,
			expect:    nil,
			queryErr:  &serviceerror.NotFound{},
			expectErr: &serviceerror.NotFound{},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			A := assert.New(t)
			mc := &mocks.Client{}
			ev := &mockEncodedValue{value: &queries.GetShopperResponse{ShopperId: id, Email: email}}
			mc.On("QueryWorkflow", mock.Anything, id, "", queryType).Return(ev, tt.queryErr)
			sut := shopper{temporal: mc}
			ctx := context.Background()
			if tt.auth != nil {
				tt.auth.On("SessionID").Return(id)
				ctx = middleware.WithAuth(ctx, tt.auth)
			}
			actual, err := sut.Shopper(ctx, nil)
			if tt.expectErr != nil {
				A.Error(err, tt.expectErr)
			} else {
				A.NoError(err)
				A.Empty(cmp.Diff(tt.expect, actual))
			}

			mc.AssertExpectations(t)
			tt.auth.AssertExpectations(t)
		})
	}
}
