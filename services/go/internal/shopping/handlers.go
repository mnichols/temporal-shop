package shopping

import (
	"context"
	"fmt"
)

var TypeHandlers *Handlers

type Handlers struct{}

func (h *Handlers) CalculateShoppingCart(ctx context.Context, cmd *commands.CalculateShoppingCartRequest) (*queries.GetCartResponse, error) {
	state := &queries.GetCartResponse{
		CartId:              cmd.CartId,
		ShopperId:           cmd.ShopperId,
		SubtotalCents:       0,
		TaxCents:            0,
		TaxRateBps:          0,
		TotalCents:          0,
		ProductIdToQuantity: make(map[string]int64),
		ProductIdToGame:     make(map[string]*values.Game),
	}

	for gid, qty := range cmd.ProductIdsToQuantity {
		if qty > 0 {
			g, exists := cmd.ProductIdToGame[gid]
			if !exists {
				return nil, fmt.Errorf("failed to find game with product id [%s] in collection", gid)
			}
			state.ProductIdToGame[g.Id] = g
			state.ProductIdToQuantity[g.Id] = qty
		}
	}

	state.TotalCents = 0
	state.SubtotalCents = 0
	if state.TaxRateBps == 0 {
		state.TaxRateBps = cmd.TaxRateBps
	}
	for k, g := range state.ProductIdToGame {
		state.SubtotalCents = state.SubtotalCents + (g.PriceCents * state.ProductIdToQuantity[k])
	}
	state.TotalCents = CalculateTotalCents(state.SubtotalCents, state.TaxRateBps)
	state.TaxCents = state.TotalCents - state.SubtotalCents
	return state, nil
}

func NewHandlers() *Handlers {
	return &Handlers{}
}
