package handler

import (
	"fmt"

	"github.com/Vinicamilotti/notification-center/shared/domain"
	"github.com/Vinicamilotti/notification-center/shared/notification"
	"github.com/gin-gonic/gin"
)

type TestHandler struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/test-webhook/:topic", h.TestWebhook)
}

func (h *TestHandler) TestWebhook(c *gin.Context) {
	topic := c.Param("topic")
	testNotification := domain.NotificationDTO{
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
				Action: fmt.Sprintf("ntfy://192.168.1.200:9999/%s", topic),
			},
			{
				Type:   domain.ActionTypeHttpCall,
				Label:  "Send another one",
				Action: fmt.Sprintf("https://192.168.1.200:9999/test-webhook/%s", topic),
				Method: "POST",
			},
		},
	}

	sender := notification.GetService()
	sender.Send(testNotification)

	c.JSON(200, gin.H{
		"message": "Test webhook handler is working!",
	})
}
