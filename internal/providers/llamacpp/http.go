package llamacpp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return fmt.Errorf(
			"decode %s response: %w\nbody:\n%s",
			path,
			err,
			string(body),
		)
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
		slog.Info(
			"llamacpp.response",
			slog.Int("status", resp.StatusCode),
			slog.String("body", string(body)),
		)

		return nil, fmt.Errorf(
			"%s %s returned %d:\n%s",
			method,
			path,
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return resp, nil
}
