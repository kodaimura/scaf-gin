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
	helper.SetAccessTokenCookie(c, "")
	helper.SetRefreshTokenCookie(c, "")
	c.Redirect(303, "/login")
}

// POST /api/accounts/signup
func (ctrl *AccountController) ApiSignup(c *gin.Context) {
	var req request.Signup
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.CreateOne(input.Account{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, response.FromModelAccount(account))
}

// POST /api/accounts/login
func (ctrl *AccountController) ApiLogin(c *gin.Context) {
	var req request.Login
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:   account.Id,
		AccountName: account.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	refreshToken, err := core.Auth.CreateRefreshToken(core.AuthPayload{
		AccountId:   account.Id,
		AccountName: account.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	helper.SetAccessTokenCookie(c, accessToken)
	helper.SetRefreshTokenCookie(c, refreshToken)

	core.Logger.Info("account login: id=%d name=%s", account.Id, account.Name)

	c.JSON(200, response.Login{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  config.AccessTokenExpiresSeconds,
		RefreshExpiresIn: config.RefreshTokenExpiresSeconds,
		Account:          response.FromModelAccount(account),
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

// POST /api/accounts/logout
func (ctrl *AccountController) ApiLogout(c *gin.Context) {
	core.Auth.RevokeRefreshToken(helper.GetRefreshToken(c))
	helper.SetAccessTokenCookie(c, "")
	helper.SetRefreshTokenCookie(c, "")
	c.JSON(200, gin.H{})
}

// GET /api/accounts/me
func (ctrl *AccountController) ApiGetOne(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	account, err := ctrl.accountService.GetOne(input.Account{Id: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.FromModelAccount(account))
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
		Id:   accountId,
		Name: req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.FromModelAccount(account))
}

// PUT /api/accounts/me/password
func (ctrl *AccountController) ApiPutPassword(c *gin.Context) {
	accountName := helper.GetAccountName(c)

	var req request.PutPassword
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		Name:     accountName,
		Password: req.OldPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	_, err = ctrl.accountService.UpdatePassword(input.UpdatePassword{
		Id:       account.Id,
		Password: req.NewPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

// DELETE /api/accounts/me
func (ctrl *AccountController) ApiDeleteOne(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	if err := ctrl.accountService.DeleteOne(input.Account{Id: accountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(204, nil)
}
