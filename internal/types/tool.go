package types

import "encoding/json"

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type Tool struct {
	Type        ToolType        `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolChoiceMode string

const (
	ToolChoiceAuto ToolChoiceMode = "auto"
	ToolChoiceNone ToolChoiceMode = "none"
	ToolChoiceTool ToolChoiceMode = "tool"
)

type ToolChoice struct {
	Mode ToolChoiceMode `json:"mode"`
	Name string         `json:"name,omitempty"`
}

type ToolCall struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Arguments        json.RawMessage `json:"arguments"`
	ThoughtSignature string          `json:"thought_signature,omitempty"`
}

func (t ToolType) String() string {
	return string(t)
}

func ParseToolType(s string) ToolType {

	switch s {

	case "function":
		return ToolTypeFunction

	default:
		return ToolTypeFunction
	}
}
