package provider

import (
	"io"

	"github.com/shubhindia/nexus/internal/types"
)

// ChatResult is the result of a chat completion.
//
// Exactly one of Response or Stream must be non-nil.
//
// For non-streaming requests, Response is populated.
// For streaming requests, Stream is populated.
type ChatResult struct {
	Response *types.ChatResponse
	Stream   io.ReadCloser
}
