package main

import (
	customHandler "github.com/Vinicamilotti/notification-center/internal/customwebhook/handler"
	grafanaApp "github.com/Vinicamilotti/notification-center/internal/grafana/application"
	grafanaHandler "github.com/Vinicamilotti/notification-center/internal/grafana/handler"
	testHandler "github.com/Vinicamilotti/notification-center/internal/testWebhook/handler"
	"github.com/Vinicamilotti/notification-center/lib/app"
	errorlib "github.com/Vinicamilotti/notification-center/lib/errorLib"
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/notification"
	"github.com/joho/godotenv"
)

func bootstrap() {
	err := errorlib.ExecMultipleCanError(func() error { return godotenv.Load() }, config.ReadConfigFile)
	if err != nil {
		panic(err)
	}

	notification.Init()
}

func main() {
	bootstrap()
	app := app.NewApp("192.168.1.200", "9999")
	grafana := grafanaHandler.NewGrafanaWebhookHandler(grafanaApp.NewGrafanaFacade())
	custom := customHandler.NewCustomWebhookHandler()
	test := testHandler.NewTestHandler()

	app.RegisterHandler(test)
	app.RegisterHandler(custom)
	app.RegisterHandler(grafana)

	if err := app.Serve(); err != nil {
		panic(err)
	}

}
