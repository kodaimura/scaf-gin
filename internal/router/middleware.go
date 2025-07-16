package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/helper"
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

// WebAuthMiddleware is a middleware that validates the JWT access token and refresh token.
// If both are invalid, the user is redirected to login,
// otherwise, a new access token is created and stored in a cookie.
func WebAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to verify the access token
		payload, err := core.Auth.VerifyAccessToken(helper.GetAccessToken(c))
		if err != nil {
			// If access token is invalid, try verifying the refresh token
			payload, err = core.Auth.VerifyRefreshToken(helper.GetRefreshToken(c))
			if err != nil {
				c.Error(core.ErrUnauthorized)
				c.Abort()
				return
			}

			// If refresh token is valid, create a new access token
			accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
				AccountId:   payload.AccountId,
				AccountName: payload.AccountName,
			})
			if err != nil {
				c.Error(core.ErrUnauthorized)
				c.Abort()
				return
			}

			// Verify the newly created access token
			payload, err := core.Auth.VerifyAccessToken(accessToken)
			if err != nil {
				c.Error(core.ErrUnauthorized)
				c.Abort()
				return
			}

			// Set the new access token in a cookie
			helper.SetAccessTokenCookie(c, accessToken)
			core.Logger.Info("access token refreshed: id=%d name=%s", payload.AccountId, payload.AccountName)
		}

		c.Set(helper.CONTEXT_KEY_PAYLOAD, payload)
		c.Next()
	}
}

// WebErrorHandler is a middleware that handles errors for web pages.
// It renders appropriate error pages based on the error type.
func WebErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		status := http.StatusInternalServerError
		template := "500.html"

		var appErr *core.AppError
		if errors.As(err, &appErr) {
			status = appErr.HTTPStatus
			switch status {
			case http.StatusUnauthorized, http.StatusForbidden:
				c.Redirect(http.StatusFound, "/login")
				return
			case http.StatusNotFound:
				template = "404.html"
			}
		} else {
			status = http.StatusInternalServerError
			template = "500.html"
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

		c.HTML(status, template, nil)
	}
}

// ApiAuthMiddleware is a middleware that validates the JWT token for API access.
// If the token is invalid, it returns an Unauthorized error in JSON format.
func ApiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAccessToken(c)
		payload, err := core.Auth.VerifyAccessToken(token)
		if err != nil {
			c.Error(core.ErrUnauthorized)
			c.Abort()
			return
		}

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
