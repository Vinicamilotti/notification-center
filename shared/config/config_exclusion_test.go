package config

import "testing"

// The excluded-topics set must be built from the JSON list the same way the
// subscribed-topics set is. If the new field is added to the struct but not
// wired into initializeTopicsSet, the map stays nil and exclusion silently
// never works (reading a nil map in Go does not panic — it just returns false).
func TestGetConfigs_ExcludedTopicsSetInitialized(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{
				Type:                 Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"*"},
				ExcludedTopicsList:   []string{"debug", "test"},
			},
		},
	})
	defer restore()

	cfgs := GetConfigs()
	excluded := cfgs.NotificationConfigs[0].ExcludedTopics

	if excluded == nil {
		t.Fatal("expected ExcludedTopics set to be initialized, got nil")
	}
	if _, ok := excluded["debug"]; !ok {
		t.Error("expected 'debug' in ExcludedTopics")
	}
	if _, ok := excluded["test"]; !ok {
		t.Error("expected 'test' in ExcludedTopics")
	}
}

// A config without excluded_topics must still produce a usable (non-nil) set,
// so the dispatch layer can look up topics unconditionally.
func TestGetConfigs_ExcludedTopicsEmptyWhenAbsent(t *testing.T) {
	restore := setupConfigDir(t, Configs{
		NotificationConfigs: []NotificationConfig{
			{
				Type:                 Ntfy,
				Enabled:              true,
				Channel:              "http://ntfy.example.com",
				SubscribedTopicsList: []string{"alerts"},
			},
		},
	})
	defer restore()

	cfgs := GetConfigs()
	excluded := cfgs.NotificationConfigs[0].ExcludedTopics

	if excluded == nil {
		t.Fatal("expected ExcludedTopics set to be non-nil even when no exclusions configured")
	}
	if len(excluded) != 0 {
		t.Errorf("expected empty ExcludedTopics, got %d entries", len(excluded))
	}
}
