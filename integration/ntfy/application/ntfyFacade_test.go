package application

import (
	"testing"

	ntfyDomain "github.com/Vinicamilotti/notification-center/integration/ntfy/domain"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

func TestNewNtfyFacade(t *testing.T) {
	f := NewNtfyFacade()
	if f == nil {
		t.Fatal("expected non-nil NtfyFacade")
	}
}

func TestProcessRequest_BasicFields(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:               "Alert Title",
		Message:             "Alert body",
		AditionalAttributes: map[string]any{},
	}

	req := f.ProcessRequest(dto)

	if req.Title != "Alert Title" {
		t.Errorf("expected title 'Alert Title', got %q", req.Title)
	}
	if req.Message != "Alert body" {
		t.Errorf("expected message 'Alert body', got %q", req.Message)
	}
}

func TestProcessRequest_WithClickAndTag(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:   "Title",
		Message: "Msg",
		AditionalAttributes: map[string]any{
			"click": "https://example.com",
			"tag":   "warning",
		},
	}

	req := f.ProcessRequest(dto)

	if req.Click != "https://example.com" {
		t.Errorf("expected click 'https://example.com', got %q", req.Click)
	}
	if len(req.Tag) != 1 || req.Tag[0] != "warning" {
		t.Errorf("expected tag ['warning'], got %v", req.Tag)
	}
}

func TestProcessRequest_MissingAttributes(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:               "T",
		Message:             "M",
		AditionalAttributes: map[string]any{},
	}

	req := f.ProcessRequest(dto)

	if req.Click != "" {
		t.Errorf("expected empty click, got %q", req.Click)
	}
}

func TestProcessRequest_NonStringAttribute(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:   "T",
		Message: "M",
		AditionalAttributes: map[string]any{
			"click": 12345,
		},
	}

	req := f.ProcessRequest(dto)
	if req.Click != "" {
		t.Errorf("expected empty click for non-string attribute, got %q", req.Click)
	}
}

func TestProcessRequest_WithActions(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:   "T",
		Message: "M",
		AditionalAttributes: map[string]any{},
		Actions: []domain.NotificationAction{
			{
				Type:   domain.ActionTypeUrl,
				Label:  "View",
				Action: "https://example.com",
			},
			{
				Type:   domain.ActionTypeHttpCall,
				Label:  "Call",
				Action: "https://api.example.com/action",
				Method: "POST",
			},
		},
	}

	req := f.ProcessRequest(dto)

	if len(req.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(req.Actions))
	}
	if req.Actions[0].Type != ntfyDomain.NtfyActionTypeView {
		t.Errorf("expected view action type, got %s", req.Actions[0].Type)
	}
	if req.Actions[0].Label != "View" {
		t.Errorf("expected label 'View', got %q", req.Actions[0].Label)
	}
	if req.Actions[1].Type != ntfyDomain.NtfyActionTypeHttp {
		t.Errorf("expected http action type, got %s", req.Actions[1].Type)
	}
	if req.Actions[1].Method != "POST" {
		t.Errorf("expected method 'POST', got %q", req.Actions[1].Method)
	}
}

func TestGetActionType(t *testing.T) {
	tests := []struct {
		input    domain.ActionType
		expected ntfyDomain.NtfyActionType
	}{
		{domain.ActionTypeUrl, ntfyDomain.NtfyActionTypeView},
		{domain.ActionTypeHttpCall, ntfyDomain.NtfyActionTypeHttp},
		{"unknown", ntfyDomain.NtfyActionTypeView},
	}

	for _, tt := range tests {
		result := getActionType(tt.input)
		if result != tt.expected {
			t.Errorf("getActionType(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestProcessRequest_EmptyActions(t *testing.T) {
	f := NewNtfyFacade()
	dto := domain.NotificationDTO{
		Title:               "T",
		Message:             "M",
		AditionalAttributes: map[string]any{},
		Actions:             []domain.NotificationAction{},
	}

	req := f.ProcessRequest(dto)

	if len(req.Actions) != 0 {
		t.Errorf("expected empty actions, got %d", len(req.Actions))
	}
}
