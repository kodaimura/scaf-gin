package main

import (
	"scaf-gin/internal/app"
)

func main() {
	application := app.New()
	if app.ShouldPrintRoutes() {
		app.PrintRoutes(application.Engine)
		return
	}
	if err := application.Run(); err != nil {
		application.Logger.Error("server stopped: %v", err)
	}
}
