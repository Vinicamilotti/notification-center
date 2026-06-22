package domain

import (
	"testing"
)

func TestGrafanaAlert_Fields(t *testing.T) {
	alert := GrafanaAlert{
		Reciver: "webhook-receiver",
		Status:  "firing",
		Title:   "High CPU",
		CommonLabels: CommonLabels{
			Alertname:     "HighCPU",
			GrafanaFolder: "Production",
		},
		CommonAnnotations: CommonAnnotations{
			Summary:     "CPU usage is critical",
			Description: "CPU exceeded 90% for 5 minutes",
		},
		Alerts: []Alert{
			{
				Status:      "firing",
				StartsAt:    "2024-01-01T00:00:00Z",
				EndsAt:      "2024-01-01T01:00:00Z",
				DashboarUrl: "https://grafana.example.com/dashboard",
				SilenceURL:  "https://grafana.example.com/silence",
				Labels:      map[string]string{"env": "prod"},
				Anotations:  map[string]string{"runbook": "https://runbook.example.com"},
				Values:      map[string]any{"cpu": 99.5},
				ValueString: "99.5",
			},
		},
	}

	if alert.Reciver != "webhook-receiver" {
		t.Errorf("expected Reciver 'webhook-receiver', got %q", alert.Reciver)
	}
	if alert.Status != "firing" {
		t.Errorf("expected Status 'firing', got %q", alert.Status)
	}
	if alert.Title != "High CPU" {
		t.Errorf("expected Title 'High CPU', got %q", alert.Title)
	}
	if alert.CommonLabels.Alertname != "HighCPU" {
		t.Errorf("expected Alertname 'HighCPU', got %q", alert.CommonLabels.Alertname)
	}
	if alert.CommonLabels.GrafanaFolder != "Production" {
		t.Errorf("expected GrafanaFolder 'Production', got %q", alert.CommonLabels.GrafanaFolder)
	}
	if alert.CommonAnnotations.Summary != "CPU usage is critical" {
		t.Errorf("expected Summary 'CPU usage is critical', got %q", alert.CommonAnnotations.Summary)
	}
	if alert.CommonAnnotations.Description != "CPU exceeded 90% for 5 minutes" {
		t.Errorf("expected Description set, got %q", alert.CommonAnnotations.Description)
	}
	if len(alert.Alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alert.Alerts))
	}
}

func TestAlert_Fields(t *testing.T) {
	alert := Alert{
		Status:      "resolved",
		StartsAt:    "2024-01-01T00:00:00Z",
		EndsAt:      "2024-01-01T01:00:00Z",
		DashboarUrl: "https://grafana.example.com/dashboard",
		SilenceURL:  "https://grafana.example.com/silence",
		Labels:      map[string]string{"severity": "critical", "env": "prod"},
		Anotations:  map[string]string{"runbook": "https://runbook.example.com"},
		Values:      map[string]any{"cpu": 95, "memory": "2GB"},
		ValueString: "95",
	}

	if alert.Status != "resolved" {
		t.Errorf("expected Status 'resolved', got %q", alert.Status)
	}
	if alert.Labels["severity"] != "critical" {
		t.Errorf("expected Labels severity 'critical', got %q", alert.Labels["severity"])
	}
	if alert.Anotations["runbook"] != "https://runbook.example.com" {
		t.Errorf("expected Anotations runbook set, got %q", alert.Anotations["runbook"])
	}
	if alert.Values["cpu"] != 95 {
		t.Errorf("expected Values cpu 95, got %v", alert.Values["cpu"])
	}
	if alert.ValueString != "95" {
		t.Errorf("expected ValueString '95', got %q", alert.ValueString)
	}
}

func TestCommonLabels_Fields(t *testing.T) {
	labels := CommonLabels{
		Alertname:     "DiskFull",
		GrafanaFolder: "Infrastructure",
	}

	if labels.Alertname != "DiskFull" {
		t.Errorf("expected Alertname 'DiskFull', got %q", labels.Alertname)
	}
	if labels.GrafanaFolder != "Infrastructure" {
		t.Errorf("expected GrafanaFolder 'Infrastructure', got %q", labels.GrafanaFolder)
	}
}

func TestCommonAnnotations_Fields(t *testing.T) {
	ann := CommonAnnotations{
		Summary:     "Disk is full",
		Description: "Disk usage exceeded 95%",
	}

	if ann.Summary != "Disk is full" {
		t.Errorf("expected Summary 'Disk is full', got %q", ann.Summary)
	}
	if ann.Description != "Disk usage exceeded 95%" {
		t.Errorf("expected Description set, got %q", ann.Description)
	}
}
