package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/shubhindia/nexus/internal/types"
)

type apiErrorResponse struct {
	Error *geminiAPIError `json:"error"`
}

type geminiAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	reqBody any,
	respBody any,
) error {
	resp, err := c.send(ctx, method, path, nil, reqBody)
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
		return fmt.Errorf("decode %s response: %w\nbody:\n%s", path, err, string(body))
	}

	return nil
}

func (c *Client) send(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	reqBody any,
) (*http.Response, error) {
	endpoint, err := url.Parse(c.url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	endpoint.Path, err = url.JoinPath(endpoint.Path, path)
	if err != nil {
		return nil, fmt.Errorf("join url path: %w", err)
	}

	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		defer func() {
			_ = resp.Body.Close()
		}()

		body, _ := io.ReadAll(resp.Body)
		return nil, toAPIError(resp.StatusCode, method, path, body)
	}

	return resp, nil
}

func toAPIError(status int, method string, path string, body []byte) error {
	message := fmt.Sprintf("%s %s returned %d", method, path, status)
	code := mapErrorCode(status, "")

	var resp apiErrorResponse
	if err := json.Unmarshal(body, &resp); err == nil && resp.Error != nil {
		if resp.Error.Message != "" {
			message = resp.Error.Message
		}

		code = mapErrorCode(status, resp.Error.Status)
	} else if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
		message = fmt.Sprintf("%s:\n%s", message, trimmed)
	}

	return &types.APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func mapErrorCode(status int, providerStatus string) string {
	switch status {
	case http.StatusBadRequest:
		return "invalid_request"
	case http.StatusUnauthorized, http.StatusForbidden:
		return "authentication_error"
	case http.StatusTooManyRequests:
		return "rate_limit_exceeded"
	case http.StatusServiceUnavailable:
		if providerStatus == "UNAVAILABLE" {
			return "service_unavailable"
		}

		return "upstream_unavailable"
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return "upstream_unavailable"
	default:
		return "upstream_error"
	}
}
