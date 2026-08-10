package router

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/core"
	account_h "scaf-gin/internal/handler/account"
	auth_h "scaf-gin/internal/handler/auth"
	"scaf-gin/internal/module"
)

type API struct {
	AccountHandler account_h.Handler
	AuthHandler    auth_h.Handler
	AccountModule  module.AccountModule
	Auth           core.AuthI
}

func SetAPI(r *gin.RouterGroup, api API) {
	r.Use(ApiErrorHandler())
	r.POST("/auth/signup", api.AuthHandler.ApiSignup)
	r.POST("/auth/login", api.AuthHandler.ApiLogin)
	r.POST("/auth/refresh", api.AuthHandler.ApiRefresh)
	r.POST("/auth/logout", api.AuthHandler.ApiLogout)
	r.POST("/auth/forgot-password", api.AuthHandler.ApiForgotPassword)
	r.GET("/auth/reset-password/verify", api.AuthHandler.ApiVerifyResetPasswordToken)
	r.POST("/auth/reset-password", api.AuthHandler.ApiResetPassword)

	auth := r.Group("", ApiAuthMiddleware(api.AccountModule, api.Auth))
	{
		auth.GET("/accounts", api.AccountHandler.ApiGetAccounts)
		auth.POST("/accounts", api.AccountHandler.ApiPostAccount)
		auth.PUT("/accounts/:target_account_id/password", api.AuthHandler.ApiPutMePassword)
		auth.PUT("/accounts/:target_account_id/disable", api.AccountHandler.ApiPutAccountDisable)
		auth.PUT("/accounts/:target_account_id/enable", api.AccountHandler.ApiPutAccountEnable)
		auth.GET("/accounts/:target_account_id", api.AccountHandler.ApiGetAccount)
		auth.PUT("/accounts/:target_account_id", api.AccountHandler.ApiPutAccount)
	}
}
