package provider

import (
	"context"

	"github.com/shubhindia/nexus/internal/types"
)

type Provider interface {
	Models(ctx context.Context) (*types.ModelsResponse, error)

	Chat(
		ctx context.Context,
		req *types.ChatRequest,
	) (*ChatResult, error)
	Capabilities() Capabilities
}
