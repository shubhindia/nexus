package gemini

import (
	"context"
	"net/http"

	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(ctx context.Context, req *types.ChatRequest) (*provider.ChatResult, error) {
	providerReq := compileChatRequest(req)
	modelPath := modelPath(req.Model)

	if req.Stream {
		stream, err := c.streamChat(
			ctx,
			modelPath,
			providerReq,
		)
		if err != nil {
			return nil, err
		}

		return &provider.ChatResult{Stream: stream}, nil
	}

	var resp generateContentResponse

	if err := c.doJSON(
		ctx,
		http.MethodPost,
		modelPath+":generateContent",
		providerReq,
		&resp,
	); err != nil {
		return nil, err
	}

	return &provider.ChatResult{Response: toAPIChatResponse(req.Model, &resp)}, nil
}

func (c *Client) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Streaming: true,
		Tools:     true,
		Vision:    false,
		JSONMode:  false,
	}
}
