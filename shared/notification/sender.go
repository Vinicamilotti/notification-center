package notification

import (
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type NotificationSender struct {
	Channels []NotificationChannel
}

func NewNotificationSender() *NotificationSender {
	return &NotificationSender{
		Channels: make([]NotificationChannel, 0),
	}
}

func (ns *NotificationSender) RegisterChannel(channel NotificationChannel) {
	ns.Channels = append(ns.Channels, channel)
}

func (ns *NotificationSender) Send(dto domain.NotificationDTO) []error {
	var errors []error
	for _, channel := range ns.Channels {
		config := channel.GetConfig()
		if !config.Enabled {
			continue
		}

		_, subscribed := config.SubscribedTopics[dto.Topic]
		_, subscribedAll := config.SubscribedTopics["*"]

		if !subscribed && !subscribedAll {
			continue
		}

		if err := channel.SendNotification(dto); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
