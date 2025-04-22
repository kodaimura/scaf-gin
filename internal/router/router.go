package router

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/controller"
	"scaf-gin/internal/infrastructure/db"
	"scaf-gin/internal/middleware"
	repository "scaf-gin/internal/repository/impl"
	"scaf-gin/internal/service"
)

var gorm = db.NewGormDB()
var sqlx = db.NewSqlxDB()

/* DI (Repository) */
var accountRepository = repository.NewGormAccountRepository(gorm)

/* DI (Query) */
//var xxxQuery = query.NewXxxQuery(sqlx)

/* DI (Service) */
var accountService = service.NewAccountService(accountRepository)

/* DI (Controller) */
var indexController = controller.NewIndexController()
var accountController = controller.NewAccountController(accountService)

func SetStatic(r *gin.Engine) {
	r.LoadHTMLGlob("web/template/*.html")
	r.Static("/css", "web/static/css")
	r.Static("/js", "web/static/js")
	r.Static("/img", "web/static/img")
	r.StaticFile("/favicon.ico", "web/static/favicon.ico")
	r.StaticFile("/manifest.json", "web/static/manifest.json")
}

func SetWeb(r *gin.RouterGroup) {
	r.GET("/signup", accountController.SignupPage)
	r.GET("/login", accountController.LoginPage)
	r.GET("/logout", accountController.Logout)

	auth := r.Group("", middleware.Auth())
	{
		auth.GET("/", indexController.IndexPage)
	}
}

func SetApi(r *gin.RouterGroup) {
	r.Use(middleware.ApiErrorHandler())
	r.POST("/accounts/signup", accountController.ApiSignup)
	r.POST("/accounts/login", accountController.ApiLogin)
	r.GET("/accounts/logout", accountController.ApiLogout)

	auth := r.Group("", middleware.ApiAuth())
	{
		auth.GET("/accounts/me", accountController.ApiGetOne)
		auth.PUT("/accounts/me", accountController.ApiPutOne)
		auth.PUT("/accounts/me/password", accountController.ApiPutPassword)
		auth.DELETE("/accounts/me", accountController.ApiDeleteOne)
	}
}
