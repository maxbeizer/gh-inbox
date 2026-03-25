package api

import (
	"fmt"

	ghAPI "github.com/cli/go-gh/v2/pkg/api"
)

// Client wraps the GitHub REST client.
type Client struct {
	rest *ghAPI.RESTClient
}

// NewClient creates a new API client using gh auth.
func NewClient() (*Client, error) {
	rest, err := ghAPI.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("creating REST client: %w", err)
	}

	return &Client{rest: rest}, nil
}
