package domain

import (
	"testing"
)

func TestActionTypeConstants(t *testing.T) {
	if ActionTypeUrl != "url" {
		t.Errorf("expected ActionTypeUrl to be 'url', got %q", ActionTypeUrl)
	}
	if ActionTypeHttpCall != "call" {
		t.Errorf("expected ActionTypeHttpCall to be 'call', got %q", ActionTypeHttpCall)
	}
}

func TestNotificationDTO_Fields(t *testing.T) {
	dto := NotificationDTO{
		Topic:               "alerts",
		Title:               "High CPU",
		Message:             "CPU usage exceeded 90%",
		AditionalAttributes: map[string]any{"severity": "critical"},
		Actions: []NotificationAction{
			{
				Type:   ActionTypeUrl,
				Label:  "Open Dashboard",
				Action: "https://grafana.example.com",
				Method: "GET",
			},
		},
	}

	if dto.Topic != "alerts" {
		t.Errorf("expected Topic 'alerts', got %q", dto.Topic)
	}
	if dto.Title != "High CPU" {
		t.Errorf("expected Title 'High CPU', got %q", dto.Title)
	}
	if dto.Message != "CPU usage exceeded 90%" {
		t.Errorf("expected Message set, got %q", dto.Message)
	}
	if dto.AditionalAttributes["severity"] != "critical" {
		t.Errorf("expected severity 'critical', got %v", dto.AditionalAttributes["severity"])
	}
	if len(dto.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(dto.Actions))
	}
	if dto.Actions[0].Type != ActionTypeUrl {
		t.Errorf("expected action Type %q, got %q", ActionTypeUrl, dto.Actions[0].Type)
	}
}

func TestNotificationAction_UrlType(t *testing.T) {
	action := NotificationAction{
		Type:   ActionTypeUrl,
		Label:  "View",
		Action: "https://example.com",
		Method: "GET",
	}

	if action.Type != ActionTypeUrl {
		t.Errorf("expected Type %q, got %q", ActionTypeUrl, action.Type)
	}
	if action.Label != "View" {
		t.Errorf("expected Label 'View', got %q", action.Label)
	}
	if action.Action != "https://example.com" {
		t.Errorf("expected Action 'https://example.com', got %q", action.Action)
	}
}

func TestNotificationAction_HttpCallType(t *testing.T) {
	action := NotificationAction{
		Type:               ActionTypeHttpCall,
		Label:              "Acknowledge",
		Action:             "https://api.example.com/ack",
		AditionalParameter: map[string]string{"token": "abc123"},
		Method:             "POST",
	}

	if action.Type != ActionTypeHttpCall {
		t.Errorf("expected Type %q, got %q", ActionTypeHttpCall, action.Type)
	}
	if action.Method != "POST" {
		t.Errorf("expected Method 'POST', got %q", action.Method)
	}
	if action.AditionalParameter == nil {
		t.Error("expected AditionalParameter to be set")
	}
}
