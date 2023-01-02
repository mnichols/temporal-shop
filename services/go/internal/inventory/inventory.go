package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/temporalio/temporal-shop/api/inventory/v1"
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
	var games []*inventory.Game

	data := []byte(gamesV1JSON)
	err := json.Unmarshal(data, &games)
	if err != nil {
		return nil, err
	}
	return &inventory.GetGamesResponse{
		Games: games,
	}, nil
}
func (i *InventoryService) Register(srv *grpc.Server) {
	inventory.RegisterInventoryServiceServer(srv, i)
}
