package openai

import (
	"encoding/json"
)

type Tool struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolChoice struct {
	Mode string `json:"mode"`
	Name string `json:"name,omitempty"`
}

const (
	ToolChoiceAuto     = "auto"
	ToolChoiceNone     = "none"
	ToolChoiceRequired = "required"
)

func (t *ToolChoice) UnmarshalJSON(data []byte) error {

	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	// Supports:
	// "auto"
	// "none"
	// "required"
	if data[0] == '"' {
		var mode string
		if err := json.Unmarshal(data, &mode); err != nil {
			return err
		}
		t.Mode = mode
		return nil
	}

	// Supports object form.
	type alias ToolChoice
	var tmp alias

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*t = ToolChoice(tmp)

	return nil
}
