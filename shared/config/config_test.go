package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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

func TestReadConfigFile_Valid(t *testing.T) {
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

	if err := ReadConfigFile(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReadConfigFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer func() {
		os.Chdir(oldDir)
		initialized = false
	}()
	os.Chdir(dir)
	initialized = false

	err := ReadConfigFile()
	if err == nil {
		t.Error("expected error for missing config.json")
	}
}

func TestReadConfigFile_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid json"), 0644)
	oldDir, _ := os.Getwd()
	defer func() {
		os.Chdir(oldDir)
		initialized = false
	}()
	os.Chdir(dir)
	initialized = false

	err := ReadConfigFile()
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestGetConfigs_ReturnsNotificationConfigs(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{Type: Discord, Enabled: false, Channel: "https://discord.com/webhook/abc"},
			{Type: Ntfy, Enabled: true, Channel: "http://ntfy.example.com"},
		},
	})
	defer restore()

	cfgs := GetConfigs()
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

	cfgs := GetConfigs()
	topics := cfgs.NotificationConfigs[0].SubscribedTopics

	if _, ok := topics["alerts"]; !ok {
		t.Error("expected 'alerts' in SubscribedTopics")
	}
	if _, ok := topics["metrics"]; !ok {
		t.Error("expected 'metrics' in SubscribedTopics")
	}
	if _, ok := topics["*"]; !ok {
		t.Error("expected '*' in SubscribedTopics")
	}
}

func TestGetConfigs_Cached(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{Type: Ntfy, Enabled: true, Channel: "http://ntfy.example.com"},
		},
	})
	defer restore()

	cfgs1 := GetConfigs()
	cfgs2 := GetConfigs()

	if len(cfgs1.NotificationConfigs) != len(cfgs2.NotificationConfigs) {
		t.Error("expected same result from cached GetConfigs calls")
	}
}

func TestReadConfigFile_ResetsCache(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{Type: Ntfy, Enabled: true, Channel: "http://ntfy.example.com"},
		},
	})
	defer restore()

	GetConfigs()

	if err := ReadConfigFile(); err != nil {
		t.Fatalf("unexpected error on second read: %v", err)
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

func TestGetConfigs_EmptyConfigs(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{},
	})
	defer restore()

	cfgs := GetConfigs()
	if len(cfgs.NotificationConfigs) != 0 {
		t.Errorf("expected 0 configs, got %d", len(cfgs.NotificationConfigs))
	}
}
