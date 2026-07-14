package api

import (
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/compiler"
	"github.com/shubhindia/nexus/internal/encoder"
	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/stream"
)

var (
	openAICompiler = compiler.NewOpenAICompiler()
	openAIEncoder  = encoder.NewOpenAIEncoder()
)

func Responses(
	gw *gateway.Gateway,
	log *slog.Logger,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req openai.ResponsesRequest

		if err := DecodeJSON(r, &req); err != nil {
			WriteError(w, err)
			return
		}

		chatReq := openAICompiler.Compile(&req)

		result, err := gw.Chat(
			r.Context(),
			chatReq,
		)
		if err != nil {
			log.Error(
				"gateway.chat",
				slog.Any("error", err),
			)

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
			w.Header().Set("X-Accel-Buffering", "no")

			rc := http.NewResponseController(w)

			if err := rc.Flush(); err != nil {
				http.Error(
					w,
					"streaming unsupported",
					http.StatusInternalServerError,
				)
				return
			}

			if err := stream.ChatCompletionToResponses(
				w,
				result.Stream,
			); err != nil {
				log.Error(
					"responses.stream",
					slog.Any("error", err),
				)
			}

			return
		}

		WriteJSON(
			w,
			http.StatusOK,
			openAIEncoder.Responses(result.Response),
		)
	}
}
