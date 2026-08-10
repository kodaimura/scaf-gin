package router

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/adapter/db"

	"scaf-gin/internal/module/account"
	"scaf-gin/internal/module/account_profile"
	"scaf-gin/internal/module/password_reset_token"

	account_uc "scaf-gin/internal/usecase/account"
	account_profile_uc "scaf-gin/internal/usecase/account_profile"
	auth_uc "scaf-gin/internal/usecase/auth"

	account_h "scaf-gin/internal/handler/account"
	account_profile_h "scaf-gin/internal/handler/account_profile"
	auth_h "scaf-gin/internal/handler/auth"
)

var gorm = db.NewGormDB()

//var sqlx = db.NewSqlxDB()

/* DI (Repository) */
var accountRepository = account.NewRepository()
var accountProfileRepository = account_profile.NewRepository()
var passwordResetTokenRepository = password_reset_token.NewRepository()

/* DI (Service) */
var accountService = account.NewService(accountRepository)
var accountProfileService = account_profile.NewService(accountProfileRepository)
var passwordResetTokenService = password_reset_token.NewService(passwordResetTokenRepository)

/* DI (Usecase) */
var authUsecase = auth_uc.NewUsecase(gorm, accountService, accountProfileService, passwordResetTokenService)
var accountUsecase = account_uc.NewUsecase(gorm, accountService)
var accountProfileUsecase = account_profile_uc.NewUsecase(gorm, accountProfileService)

/* DI (Handler) */
var accountHandler = account_h.NewHandler(accountUsecase)
var accountProfileHandler = account_profile_h.NewHandler(accountProfileUsecase)
var authHandler = auth_h.NewHandler(authUsecase)

func SetApi(r *gin.RouterGroup) {
	r.Use(ApiErrorHandler())
	r.POST("/auth/signup", authHandler.ApiSignup)
	r.POST("/auth/login", authHandler.ApiLogin)
	r.POST("/auth/refresh", authHandler.ApiRefresh)
	r.POST("/auth/logout", authHandler.ApiLogout)
	r.POST("/auth/forgot-password", authHandler.ApiForgotPassword)
	r.GET("/auth/reset-password/verify", authHandler.ApiVerifyResetPasswordToken)
	r.POST("/auth/reset-password", authHandler.ApiResetPassword)

	auth := r.Group("", ApiAuthMiddleware())
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
