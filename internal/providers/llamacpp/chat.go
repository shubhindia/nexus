package llamacpp

import (
	"context"

	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(
	ctx context.Context,
	req *types.ChatRequest,
) (*provider.ChatResult, error) {

	return &provider.ChatResult{
		Response: &types.ChatResponse{
			ID:     "chatcmpl-test",
			Object: "chat.completion",
			Model:  req.Model,
			Choices: []types.Choice{
				{
					Index: 0,
					Message: types.Message{
						Role:    "assistant",
						Content: "Hello from Nexus!",
					},
					FinishReason: "stop",
				},
			},
		},
	}, nil
}
