package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	handlerutil "scaf-gin/internal/handler"
)

type accessTokenAuthorizer interface {
	AuthorizeAccessToken(payload core.AuthPayload) (int64, error)
}

func apiAuthMiddleware(authorizer accessTokenAuthorizer, authService core.Auth) gin.HandlerFunc {
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

		accountID, err := authorizer.AuthorizeAccessToken(payload)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		payload.AccountID = accountID

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
			"message": "Internal server error",
		}

		var appErr *core.AppError
		if errors.As(err, &appErr) {
			status = appErr.HTTPStatus
			if appErr.ErrorCode == core.ErrCodeUnprocessable {
				resp = gin.H{
					"message": "Validation error",
					"errors":  appErr.ErrorDetails,
				}
			} else {
				resp = gin.H{
					"code": appErr.ErrorCode,
				}
				if len(appErr.ErrorDetails) > 0 {
					resp["details"] = appErr.ErrorDetails
				}
			}
		}

		if status >= 500 {
			log.ErrorFields("unexpected error occurred", errorLogFields(c, err, status))
		} else if isValidationError(err) {
			log.WarnFields("validation error occurred", errorLogFields(c, err, status))
		} else {
			log.WarnFields("application error occurred", errorLogFields(c, err, status))
		}

		handlerutil.Error(c, status, resp)
	}
}

func accessLogMiddleware(log core.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		if c.Request.Method == http.MethodOptions {
			return
		}

		var account any
		if accountID := handlerutil.GetAccountID(c); accountID != 0 {
			account = map[string]any{"id": accountID}
		}

		log.InfoFields("access", map[string]any{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": time.Since(startedAt).Milliseconds(),
			"client":      c.ClientIP(),
			"account":     account,
		})
	}
}

func securityHeadersMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin")

		if cfg.AppEnv == "production" {
			c.Header("Content-Security-Policy", "default-src 'self'")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

func errorLogFields(c *gin.Context, err error, status int) map[string]any {
	fields := map[string]any{
		"status_code": status,
		"path":        c.Request.URL.String(),
		"method":      c.Request.Method,
		"error":       err.Error(),
	}
	if accountID := handlerutil.GetAccountID(c); accountID != 0 {
		fields["account_id"] = accountID
	}

	var appErr *core.AppError
	if errors.As(err, &appErr) {
		fields["error_code"] = appErr.ErrorCode
		if appErr.ErrorCode == core.ErrCodeUnprocessable {
			fields["error_type"] = "validation_error"
			fields["errors"] = appErr.ErrorDetails
			return fields
		}
		fields["error_type"] = "app_error"
		return fields
	}

	fields["error_type"] = "unexpected_error"
	fields["error_class"] = fmt.Sprintf("%T", err)
	return fields
}

func isValidationError(err error) bool {
	var appErr *core.AppError
	return errors.As(err, &appErr) && appErr.ErrorCode == core.ErrCodeUnprocessable
}
