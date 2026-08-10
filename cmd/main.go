package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/adapter/auth"
	"scaf-gin/internal/adapter/logger"
	"scaf-gin/internal/adapter/mailer"
	"scaf-gin/internal/core"
	"scaf-gin/internal/router"
)

func main() {
	core.SetLogger(logger.NewJSONLogger())
	core.SetMailer(mailer.NewMailer())
	core.SetAuth(auth.NewJwtAuth())
	if config.AppEnv == "production" || os.Getenv("PRINT_ROUTES") == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(accessLogMiddleware())
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		core.Logger.Error("panic recovered: %v", recovered)
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.FrontendOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.SetApi(r.Group("/api"))
	if os.Getenv("PRINT_ROUTES") == "true" {
		printRoutes(r)
		return
	}
	r.Run(":" + config.AppPort)
}

func accessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		core.Logger.Info(
			"access method=%s path=%s status=%d latency_ms=%d client_ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startedAt).Milliseconds(),
			c.ClientIP(),
		)
		if len(c.Errors) > 0 {
			core.Logger.Error("request errors: %s", fmt.Sprint(c.Errors))
		}
	}
}

func printRoutes(r *gin.Engine) {
	routes := r.Routes()
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})
	for _, route := range routes {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
}
