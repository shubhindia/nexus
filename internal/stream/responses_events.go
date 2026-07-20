package stream

type ResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type ResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseOutputItem struct {
	ID               string                  `json:"id,omitempty"`
	Type             string                  `json:"type"`
	Status           string                  `json:"status,omitempty"`
	Role             string                  `json:"role,omitempty"`
	Content          []ResponseOutputContent `json:"content,omitempty"`
	CallID           string                  `json:"call_id,omitempty"`
	Name             string                  `json:"name,omitempty"`
	Arguments        string                  `json:"arguments,omitempty"`
	ThoughtSignature string                  `json:"thought_signature,omitempty"`
}

type Response struct {
	ID        string               `json:"id"`
	Object    string               `json:"object"`
	CreatedAt int64                `json:"created_at"`
	Status    string               `json:"status"`
	Model     string               `json:"model"`
	Output    []ResponseOutputItem `json:"output"`
	Usage     ResponseUsage        `json:"usage"`
}

type ResponseCreatedEvent struct {
	Type           string   `json:"type"`
	SequenceNumber int      `json:"sequence_number,omitempty"`
	Response       Response `json:"response"`
}

type ResponseOutputItemAddedEvent struct {
	Type           string             `json:"type"`
	SequenceNumber int                `json:"sequence_number,omitempty"`
	OutputIndex    int                `json:"output_index"`
	Item           ResponseOutputItem `json:"item"`
}

type ResponseOutputTextDeltaEvent struct {
	Type           string `json:"type"`
	SequenceNumber int    `json:"sequence_number,omitempty"`
	OutputIndex    int    `json:"output_index"`
	ContentIndex   int    `json:"content_index"`
	Delta          string `json:"delta"`
}

type ResponseOutputItemDoneEvent struct {
	Type           string             `json:"type"`
	SequenceNumber int                `json:"sequence_number,omitempty"`
	OutputIndex    int                `json:"output_index"`
	Item           ResponseOutputItem `json:"item"`
}

type ResponseCompletedEvent struct {
	Type           string   `json:"type"`
	SequenceNumber int      `json:"sequence_number,omitempty"`
	Response       Response `json:"response"`
}
