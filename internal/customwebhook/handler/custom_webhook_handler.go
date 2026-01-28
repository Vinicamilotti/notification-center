package handler

import (
	"encoding/json"

	"github.com/Vinicamilotti/notification-center/shared/domain"
	"github.com/Vinicamilotti/notification-center/shared/notification"
	"github.com/gin-gonic/gin"
)

type CustomWebhookHandler struct {
}

func NewCustomWebhookHandler() *CustomWebhookHandler {
	return &CustomWebhookHandler{}
}

func (h *CustomWebhookHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/custom-webhook/:topic", h.HandleCustomWebhook)
}

func (h *CustomWebhookHandler) HandleCustomWebhook(c *gin.Context) {
	topic := c.Param("topic")
	getBody := c.Request.Body

	var payload domain.NotificationDTO

	err := json.NewDecoder(getBody).Decode(&payload)
	payload.Topic = topic
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}

	sender := notification.GetService()
	sender.Send(payload)

	c.JSON(200, gin.H{
		"message": "Custom webhook received",
	})

}
