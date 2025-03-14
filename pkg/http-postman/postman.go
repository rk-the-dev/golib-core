package httppostman

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// PostmanClient defines the interface for HTTP requests
type PostmanClient interface {
	GET(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	POST(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error)
	PUT(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error)
	DELETE(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Close() error
}

// postmanClient implements PostmanClient
type postmanClient struct {
	client *http.Client
}

var (
	instance *postmanClient
	once     sync.Once
)

// NewPostmanClient initializes and returns an HTTP client with timeout
func NewPostmanClient(timeout time.Duration) PostmanClient {
	once.Do(func() {
		instance = &postmanClient{
			client: &http.Client{
				Timeout: timeout,
			},
		}
	})
	return instance
}

// GET makes a GET request
func (p *postmanClient) GET(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return p.doRequest(ctx, http.MethodGet, url, nil, headers)
}

// POST makes a POST request
func (p *postmanClient) POST(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return p.doRequest(ctx, http.MethodPost, url, body, headers)
}

// PUT makes a PUT request
func (p *postmanClient) PUT(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return p.doRequest(ctx, http.MethodPut, url, body, headers)
}

// DELETE makes a DELETE request
func (p *postmanClient) DELETE(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return p.doRequest(ctx, http.MethodDelete, url, nil, headers)
}

// doRequest handles the actual HTTP request execution
func (p *postmanClient) doRequest(ctx context.Context, method, url string, body []byte, headers map[string]string) (*http.Response, error) {
	reqBody := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

// Close closes the HTTP client (useful for persistent connections)
func (p *postmanClient) Close() error {
	// No explicit close operation needed for http.Client
	return nil
}
