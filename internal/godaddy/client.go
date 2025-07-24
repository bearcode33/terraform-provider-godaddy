package godaddy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RateLimitResponse represents the structure of rate limit error response
type RateLimitResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	RetryAfterSec int    `json:"retryAfterSec"`
}

const (
	productionURL = "https://api.godaddy.com"
	testURL       = "https://api.ote-godaddy.com"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	apiSecret  string
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTestEnvironment() ClientOption {
	return func(c *Client) {
		c.baseURL = testURL
	}
}

func NewClient(apiKey, apiSecret string, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   productionURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	return c.doRequestWithRetry(ctx, method, path, body, 0)
}

func (c *Client) doRequestWithRetry(ctx context.Context, method, path string, body interface{}, retryCount int) (*http.Response, error) {
	const maxRetries = 3
	
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", c.apiKey, c.apiSecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		// Handle rate limiting with intelligent retry
		if resp.StatusCode == 429 && retryCount < maxRetries {
			fmt.Printf("RATE LIMIT HIT: %s %s - Status: %d, Body: %s (Retry %d/%d)\n",
				req.Method, req.URL.Path, resp.StatusCode, string(bodyBytes), retryCount+1, maxRetries)
			
			// Parse retryAfterSec from response body
			waitTime := c.parseRetryAfter(bodyBytes)
			if waitTime > 0 {
				fmt.Printf("Waiting %d seconds before retry as specified by GoDaddy API\n", waitTime)
				select {
				case <-time.After(time.Duration(waitTime) * time.Second):
					// Continue with retry
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			} else {
				// Fallback exponential backoff if retryAfterSec not found
				backoffTime := time.Duration(1<<uint(retryCount)) * time.Second
				fmt.Printf("Using exponential backoff: waiting %v before retry\n", backoffTime)
				select {
				case <-time.After(backoffTime):
					// Continue with retry
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			
			// Recursive retry
			return c.doRequestWithRetry(ctx, method, path, body, retryCount+1)
		}
		
		// DIAGNOSTIC LOG for non-429 errors or max retries reached
		if resp.StatusCode == 429 {
			fmt.Printf("RATE LIMIT EXHAUSTED: Max retries (%d) reached for %s %s\n",
				maxRetries, req.Method, req.URL.Path)
		}
		
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// parseRetryAfter extracts retryAfterSec from GoDaddy's rate limit response
func (c *Client) parseRetryAfter(bodyBytes []byte) int {
	var rateLimitResp RateLimitResponse
	if err := json.Unmarshal(bodyBytes, &rateLimitResp); err == nil && rateLimitResp.RetryAfterSec > 0 {
		return rateLimitResp.RetryAfterSec
	}
	
	// Fallback: try to extract from message string
	bodyStr := string(bodyBytes)
	if strings.Contains(bodyStr, "retryAfterSec") {
		// Look for "retryAfterSec":51 pattern
		start := strings.Index(bodyStr, "retryAfterSec\":")
		if start != -1 {
			start += len("retryAfterSec\":")
			end := start
			for end < len(bodyStr) && bodyStr[end] >= '0' && bodyStr[end] <= '9' {
				end++
			}
			if waitTimeStr := bodyStr[start:end]; len(waitTimeStr) > 0 {
				if waitTime, err := strconv.Atoi(waitTimeStr); err == nil {
					return waitTime
				}
			}
		}
	}
	
	return 0 // No valid retry time found
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
