package router

import (
	"github.com/gin-gonic/gin"

	"goscaf/internal/middleware"
	"goscaf/internal/infrastructure/db"
	"goscaf/internal/controller"
	"goscaf/internal/service"
	repository "goscaf/internal/repository/impl"
)

var gorm = db.NewGormDB()
var sqlx = db.NewSqlxDB()

/* DI (Repository) */
var accountRepository = repository.NewGormAccountRepository(gorm)

/* DI (Query) */
//var xxxQuery = query.NewXxxQuery(gorm)

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

	auth := r.Group("", middleware.JwtAuth())
	{
		auth.GET("/", indexController.IndexPage)
	}
}

func SetApi(r *gin.RouterGroup) {
	r.Use(middleware.ApiErrorHandler())
	r.POST("/signup", accountController.ApiSignup)
	r.POST("/login", accountController.ApiLogin)
	r.GET("/logout", accountController.ApiLogout)

	auth := r.Group("", middleware.ApiJwtAuth())
	{
		auth.GET("/accounts/me", accountController.ApiGetOne)
		auth.PUT("/accounts/me", accountController.ApiPutOne)
		auth.PUT("/accounts/me/password", accountController.ApiPutPassword)
		auth.DELETE("/accounts/me", accountController.ApiDeleteOne)
	}
}