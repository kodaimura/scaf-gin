package main

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/adapter/auth"
	"scaf-gin/internal/adapter/file"
	"scaf-gin/internal/adapter/logger"
	"scaf-gin/internal/adapter/mailer"
	"scaf-gin/internal/core"
	"scaf-gin/internal/router"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.SetStatic(r)
	router.SetWeb(r.Group("/"))
	router.SetApi(r.Group("/api"))
	r.Run(":" + config.AppPort)
}
