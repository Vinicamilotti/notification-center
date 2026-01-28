package channel

import (
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type DiscordChannel struct {
	Config config.NotificationConfig
}

func NewDiscordChannel(config config.NotificationConfig) *DiscordChannel {
	return &DiscordChannel{
		Config: config,
	}
}

func (nc *DiscordChannel) GetConfig() config.NotificationConfig {
	return nc.Config
}

func (nc *DiscordChannel) SendNotification(dto domain.NotificationDTO) error {
	return nil
}
