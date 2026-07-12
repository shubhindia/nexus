package provider

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/shubhindia/nexus/internal/types"
)

type LlamaCPP struct {
	baseURL string
	client  *http.Client
}

func NewLlamaCPP(baseURL string) *LlamaCPP {
	return &LlamaCPP{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (l *LlamaCPP) Models(ctx context.Context) (*types.ModelsResponse, error) {
	return nil, errors.New("not implemented")
}

func (l *LlamaCPP) Chat(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	return nil, errors.New("not implemented")
}
