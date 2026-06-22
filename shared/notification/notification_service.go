package notification

import (
	ntfyChannel "github.com/Vinicamilotti/notification-center/integration/ntfy/channel"
	"github.com/Vinicamilotti/notification-center/shared/config"
)

type NotificationChannelBuilder func(config config.NotificationConfig) NotificationChannel

var notificationService NotificationSender

var initialized = false

func Init(cfg config.Configs) {
	if initialized {
		return
	}

	availableChannels := map[config.ConfigType]NotificationChannelBuilder{
		config.Ntfy: func(cfg config.NotificationConfig) NotificationChannel {
			return ntfyChannel.NewNtfyChannel(cfg)
		},
		config.Discord: func(cfg config.NotificationConfig) NotificationChannel {
			return nil
		},
	}

	notificationService = *NewNotificationSender()

	for _, channelConfig := range cfg.NotificationConfigs {
		if channel, ok := availableChannels[channelConfig.Type]; ok {
			notificationService.RegisterChannel(channel(channelConfig))
		}
	}

	initialized = true

}

func GetService() NotificationSender {
	return notificationService
}
