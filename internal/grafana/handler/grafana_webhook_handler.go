package handler

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Vinicamilotti/notification-center/internal/grafana/application"
	"github.com/Vinicamilotti/notification-center/internal/grafana/domain"
	"github.com/Vinicamilotti/notification-center/shared/notification"
	"github.com/gin-gonic/gin"
)

type GrafanaWebhookHandler struct {
	sender        notification.NotificationSender
	grafanaFacade *application.GrafanaFacade
}

func NewGrafanaWebhookHandler(grafanaFacade *application.GrafanaFacade) *GrafanaWebhookHandler {
	return &GrafanaWebhookHandler{
		grafanaFacade: grafanaFacade,
	}
}

func (h *GrafanaWebhookHandler) HandleWebhook(ctx *gin.Context) {
	getTopic := ctx.Param("topic")
	getBody := ctx.Request.Body
	bodyAsByteArray, err := io.ReadAll(getBody)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	var alert domain.GrafanaAlert
	if err = json.Unmarshal(bodyAsByteArray, &alert); err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	dto, err := h.grafanaFacade.ProcessAlert(getTopic, alert)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to process alert"})
		return
	}

	sender := notification.GetService()
	errors := sender.Send(dto)

	for _, sendErr := range errors {
		fmt.Printf("Error sending notification: %v\n", sendErr)
	}

	ctx.JSON(200, gin.H{"status": "success", "notification": dto})

}

func (h *GrafanaWebhookHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/grafana/:topic", h.HandleWebhook)
}
