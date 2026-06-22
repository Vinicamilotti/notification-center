package notification

import (
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
)

func ntfyConfig(enabled bool, topics ...string) config.NotificationConfig {
	return config.NotificationConfig{
		Type:                 config.Ntfy,
		Enabled:              enabled,
		Channel:              "http://ntfy.example.com",
		SubscribedTopicsList: topics,
	}
}

func TestGetService_ReturnsCurrentService(t *testing.T) {
	notificationService = *NewNotificationSender()

	svc := GetService()
	if len(svc.Channels) != 0 {
		t.Errorf("expected 0 channels on fresh sender, got %d", len(svc.Channels))
	}
}

func TestInit_RegistersNtfyChannel(t *testing.T) {
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			ntfyConfig(true, "*"),
		},
	})

	svc := GetService()
	if len(svc.Channels) != 1 {
		t.Fatalf("expected 1 channel after Init with Ntfy config, got %d", len(svc.Channels))
	}
	if svc.Channels[0] == nil {
		t.Error("expected a non-nil ntfy channel")
	}
}

func TestInit_EmptyConfig(t *testing.T) {
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{},
	})

	svc := GetService()
	if len(svc.Channels) != 0 {
		t.Errorf("expected 0 channels with empty config, got %d", len(svc.Channels))
	}
}

func TestInit_MultipleNtfyChannels(t *testing.T) {
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			ntfyConfig(true, "alerts"),
			ntfyConfig(false, "metrics"),
		},
	})

	svc := GetService()
	if len(svc.Channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(svc.Channels))
	}
}

func TestInit_IgnoresUnknownType(t *testing.T) {
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			{Type: config.ConfigType("unknown"), Enabled: true},
			ntfyConfig(true, "*"),
		},
	})

	svc := GetService()
	if len(svc.Channels) != 1 {
		t.Errorf("expected unknown type to be ignored, got %d channels", len(svc.Channels))
	}
}

func TestInit_ResetsPreviousChannels(t *testing.T) {
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			ntfyConfig(true, "*"),
			ntfyConfig(true, "alerts"),
		},
	})
	if got := len(GetService().Channels); got != 2 {
		t.Fatalf("expected 2 channels, got %d", got)
	}

	// A second Init with a smaller config must not accumulate channels.
	Init(config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			ntfyConfig(true, "*"),
		},
	})
	if got := len(GetService().Channels); got != 1 {
		t.Errorf("expected channels to be reset to 1, got %d", got)
	}
}
