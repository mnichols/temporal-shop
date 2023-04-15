package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/temporalio/temporal-shop/api/inventory/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	"google.golang.org/grpc"
)

func InventorySessionID(sessionID string) string {
	return fmt.Sprintf("inv_%s", sessionID)
}
func NewInventoryService() (*InventoryService, error) {
	return &InventoryService{}, nil
}

type InventoryService struct {
	inventory.UnimplementedInventoryServiceServer
}

func (i *InventoryService) GetGames(ctx context.Context, request *inventory.GetGamesRequest) (*inventory.GetGamesResponse, error) {
	var games []*values.Game

	data := []byte(gamesV1JSON)
	err := json.Unmarshal(data, &games)
	if err != nil {
		return nil, err
	}
	var include []*values.Game
	if len(request.IncludeProductIds) == 0 {
		include = games
	} else {
		incMap := map[string]struct{}{}
		for _, id := range request.IncludeProductIds {
			incMap[id] = struct{}{}
		}
		for _, g := range games {
			if _, exists := incMap[g.Id]; exists {
				include = append(include, g)
			}
		}
	}
	return &inventory.GetGamesResponse{
		Games: include,
	}, nil
}
func (i *InventoryService) Register(srv *grpc.Server) {
	inventory.RegisterInventoryServiceServer(srv, i)
}
