package llamacpp

import (
	"context"
	"net/http"

	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(
	ctx context.Context,
	req *types.ChatRequest,
) (*provider.ChatResult, error) {

	if req.Stream {
		stream, err := c.doStream(
			ctx,
			http.MethodPost,
			"/v1/chat/completions",
			toProviderChatRequest(req),
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
		toProviderChatRequest(req),
		&resp,
	); err != nil {
		return nil, err
	}

	return &provider.ChatResult{
		Response: toAPIChatResponse(&resp),
	}, nil
}
