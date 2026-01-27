package application

import (
	"github.com/Vinicamilotti/notification-center/internal/grafana/domain"
	notificationEntities "github.com/Vinicamilotti/notification-center/shared/domain"
)

type GrafanaFacade struct {
	// Add fields if necessary
}

func NewGrafanaFacade() *GrafanaFacade {
	return &GrafanaFacade{}
}

func (f *GrafanaFacade) ProcessAlert(topic string, alert domain.GrafanaAlert) (notificationEntities.NotificationDTO, error) {
	notification := notificationEntities.NotificationDTO{
		Channel: topic, // Example channel
		Title:   alert.Title,
		Message: "test",
	}

	return notification, nil
}
