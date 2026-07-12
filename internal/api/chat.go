package api

import (
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
			http.Error(
				w,
				"streaming not implemented",
				http.StatusNotImplemented,
			)
			return
		}

		WriteJSON(w, http.StatusOK, result.Response)
	}
}
