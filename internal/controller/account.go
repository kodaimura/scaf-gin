package controller

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/dto/input"
	"scaf-gin/internal/dto/request"
	"scaf-gin/internal/dto/response"
	"scaf-gin/internal/helper"
	"scaf-gin/internal/service"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// GET /signup
func (ctrl *AccountController) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

// GET /login
func (ctrl *AccountController) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

// GET /logout
func (ctrl *AccountController) Logout(c *gin.Context) {
	core.Auth.RevokeRefreshToken(helper.GetRefreshToken(c))
	helper.RemoveAccessTokenCookie(c)
	helper.RemoveRefreshTokenCookie(c)
	c.Redirect(303, "/login")
}

// POST /api/accounts/signup
func (ctrl *AccountController) ApiSignup(c *gin.Context) {
	var req request.Signup
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.Signup(input.Signup{
		AccountName:     req.AccountName,
		AccountPassword: req.AccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	})
}

// POST /api/accounts/login
func (ctrl *AccountController) ApiLogin(c *gin.Context) {
	var req request.Login
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		AccountName:     req.AccountName,
		AccountPassword: req.AccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := core.Auth.CreateRefreshToken(core.AuthPayload{
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	helper.SetAccessTokenCookie(c, accessToken)
	helper.SetRefreshTokenCookie(c, refreshToken)

	core.Logger.Info("account login: id=%d name=%s", account.AccountId, account.AccountName)

	c.JSON(200, response.Login{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  config.AccessTokenExpiresSeconds,
		RefreshExpiresIn: config.RefreshTokenExpiresSeconds,
		Account: response.Account{
			AccountId:   account.AccountId,
			AccountName: account.AccountName,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
		},
	})
}

// POST /api/accounts/refresh
func (ctrl *AccountController) ApiRefresh(c *gin.Context) {
	refreshToken := helper.GetRefreshToken(c)

	payload, err := core.Auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.Error(core.NewAppError("invalid or expired refresh token", core.ErrCodeUnauthorized))
		return
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:   payload.AccountId,
		AccountName: payload.AccountName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	helper.SetAccessTokenCookie(c, accessToken)

	core.Logger.Info("access token refreshed: id=%d name=%s", payload.AccountId, payload.AccountName)

	c.JSON(200, response.Refresh{
		AccessToken: accessToken,
		ExpiresIn:   config.AccessTokenExpiresSeconds,
	})
}

// GET /api/accounts/logout
func (ctrl *AccountController) ApiLogout(c *gin.Context) {
	core.Auth.RevokeRefreshToken(helper.GetRefreshToken(c))
	helper.RemoveAccessTokenCookie(c)
	helper.RemoveRefreshTokenCookie(c)
	c.JSON(200, gin.H{})
}

// GET /api/accounts/me
func (ctrl *AccountController) ApiGetOne(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	account, err := ctrl.accountService.GetOne(input.Account{AccountId: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	})
}

// PUT /api/accounts/me/password
func (ctrl *AccountController) ApiPutPassword(c *gin.Context) {
	accountName := helper.GetAccountName(c)

	var req request.PutAccountPassword
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		AccountName:     accountName,
		AccountPassword: req.OldAccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	_, err = ctrl.accountService.UpdateOne(input.Account{
		AccountId:       account.AccountId,
		AccountPassword: req.NewAccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

// PUT /api/accounts/me
func (ctrl *AccountController) ApiPutOne(c *gin.Context) {
	accountId := helper.GetAccountId(c)

	var req request.PutAccount
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.UpdateOne(input.Account{
		AccountId:   accountId,
		AccountName: req.AccountName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	})
}

// DELETE /api/accounts/me
func (ctrl *AccountController) ApiDeleteOne(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	if err := ctrl.accountService.DeleteOne(input.Account{AccountId: accountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}
