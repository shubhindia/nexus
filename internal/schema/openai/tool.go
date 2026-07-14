package openai

import "encoding/json"

type Tool struct {
	Type string `json:"type"`

	Name string `json:"name"`

	Description string `json:"description,omitempty"`

	Parameters json.RawMessage `json:"parameters"`
}

type ToolChoice struct {
	Mode string `json:"mode"`

	Name string `json:"name,omitempty"`
}
