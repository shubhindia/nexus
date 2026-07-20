package gemini

import (
	"net/http"
	"time"
)

type Client struct {
	url    string
	apiKey string
	client *http.Client
}

type Option func(*Client)

func New(url string, apiKey string, opts ...Option) *Client {
	c := &Client{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Name() string {
	return "gemini"
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}
