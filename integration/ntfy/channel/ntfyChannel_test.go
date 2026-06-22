package channel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ntfyDomain "github.com/Vinicamilotti/notification-center/integration/ntfy/domain"
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

func makeNtfyConfig(baseURL string) config.NotificationConfig {
	return config.NotificationConfig{
		Type:    config.Ntfy,
		Enabled: true,
		Channel: baseURL,
		SubscribedTopics: map[string]any{
			"*": nil,
		},
	}
}

func TestNewNtfyChannel(t *testing.T) {
	cfg := makeNtfyConfig("http://ntfy.example.com")
	ch := NewNtfyChannel(cfg)
	if ch == nil {
		t.Fatal("expected non-nil NtfyChannel")
	}
}

func TestGetConfig(t *testing.T) {
	cfg := makeNtfyConfig("http://ntfy.example.com")
	ch := NewNtfyChannel(cfg)
	got := ch.GetConfig()
	if got.Channel != "http://ntfy.example.com" {
		t.Errorf("expected channel URL 'http://ntfy.example.com', got %q", got.Channel)
	}
}

func TestParseTags(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))

	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"tag1", "tag2"}, "tag1,tag2"},
		{[]string{"only"}, "only"},
		{[]string{}, ""},
	}

	for _, tt := range tests {
		result := ch.parseTags(tt.input)
		if result != tt.expected {
			t.Errorf("parseTags(%v): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestDeterminateAction_BasicView(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	actions := []ntfyDomain.NtfyAction{
		{Type: ntfyDomain.NtfyActionTypeView, Label: "Open", Url: "https://example.com"},
	}

	result := ch.determinateAction(actions)
	expected := "view, Open, https://example.com"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeterminateAction_WithClear(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	actions := []ntfyDomain.NtfyAction{
		{Type: ntfyDomain.NtfyActionTypeView, Label: "Dismiss", Url: "https://example.com", Clear: true},
	}

	result := ch.determinateAction(actions)
	expected := "view, Dismiss, https://example.com, clear=true"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeterminateAction_WithMethod(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	actions := []ntfyDomain.NtfyAction{
		{Type: ntfyDomain.NtfyActionTypeHttp, Label: "Confirm", Url: "https://api.example.com", Method: "POST"},
	}

	result := ch.determinateAction(actions)
	expected := "http, Confirm, https://api.example.com, method=POST"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeterminateAction_WithBody(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	actions := []ntfyDomain.NtfyAction{
		{Type: ntfyDomain.NtfyActionTypeHttp, Label: "Act", Url: "https://example.com", Body: `{"k":"v"}`},
	}

	result := ch.determinateAction(actions)
	expected := `http, Act, https://example.com, body={"k":"v"}`
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeterminateAction_MultipleActions(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	actions := []ntfyDomain.NtfyAction{
		{Type: ntfyDomain.NtfyActionTypeView, Label: "A", Url: "https://a.com"},
		{Type: ntfyDomain.NtfyActionTypeHttp, Label: "B", Url: "https://b.com"},
	}

	result := ch.determinateAction(actions)
	expected := "view, A, https://a.com;http, B, https://b.com"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeterminateAction_Empty(t *testing.T) {
	ch := NewNtfyChannel(makeNtfyConfig("http://example.com"))
	result := ch.determinateAction([]ntfyDomain.NtfyAction{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestSendNotification_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ch := NewNtfyChannel(makeNtfyConfig(server.URL))
	dto := domain.NotificationDTO{
		Topic:               "test",
		Title:               "Test Title",
		Message:             "Test Message",
		AditionalAttributes: map[string]any{},
	}

	err := ch.SendNotification(dto)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSendNotification_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	ch := NewNtfyChannel(makeNtfyConfig(server.URL))
	dto := domain.NotificationDTO{
		Topic:               "test",
		Title:               "Test Title",
		Message:             "Test Message",
		AditionalAttributes: map[string]any{},
	}

	err := ch.SendNotification(dto)
	if err == nil {
		t.Error("expected error for 5xx response")
	}
}
