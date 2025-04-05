package controller

import (
	"github.com/gin-gonic/gin"

	"goscaf/config"
	"goscaf/pkg/jwt"
	"goscaf/pkg/errs"
	"goscaf/internal/common"
	"goscaf/internal/service"
	"goscaf/internal/dto/input"
	"goscaf/internal/dto/request"
	"goscaf/internal/dto/response"
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
	c.SetCookie(common.COOKIE_KEY_ACCESS_TOKEN, "", 0, "/", config.AppHost, config.SecureCookie, true)
	c.Redirect(303, "/login")
}

// POST /api/signup
func (ctrl *AccountController) ApiSignup(c *gin.Context) {
	var req request.Signup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctrl.accountService.Signup(input.Signup{
		AccountName: req.AccountName,
		AccountPassword: req.AccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId: account.AccountId,
		AccountName: account.AccountName,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

// POST /api/login
func (ctrl *AccountController) ApiLogin(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		AccountName: req.AccountName,
		AccountPassword: req.AccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	claims := map[string]interface{}{
		"account_id":  account.AccountId,
		"account_name": account.AccountName,
	}
	pl := jwt.NewPayload(account.AccountId, config.JwtExpiresSeconds, claims)
	encoded, err := jwt.EncodeToken(pl, config.JwtSecretKey)
	if err != nil {
		c.Error(err)
	}

	c.SetCookie(common.COOKIE_KEY_ACCESS_TOKEN, encoded, config.JwtExpiresSeconds, "/", config.AppHost, config.SecureCookie, true)
	c.JSON(200, response.Login{
		AccessToken: encoded,
		ExpiresIn: config.JwtExpiresSeconds,
		Account: response.Account{
			AccountId: account.AccountId,
			AccountName: account.AccountName,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
	})
}

// GET /api/logout
func (ctrl *AccountController) ApiLogout(c *gin.Context) {
	c.SetCookie(common.COOKIE_KEY_ACCESS_TOKEN, "", 0, "/", config.AppHost, config.SecureCookie, true)
	c.JSON(200, gin.H{})
}

// GET /api/accounts/me
func (ctrl *AccountController) ApiGetOne(c *gin.Context) {
	accountId := common.GetAccountId(c)
	account, err := ctrl.accountService.GetOne(input.Account{AccountId: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId: account.AccountId,
		AccountName: account.AccountName,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

// PUT /api/accounts/me/password
func (ctrl *AccountController) ApiPutPassword(c *gin.Context) {
	accountName := common.GetAccountName(c)

	var req request.PutAccountPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctrl.accountService.Login(input.Login{
		AccountName: accountName, 
		AccountPassword: req.OldAccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	_, err = ctrl.accountService.UpdateOne(input.Account{
		AccountId: account.AccountId,
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
	accountId := common.GetAccountId(c)

	var req request.PutAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctrl.accountService.UpdateOne(input.Account{
		AccountId: accountId,
		AccountName: req.AccountName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, response.Account{
		AccountId: account.AccountId,
		AccountName: account.AccountName,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

// DELETE /api/accounts/me
func (ctrl *AccountController) ApiDeleteOne(c *gin.Context) {
	accountId := common.GetAccountId(c)
	if err := ctrl.accountService.DeleteOne(input.Account{AccountId: accountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}