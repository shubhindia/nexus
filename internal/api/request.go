package api

import (
	"encoding/json"
	"io"
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
	defer func() {
		_ = body.Close()
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)
}
