package llamacpp

import (
	"context"
	"errors"

	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(
	ctx context.Context,
	req *types.ChatRequest,
) (*types.ChatResponse, error) {
	return nil, errors.New("not implemented")
}
