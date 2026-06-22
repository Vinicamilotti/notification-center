package notification

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
)

func setupServiceConfig(t *testing.T, cfg config.Configs) (restore func()) {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd error: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir error: %v", err)
	}

	// force config package to reload from the new directory
	config.ReadConfigFile()

	return func() {
		os.Chdir(oldDir)
		notificationService = *NewNotificationSender()
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
	restore := setupServiceConfig(t, config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			{
				Type:                 config.Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"*"},
			},
		},
	})
	defer restore()

	Init()

	svc := GetService()
	if len(svc.Channels) != 1 {
		t.Errorf("expected 1 channel after Init with Ntfy config, got %d", len(svc.Channels))
	}
}

func TestInit_EmptyConfig(t *testing.T) {
	restore := setupServiceConfig(t, config.Configs{
		NotificationConfigs: []config.NotificationConfig{},
	})
	defer restore()

	Init()

	svc := GetService()
	if len(svc.Channels) != 0 {
		t.Errorf("expected 0 channels with empty config, got %d", len(svc.Channels))
	}
}

func TestInit_MultipleChannels(t *testing.T) {
	restore := setupServiceConfig(t, config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			{
				Type:                 config.Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"alerts"},
			},
			{
				Type:                 config.Ntfy,
				Enabled:              false,
				Channel:              "http://ntfy2.example.com",
				SubscribedTopicsList: []string{"metrics"},
			},
		},
	})
	defer restore()

	Init()

	svc := GetService()
	if len(svc.Channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(svc.Channels))
	}
}

func TestReloadNotificationService(t *testing.T) {
	restore := setupServiceConfig(t, config.Configs{
		NotificationConfigs: []config.NotificationConfig{
			{
				Type:                 config.Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"*"},
			},
		},
	})
	defer restore()

	Init()
	before := len(GetService().Channels)

	ReloadNotificationService()
	after := len(GetService().Channels)

	if before != after {
		t.Errorf("expected same channel count after reload: before=%d after=%d", before, after)
	}
}
