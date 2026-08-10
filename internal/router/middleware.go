package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/helper"
	"scaf-gin/internal/service"
)

// BasicAuthMiddleware is a middleware that checks for Basic Authentication credentials.
// If the credentials are incorrect, it returns an Unauthorized status.
func BasicAuthMiddleware() gin.HandlerFunc {
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

// ApiAuthMiddleware is a middleware that validates the JWT token for API access.
// If the token is invalid, it returns an Unauthorized error in JSON format.
func ApiAuthMiddleware(dbConn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAccessToken(c)
		if token == "" {
			c.Error(core.ErrAuthMissing)
			c.Abort()
			return
		}
		payload, err := core.Auth.VerifyAccessToken(token)
		if err != nil {
			c.Error(core.ErrAuthInvalid)
			c.Abort()
			return
		}

		accountID, err := service.ValidateAccessTokenAccount(dbConn, payload)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		payload.AccountId = accountID

		c.Set(helper.CONTEXT_KEY_PAYLOAD, payload)
		c.Next()
	}
}

// ApiErrorHandler is a middleware that handles API errors.
// It checks for specific error types and returns the appropriate HTTP status and message.
func ApiErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		status := http.StatusInternalServerError
		resp := gin.H{
			"message": "Unexpected Error",
		}

		var appErr *core.AppError
		if errors.As(err, &appErr) {
			status = appErr.HTTPStatus
			resp = gin.H{
				"code":    appErr.ErrorCode,
				"message": appErr.Message,
			}
			if len(appErr.ErrorDetails) > 0 {
				resp["details"] = appErr.ErrorDetails
			}
		}

		if status >= 500 {
			core.Logger.Error(
				"Error: %v method=%s url=%s headers=%v",
				err,
				c.Request.Method,
				c.Request.URL.String(),
				c.Request.Header,
			)
		}

		c.JSON(status, resp)
	}
}
