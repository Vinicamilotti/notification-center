package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://example.com")
	if client.BaseURL != "https://example.com" {
		t.Errorf("expected BaseURL 'https://example.com', got %q", client.BaseURL)
	}
}

func TestClient_Post(t *testing.T) {
	client := NewClient("https://example.com")
	req, err := client.Post("/path", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "POST" {
		t.Errorf("expected POST, got %s", req.Method)
	}
	if req.URL.String() != "https://example.com/path" {
		t.Errorf("unexpected URL: %s", req.URL.String())
	}
}

func TestClient_Get(t *testing.T) {
	client := NewClient("https://example.com")
	req, err := client.Get("/resource", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "GET" {
		t.Errorf("expected GET, got %s", req.Method)
	}
}

func TestClient_Delete(t *testing.T) {
	client := NewClient("https://example.com")
	req, err := client.Delete("/resource", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", req.Method)
	}
}

func TestClient_Put(t *testing.T) {
	client := NewClient("https://example.com")
	req, err := client.Put("/resource", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "PUT" {
		t.Errorf("expected PUT, got %s", req.Method)
	}
}

func TestClient_Patch(t *testing.T) {
	client := NewClient("https://example.com")
	req, err := client.Patch("/resource", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "PATCH" {
		t.Errorf("expected PATCH, got %s", req.Method)
	}
}

func TestClient_URLConcatenation(t *testing.T) {
	client := NewClient("http://base.com")
	req, err := client.Get("/some/path", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "http://base.com/some/path"
	if req.URL.String() != expected {
		t.Errorf("expected %q, got %q", expected, req.URL.String())
	}
}

func TestClient_PostWithBody(t *testing.T) {
	client := NewClient("http://example.com")
	body := strings.NewReader(`{"key":"value"}`)
	req, err := client.Post("/endpoint", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Body == nil {
		t.Error("expected non-nil request body")
	}
}

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	req, err := client.Get("/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
