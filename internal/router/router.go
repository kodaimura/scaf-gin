package router

import (
	"net/http"	
	"github.com/gin-gonic/gin"

	"goscaf/internal/middleware"
)

func SetStatic(r *gin.Engine) {	
	r.LoadHTMLGlob("web/template/*.html")
	r.Static("/css", "web/static/css")
	r.Static("/js", "web/static/js")
	r.Static("/img", "web/static/img")
	r.StaticFile("/favicon.ico", "web/static/favicon.ico")
	r.StaticFile("/manifest.json", "web/static/manifest.json")
}

var ic = controller.NewIndexController()
var ac = controller.NewAccountController()

func SetWeb(r *gin.Engine) {
	r.GET("/signup", ac.SignupPage)
	r.GET("/login", ac.LoginPage)
	r.GET("/logout", ac.Logout)

	auth := r.Group("", middleware.JwtAuth())
	{
		auth.GET("/", ic.IndexPage)
	}
}

func SetApi(r *gin.Engine) {
	r.Use(middleware.ApiErrorHandler())
	r.POST("/signup", ac.ApiSignup)
	r.POST("/login", ac.ApiLogin)

	auth := r.Group("", middleware.ApiJwtAuth())
	{
		auth.GET("/accounts/me", ac.ApiGetOne)
		auth.PUT("/accounts/me", ac.ApiPutOne)
		auth.PUT("/accounts/me/password", ac.ApiPutPassword)
		auth.DELETE("/accounts/me", ac.ApiDeleteOne)
	}
}