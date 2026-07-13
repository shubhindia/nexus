package api

import (
	"github.com/shubhindia/nexus/internal/prompt"
	"github.com/shubhindia/nexus/internal/types"
)

var codexCompiler = prompt.NewCodexCompiler()

func toChatRequest(
	req *types.ResponsesRequest,
) *types.ChatRequest {

	return &types.ChatRequest{
		Model:    req.Model,
		Messages: codexCompiler.Compile(req),
		Stream:   req.Stream,
	}
}

func toResponsesResponse(
	chat *types.ChatResponse,
) *types.ResponsesResponse {

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
