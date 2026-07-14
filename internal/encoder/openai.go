package encoder

import (
	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/types"
)

type OpenAIEncoder struct{}

func NewOpenAIEncoder() *OpenAIEncoder {
	return &OpenAIEncoder{}
}

func (e *OpenAIEncoder) Responses(
	resp *types.ChatResponse,
) *openai.ResponsesResponse {

	output := make([]openai.ResponseOutput, 0, len(resp.Choices))

	for _, choice := range resp.Choices {

		output = append(output, openai.ResponseOutput{
			Type: "message",
			Role: choice.Message.Role.String(),
			Content: []openai.ResponseOutputContent{
				{
					Type: "output_text",
					Text: choice.Message.Content,
				},
			},
		})
	}

	return &openai.ResponsesResponse{
		ID:     resp.ID,
		Object: resp.Object,
		Model:  resp.Model,
		Output: output,
	}
}
