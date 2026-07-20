package gemini

import (
	"context"
	"net/http"
	"net/url"

	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Models(ctx context.Context) (*types.ModelsResponse, error) {
	var resp listModelsResponse

	query := url.Values{}
	query.Set("pageSize", "100")

	if err := c.doJSONWithQuery(ctx, http.MethodGet, "models", query, nil, &resp); err != nil {
		return nil, err
	}

	return toAPIModelsResponse(&resp), nil
}
