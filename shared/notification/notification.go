package notification

import (
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type NotificationChannel interface {
	GetConfig() config.NotificationConfig
	SendNotification(dto domain.NotificationDTO) error
}
