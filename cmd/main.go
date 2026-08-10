package main

import (
	"scaf-gin/config"
	"scaf-gin/internal/app"
	"scaf-gin/internal/core"
)

func main() {
	cfg := config.Current
	application, err := app.New(cfg)
	if err != nil {
		core.NewJSONLogger(cfg.LogLevel).Error("startup error: %v", err)
		return
	}
	if app.ShouldPrintRoutes() {
		app.PrintRoutes(application.Engine)
		return
	}
	if err := application.Run(); err != nil {
		application.Logger.Error("server stopped: %v", err)
	}
}
