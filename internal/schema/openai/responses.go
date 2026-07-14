package openai

type ResponsesRequest struct {
	Model        string          `json:"model"`
	Instructions string          `json:"instructions,omitempty"`
	Input        []ResponseInput `json:"input"`
	Stream       bool            `json:"stream"`

	Tools      []Tool      `json:"tools,omitempty"`
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
}

type ResponseInput struct {
	Role    string                 `json:"role"`
	Content []ResponseInputContent `json:"content"`
}

type ResponseInputContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type ResponsesResponse struct {
	ID     string               `json:"id"`
	Object string               `json:"object"`
	Model  string               `json:"model"`
	Output []ResponseOutputItem `json:"output"`
}

type ResponseOutputItem struct {
	Type string `json:"type"`

	// message
	Role    string                  `json:"role,omitempty"`
	Content []ResponseOutputContent `json:"content,omitempty"`

	// function_call
	CallID    string `json:"call_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type ResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
