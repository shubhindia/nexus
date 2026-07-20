package llamacpp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(
	ctx context.Context,
	req *types.ChatRequest,
) (*provider.ChatResult, error) {

	providerReq := compileChatRequest(req)
	body, _ := json.MarshalIndent(providerReq, "", "  ")

	slog.Info(
		"llamacpp.request",
		slog.Int("tools", len(providerReq.Tools)),
		slog.String("body", string(body)),
	)

	if req.Stream {
		stream, err := c.doStream(
			ctx,
			http.MethodPost,
			"/v1/chat/completions",
			providerReq,
		)
		if err != nil {
			return nil, err
		}

		return &provider.ChatResult{
			Stream: stream,
		}, nil
	}

	var resp chatResponse

	if err := c.doJSON(
		ctx,
		http.MethodPost,
		"/v1/chat/completions",
		providerReq,
		&resp,
	); err != nil {
		return nil, err
	}

	return &provider.ChatResult{
		Response: toAPIChatResponse(&resp),
	}, nil
}

func (c *Client) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Streaming: true,
		Tools:     false,
		Vision:    false,
		JSONMode:  false,
	}
}
