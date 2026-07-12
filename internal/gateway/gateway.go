package gateway

import (
	"context"

	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

type Gateway struct {
	provider provider.Provider
}

func New(p provider.Provider) *Gateway {
	return &Gateway{
		provider: p,
	}
}

func (g *Gateway) Chat(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	return g.provider.Chat(ctx, req)
}

func (g *Gateway) Models(ctx context.Context) (*types.ModelsResponse, error) {
	return g.provider.Models(ctx)
}
