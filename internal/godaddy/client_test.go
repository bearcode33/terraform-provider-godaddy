package godaddy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		apiSecret   string
		opts        []ClientOption
		wantBaseURL string
	}{
		{
			name:        "Default production client",
			apiKey:      "test-key",
			apiSecret:   "test-secret",
			opts:        nil,
			wantBaseURL: productionURL,
		},
		{
			name:      "Test environment client",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			opts: []ClientOption{
				WithTestEnvironment(),
			},
			wantBaseURL: testURL,
		},
		{
			name:      "Custom HTTP client",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			opts: []ClientOption{
				WithHTTPClient(&http.Client{Timeout: 60}),
			},
			wantBaseURL: productionURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.apiSecret, tt.opts...)

			if client.apiKey != tt.apiKey {
				t.Errorf("NewClient() apiKey = %v, want %v", client.apiKey, tt.apiKey)
			}
			if client.apiSecret != tt.apiSecret {
				t.Errorf("NewClient() apiSecret = %v, want %v", client.apiSecret, tt.apiSecret)
			}
			if client.baseURL != tt.wantBaseURL {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.wantBaseURL)
			}
		})
	}
}

func TestClient_doRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		auth := r.Header.Get("Authorization")
		expectedAuth := "sso-key test-key:test-secret"
		if auth != expectedAuth {
			t.Errorf("Authorization header = %v, want %v", auth, expectedAuth)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type header = %v, want application/json", contentType)
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-secret")
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.doRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("doRequest() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code": "NOT_FOUND", "message": "Domain not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-secret")
	client.baseURL = server.URL

	ctx := context.Background()
	_, err := client.doRequest(ctx, http.MethodGet, "/test", nil)
	if err == nil {
		t.Error("doRequest() expected error for 404 response")
	}
}
