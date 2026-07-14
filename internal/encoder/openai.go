package encoder

import (
	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/types"
)

type OpenAIEncoder struct{}

func NewOpenAIEncoder() *OpenAIEncoder {
	return &OpenAIEncoder{}
}

func (e *OpenAIEncoder) encodeMessage(
	msg types.Message,
) openai.ResponseOutputItem {

	return openai.ResponseOutputItem{
		Type: "message",
		Role: msg.Role.String(),
		Content: []openai.ResponseOutputContent{
			{
				Type: "output_text",
				Text: msg.Content,
			},
		},
	}
}

func (e *OpenAIEncoder) encodeFunctionCall(
	call types.ToolCall,
) openai.ResponseOutputItem {

	return openai.ResponseOutputItem{
		Type:      "function_call",
		CallID:    call.ID,
		Name:      call.Name,
		Arguments: string(call.Arguments),
	}
}

func (e *OpenAIEncoder) Responses(
	resp *types.ChatResponse,
) *openai.ResponsesResponse {

	output := make([]openai.ResponseOutputItem, 0)

	for _, choice := range resp.Choices {

		msg := choice.Message

		if len(msg.ToolCalls) > 0 {
			for _, call := range msg.ToolCalls {
				output = append(output, e.encodeFunctionCall(call))
			}

			continue
		}

		output = append(output, e.encodeMessage(msg))
	}

	return &openai.ResponsesResponse{
		ID:     resp.ID,
		Object: resp.Object,
		Model:  resp.Model,
		Output: output,
	}
}
