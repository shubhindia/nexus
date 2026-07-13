package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type chatCompletionChunk struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`

	Choices []struct {
		Index int `json:"index"`

		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`

		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

func ChatCompletionToResponses(
	w io.Writer,
	r io.Reader,
) error {

	scanner := bufio.NewScanner(r)

	rw := NewResponseWriter(w)

	created := false

	for scanner.Scan() {

		line := scanner.Text()

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			return rw.Completed()
		}

		var chunk chatCompletionChunk

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf(
				"failed to decode SSE chunk:\n\n%s\n\nerror: %w",
				data,
				err,
			)
		}

		if !created {
			if err := rw.Created(
				chunk.ID,
				chunk.Model,
			); err != nil {
				return err
			}

			created = true
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta.Content
		if delta == "" {
			continue
		}

		if err := rw.Delta(delta); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
