package router

import (
	"net/http"	
	"github.com/gin-gonic/gin"

	"goscaf/internal/middleware"
)

func SetStatic(r *gin.Engine) {	
	//r.LoadHTMLGlob("web/template/*.html")
	//r.Static("/css", "web/static/css")
	//r.Static("/js", "web/static/js")
	//r.Static("/img", "web/static/img")
	//r.StaticFile("/favicon.ico", "web/static/favicon.ico")
	//r.StaticFile("/manifest.json", "web/static/manifest.json")

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
}

func SetWeb(r *gin.Engine) {

}

func SetApi(r *gin.Engine) {
	r.Use(middleware.ApiErrorHandler())

}