package domain

type MessageData struct {
	Sumary      string
	Status      string
	StartDate   string
	EndDate     string
	Description string
}

func NewMessageData(alert GrafanaAlert) MessageData {
	return MessageData{
		Sumary:      alert.CommonAnnotations.Summary,
		Status:      alert.Status,
		StartDate:   alert.Alerts[0].StartsAt,
		EndDate:     alert.Alerts[0].EndsAt,
		Description: alert.CommonAnnotations.Description,
	}
}
