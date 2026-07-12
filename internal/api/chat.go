package api

import (
	"io"
	"net/http"

	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/types"
)

func Chat(gw *gateway.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.ChatRequest

		if err := DecodeJSON(r, &req); err != nil {
			WriteError(w, &types.APIError{
				Status:  http.StatusBadRequest,
				Code:    "invalid_request",
				Message: err.Error(),
			})
			return
		}

		result, err := gw.Chat(r.Context(), &req)
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

			if _, err := io.Copy(w, result.Stream); err != nil {
				WriteError(w, err)
			}

			return
		}

		WriteJSON(w, http.StatusOK, result.Response)
	}
}
