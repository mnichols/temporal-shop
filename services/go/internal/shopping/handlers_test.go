package shopping

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/assert"
	commands "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/commands/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	values "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

type cartTestCase struct {
	description string
	expectState *queries.GetCartResponse
	expectError error
	cmd         *commands.CalculateShoppingCartRequest
}

func Test_CalculateShoppingCart(t *testing.T) {
	args := ShoppingCartArgs{
		ID:         cuid.New(),
		ShopperID:  cuid.New(),
		TaxRateBPS: DefaultTaxRateBPS,
	}
	games := []*values.Game{}
	gamesMap := make(map[string]*values.Game)

	for i := 0; i < 20; i++ {
		game := &values.Game{
			Id:         fmt.Sprintf("%d", i),
			Title:      cuid.New(),
			Category:   fmt.Sprintf("category_%d", i),
			ImageUrl:   "",
			PriceCents: int64(50 + i),
		}
		games = append(games, game)
		gamesMap[game.Id] = game
	}

	cmdWithZeroQuantity := &commands.CalculateShoppingCartRequest{
		CartId:     args.ID,
		TaxRateBps: DefaultTaxRateBPS,
		ShopperId:  args.ShopperID,
		ProductIdsToQuantity: map[string]int64{
			games[0].Id: 1,
			games[1].Id: 0,
			games[2].Id: 2,
		},
		ProductIdToGame: map[string]*values.Game{
			games[0].Id: gamesMap[games[0].Id],
			games[1].Id: gamesMap[games[1].Id],
			games[2].Id: gamesMap[games[2].Id],
		},
	}
	cmdWithMissingGame := &commands.CalculateShoppingCartRequest{
		CartId:     args.ID,
		TaxRateBps: DefaultTaxRateBPS,
		ShopperId:  args.ShopperID,
		ProductIdsToQuantity: map[string]int64{
			games[0].Id: 1,
			games[1].Id: 0,
			games[2].Id: 2,
		},
		ProductIdToGame: map[string]*values.Game{
			games[0].Id: gamesMap[games[0].Id],
			games[1].Id: gamesMap[games[1].Id],
		},
	}

	//cmdWithLotsOfThings := &commands.SetCartItemsRequest{
	//	CartId: args.ID,
	//	ProductIdsToQuantity: map[string]int64{
	//		games[4].Id: 1,
	//		games[8].Id: 3,
	//		games[3].Id: 2,
	//		games[9].Id: 23,
	//	},
	//}
	cases := []cartTestCase{{
		description: "zero quantities are pruned",
		expectState: &queries.GetCartResponse{
			CartId:        args.ID,
			ShopperId:     args.ShopperID,
			SubtotalCents: 154,
			TaxCents:      6,
			TaxRateBps:    args.TaxRateBPS,
			TotalCents:    160,
			ProductIdToQuantity: map[string]int64{
				games[0].Id: 1,
				games[2].Id: 2,
			},
			ProductIdToGame: map[string]*values.Game{
				games[0].Id: gamesMap[games[0].Id],
				games[2].Id: gamesMap[games[2].Id],
			},
		},
		cmd: cmdWithZeroQuantity,
	}, {
		description: "zero quantities are pruned regardless of query",
		expectState: &queries.GetCartResponse{
			CartId:        args.ID,
			ShopperId:     args.ShopperID,
			SubtotalCents: 154,
			TaxCents:      6,
			TaxRateBps:    args.TaxRateBPS,
			TotalCents:    160,
			ProductIdToQuantity: map[string]int64{
				games[0].Id: 1,
				games[2].Id: 2,
			},
			ProductIdToGame: map[string]*values.Game{
				games[0].Id: gamesMap[games[0].Id],
				games[2].Id: gamesMap[games[2].Id],
			},
		},
		cmd: cmdWithZeroQuantity,
	}, {
		description: "quantity provided for unknown game errors",
		expectState: nil,
		expectError: fmt.Errorf("failed to find game"),
		cmd:         cmdWithMissingGame,
	},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			A := assert.New(t)
			sut := NewHandlers()

			actual, err := sut.CalculateShoppingCart(context.Background(), tc.cmd)
			if tc.expectError != nil {
				A.ErrorContains(err, tc.expectError.Error())
			} else {
				A.NoError(err)
				A.Empty(cmp.Diff(tc.expectState, actual, protocmp.Transform()))
			}
		})
	}
}
