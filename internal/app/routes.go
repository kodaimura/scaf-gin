package app

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/core"
	account_handler "scaf-gin/internal/handler/account"
	auth_handler "scaf-gin/internal/handler/auth"
	"scaf-gin/internal/module"
)

func registerAPIRoutes(
	r *gin.RouterGroup,
	accountHandler account_handler.Handler,
	authHandler auth_handler.Handler,
	accountModule module.AccountModule,
	authService core.Auth,
	log core.Logger,
) {
	r.Use(apiErrorHandler(log))
	r.POST("/auth/signup", authHandler.ApiSignup)
	r.POST("/auth/login", authHandler.ApiLogin)
	r.POST("/auth/refresh", authHandler.ApiRefresh)
	r.POST("/auth/logout", authHandler.ApiLogout)
	r.POST("/auth/forgot-password", authHandler.ApiForgotPassword)
	r.GET("/auth/reset-password/verify", authHandler.ApiVerifyResetPasswordToken)
	r.POST("/auth/reset-password", authHandler.ApiResetPassword)

	auth := r.Group("", apiAuthMiddleware(accountModule, authService))
	{
		auth.GET("/accounts", accountHandler.ApiGetAccounts)
		auth.POST("/accounts", accountHandler.ApiPostAccount)
		auth.GET("/accounts/me", accountHandler.ApiGetCurrentAccount)
		auth.PUT("/accounts/:target_account_id/password", authHandler.ApiPutMePassword)
		auth.PUT("/accounts/:target_account_id/disable", accountHandler.ApiPutAccountDisable)
		auth.PUT("/accounts/:target_account_id/enable", accountHandler.ApiPutAccountEnable)
		auth.GET("/accounts/:target_account_id", accountHandler.ApiGetAccount)
		auth.PUT("/accounts/:target_account_id", accountHandler.ApiPutAccount)
	}
}
