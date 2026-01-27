package handler

import (
	"encoding/json"
	"io"

	"github.com/Vinicamilotti/notification-center/internal/grafana/application"
	"github.com/Vinicamilotti/notification-center/internal/grafana/domain"
	"github.com/gin-gonic/gin"
)

type GrafanaWebhookHandler struct {
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

	notification, err := h.grafanaFacade.ProcessAlert(getTopic, alert)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to process alert"})
		return
	}

	ctx.JSON(200, gin.H{"status": "success", "notification": notification})

}

func (h *GrafanaWebhookHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/grafana/:topic", h.HandleWebhook)
}
