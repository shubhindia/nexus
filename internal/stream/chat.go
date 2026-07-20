package stream

type ChatCompletionChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index int   `json:"index"`
	Delta Delta `json:"delta"`

	FinishReason *string `json:"finish_reason"`
}

type Delta struct {
	Role      string          `json:"role,omitempty"`
	Content   string          `json:"content,omitempty"`
	ToolCalls []ToolCallDelta `json:"tool_calls,omitempty"`
}

type ToolCallDelta struct {
	Index    int               `json:"index,omitempty"`
	ID       string            `json:"id,omitempty"`
	Type     string            `json:"type,omitempty"`
	Function ToolFunctionDelta `json:"function,omitempty"`
}

type ToolFunctionDelta struct {
	Name             string `json:"name,omitempty"`
	Arguments        string `json:"arguments,omitempty"`
	ThoughtSignature string `json:"thought_signature,omitempty"`
}
