package main

import (
	"io"
	"os"
	"github.com/gin-gonic/gin"

	"goscaf/config"
	"goscaf/internal/infrastructure/logger"
	"goscaf/internal/infrastructure/mailer"
	"goscaf/internal/infrastructure/auth"
	"goscaf/internal/infrastructure/file"
	"goscaf/internal/core"
	"goscaf/internal/router"
)

func main() {
	f1 := file.GetAccessLogFile()
	f2 := file.GetAppLogFile()
	defer f1.Close()
	defer f2.Close()

    gin.DefaultWriter = io.MultiWriter(os.Stdout, f1)

	core.SetLogger(logger.NewMultiLogger(f2))
	core.SetMailer(mailer.NewSmtpMailer())
	core.SetAuth(auth.NewJwtAuth())

	r := gin.Default()
	router.SetStatic(r)
	router.SetWeb(r.Group("/"))
	router.SetApi(r.Group("/api"))
	r.Run(":" + config.AppPort)
}