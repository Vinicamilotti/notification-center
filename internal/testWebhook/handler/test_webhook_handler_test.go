package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestWebhookRouter() *gin.Engine {
	router := gin.New()
	handler := NewTestHandler()
	handler.RegisterRoutes(router)
	return router
}

func TestTestWebhook_ReturnsOK(t *testing.T) {
	router := setupTestWebhookRouter()

	req := httptest.NewRequest(http.MethodPost, "/test-webhook/infra", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestTestWebhook_ResponseBody(t *testing.T) {
	router := setupTestWebhookRouter()

	req := httptest.NewRequest(http.MethodPost, "/test-webhook/infra", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if response["message"] != "Test webhook handler is working!" {
		t.Errorf("unexpected message: %v", response["message"])
	}
}

func TestTestWebhook_DifferentTopics(t *testing.T) {
	router := setupTestWebhookRouter()

	topics := []string{"alerts", "metrics", "db-errors"}
	for _, topic := range topics {
		req := httptest.NewRequest(http.MethodPost, "/test-webhook/"+topic, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("topic %q: expected status 200, got %d", topic, w.Code)
		}
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	handler := NewTestHandler()
	handler.RegisterRoutes(router)

	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/test-webhook/:topic" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected POST /test-webhook/:topic route to be registered")
	}
}
