package llamacpp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/logger"
	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Chat(
	ctx context.Context,
	req *types.ChatRequest,
) (*provider.ChatResult, error) {
	log := logger.FromContext(ctx)

	log.Info(
		"provider.request",
		slog.String("provider", "llamacpp"),
		slog.String("operation", "chat"),
	)
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
	log.Info(
		"provider.response",
		slog.String("provider", "llamacpp"),
	)

	return &provider.ChatResult{
		Response: toAPIChatResponse(&resp),
	}, nil
}
