package app

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	account_handler "scaf-gin/internal/handler/account"
	auth_handler "scaf-gin/internal/handler/auth"
	"scaf-gin/internal/module"
	account_usecase "scaf-gin/internal/usecase/account"
	auth_usecase "scaf-gin/internal/usecase/auth"
)

type App struct {
	Engine *gin.Engine
	Logger core.LoggerI
}

func New() *App {
	log := core.NewJSONLogger()
	core.SetLogger(log)

	if config.AppEnv == "production" || ShouldPrintRoutes() {
		gin.SetMode(gin.ReleaseMode)
	}

	authService := core.NewJwtAuth()
	mailerService := core.NewMailer()

	dbConn := core.NewGormDB()
	accountModule := module.NewAccountModule(dbConn)
	passwordResetTokenModule := module.NewPasswordResetTokenModule(dbConn)

	authUsecase := auth_usecase.NewUsecase(
		dbConn,
		accountModule,
		passwordResetTokenModule,
		authService,
		mailerService,
	)
	accountUsecase := account_usecase.NewUsecase(accountModule)

	accountHandler := account_handler.NewHandler(accountUsecase)
	authHandler := auth_handler.NewHandler(authUsecase)

	engine := gin.New()
	engine.Use(accessLogMiddleware(log))
	engine.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered: %v", recovered)
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     config.FrontendOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	registerAPIRoutes(engine.Group("/api"), accountHandler, authHandler, accountModule, authService)

	return &App{
		Engine: engine,
		Logger: log,
	}
}

func (a *App) Run() error {
	return a.Engine.Run(":" + config.AppPort)
}

func ShouldPrintRoutes() bool {
	return os.Getenv("PRINT_ROUTES") == "true"
}

func PrintRoutes(r *gin.Engine) {
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

func accessLogMiddleware(log core.LoggerI) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		log.Info(
			"access method=%s path=%s status=%d latency_ms=%d client_ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startedAt).Milliseconds(),
			c.ClientIP(),
		)
		if len(c.Errors) > 0 {
			log.Error("request errors: %s", fmt.Sprint(c.Errors))
		}
	}
}
