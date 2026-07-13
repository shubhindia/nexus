package api

import "github.com/shubhindia/nexus/internal/types"

func toChatRequest(req *types.ResponsesRequest) *types.ChatRequest {
	messages := make([]types.Message, 0)

	if req.Instructions != "" {
		messages = append(messages, types.Message{
			Role:    "system",
			Content: req.Instructions,
		})
	}

	for _, input := range req.Input {
		var text string

		for _, content := range input.Content {
			if content.Type == "input_text" {
				text += content.Text
			}
		}

		if text == "" {
			continue
		}

		messages = append(messages, types.Message{
			Role:    input.Role,
			Content: text,
		})
	}

	return &types.ChatRequest{
		Model:    req.Model,
		Messages: messages,
		Stream:   req.Stream,
	}
}

func toResponsesResponse(chat *types.ChatResponse) *types.ResponsesResponse {
	text := ""

	if len(chat.Choices) > 0 {
		text = chat.Choices[0].Message.Content
	}

	return &types.ResponsesResponse{
		ID:     chat.ID,
		Object: "response",
		Model:  chat.Model,
		Output: []types.ResponseOutput{
			{
				Type: "message",
				Role: "assistant",
				Content: []types.ResponseOutputContent{
					{
						Type: "output_text",
						Text: text,
					},
				},
			},
		},
	}
}
