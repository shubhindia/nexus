package api

import (
	"encoding/json"
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

	decoder := json.NewDecoder(body)

	return decoder.Decode(dst)
}
