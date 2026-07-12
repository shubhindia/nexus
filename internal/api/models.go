package api

import (
	"net/http"

	"github.com/shubhindia/nexus/internal/gateway"
)

func Models(gw *gateway.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models, err := gw.Models(r.Context())
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, models)
	}
}
