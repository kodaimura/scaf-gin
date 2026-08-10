package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	handlerutil "scaf-gin/internal/handler"
	"scaf-gin/internal/module"
	"scaf-gin/internal/service"
)

func basicAuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != cfg.BasicAuthUser || pass != cfg.BasicAuthPass {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func apiAuthMiddleware(accountModule module.AccountModule, authService core.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := handlerutil.GetAccessToken(c)
		if token == "" {
			c.Error(core.ErrAuthMissing)
			c.Abort()
			return
		}
		payload, err := authService.VerifyAccessToken(token)
		if err != nil {
			c.Error(core.ErrAuthInvalid)
			c.Abort()
			return
		}

		accountID, err := service.ValidateAccessTokenAccount(accountModule, payload)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		payload.AccountId = accountID

		handlerutil.SetPayload(c, payload)
		c.Next()
	}
}

func apiErrorHandler(log core.Logger) gin.HandlerFunc {
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
			log.Error(
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
