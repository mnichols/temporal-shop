package inventory

import (
	"context"
	inventory "github.com/temporalio/temporal-shop/services/go/api/generated/inventory/v1"
)

var TypeHandlers *Handlers

type Handlers struct {
	inner inventory.InventoryServiceClient
}

func NewHandlers(c inventory.InventoryServiceClient) *Handlers {
	return &Handlers{c}
}
func (h *Handlers) GetGames(ctx context.Context, in *inventory.GetGamesRequest) (*inventory.GetGamesResponse, error) {
	return h.inner.GetGames(ctx, in)
}
