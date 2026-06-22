package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// setupConfigDir writes the given config as config.json into a temp dir, makes
// it the working directory and resets the package cache. The returned function
// restores the previous working directory and the cache state.
func setupConfigDir(t *testing.T, cfg Configs) (restore func()) {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644); err != nil {
		t.Fatalf("write config error: %v", err)
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd error: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir error: %v", err)
	}

	initialized = false
	return func() {
		os.Chdir(oldDir)
		initialized = false
	}
}

func TestGetConfigs_Valid(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{
				Type:                 Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"alerts", "metrics"},
			},
		},
	})
	defer restore()

	cfgs, err := GetConfigs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs.NotificationConfigs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(cfgs.NotificationConfigs))
	}
	if cfgs.NotificationConfigs[0].Type != Ntfy {
		t.Errorf("expected type Ntfy, got %q", cfgs.NotificationConfigs[0].Type)
	}
}

func TestGetConfigs_MissingFile(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer func() {
		os.Chdir(oldDir)
		initialized = false
	}()
	os.Chdir(dir)
	initialized = false

	_, err := GetConfigs()
	if err == nil {
		t.Error("expected error for missing config.json")
	}
}

func TestGetConfigs_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("write config error: %v", err)
	}
	oldDir, _ := os.Getwd()
	defer func() {
		os.Chdir(oldDir)
		initialized = false
	}()
	os.Chdir(dir)
	initialized = false

	_, err := GetConfigs()
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestGetConfigs_ReturnsAllNotificationConfigs(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{Type: Discord, Enabled: false, Channel: "https://discord.com/webhook/abc"},
			{Type: Ntfy, Enabled: true, Channel: "http://ntfy.example.com"},
		},
	})
	defer restore()

	cfgs, err := GetConfigs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs.NotificationConfigs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(cfgs.NotificationConfigs))
	}
	if cfgs.NotificationConfigs[0].Type != Discord {
		t.Errorf("expected first config type Discord, got %q", cfgs.NotificationConfigs[0].Type)
	}
	if cfgs.NotificationConfigs[1].Type != Ntfy {
		t.Errorf("expected second config type Ntfy, got %q", cfgs.NotificationConfigs[1].Type)
	}
}

func TestGetConfigs_TopicsSetInitialized(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{
				Type:                 Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"alerts", "metrics", "*"},
			},
		},
	})
	defer restore()

	cfgs, err := GetConfigs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	topics := cfgs.NotificationConfigs[0].SubscribedTopics

	for _, want := range []string{"alerts", "metrics", "*"} {
		if _, ok := topics[want]; !ok {
			t.Errorf("expected %q in SubscribedTopics", want)
		}
	}
	if len(topics) != 3 {
		t.Errorf("expected 3 topics, got %d", len(topics))
	}
}

func TestGetConfigs_Cached(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{Type: Ntfy, Enabled: true, Channel: "http://ntfy.example.com"},
		},
	})
	defer restore()

	if _, err := GetConfigs(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Remove the file: a second call must succeed from cache without re-reading.
	if err := os.Remove("config.json"); err != nil {
		t.Fatalf("remove config error: %v", err)
	}

	cfgs, err := GetConfigs()
	if err != nil {
		t.Fatalf("expected cached result, got error: %v", err)
	}
	if len(cfgs.NotificationConfigs) != 1 {
		t.Errorf("expected 1 cached config, got %d", len(cfgs.NotificationConfigs))
	}
}

func TestGetConfigs_EmptyConfigs(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{},
	})
	defer restore()

	cfgs, err := GetConfigs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs.NotificationConfigs) != 0 {
		t.Errorf("expected 0 configs, got %d", len(cfgs.NotificationConfigs))
	}
}

func TestConfigTypeConstants(t *testing.T) {
	if Ntfy != "ntfy" {
		t.Errorf("expected Ntfy to be 'ntfy', got %q", Ntfy)
	}
	if Discord != "discord" {
		t.Errorf("expected Discord to be 'discord', got %q", Discord)
	}
}
