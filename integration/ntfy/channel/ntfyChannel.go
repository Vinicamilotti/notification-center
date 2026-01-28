package channel

import (
	"net/http"

	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type NtfyChannel struct {
	client http.Client
	Config config.NotificationConfig
}

func NewNtfyChannel(config config.NotificationConfig) *NtfyChannel {
	return &NtfyChannel{
		client: http.Client{},
		Config: config,
	}
}

func (nc *NtfyChannel) GetConfig() config.NotificationConfig {
	return nc.Config
}

func (nc *NtfyChannel) SendNotification(dto domain.NotificationDTO) error {
	return nil
}
