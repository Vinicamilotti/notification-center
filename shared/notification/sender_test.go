package notification

import (
	"errors"
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type mockChannel struct {
	cfg    config.NotificationConfig
	err    error
	called bool
	lastDTO domain.NotificationDTO
}

func (m *mockChannel) GetConfig() config.NotificationConfig {
	return m.cfg
}

func (m *mockChannel) SendNotification(dto domain.NotificationDTO) error {
	m.called = true
	m.lastDTO = dto
	return m.err
}

func makeConfig(enabled bool, topics []string) config.NotificationConfig {
	topicSet := make(map[string]any)
	for _, t := range topics {
		topicSet[t] = nil
	}
	return config.NotificationConfig{
		Enabled:          enabled,
		SubscribedTopics: topicSet,
	}
}

func TestNewNotificationSender(t *testing.T) {
	ns := NewNotificationSender()
	if ns == nil {
		t.Fatal("expected non-nil NotificationSender")
	}
	if len(ns.Channels) != 0 {
		t.Errorf("expected empty channels, got %d", len(ns.Channels))
	}
}

func TestRegisterChannel(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{}
	ns.RegisterChannel(ch)
	if len(ns.Channels) != 1 {
		t.Errorf("expected 1 channel, got %d", len(ns.Channels))
	}
}

func TestSend_EnabledAndSubscribed(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfig(true, []string{"alerts"})}
	ns.RegisterChannel(ch)

	dto := domain.NotificationDTO{Topic: "alerts", Title: "Test"}
	errs := ns.Send(dto)

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if !ch.called {
		t.Error("expected channel to be called")
	}
}

func TestSend_DisabledChannel(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfig(false, []string{"alerts"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "alerts"})

	if ch.called {
		t.Error("expected disabled channel to not be called")
	}
}

func TestSend_NotSubscribedToTopic(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfig(true, []string{"metrics"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "alerts"})

	if ch.called {
		t.Error("expected channel to not be called for unsubscribed topic")
	}
}

func TestSend_WildcardSubscription(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfig(true, []string{"*"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "any-topic"})

	if !ch.called {
		t.Error("expected wildcard channel to be called for any topic")
	}
}

func TestSend_ChannelReturnsError(t *testing.T) {
	ns := NewNotificationSender()
	sentinel := errors.New("send failed")
	ch := &mockChannel{cfg: makeConfig(true, []string{"alerts"}), err: sentinel}
	ns.RegisterChannel(ch)

	errs := ns.Send(domain.NotificationDTO{Topic: "alerts"})

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0] != sentinel {
		t.Errorf("expected sentinel error, got %v", errs[0])
	}
}

func TestSend_MultipleChannels(t *testing.T) {
	ns := NewNotificationSender()
	ch1 := &mockChannel{cfg: makeConfig(true, []string{"alerts"})}
	ch2 := &mockChannel{cfg: makeConfig(true, []string{"*"})}
	ch3 := &mockChannel{cfg: makeConfig(false, []string{"alerts"})}
	ns.RegisterChannel(ch1)
	ns.RegisterChannel(ch2)
	ns.RegisterChannel(ch3)

	ns.Send(domain.NotificationDTO{Topic: "alerts"})

	if !ch1.called {
		t.Error("expected ch1 to be called")
	}
	if !ch2.called {
		t.Error("expected ch2 (wildcard) to be called")
	}
	if ch3.called {
		t.Error("expected ch3 (disabled) to not be called")
	}
}

func TestSend_CollectsAllErrors(t *testing.T) {
	ns := NewNotificationSender()
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	ch1 := &mockChannel{cfg: makeConfig(true, []string{"*"}), err: err1}
	ch2 := &mockChannel{cfg: makeConfig(true, []string{"*"}), err: err2}
	ns.RegisterChannel(ch1)
	ns.RegisterChannel(ch2)

	errs := ns.Send(domain.NotificationDTO{Topic: "test"})

	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestSend_NoChannels(t *testing.T) {
	ns := NewNotificationSender()
	errs := ns.Send(domain.NotificationDTO{Topic: "alerts"})
	if len(errs) != 0 {
		t.Errorf("expected no errors with no channels, got %v", errs)
	}
}
