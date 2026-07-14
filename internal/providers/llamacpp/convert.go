package llamacpp

import (
	"path/filepath"

	"github.com/shubhindia/nexus/internal/types"
)

func toAPIModelsResponse(resp *modelsResponse) *types.ModelsResponse {
	r := &types.ModelsResponse{
		Object: resp.Object,
		Data:   make([]types.Model, len(resp.Data)),
	}

	for i, model := range resp.Data {
		r.Data[i] = types.Model{
			ID:      filepath.Base(model.ID),
			Object:  model.Object,
			OwnedBy: model.OwnedBy,
		}
	}

	return r
}

func toProviderChatRequest(req *types.ChatRequest) *chatRequest {
	r := &chatRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
		Messages:    make([]chatMessage, len(req.Messages)),
	}

	for i, msg := range req.Messages {
		r.Messages[i] = chatMessage{
			Role:    msg.Role.String(),
			Content: msg.Content,
		}
	}

	return r
}

func toAPIChatResponse(resp *chatResponse) *types.ChatResponse {
	r := &types.ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   filepath.Base(resp.Model),
		Choices: make([]types.Choice, len(resp.Choices)),
	}

	if resp.Usage != nil {
		r.Usage = &types.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}

	for i, choice := range resp.Choices {
		r.Choices[i] = types.Choice{
			Index: choice.Index,
			Message: types.Message{
				Role:    types.ParseRole(choice.Message.Role),
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}

	return r
}
