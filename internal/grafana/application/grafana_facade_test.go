package application

import (
	"strings"
	"testing"

	"github.com/Vinicamilotti/notification-center/internal/grafana/domain"
)

func makeAlertWithStatus(status string) domain.GrafanaAlert {
	return domain.GrafanaAlert{
		Status: status,
		Title:  "Test Alert",
		CommonAnnotations: domain.CommonAnnotations{
			Summary:     "Summary text",
			Description: "Description text",
		},
		Alerts: []domain.Alert{
			{
				StartsAt:    "2024-01-01T00:00:00Z",
				EndsAt:      "2024-01-01T01:00:00Z",
				DashboarUrl: "https://grafana.example.com/dashboard",
				Values:      map[string]any{"cpu": 80},
			},
		},
	}
}

func TestNewGrafanaFacade(t *testing.T) {
	f := NewGrafanaFacade()
	if f == nil {
		t.Fatal("expected non-nil GrafanaFacade")
	}
}

func TestDeterminateTags_Firing(t *testing.T) {
	f := NewGrafanaFacade()
	result := f.determinateTags("firing")
	if result != "rotating_light" {
		t.Errorf("expected 'rotating_light', got %q", result)
	}
}

func TestDeterminateTags_FiringUppercase(t *testing.T) {
	f := NewGrafanaFacade()
	result := f.determinateTags("FIRING")
	if result != "rotating_light" {
		t.Errorf("expected 'rotating_light', got %q", result)
	}
}

func TestDeterminateTags_Resolved(t *testing.T) {
	f := NewGrafanaFacade()
	result := f.determinateTags("resolved")
	if result != "heavy_check_mark" {
		t.Errorf("expected 'heavy_check_mark', got %q", result)
	}
}

func TestDeterminateTags_ResolvedUppercase(t *testing.T) {
	f := NewGrafanaFacade()
	result := f.determinateTags("RESOLVED")
	if result != "heavy_check_mark" {
		t.Errorf("expected 'heavy_check_mark', got %q", result)
	}
}

func TestDeterminateTags_Unknown(t *testing.T) {
	f := NewGrafanaFacade()
	result := f.determinateTags("pending")
	if result != "warning" {
		t.Errorf("expected 'warning' for unknown status, got %q", result)
	}
}

func TestCreateAttributes_WithAlerts(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("firing")

	attrs := f.createAttributes(alert)

	if attrs["tag"] != "rotating_light" {
		t.Errorf("expected tag 'rotating_light', got %v", attrs["tag"])
	}
	if attrs["click"] != "https://grafana.example.com/dashboard" {
		t.Errorf("expected click URL, got %v", attrs["click"])
	}
}

func TestCreateAttributes_NoAlerts(t *testing.T) {
	f := NewGrafanaFacade()
	alert := domain.GrafanaAlert{
		Status: "resolved",
		CommonAnnotations: domain.CommonAnnotations{
			Summary:     "S",
			Description: "D",
		},
		Alerts: []domain.Alert{},
	}

	attrs := f.createAttributes(alert)

	if attrs["tag"] != "heavy_check_mark" {
		t.Errorf("expected tag 'heavy_check_mark', got %v", attrs["tag"])
	}
	if _, exists := attrs["click"]; exists {
		t.Error("expected no 'click' attribute when there are no alerts")
	}
}

func TestProcessAlert_ReturnsTopic(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("firing")

	dto, err := f.ProcessAlert("infra", alert)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Topic != "infra" {
		t.Errorf("expected topic 'infra', got %q", dto.Topic)
	}
}

func TestProcessAlert_ReturnsTitle(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("firing")

	dto, err := f.ProcessAlert("infra", alert)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Title != "Test Alert" {
		t.Errorf("expected title 'Test Alert', got %q", dto.Title)
	}
}

func TestProcessAlert_MessageContainsSummary(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("firing")

	dto, err := f.ProcessAlert("infra", alert)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(dto.Message, "Summary text") {
		t.Errorf("expected message to contain summary, got %q", dto.Message)
	}
}

func TestProcessAlert_MessageContainsDescription(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("resolved")

	dto, err := f.ProcessAlert("db", alert)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(dto.Message, "Description text") {
		t.Errorf("expected message to contain description, got %q", dto.Message)
	}
}

func TestProcessAlert_SetsAdditionalAttributes(t *testing.T) {
	f := NewGrafanaFacade()
	alert := makeAlertWithStatus("firing")

	dto, err := f.ProcessAlert("infra", alert)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.AditionalAttributes == nil {
		t.Fatal("expected non-nil AditionalAttributes")
	}
	if dto.AditionalAttributes["tag"] != "rotating_light" {
		t.Errorf("expected tag 'rotating_light', got %v", dto.AditionalAttributes["tag"])
	}
}
