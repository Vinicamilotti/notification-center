package application

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Vinicamilotti/notification-center/internal/grafana/domain"
	notificationEntities "github.com/Vinicamilotti/notification-center/shared/domain"
)

type GrafanaFacade struct {
	// Add fields if necessary
}

func NewGrafanaFacade() *GrafanaFacade {
	return &GrafanaFacade{}
}

func (f *GrafanaFacade) determinateTags(status string) string {
	stUpper := strings.ToUpper(status)
	if stUpper == "FIRING" {
		return "rotating_light"
	}

	if stUpper == "RESOLVED" {
		return "heavy_check_mark"
	}
	return "warning"
}

func (f *GrafanaFacade) createAttributes(alert domain.GrafanaAlert) map[string]any {
	att := make(map[string]any)
	att["tag"] = f.determinateTags(alert.Status)
	if len(alert.Alerts) > 0 {
		att["click"] = alert.Alerts[0].DashboarUrl
	}

	return att
}

func (f *GrafanaFacade) createMessage(alert domain.GrafanaAlert) string {
	sumary := alert.CommonAnnotations.Summary
	description := alert.CommonAnnotations.Description

	messageTemplate := `
**{{.Sumary}}**
{{.Description}}
- Started at: {{.StartedAt}} 
- Ended at: {{.EndedAt}}
Values: 
{{.Values}}
`

	t, err := template.New("markdown").Parse(messageTemplate)

	if err != nil {
		return fmt.Sprintf("%s %s", sumary, description)
	}
	var msgBuff bytes.Buffer

	data := domain.NewMessageData(alert)
	t.Execute(&msgBuff, data)

	return msgBuff.String()

}

func (f *GrafanaFacade) ProcessAlert(topic string, alert domain.GrafanaAlert) (notificationEntities.NotificationDTO, error) {
	msg := f.createMessage(alert)
	notification := notificationEntities.NotificationDTO{
		Topic:               topic, // Example channel
		Title:               alert.Title,
		AditionalAttributes: f.createAttributes(alert),
		Message:             msg,
	}

	return notification, nil
}
