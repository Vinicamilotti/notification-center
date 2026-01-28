package notification

import "github.com/Vinicamilotti/notification-center/shared/domain"

type NotificationChannel interface {
	SendNotification(dto domain.NotificationDTO) error
}
