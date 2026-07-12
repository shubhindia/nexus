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
	resp, err := c.send(
		ctx,
		method,
		path,
		reqBody,
	)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *Client) doStream(
	ctx context.Context,
	method string,
	path string,
	reqBody any,
) (io.ReadCloser, error) {
	resp, err := c.send(
		ctx,
		method,
		path,
		reqBody,
	)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *Client) send(
	ctx context.Context,
	method string,
	path string,
	reqBody any,
) (*http.Response, error) {
	endpoint, err := url.JoinPath(c.url, path)
	if err != nil {
		return nil, fmt.Errorf("join url: %w", err)
	}

	var body io.Reader

	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer func() {
			_ = resp.Body.Close()
		}()

		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"%s %s returned %d: %s",
			method,
			path,
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return resp, nil
}
