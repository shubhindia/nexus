package gemini

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shubhindia/nexus/internal/types"
)

func toAPIModelsResponse(resp *listModelsResponse) *types.ModelsResponse {
	data := make([]types.Model, 0, len(resp.Models))

	for _, model := range resp.Models {
		if !supportsGenerateContent(model) {
			continue
		}

		data = append(data, types.Model{
			ID:      trimModelName(model.Name),
			Object:  "model",
			OwnedBy: "google",
		})
	}

	return &types.ModelsResponse{
		Object: "list",
		Data:   data,
	}
}

func compileChatRequest(req *types.ChatRequest) *generateContentRequest {
	r := &generateContentRequest{
		Contents:   make([]geminiContent, 0, len(req.Messages)),
		Tools:      convertTools(req.Tools),
		ToolConfig: convertToolChoice(req.ToolChoice),
	}

	if cfg := generationConfig(req); cfg != nil {
		r.GenerationConfig = cfg
	}

	var systemParts []geminiPart

	for _, msg := range req.Messages {
		content := compileMessage(msg)

		switch msg.Role {
		case types.RoleSystem, types.RoleDeveloper:
			systemParts = append(systemParts, content.Parts...)
		default:
			if len(content.Parts) == 0 {
				continue
			}

			r.Contents = append(r.Contents, content)
		}
	}

	if len(systemParts) > 0 {
		r.SystemInstruction = &geminiContent{Parts: systemParts}
	}

	return r
}

func compileMessage(msg types.Message) geminiContent {
	content := geminiContent{}
	if role := geminiRole(msg.Role); role != "" {
		content.Role = role
	}

	if msg.Role == types.RoleTool && msg.Name != "" {
		content.Parts = append(content.Parts, geminiPart{
			FunctionResponse: &geminiFunctionResult{
				Name:     msg.Name,
				Response: compileToolResponse(msg.Content),
			},
		})
		return content
	}

	if msg.Content != "" {
		content.Parts = append(content.Parts, geminiPart{Text: msg.Content})
	}

	for _, call := range msg.ToolCalls {
		content.Parts = append(content.Parts, geminiPart{
			FunctionCall: &geminiFunctionCall{
				Name:             call.Name,
				Args:             call.Arguments,
				ThoughtSignature: call.ThoughtSignature,
			},
		})
	}

	return content
}

func compileToolResponse(content string) json.RawMessage {
	content = strings.TrimSpace(content)
	if content == "" {
		return json.RawMessage(`{}`)
	}

	raw := json.RawMessage(content)
	if json.Valid(raw) {
		return raw
	}

	wrapped, err := json.Marshal(map[string]string{"output": content})
	if err != nil {
		return json.RawMessage(`{}`)
	}

	return wrapped
}

func generationConfig(req *types.ChatRequest) *geminiGenerationConfig {
	if req.Temperature == nil && req.TopP == nil && req.MaxTokens == nil {
		return nil
	}

	return &geminiGenerationConfig{
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		MaxOutputTokens: req.MaxTokens,
	}
}

func toAPIChatResponse(model string, resp *generateContentResponse) *types.ChatResponse {
	out := &types.ChatResponse{
		ID:      "chatcmpl-" + uuid.NewString(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   responseModel(model, resp.ModelVersion),
		Choices: make([]types.Choice, 0, len(resp.Candidates)),
	}

	if resp.UsageMetadata != nil {
		out.Usage = &types.Usage{
			PromptTokens:     resp.UsageMetadata.PromptTokenCount,
			CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      resp.UsageMetadata.TotalTokenCount,
		}
	}

	for i, candidate := range resp.Candidates {
		content, toolCalls := candidateMessage(candidate.Content)

		out.Choices = append(out.Choices, types.Choice{
			Index: i,
			Message: types.Message{
				Role:      types.RoleAssistant,
				Content:   content,
				ToolCalls: toolCalls,
			},
			FinishReason: finishReason(candidate.FinishReason, len(toolCalls) > 0),
		})
	}

	return out
}

func candidateMessage(content geminiContent) (string, []types.ToolCall) {
	var textParts []string
	toolCalls := make([]types.ToolCall, 0)

	for _, part := range content.Parts {
		if part.Text != "" {
			textParts = append(textParts, part.Text)
		}

		if part.FunctionCall != nil {
			toolCalls = append(toolCalls, types.ToolCall{
				ID:               uuid.NewString(),
				Name:             part.FunctionCall.Name,
				Arguments:        part.FunctionCall.Args,
				ThoughtSignature: part.FunctionCall.ThoughtSignature,
			})
		}
	}

	if len(toolCalls) == 0 {
		toolCalls = nil
	}

	return strings.Join(textParts, ""), toolCalls
}

func responseModel(requestModel, modelVersion string) string {
	if modelVersion != "" {
		return trimModelName(modelVersion)
	}

	return trimModelName(requestModel)
}

func geminiRole(role types.Role) string {
	switch role {
	case types.RoleAssistant:
		return "model"
	case types.RoleTool:
		return "user"
	default:
		return ""
	}
}

func supportsGenerateContent(model geminiModel) bool {
	for _, method := range model.SupportedGenerationMethods {
		if method == "generateContent" || method == "streamGenerateContent" {
			return true
		}
	}

	return false
}

func trimModelName(name string) string {
	return strings.TrimPrefix(name, "models/")
}

func finishReason(reason string, hasToolCalls bool) string {
	if hasToolCalls {
		return "tool_calls"
	}

	switch reason {
	case "MAX_TOKENS":
		return "length"
	case "STOP", "":
		return "stop"
	case "SAFETY", "RECITATION", "BLOCKLIST", "PROHIBITED_CONTENT", "SPII":
		return "content_filter"
	default:
		return strings.ToLower(reason)
	}
}
