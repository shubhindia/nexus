package llamacpp

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/shubhindia/nexus/internal/types"
)

func (c *Client) Models(
	ctx context.Context,
) (*types.ModelsResponse, error) {

	var resp modelsResponse

	if err := c.do(
		ctx,
		http.MethodGet,
		"/v1/models",
		nil,
		&resp,
	); err != nil {
		return nil, err
	}

	models := &types.ModelsResponse{
		Object: resp.Object,
		Data:   make([]types.Model, len(resp.Data)),
	}

	for i, model := range resp.Data {
		models.Data[i] = types.Model{
			ID:      filepath.Base(model.ID),
			Object:  model.Object,
			OwnedBy: model.OwnedBy,
		}
	}

	return models, nil
}
