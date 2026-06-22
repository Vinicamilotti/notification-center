package domain

import (
	"testing"
)

func TestNtfyActionTypeConstants(t *testing.T) {
	if NtfyActionTypeView != "view" {
		t.Errorf("expected NtfyActionTypeView to be 'view', got %q", NtfyActionTypeView)
	}
	if NtfyActionTypeHttp != "http" {
		t.Errorf("expected NtfyActionTypeHttp to be 'http', got %q", NtfyActionTypeHttp)
	}
}

func TestNtfyRequest_Fields(t *testing.T) {
	req := NtfyRequest{
		Title:   "Alert Title",
		Message: "Alert body",
		Click:   "https://example.com",
		Tag:     []string{"warning", "cpu"},
		Actions: []NtfyAction{
			{Type: NtfyActionTypeView, Label: "Open", Url: "https://example.com"},
		},
	}

	if req.Title != "Alert Title" {
		t.Errorf("expected Title 'Alert Title', got %q", req.Title)
	}
	if req.Message != "Alert body" {
		t.Errorf("expected Message 'Alert body', got %q", req.Message)
	}
	if req.Click != "https://example.com" {
		t.Errorf("expected Click 'https://example.com', got %q", req.Click)
	}
	if len(req.Tag) != 2 {
		t.Errorf("expected 2 tags, got %d", len(req.Tag))
	}
	if len(req.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(req.Actions))
	}
}

func TestNtfyAction_ViewType(t *testing.T) {
	action := NtfyAction{
		Type:  NtfyActionTypeView,
		Label: "Open Dashboard",
		Url:   "https://grafana.example.com",
		Clear: true,
	}

	if action.Type != NtfyActionTypeView {
		t.Errorf("expected Type %q, got %q", NtfyActionTypeView, action.Type)
	}
	if action.Label != "Open Dashboard" {
		t.Errorf("expected Label 'Open Dashboard', got %q", action.Label)
	}
	if !action.Clear {
		t.Error("expected Clear to be true")
	}
}

func TestNtfyAction_HttpType(t *testing.T) {
	action := NtfyAction{
		Type:   NtfyActionTypeHttp,
		Label:  "Acknowledge",
		Url:    "https://api.example.com/ack",
		Method: "POST",
		Body:   `{"status":"ack"}`,
	}

	if action.Type != NtfyActionTypeHttp {
		t.Errorf("expected Type %q, got %q", NtfyActionTypeHttp, action.Type)
	}
	if action.Method != "POST" {
		t.Errorf("expected Method 'POST', got %q", action.Method)
	}
	if action.Body != `{"status":"ack"}` {
		t.Errorf("expected Body to be set, got %q", action.Body)
	}
}
