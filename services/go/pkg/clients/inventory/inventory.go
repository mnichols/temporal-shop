package inventory

import (
	"context"
	"google.golang.org/grpc"
)

type Client struct {
	remote inventory.InventoryServiceClient
	conn   *grpc.ClientConn
}

func (c *Client) GetGames(ctx context.Context, in *inventory.GetGamesRequest, opts ...grpc.CallOption) (*inventory.GetGamesResponse, error) {
	return c.remote.GetGames(ctx, in)
}

func NewClient(ctx context.Context, conn *grpc.ClientConn) (*Client, error) {
	return &Client{
		remote: inventory.NewInventoryServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
