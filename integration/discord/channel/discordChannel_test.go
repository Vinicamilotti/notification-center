package channel

import (
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

func makeDiscordConfig() config.NotificationConfig {
	return config.NotificationConfig{
		Type:    config.Discord,
		Enabled: true,
		Channel: "https://discord.com/api/webhooks/123/abc",
		SubscribedTopics: map[string]any{
			"alerts": nil,
		},
	}
}

func TestNewDiscordChannel(t *testing.T) {
	ch := NewDiscordChannel(makeDiscordConfig())
	if ch == nil {
		t.Fatal("expected non-nil DiscordChannel")
	}
}

func TestDiscordChannel_GetConfig(t *testing.T) {
	cfg := makeDiscordConfig()
	ch := NewDiscordChannel(cfg)
	got := ch.GetConfig()

	if got.Type != config.Discord {
		t.Errorf("expected Type %q, got %q", config.Discord, got.Type)
	}
	if got.Channel != cfg.Channel {
		t.Errorf("expected Channel %q, got %q", cfg.Channel, got.Channel)
	}
	if !got.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestDiscordChannel_SendNotification(t *testing.T) {
	ch := NewDiscordChannel(makeDiscordConfig())
	dto := domain.NotificationDTO{
		Topic:   "alerts",
		Title:   "Test Alert",
		Message: "Something happened",
	}
	err := ch.SendNotification(dto)
	if err != nil {
		t.Errorf("expected no error from stub SendNotification, got %v", err)
	}
}
