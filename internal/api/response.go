package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/types"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	slog.Error(
		"api.error",
		slog.String("type", fmt.Sprintf("%T", err)),
		slog.Any("error", err),
	)

	apiErr, ok := err.(*types.APIError)
	if !ok {
		apiErr = &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "internal_error",
			Message: err.Error(),
		}
	}

	WriteJSON(w, apiErr.Status, map[string]any{
		"error": apiErr,
	})
}
