package providers

import (
	"context"

	"github.com/shubhindia/nexus/internal/types"
)

type Provider interface {
	Name() string

	Chat(
		ctx context.Context,
		req *types.ChatRequest,
	) (*types.ChatResponse, error)
}
