package router

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/adapter/db"

	"scaf-gin/internal/module"

	account_uc "scaf-gin/internal/usecase/account"
	auth_uc "scaf-gin/internal/usecase/auth"

	account_h "scaf-gin/internal/handler/account"
	auth_h "scaf-gin/internal/handler/auth"
)

func SetApi(r *gin.RouterGroup) {
	dbConn := db.NewGormDB()

	accountModule := module.NewAccountModule(dbConn)
	passwordResetTokenModule := module.NewPasswordResetTokenModule(dbConn)

	authUsecase := auth_uc.NewUsecase(dbConn, accountModule, passwordResetTokenModule)
	accountUsecase := account_uc.NewUsecase(accountModule)

	accountHandler := account_h.NewHandler(accountUsecase)
	authHandler := auth_h.NewHandler(authUsecase)

	r.Use(ApiErrorHandler())
	r.POST("/auth/signup", authHandler.ApiSignup)
	r.POST("/auth/login", authHandler.ApiLogin)
	r.POST("/auth/refresh", authHandler.ApiRefresh)
	r.POST("/auth/logout", authHandler.ApiLogout)
	r.POST("/auth/forgot-password", authHandler.ApiForgotPassword)
	r.GET("/auth/reset-password/verify", authHandler.ApiVerifyResetPasswordToken)
	r.POST("/auth/reset-password", authHandler.ApiResetPassword)

	auth := r.Group("", ApiAuthMiddleware(dbConn))
	{
		auth.GET("/accounts", accountHandler.ApiGetAccounts)
		auth.POST("/accounts", accountHandler.ApiPostAccount)
		auth.PUT("/accounts/:target_account_id/password", authHandler.ApiPutMePassword)
		auth.PUT("/accounts/:target_account_id/disable", accountHandler.ApiPutAccountDisable)
		auth.PUT("/accounts/:target_account_id/enable", accountHandler.ApiPutAccountEnable)
		auth.GET("/accounts/:target_account_id", accountHandler.ApiGetAccount)
		auth.PUT("/accounts/:target_account_id", accountHandler.ApiPutAccount)
	}
}
