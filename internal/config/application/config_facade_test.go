package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
)

func setupConfigDir(t *testing.T, cfg config.Configs) {
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
	t.Cleanup(func() { os.Chdir(oldDir) })
}

func TestAddConfigure_ConfigIsReturnedByGetConfigs(t *testing.T) {
	setupConfigDir(t, config.Configs{NotificationConfigs: []config.NotificationConfig{}})

	facade := NewConfigFacade()

	newCfg := config.NotificationConfig{
		Type:    config.Ntfy,
		Enabled: true,
		Channel: "added-channel",
	}

	id, err := facade.AddConfigure(newCfg)
	if err != nil {
		t.Fatalf("AddConfigure returned error: %v", err)
	}
	if id == "" {
		t.Fatal("expected AddConfigure to return a non-empty UUID")
	}

	cfgs, err := config.GetConfigs()
	if err != nil {
		t.Fatalf("GetConfigs returned error: %v", err)
	}

	found := false
	for _, c := range cfgs.NotificationConfigs {
		if c.ID == id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected GetConfigs to contain a config with ID %q, got %+v",
			id, cfgs.NotificationConfigs)
	}
}
