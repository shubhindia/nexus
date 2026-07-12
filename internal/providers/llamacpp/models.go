package llamacpp

import (
	"context"
	"net/http"

	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Models(
	ctx context.Context,
) (*types.ModelsResponse, error) {

	var resp modelsResponse

	if err := c.doJSON(
		ctx,
		http.MethodGet,
		"/v1/models",
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	return toAPIModelsResponse(&resp), nil
}
