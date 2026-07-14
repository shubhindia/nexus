package llamacpp

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/shubhindia/nexus/internal/types"
)

type textualToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

var toolTags = []struct {
	Start string
	End   string
}{
	{"<function>", "</function>"},
	{"<function_call>", "</function_call>"},
	{"<tools>", "</tools>"},
}

func normalizeToolCalls(
	msg *chatMessage,
) []types.ToolCall {

	if len(msg.ToolCalls) > 0 {
		return convertToolCalls(msg.ToolCalls)
	}

	content := strings.TrimSpace(msg.Content)

	for _, tag := range toolTags {

		if !strings.HasPrefix(content, tag.Start) {
			continue
		}

		if !strings.HasSuffix(content, tag.End) {
			continue
		}

		jsonBody := strings.TrimSpace(
			strings.TrimSuffix(
				strings.TrimPrefix(content, tag.Start),
				tag.End,
			),
		)

		var call textualToolCall

		if err := json.Unmarshal(
			[]byte(jsonBody),
			&call,
		); err != nil {
			return nil
		}

		msg.Content = ""

		return []types.ToolCall{
			{
				ID:        uuid.NewString(),
				Name:      call.Name,
				Arguments: call.Arguments,
			},
		}
	}

	return nil
}
