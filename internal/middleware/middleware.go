package middleware

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"

	"goscaf/config"
	"goscaf/pkg/logger"
	"goscaf/internal/commom"
)


func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != config.BasicAuthUser || pass != config.BasicAuthPass {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}



func getAccessToken (c *gin.Context) string {
	token, err := c.Cookie(common.COOKIE_KEY_ACCESS_TOKEN)
	if err == nil {
		return token
	}

	bearer := c.Request.Header.Get("Authorization")
	if bearer != "" && !strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimSpace(bearer[7:])
	}

	return ""
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAccessToken(c)
		pl, err := jwt.DecodeToken(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set(common.CONTEXT_KEY_JWT_PAYLOAD, pl)
		c.Next()
	}
}


func ApiJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {	
		token := getAccessToken(c)
		pl, err := jwt.DecodeToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		
		c.Set(common.CONTEXT_KEY_JWT_PAYLOAD, pl)
		c.Next()
	}
}

func ApiErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case errs.BadRequestError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(), 
				})
			case errs.UnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": e.Error(),
				})
			case errs.ForbiddenError:
				c.JSON(http.StatusForbidden, gin.H{
					"error": e.Error(),
				})
			case errs.NotFoundError:
				c.JSON(http.StatusNotFound, gin.H{
					"error": e.Error(),
				})
			case errs.ConflictError:
				c.JSON(http.StatusConflict, gin.H{
					"error": e.Error(),
				})
			default:
				logger.Error(e.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": e.Error(),
				})
			}
		}
	}
}