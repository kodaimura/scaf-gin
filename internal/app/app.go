package app

import (
	"fmt"
	"net/http"
	"os"
	"sort"

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
	Logger core.Logger
}

func New(cfg config.Config) (*App, error) {
	log := core.NewJSONLogger(cfg.LogLevel)

	if cfg.AppEnv == "production" || ShouldPrintRoutes() {
		gin.SetMode(gin.ReleaseMode)
	}

	authService := core.NewJwtAuth(cfg)
	mailerService := core.NewMailer(cfg)

	dbConn, err := core.NewGormDB(cfg, log)
	if err != nil {
		return nil, err
	}
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
	authHandler := auth_handler.NewHandler(authUsecase, log)

	engine := gin.New()
	engine.Use(accessLogMiddleware(log))
	engine.Use(securityHeadersMiddleware(cfg))
	engine.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.ErrorFields("panic recovered", map[string]any{
			"error_type": "panic",
			"recovered":  fmt.Sprint(recovered),
			"path":       c.Request.URL.String(),
			"method":     c.Request.Method,
		})
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
	}))
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	registerAPIRoutes(engine.Group("/api"), accountHandler, authHandler, accountModule, authService, log)

	return &App{
		Engine: engine,
		Logger: log,
	}, nil
}

func (a *App) Run() error {
	return a.Engine.Run(":8000")
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
