package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientConnection(ctx context.Context, cfg *Config) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
