package gemini

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (c *Client) doJSONWithQuery(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	reqBody any,
	respBody any,
) error {
	resp, err := c.send(ctx, method, path, query, reqBody)
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

func (c *Client) streamChat(
	ctx context.Context,
	path string,
	reqBody any,
) (io.ReadCloser, error) {
	query := url.Values{}
	query.Set("alt", "sse")

	resp, err := c.send(ctx, http.MethodPost, path+":streamGenerateContent", query, reqBody)
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()

	go func() {
		defer func() {
			_ = resp.Body.Close()
		}()

		err := translateStream(pw, resp.Body)
		_ = pw.CloseWithError(err)
	}()

	return pr, nil
}

func translateStream(w *io.PipeWriter, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	chunkID := "chatcmpl-" + uuid.NewString()
	created := time.Now().Unix()
	started := false
	finished := false
	model := "gemini"

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}

		var resp generateContentResponse
		if err := json.Unmarshal([]byte(payload), &resp); err != nil {
			return fmt.Errorf("decode Gemini SSE chunk: %w", err)
		}

		if len(resp.Candidates) == 0 {
			continue
		}

		candidate := resp.Candidates[0]
		model = responseModel(model, resp.ModelVersion)
		content, toolCalls := candidateMessage(candidate.Content)

		if !started {
			if err := writeChatChunk(w, map[string]any{
				"id":      chunkID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   model,
				"choices": []map[string]any{{
					"index": 0,
					"delta": map[string]any{"role": "assistant"},
				}},
			}); err != nil {
				return err
			}
			started = true
		}

		if content != "" {
			if err := writeChatChunk(w, map[string]any{
				"id":      chunkID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   model,
				"choices": []map[string]any{{
					"index": 0,
					"delta": map[string]any{"content": content},
				}},
			}); err != nil {
				return err
			}
		}

		if len(toolCalls) > 0 {
			deltaToolCalls := make([]map[string]any, 0, len(toolCalls))

			for index, call := range toolCalls {
				deltaToolCalls = append(deltaToolCalls, map[string]any{
					"index": index,
					"id":    call.ID,
					"type":  "function",
					"function": map[string]any{
						"name":              call.Name,
						"arguments":         string(call.Arguments),
						"thought_signature": call.ThoughtSignature,
					},
				})
			}

			if err := writeChatChunk(w, map[string]any{
				"id":      chunkID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   model,
				"choices": []map[string]any{{
					"index": 0,
					"delta": map[string]any{"tool_calls": deltaToolCalls},
				}},
			}); err != nil {
				return err
			}
		}

		if candidate.FinishReason != "" {
			reason := finishReason(candidate.FinishReason, len(toolCalls) > 0)
			if err := writeChatChunk(w, map[string]any{
				"id":      chunkID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   model,
				"choices": []map[string]any{{
					"index":         0,
					"delta":         map[string]any{},
					"finish_reason": reason,
				}},
			}); err != nil {
				return err
			}
			finished = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !finished {
		if err := writeChatChunk(w, map[string]any{
			"id":      chunkID,
			"object":  "chat.completion.chunk",
			"created": created,
			"model":   model,
			"choices": []map[string]any{{
				"index":         0,
				"delta":         map[string]any{},
				"finish_reason": "stop",
			}},
		}); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "data: [DONE]\n\n")
	return err
}

func writeChatChunk(w io.Writer, payload any) error {
	return streamWriteSSE(w, payload)
}

func streamWriteSSE(w io.Writer, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", data)
	return err
}

func modelPath(model string) string {
	model = trimModelName(model)
	return "models/" + model
}
