package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/types"
)

func Responses(
	gw *gateway.Gateway,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req types.ResponsesRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, err)
			return
		}

		result, err := gw.Chat(
			r.Context(),
			toChatRequest(&req),
		)
		if err != nil {
			WriteError(w, err)
			return
		}

		if result.Stream != nil {

			defer func() {
				_ = result.Stream.Close()
			}()

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			_, _ = io.Copy(w, result.Stream)
			return
		}

		WriteJSON(
			w,
			http.StatusOK,
			toResponsesResponse(result.Response),
		)
	}
}
