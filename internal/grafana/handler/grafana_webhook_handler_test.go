package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vinicamilotti/notification-center/internal/grafana/application"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupGrafanaRouter() *gin.Engine {
	router := gin.New()
	facade := application.NewGrafanaFacade()
	handler := NewGrafanaWebhookHandler(facade)
	handler.RegisterRoutes(router)
	return router
}

func validGrafanaPayload() map[string]any {
	return map[string]any{
		"receiver": "test-receiver",
		"status":   "firing",
		"title":    "Test Alert",
		"commonLabels": map[string]any{
			"alertname":      "HighCPU",
			"grafana_folder": "Infrastructure",
		},
		"commonAnnotations": map[string]any{
			"summary":     "CPU usage high",
			"description": "CPU has been above 90% for 5 minutes",
		},
		"alerts": []map[string]any{
			{
				"status":       "firing",
				"startsAt":     "2024-01-01T00:00:00Z",
				"endsAt":       "2024-01-01T01:00:00Z",
				"dashboardUrl": "https://grafana.example.com/d/abc",
				"silenceUrl":   "https://grafana.example.com/silence",
				"values":       map[string]any{"cpu": 95},
				"valueString":  "cpu=95",
			},
		},
	}
}

func TestHandleWebhook_ValidPayload(t *testing.T) {
	router := setupGrafanaRouter()
	body, _ := json.Marshal(validGrafanaPayload())

	req := httptest.NewRequest(http.MethodPost, "/grafana/infra", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	router := setupGrafanaRouter()

	req := httptest.NewRequest(http.MethodPost, "/grafana/infra", bytes.NewReader([]byte(`not-json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleWebhook_ResponseContainsStatus(t *testing.T) {
	router := setupGrafanaRouter()
	body, _ := json.Marshal(validGrafanaPayload())

	req := httptest.NewRequest(http.MethodPost, "/grafana/alerts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if response["status"] != "success" {
		t.Errorf("expected status 'success' in response, got %v", response["status"])
	}
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()
	facade := application.NewGrafanaFacade()
	handler := NewGrafanaWebhookHandler(facade)
	handler.RegisterRoutes(router)

	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/grafana/:topic" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected POST /grafana/:topic route to be registered")
	}
}
