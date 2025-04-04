package middleware

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"

	"goscaf/pkg/logger"
)


//func BasicAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		cf := config.GetConfig()
//
//		user, pass, ok := c.Request.BasicAuth()
//		if !ok || user != cf.BasicAuthUser || pass != cf.BasicAuthPass {
//			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
//			c.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//		c.Next()
//	}
//}
//
//
//func JwtAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if err := jwt.Auth(c); err != nil {
//			c.Redirect(http.StatusSeeOther, "/login")
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
//
//
//func ApiJwtAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if err := jwt.Auth(c); err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}


func ApiErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error(err.Error())
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
				"error": fmt.Sprintf("%v", c.Errors[0].Error()),
			})
		}
	}
}