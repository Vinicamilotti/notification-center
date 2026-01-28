package main

import (
	grafanaApp "github.com/Vinicamilotti/notification-center/internal/grafana/application"
	grafanaHandler "github.com/Vinicamilotti/notification-center/internal/grafana/handler"
	"github.com/Vinicamilotti/notification-center/lib/app"
	errorlib "github.com/Vinicamilotti/notification-center/lib/errorLib"
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/joho/godotenv"
)

func bootstrap() {
	err := errorlib.ExecMultipleCanError(func() error { return godotenv.Load() }, config.ReadConfigFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	bootstrap()
	app := app.NewApp("192.168.1.200", "9999")
	grafana := grafanaHandler.NewGrafanaWebhookHandler(grafanaApp.NewGrafanaFacade())
	app.RegisterHandler(grafana)

	if err := app.Serve(); err != nil {
		panic(err)
	}

}
