package provider

import (
	"context"

	"github.com/shubhindia/nexus/internal/types"
)

type Provider interface {
	Chat(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error)
	Models(ctx context.Context) (*types.ModelsResponse, error)
}
