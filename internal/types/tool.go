package types

import "encoding/json"

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolChoice struct {
	Mode string `json:"mode"`
	Name string `json:"name,omitempty"`
}

const (
	ToolChoiceAuto = "auto"
	ToolChoiceNone = "none"
)

type ToolCall struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Arguments json.RawMessage `json:"arguments"`
}

type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`

	Output string `json:"output"`
}
