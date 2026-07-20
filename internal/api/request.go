package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func DecodeJSON(
	r *http.Request,
	dst any,
) error {

	body, err := RequestBody(r)
	if err != nil {
		return err
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	// TEMPORARY
	slog.Info(
		"http.request.body",
		slog.String("body", string(data)),
	)

	return json.Unmarshal(data, dst)
}
