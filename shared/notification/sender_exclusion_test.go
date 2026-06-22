package notification

import (
	"testing"

	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

func makeConfigWithExclusions(enabled bool, topics []string, excluded []string) config.NotificationConfig {
	topicSet := make(map[string]any)
	for _, t := range topics {
		topicSet[t] = nil
	}
	excludedSet := make(map[string]any)
	for _, t := range excluded {
		excludedSet[t] = nil
	}
	return config.NotificationConfig{
		Enabled:          enabled,
		SubscribedTopics: topicSet,
		ExcludedTopics:   excludedSet,
	}
}

// Exclusion must win over a wildcard ("*") subscription. This is the case a
// naive implementation overlooks: it sees "*" is subscribed and sends anyway.
func TestSend_ExcludedTopicUnderWildcard(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfigWithExclusions(true, []string{"*"}, []string{"debug"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "debug"})

	if ch.called {
		t.Error("excluded topic 'debug' must NOT be sent even under wildcard subscription")
	}
}

// Exclusion must also win over an explicit subscription to the same topic.
func TestSend_ExcludedTopicExplicitlySubscribed(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfigWithExclusions(true, []string{"debug", "alerts"}, []string{"debug"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "debug"})

	if ch.called {
		t.Error("excluded topic must NOT be sent even when explicitly subscribed")
	}
}

// Existing behavior must be preserved: non-excluded topics still flow through
// the wildcard.
func TestSend_NonExcludedTopicUnderWildcard(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfigWithExclusions(true, []string{"*"}, []string{"debug"})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "alerts"})

	if !ch.called {
		t.Error("non-excluded topic 'alerts' must still be sent under wildcard")
	}
}

// With no exclusions configured the dispatch behaves exactly as before.
func TestSend_NoExclusionsBackwardCompatible(t *testing.T) {
	ns := NewNotificationSender()
	ch := &mockChannel{cfg: makeConfigWithExclusions(true, []string{"*"}, []string{})}
	ns.RegisterChannel(ch)

	ns.Send(domain.NotificationDTO{Topic: "anything"})

	if !ch.called {
		t.Error("wildcard with no exclusions must still send")
	}
}
