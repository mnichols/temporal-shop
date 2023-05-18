package inventory

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	inventory "github.com/temporalio/temporal-shop/services/go/api/generated/inventory/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

func TestInventoryService_GetGames(t *testing.T) {

	allGames, err := getAllGames(gamesV1JSON)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx     context.Context
		request *inventory.GetGamesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *inventory.GetGamesResponse
		wantErr error
	}{
		{

			name:    "happy path",
			args:    args{ctx: context.Background(), request: &inventory.GetGamesRequest{Version: "1"}},
			want:    &inventory.GetGamesResponse{Games: allGames},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			A := assert.New(t)
			i := &InventoryService{}
			actual, err := i.GetGames(tt.args.ctx, tt.args.request)
			if tt.wantErr != nil {
				A.EqualError(err, tt.wantErr.Error())
			} else {
				A.NoError(err)
				A.Empty(cmp.Diff(tt.want, actual, protocmp.Transform()))
			}

		})
	}
}
