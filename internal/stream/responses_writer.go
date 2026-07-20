package stream

import (
	"io"
	"strings"
	"time"

	"github.com/shubhindia/nexus/internal/toolstate"
)

type ResponseWriter struct {
	w io.Writer

	sequence int

	responseID string
	model      string
	createdAt  int64

	text strings.Builder

	messageIndex int
	output       []ResponseOutputItem
}

func NewResponseWriter(
	w io.Writer,
) *ResponseWriter {
	return &ResponseWriter{
		w:            w,
		createdAt:    time.Now().Unix(),
		messageIndex: -1,
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
	if rw.messageIndex == -1 {
		rw.messageIndex = len(rw.output)
		rw.output = append(rw.output, ResponseOutputItem{
			Type:   "message",
			Status: "in_progress",
			Role:   "assistant",
			Content: []ResponseOutputContent{{
				Type: "output_text",
				Text: "",
			}},
		})

		if err := writeSSE(
			rw.w,
			ResponseOutputItemAddedEvent{
				Type:           "response.output_item.added",
				SequenceNumber: rw.next(),
				OutputIndex:    rw.messageIndex,
				Item:           rw.output[rw.messageIndex],
			},
		); err != nil {
			return err
		}
	}

	rw.text.WriteString(delta)

	return writeSSE(
		rw.w,
		ResponseOutputTextDeltaEvent{
			Type:           "response.output_text.delta",
			SequenceNumber: rw.next(),
			OutputIndex:    rw.messageIndex,
			ContentIndex:   0,
			Delta:          delta,
		},
	)
}

func (rw *ResponseWriter) FunctionCalls(calls []ResponseOutputItem) error {
	for _, call := range calls {
		toolstate.StoreThoughtSignature(call.CallID, call.ThoughtSignature)

		index := len(rw.output)
		rw.output = append(rw.output, call)

		if err := writeSSE(
			rw.w,
			ResponseOutputItemAddedEvent{
				Type:           "response.output_item.added",
				SequenceNumber: rw.next(),
				OutputIndex:    index,
				Item:           call,
			},
		); err != nil {
			return err
		}

		if err := writeSSE(
			rw.w,
			ResponseOutputItemDoneEvent{
				Type:           "response.output_item.done",
				SequenceNumber: rw.next(),
				OutputIndex:    index,
				Item:           call,
			},
		); err != nil {
			return err
		}
	}

	return nil
}

func (rw *ResponseWriter) Completed() error {
	if rw.messageIndex != -1 {
		item := ResponseOutputItem{
			Type:   "message",
			Status: "completed",
			Role:   "assistant",
			Content: []ResponseOutputContent{{
				Type: "output_text",
				Text: rw.text.String(),
			}},
		}

		rw.output[rw.messageIndex] = item

		if err := writeSSE(
			rw.w,
			ResponseOutputItemDoneEvent{
				Type:           "response.output_item.done",
				SequenceNumber: rw.next(),
				OutputIndex:    rw.messageIndex,
				Item:           item,
			},
		); err != nil {
			return err
		}
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
				Output:    rw.output,
				Usage:     ResponseUsage{},
			},
		},
	)
}
