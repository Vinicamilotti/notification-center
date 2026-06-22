package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupCustomWebhookRouter() *gin.Engine {
	router := gin.New()
	handler := NewCustomWebhookHandler()
	handler.RegisterRoutes(router)
	return router
}

func TestHandleCustomWebhook_ValidPayload(t *testing.T) {
	router := setupCustomWebhookRouter()
	payload := map[string]any{
		"title":   "Test Notification",
		"message": "Hello from test",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/custom-webhook/my-topic", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestHandleCustomWebhook_InvalidJSON(t *testing.T) {
	router := setupCustomWebhookRouter()

	req := httptest.NewRequest(http.MethodPost, "/custom-webhook/my-topic", bytes.NewReader([]byte(`not-json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleCustomWebhook_ResponseBody(t *testing.T) {
	router := setupCustomWebhookRouter()
	payload := map[string]any{"title": "T", "message": "M"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/custom-webhook/alerts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if response["message"] != "Custom webhook received" {
		t.Errorf("expected message 'Custom webhook received', got %v", response["message"])
	}
}

func TestHandleCustomWebhook_EmptyJSON(t *testing.T) {
	router := setupCustomWebhookRouter()

	req := httptest.NewRequest(http.MethodPost, "/custom-webhook/topic", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty JSON object, got %d", w.Code)
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	handler := NewCustomWebhookHandler()
	handler.RegisterRoutes(router)

	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/custom-webhook/:topic" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected POST /custom-webhook/:topic route to be registered")
	}
}
