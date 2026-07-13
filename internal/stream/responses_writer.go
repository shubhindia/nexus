package stream

import (
	"io"
	"strings"
	"time"
)

type ResponseWriter struct {
	w io.Writer

	sequence int

	responseID string
	model      string
	createdAt  int64

	text strings.Builder
}

func NewResponseWriter(
	w io.Writer,
) *ResponseWriter {
	return &ResponseWriter{
		w:         w,
		createdAt: time.Now().Unix(),
	}
}

func (rw *ResponseWriter) next() int {
	n := rw.sequence
	rw.sequence++
	return n
}

func (rw *ResponseWriter) Created(
	responseID string,
	model string,
) error {

	rw.responseID = responseID
	rw.model = model

	return writeSSE(
		rw.w,
		ResponseCreatedEvent{
			Type:           "response.created",
			SequenceNumber: rw.next(),
			Response: Response{
				ID:        rw.responseID,
				Object:    "response",
				CreatedAt: rw.createdAt,
				Status:    "in_progress",
				Model:     rw.model,
				Output:    []ResponseOutputItem{},
				Usage:     ResponseUsage{},
			},
		},
	)
}

func (rw *ResponseWriter) Delta(
	delta string,
) error {

	rw.text.WriteString(delta)

	return writeSSE(
		rw.w,
		ResponseOutputTextDeltaEvent{
			Type:           "response.output_text.delta",
			SequenceNumber: rw.next(),
			OutputIndex:    0,
			ContentIndex:   0,
			Delta:          delta,
		},
	)
}

func (rw *ResponseWriter) Completed() error {

	text := rw.text.String()

	item := ResponseOutputItem{
		Type: "message",
		Role: "assistant",
		Content: []ResponseOutputContent{
			{
				Type: "output_text",
				Text: text,
			},
		},
	}

	if err := writeSSE(
		rw.w,
		ResponseOutputItemAddedEvent{
			Type:           "response.output_item.added",
			SequenceNumber: rw.next(),
			OutputIndex:    0,
			Item: ResponseOutputItem{
				Type: "message",
				Role: "assistant",
				Content: []ResponseOutputContent{
					{
						Type: "output_text",
						Text: "",
					},
				},
			},
		},
	); err != nil {
		return err
	}

	if err := writeSSE(
		rw.w,
		ResponseOutputItemDoneEvent{
			Type:           "response.output_item.done",
			SequenceNumber: rw.next(),
			OutputIndex:    0,
			Item:           item,
		},
	); err != nil {
		return err
	}

	return writeSSE(
		rw.w,
		ResponseCompletedEvent{
			Type:           "response.completed",
			SequenceNumber: rw.next(),
			Response: Response{
				ID:        rw.responseID,
				Object:    "response",
				CreatedAt: rw.createdAt,
				Status:    "completed",
				Model:     rw.model,
				Output: []ResponseOutputItem{
					item,
				},
				Usage: ResponseUsage{},
			},
		},
	)
}
