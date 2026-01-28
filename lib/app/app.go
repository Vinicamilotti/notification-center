package app

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	Port     string
	Listen   string
	Handlers []Handler
}

func NewApp(listen string, port string) *App {
	return &App{
		Port:     port,
		Listen:   listen,
		Handlers: []Handler{},
	}
}

func (a *App) RegisterHandler(h Handler) {
	a.Handlers = append(a.Handlers, h)
}

func (a *App) Serve() error {
	r := gin.Default()

	for _, h := range a.Handlers {
		h.RegisterRoutes(r)
	}

	return r.Run(a.Listen + ":" + a.Port)
}
