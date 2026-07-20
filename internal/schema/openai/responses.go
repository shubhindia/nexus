package openai

import "github.com/google/uuid"

type ResponsesRequest struct {
	Model        string          `json:"model"`
	Instructions string          `json:"instructions,omitempty"`
	Input        []ResponseInput `json:"input"`
	Stream       bool            `json:"stream"`

	Tools      []Tool      `json:"tools,omitempty"`
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
}

type ResponseInput struct {
	Type             string                 `json:"type,omitempty"`
	Role             string                 `json:"role,omitempty"`
	Content          []ResponseInputContent `json:"content,omitempty"`
	CallID           string                 `json:"call_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Arguments        string                 `json:"arguments,omitempty"`
	ThoughtSignature string                 `json:"thought_signature,omitempty"`
	Output           string                 `json:"output,omitempty"`
	Status           string                 `json:"status,omitempty"`
}

type ResponseInputContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type ResponsesResponse struct {
	ID                 string               `json:"id"`
	Object             string               `json:"object"`
	CreatedAt          int64                `json:"created_at,omitempty"`
	Status             string               `json:"status,omitempty"`
	Background         bool                 `json:"background"`
	Error              any                  `json:"error"`
	IncompleteDetails  any                  `json:"incomplete_details"`
	Instructions       *string              `json:"instructions"`
	MaxOutputTokens    *int                 `json:"max_output_tokens"`
	Model              string               `json:"model"`
	Output             []ResponseOutputItem `json:"output"`
	ParallelToolCalls  bool                 `json:"parallel_tool_calls"`
	PreviousResponseID any                  `json:"previous_response_id"`
	Reasoning          ResponseReasoning    `json:"reasoning"`
	Store              bool                 `json:"store"`
	Text               ResponseText         `json:"text"`
	ToolChoice         string               `json:"tool_choice"`
	Tools              []any                `json:"tools"`
	TopP               *float64             `json:"top_p"`
	Temperature        *float64             `json:"temperature"`
	Truncation         string               `json:"truncation"`
	Usage              *ResponseUsage       `json:"usage,omitempty"`
	User               any                  `json:"user"`
	Metadata           map[string]any       `json:"metadata"`
}

type ResponseUsage struct {
	InputTokens         int                        `json:"input_tokens"`
	InputTokensDetails  ResponseInputTokenDetails  `json:"input_tokens_details"`
	OutputTokens        int                        `json:"output_tokens"`
	OutputTokensDetails ResponseOutputTokenDetails `json:"output_tokens_details"`
	TotalTokens         int                        `json:"total_tokens"`
}

type ResponseInputTokenDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type ResponseOutputTokenDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type ResponseReasoning struct {
	Effort  any   `json:"effort"`
	Summary []any `json:"summary"`
}

type ResponseText struct {
	Format ResponseTextFormat `json:"format"`
}

type ResponseTextFormat struct {
	Type string `json:"type"`
}

type ResponseOutputItem struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Status string `json:"status,omitempty"`

	// message
	Role    string                  `json:"role,omitempty"`
	Content []ResponseOutputContent `json:"content,omitempty"`

	// function_call
	CallID           string `json:"call_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Arguments        string `json:"arguments,omitempty"`
	ThoughtSignature string `json:"thought_signature,omitempty"`
}

type ResponseOutputContent struct {
	Type        string `json:"type"`
	Text        string `json:"text"`
	Annotations []any  `json:"annotations,omitempty"`
}

func NewResponseItemID() string {
	return "msg_" + uuid.NewString()
}
