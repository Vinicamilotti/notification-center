package domain

import (
	"strings"
	"testing"
)

func TestHandleValues_Empty(t *testing.T) {
	result := handleValues(map[string]any{})
	if result != "" {
		t.Errorf("expected empty string for empty map, got %q", result)
	}
}

func TestHandleValues_SingleEntry(t *testing.T) {
	result := handleValues(map[string]any{"cpu": 99.5})
	if !strings.Contains(result, "cpu") {
		t.Errorf("expected result to contain 'cpu', got %q", result)
	}
	if !strings.Contains(result, "99.5") {
		t.Errorf("expected result to contain '99.5', got %q", result)
	}
}

func TestHandleValues_MultipleEntries(t *testing.T) {
	values := map[string]any{
		"cpu":    80,
		"memory": "2GB",
	}
	result := handleValues(values)
	if !strings.Contains(result, "cpu") {
		t.Errorf("expected result to contain 'cpu', got %q", result)
	}
	if !strings.Contains(result, "memory") {
		t.Errorf("expected result to contain 'memory', got %q", result)
	}
}

func TestNewMessageData(t *testing.T) {
	alert := GrafanaAlert{
		Status: "firing",
		CommonAnnotations: CommonAnnotations{
			Summary:     "High CPU",
			Description: "CPU usage is too high",
		},
		Alerts: []Alert{
			{
				StartsAt: "2024-01-01T00:00:00Z",
				EndsAt:   "2024-01-01T01:00:00Z",
				Values:   map[string]any{"cpu": 95},
			},
		},
	}

	data := NewMessageData(alert)

	if data.Sumary != "High CPU" {
		t.Errorf("expected Summary 'High CPU', got %q", data.Sumary)
	}
	if data.Description != "CPU usage is too high" {
		t.Errorf("expected Description 'CPU usage is too high', got %q", data.Description)
	}
	if data.Status != "firing" {
		t.Errorf("expected Status 'firing', got %q", data.Status)
	}
	if data.StartedAt != "2024-01-01T00:00:00Z" {
		t.Errorf("expected StartedAt '2024-01-01T00:00:00Z', got %q", data.StartedAt)
	}
	if data.EndedAt != "2024-01-01T01:00:00Z" {
		t.Errorf("expected EndedAt '2024-01-01T01:00:00Z', got %q", data.EndedAt)
	}
	if !strings.Contains(data.Values, "cpu") {
		t.Errorf("expected Values to contain 'cpu', got %q", data.Values)
	}
}

func TestNewMessageData_EmptyValues(t *testing.T) {
	alert := GrafanaAlert{
		Status:            "resolved",
		CommonAnnotations: CommonAnnotations{Summary: "S", Description: "D"},
		Alerts: []Alert{
			{
				StartsAt: "2024-01-01T00:00:00Z",
				EndsAt:   "2024-01-01T01:00:00Z",
				Values:   map[string]any{},
			},
		},
	}

	data := NewMessageData(alert)

	if data.Values != "" {
		t.Errorf("expected empty Values for empty map, got %q", data.Values)
	}
}
