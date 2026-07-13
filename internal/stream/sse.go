package stream

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeSSE(
	w io.Writer,
	v any,
) error {

	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		w,
		"data: %s\n\n",
		payload,
	)

	return err
}
