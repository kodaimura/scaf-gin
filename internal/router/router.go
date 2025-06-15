package router

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/domain/account"
	"scaf-gin/internal/domain/auth"
	"scaf-gin/internal/domain/general"
	"scaf-gin/internal/infrastructure/db"
	"scaf-gin/internal/middleware"
)

var gorm = db.NewGormDB()
//var sqlx = db.NewSqlxDB()

/* DI (Repository) */
var accountRepository = account.NewRepository()

/* DI (Service) */
var authService = auth.NewService(accountRepository)
var accountService = account.NewService(accountRepository)

/* DI (Controller) */
var authController = auth.NewController(gorm, authService)
var accountController = account.NewController(gorm, accountService)
var generalController = general.NewController()

func SetStatic(r *gin.Engine) {
	r.LoadHTMLGlob("web/template/*.html")
	r.Static("/css", "web/static/css")
	r.Static("/js", "web/static/js")
	r.Static("/img", "web/static/img")
	r.StaticFile("/favicon.ico", "web/static/favicon.ico")
	r.StaticFile("/manifest.json", "web/static/manifest.json")
	r.NoRoute(func(c *gin.Context) { c.HTML(404, "404.html", nil) })
}

func SetWeb(r *gin.RouterGroup) {
	r.Use(middleware.WebErrorHandler())
	r.GET("/signup", authController.SignupPage)
	r.GET("/login", authController.LoginPage)
	r.GET("/logout", authController.Logout)

	auth := r.Group("", middleware.WebAuth())
	{
		auth.GET("/", generalController.IndexPage)
	}
}

func SetApi(r *gin.RouterGroup) {
	r.Use(middleware.ApiErrorHandler())
	r.POST("/accounts/signup", authController.ApiSignup)
	r.POST("/accounts/login", authController.ApiLogin)
	r.POST("/accounts/refresh", authController.ApiRefresh)
	r.POST("/accounts/logout", authController.ApiLogout)

	auth := r.Group("", middleware.ApiAuth())
	{
		auth.GET("/accounts/me", accountController.ApiGetMe)
		auth.PUT("/accounts/me", accountController.ApiPutMe)
		auth.PUT("/accounts/me/password", authController.ApiPutMePassword)
		auth.DELETE("/accounts/me", accountController.ApiDeleteMe)
	}
}
