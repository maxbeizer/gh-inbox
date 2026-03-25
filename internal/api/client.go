package api

import (
	"fmt"

	ghAPI "github.com/cli/go-gh/v2/pkg/api"
)

// Client wraps GitHub REST and GraphQL clients.
type Client struct {
	rest *ghAPI.RESTClient
	gql  *ghAPI.GraphQLClient
}

// NewClient creates a new API client using gh auth.
func NewClient() (*Client, error) {
	rest, err := ghAPI.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("creating REST client: %w", err)
	}

	gql, err := ghAPI.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("creating GraphQL client: %w", err)
	}

	return &Client{rest: rest, gql: gql}, nil
}
