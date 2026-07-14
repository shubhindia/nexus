package llamacpp

import (
	"encoding/json"

	"github.com/shubhindia/nexus/internal/types"
)

type chatTool struct {
	Type     string       `json:"type"`
	Function chatFunction `json:"function"`
}

type chatFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"`
}

type chatToolChoice struct {
	Type     string            `json:"type,omitempty"`
	Function *chatToolFunction `json:"function,omitempty"`
}

type chatToolFunction struct {
	Name string `json:"name"`
}

type chatToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function chatFunctionCall `json:"function"`
}

type chatFunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func convertTools(
	tools []types.Tool,
) []chatTool {

	if len(tools) == 0 {
		return nil
	}

	out := make([]chatTool, 0, len(tools))

	for _, tool := range tools {

		out = append(out, chatTool{
			Type: tool.Type.String(),

			Function: chatFunction{
				Name: tool.Name,

				Description: tool.Description,

				Parameters: tool.Parameters,
			},
		})
	}

	return out
}

func convertToolChoice(
	choice *types.ToolChoice,
) *chatToolChoice {

	if choice == nil {
		return nil
	}

	switch choice.Mode {

	case types.ToolChoiceAuto:
		return &chatToolChoice{
			Type: "auto",
		}

	case types.ToolChoiceNone:
		return &chatToolChoice{
			Type: "none",
		}

	case types.ToolChoiceTool:
		return &chatToolChoice{
			Type: "function",
			Function: &chatToolFunction{
				Name: choice.Name,
			},
		}
	}

	return nil
}

func convertToolCalls(
	calls []chatToolCall,
) []types.ToolCall {

	if len(calls) == 0 {
		return nil
	}

	out := make([]types.ToolCall, 0, len(calls))
	for _, call := range calls {
		out = append(out, types.ToolCall{
			ID:        call.ID,
			Name:      call.Function.Name,
			Arguments: call.Function.Arguments,
		})
	}

	return out
}
