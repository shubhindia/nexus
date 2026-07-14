package openai

import "github.com/shubhindia/nexus/internal/types"

type ResponsesRequest struct {
	Model        string          `json:"model"`
	Instructions string          `json:"instructions,omitempty"`
	Input        []ResponseInput `json:"input"`
	Stream       bool            `json:"stream"`

	Tools      []types.Tool      `json:"tools,omitempty"`
	ToolChoice *types.ToolChoice `json:"tool_choice,omitempty"`
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
	ID     string           `json:"id"`
	Object string           `json:"object"`
	Model  string           `json:"model"`
	Output []ResponseOutput `json:"output"`
}

type ResponseOutput struct {
	Type    string                  `json:"type"`
	Role    string                  `json:"role"`
	Content []ResponseOutputContent `json:"content"`
}

type ResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
