package llamacpp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	reqBody any,
	respBody any,
) error {
	endpoint, err := url.JoinPath(c.url, path)
	if err != nil {
		return err
	}

	var body *bytes.Reader

	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"llamacpp: %s %s returned %d: %s",
			method,
			path,
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return json.NewDecoder(resp.Body).Decode(respBody)
}
