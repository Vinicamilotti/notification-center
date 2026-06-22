package notification

import (
	ntfyChannel "github.com/Vinicamilotti/notification-center/integration/ntfy/channel"
	"github.com/Vinicamilotti/notification-center/shared/config"
)

type NotificationChannelBuilder func(config config.NotificationConfig) NotificationChannel

var notificationService NotificationSender

func Init() {
	availableChannels := map[config.ConfigType]NotificationChannelBuilder{
		config.Ntfy: func(cfg config.NotificationConfig) NotificationChannel {
			return ntfyChannel.NewNtfyChannel(cfg)
		},
		config.Discord: func(cfg config.NotificationConfig) NotificationChannel {
			return nil
		},
	}

	notificationService = *NewNotificationSender()
	notificationChannels := config.GetConfigs().NotificationConfigs
	for _, channelConfig := range notificationChannels {
		if channel, ok := availableChannels[channelConfig.Type]; ok {
			notificationService.RegisterChannel(channel(channelConfig))
		}
	}

}

func GetService() NotificationSender {
	return notificationService
}

func ReloadNotificationService() {
	config.ReadConfigFile()
	Init()
}
