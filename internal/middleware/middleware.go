package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/helper"
)

// BasicAuth is a middleware that checks for Basic Authentication credentials.
// If the credentials are incorrect, it returns an Unauthorized status.
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

// Auth is a middleware that validates the JWT token for authenticated access.
// If the token is invalid, it redirects to the login page.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAccessToken(c)
		pl, err := core.Auth.ValidateToken(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set(helper.CONTEXT_KEY_PAYLOAD, pl)
		c.Next()
	}
}

// ApiAuth is a middleware that validates the JWT token for API access.
// If the token is invalid, it returns an Unauthorized error in JSON format.
func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAccessToken(c)
		pl, err := core.Auth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set(helper.CONTEXT_KEY_PAYLOAD, pl)
		c.Next()
	}
}

// ApiErrorHandler is a middleware that handles API errors.
// It checks for specific error types and returns the appropriate HTTP status and message.
func ApiErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch {
			case errors.Is(err, core.ErrBadRequest):
				if appErr, ok := err.(*core.AppError); ok {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": appErr.Error(),
						"details": appErr.Details(),
					})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": err.Error(),
					})
				}
				break
			case errors.Is(err, core.ErrUnauthorized):
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": err.Error(),
				})
				break
			case errors.Is(err, core.ErrForbidden):
				c.JSON(http.StatusForbidden, gin.H{
					"message": err.Error(),
				})
				break
			case errors.Is(err, core.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{
					"message": err.Error(),
				})
				break
			case errors.Is(err, core.ErrConflict):
				c.JSON(http.StatusConflict, gin.H{
					"message": err.Error(),
				})
				break
			case errors.Is(err, core.ErrUnexpected):
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				break
			default:
				core.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
		}
	}
}
