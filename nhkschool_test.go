package nhkschool

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("my-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.apiKey != "my-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "my-key")
	}
	if c.baseURL != "https://api.nhk.or.jp/school/v2" {
		t.Errorf("baseURL = %q, want default", c.baseURL)
	}
	if c.httpClient != http.DefaultClient {
		t.Error("httpClient should default to http.DefaultClient")
	}
}

func TestNewClient_EmptyAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("NewClient('') should return an error")
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	c, err := NewClient("key", WithBaseURL("http://localhost:9999"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.baseURL != "http://localhost:9999" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "http://localhost:9999")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c, err := NewClient("key", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.httpClient != custom {
		t.Error("httpClient should be the custom client")
	}
}
