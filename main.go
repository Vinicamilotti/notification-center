package main

import (
	"github.com/Vinicamilotti/notification-center/internal/grafana/application"
	"github.com/Vinicamilotti/notification-center/internal/grafana/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	handler := handler.NewGrafanaWebhookHandler(application.NewGrafanaFacade())
	handler.RegisterRoutes(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run(":9999")
}
