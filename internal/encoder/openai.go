package encoder

import (
	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/toolstate"
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
		ID:     openai.NewResponseItemID(),
		Type:   "message",
		Status: "completed",
		Role:   msg.Role.String(),
		Content: []openai.ResponseOutputContent{
			{
				Type:        "output_text",
				Text:        msg.Content,
				Annotations: []any{},
			},
		},
	}
}

func (e *OpenAIEncoder) encodeFunctionCall(
	call types.ToolCall,
) openai.ResponseOutputItem {
	toolstate.StoreThoughtSignature(call.ID, call.ThoughtSignature)

	return openai.ResponseOutputItem{
		ID:               openai.NewResponseItemID(),
		Type:             "function_call",
		Status:           "completed",
		CallID:           call.ID,
		Name:             call.Name,
		Arguments:        string(call.Arguments),
		ThoughtSignature: call.ThoughtSignature,
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

	var usage *openai.ResponseUsage
	if resp.Usage != nil {
		usage = &openai.ResponseUsage{
			InputTokens: resp.Usage.PromptTokens,
			InputTokensDetails: openai.ResponseInputTokenDetails{
				CachedTokens: 0,
			},
			OutputTokens: resp.Usage.CompletionTokens,
			OutputTokensDetails: openai.ResponseOutputTokenDetails{
				ReasoningTokens: 0,
			},
			TotalTokens: resp.Usage.TotalTokens,
		}
	}

	return &openai.ResponsesResponse{
		ID:                 resp.ID,
		Object:             "response",
		CreatedAt:          resp.Created,
		Status:             "completed",
		Background:         false,
		Error:              nil,
		IncompleteDetails:  nil,
		Instructions:       nil,
		MaxOutputTokens:    nil,
		Model:              resp.Model,
		Output:             output,
		ParallelToolCalls:  true,
		PreviousResponseID: nil,
		Reasoning: openai.ResponseReasoning{
			Effort:  nil,
			Summary: []any{},
		},
		Store: true,
		Text: openai.ResponseText{
			Format: openai.ResponseTextFormat{Type: "text"},
		},
		ToolChoice:  "auto",
		Tools:       []any{},
		TopP:        nil,
		Temperature: nil,
		Truncation:  "disabled",
		Usage:       usage,
		User:        nil,
		Metadata:    map[string]any{},
	}
}
