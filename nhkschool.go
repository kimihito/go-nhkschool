// Package nhkschool provides a client for the NHK for School API v2.
//
// It supports all four API endpoints: content detail, content list by curriculum
// standard code, keyword search, and bulk TSV data retrieval.
//
// # Getting Started
//
// Create a client with your API key:
//
//	client, err := nhkschool.NewClient("YOUR_API_KEY")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Retrieve a video by ID:
//
//	video, err := client.GetVideo(ctx, "D0005110412_00000")
//
// # Error Handling
//
// Non-2xx responses return an *APIError:
//
//	var apiErr *nhkschool.APIError
//	if errors.As(err, &apiErr) {
//		// apiErr.StatusCode, apiErr.Body
//	}
package nhkschool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.nhk.or.jp/school/v2"

// Client is an NHK for School API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) {
		cl.httpClient = c
	}
}

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(cl *Client) {
		cl.baseURL = url
	}
}

// NewClient creates a new NHK for School API client.
// apiKey is required and must not be empty.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("nhkschool: apiKey must not be empty")
	}
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// do executes an HTTP GET request and returns the response body.
// It returns an *APIError for non-2xx status codes.
func (c *Client) do(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	return body, nil
}
