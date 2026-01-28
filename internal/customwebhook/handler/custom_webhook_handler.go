package handler

import (
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
	sender := notification.GetService()
	sender.Send(domain.NotificationDTO{
		Title:   "Custom Webhook",
		Message: "A custom webhook was received.",
		Topic:   topic,
		AditionalAttributes: map[string]any{
			"tag":   "warning",
			"click": "https://example.com",
		},
		Actions: []domain.NotificationAction{
			{
				Type:   domain.ActionTypeUrl,
				Label:  "View",
				Action: "https://example.com/view",
			},
			{
				Type:   domain.ActionTypeHttpCall,
				Label:  "Acknowledge",
				Action: "https://example.com/acknowledge",
			},
		},
	})

	c.JSON(200, gin.H{
		"message": "Custom webhook received",
	})
}
