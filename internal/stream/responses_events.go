package stream

type Response struct {
	ID     string `json:"id"`
	Object string `json:"object,omitempty"`
	Model  string `json:"model,omitempty"`
	Status string `json:"status,omitempty"`
}

type ResponseItem struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Role   string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
}

type ResponseContentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type ResponseCreatedEvent struct {
	Type           string   `json:"type"`
	SequenceNumber int      `json:"sequence_number"`
	Response       Response `json:"response"`
}

type ResponseOutputItemAddedEvent struct {
	Type           string       `json:"type"`
	SequenceNumber int          `json:"sequence_number"`
	OutputIndex    int          `json:"output_index"`
	Item           ResponseItem `json:"item"`
}

type ResponseContentPartAddedEvent struct {
	Type           string              `json:"type"`
	SequenceNumber int                 `json:"sequence_number"`
	OutputIndex    int                 `json:"output_index"`
	ContentIndex   int                 `json:"content_index"`
	ItemID         string              `json:"item_id"`
	Part           ResponseContentPart `json:"part"`
}

type ResponseOutputTextDeltaEvent struct {
	Type           string `json:"type"`
	SequenceNumber int    `json:"sequence_number"`
	OutputIndex    int    `json:"output_index"`
	ContentIndex   int    `json:"content_index"`
	ItemID         string `json:"item_id"`
	Delta          string `json:"delta"`
}

type ResponseOutputTextDoneEvent struct {
	Type           string `json:"type"`
	SequenceNumber int    `json:"sequence_number"`
	OutputIndex    int    `json:"output_index"`
	ContentIndex   int    `json:"content_index"`
	ItemID         string `json:"item_id"`
	Text           string `json:"text"`
}

type ResponseContentPartDoneEvent struct {
	Type           string              `json:"type"`
	SequenceNumber int                 `json:"sequence_number"`
	OutputIndex    int                 `json:"output_index"`
	ContentIndex   int                 `json:"content_index"`
	ItemID         string              `json:"item_id"`
	Part           ResponseContentPart `json:"part"`
}

type ResponseOutputItemDoneEvent struct {
	Type           string       `json:"type"`
	SequenceNumber int          `json:"sequence_number"`
	OutputIndex    int          `json:"output_index"`
	Item           ResponseItem `json:"item"`
}

type ResponseCompletedEvent struct {
	Type           string   `json:"type"`
	SequenceNumber int      `json:"sequence_number"`
	Response       Response `json:"response"`
}
