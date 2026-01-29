package domain

import (
	"fmt"
)

type MessageData struct {
	Sumary      string
	Status      string
	StartedAt   string
	EndedAt     string
	Description string
	Values      string
}

func handleValues(values map[string]any) string {
	valuesString := ""
	for k, v := range values {
		valuesString = fmt.Sprintf("%s\n- %s = %v", valuesString, k, v)

	}

	return valuesString
}

func NewMessageData(alert GrafanaAlert) MessageData {
	return MessageData{
		Sumary:      alert.CommonAnnotations.Summary,
		Status:      alert.Status,
		StartedAt:   alert.Alerts[0].StartsAt,
		EndedAt:     alert.Alerts[0].EndsAt,
		Description: alert.CommonAnnotations.Description,
		Values:      handleValues(alert.Alerts[0].Values),
	}
}
