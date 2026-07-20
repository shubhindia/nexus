package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

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

		var chunk ChatCompletionChunk

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

		toolCalls := chunk.Choices[0].Delta.ToolCalls
		if len(toolCalls) > 0 {
			items := make([]ResponseOutputItem, 0, len(toolCalls))

			for _, toolCall := range toolCalls {
				items = append(items, ResponseOutputItem{
					Type:             "function_call",
					Status:           "completed",
					CallID:           toolCall.ID,
					Name:             toolCall.Function.Name,
					Arguments:        toolCall.Function.Arguments,
					ThoughtSignature: toolCall.Function.ThoughtSignature,
				})
			}

			if err := rw.FunctionCalls(items); err != nil {
				return err
			}
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
