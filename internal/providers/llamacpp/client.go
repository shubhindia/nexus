package llamacpp

import (
	"net/http"
	"time"
)

type Client struct {
	url    string
	client *http.Client
}

type Option func(*Client)

func (c *Client) Name() string {
	return "llama.cpp"
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func New(url string, opts ...Option) *Client {
	c := &Client{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
