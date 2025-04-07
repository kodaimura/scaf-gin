package main

import (
	"io"
	"os"
	"github.com/gin-gonic/gin"

	"goscaf/config"
	"goscaf/internal/infrastructure/logger"
	"goscaf/internal/infrastructure/mailer"
	"goscaf/internal/infrastructure/auth"
	"goscaf/internal/core"
	"goscaf/internal/router"
)

func main() {
    f, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
    gin.DefaultWriter = io.MultiWriter(os.Stdout, f)

	core.SetLogger(logger.NewMultiLogger(f))
	core.SetMailer(mailer.NewMockMailer())
	core.SetAuth(auth.NewJwtAuth())

	r := gin.Default()
	router.SetStatic(r)
	router.SetWeb(r.Group("/"))
	router.SetApi(r.Group("/api"))
	r.Run(":" + config.AppPort)
}