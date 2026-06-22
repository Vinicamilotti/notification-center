package app

import (
	"testing"

	"github.com/gin-gonic/gin"
)

type mockHandler struct {
	called bool
}

func (m *mockHandler) RegisterRoutes(r *gin.Engine) {
	m.called = true
}

func TestNewApp(t *testing.T) {
	a := NewApp("0.0.0.0", "8080")
	if a == nil {
		t.Fatal("expected non-nil App")
	}
	if a.Listen != "0.0.0.0" {
		t.Errorf("expected Listen '0.0.0.0', got %q", a.Listen)
	}
	if a.Port != "8080" {
		t.Errorf("expected Port '8080', got %q", a.Port)
	}
	if len(a.Handlers) != 0 {
		t.Errorf("expected empty Handlers, got %d", len(a.Handlers))
	}
}

func TestApp_RegisterHandler(t *testing.T) {
	a := NewApp("0.0.0.0", "8080")
	h := &mockHandler{}
	a.RegisterHandler(h)

	if len(a.Handlers) != 1 {
		t.Fatalf("expected 1 handler, got %d", len(a.Handlers))
	}
}

func TestApp_RegisterHandler_Multiple(t *testing.T) {
	a := NewApp("0.0.0.0", "8080")
	h1 := &mockHandler{}
	h2 := &mockHandler{}
	h3 := &mockHandler{}
	a.RegisterHandler(h1)
	a.RegisterHandler(h2)
	a.RegisterHandler(h3)

	if len(a.Handlers) != 3 {
		t.Errorf("expected 3 handlers, got %d", len(a.Handlers))
	}
}

func TestApp_RegisterHandler_Order(t *testing.T) {
	a := NewApp("127.0.0.1", "9090")
	h1 := &mockHandler{}
	h2 := &mockHandler{}
	a.RegisterHandler(h1)
	a.RegisterHandler(h2)

	if a.Handlers[0] != h1 {
		t.Error("expected first handler to be h1")
	}
	if a.Handlers[1] != h2 {
		t.Error("expected second handler to be h2")
	}
}
