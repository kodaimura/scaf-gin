package controller

import (
	"github.com/gin-gonic/gin"

	"goscaf/pkg/jwt"
	"goscaf/pkg/utils"
	"goscaf/internal/service"
	"goscaf/internal/dto/input"
	"goscaf/internal/dto/request"
	"goscaf/internal/dto/response"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController() *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(),
	}
}

// GET /signup
func (ctr *AccountController) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

// GET /login
func (ctr *AccountController) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

// GET /logout
func (ctr *AccountController) Logout(c *gin.Context) {
	c.SetCookie(common.COOKIE_KEY_ACCESS_TOKEN, "", 0, "/", config.AppHost, false, true)
	c.Redirect(303, "/login")
}

// POST /api/signup
func (ctr *AccountController) Signup(c *gin.Context) {
	var req request.Signup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctr.accountService.Signup(input.Signup{
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
func (ctr *AccountController) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctr.accountService.Login(input.Login{
		AccountName: req.AccountName,
		AccountPassword: req.AccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	claims := map[string]interface{}{
		"account_id":  account.AccountId,
		"account_nme": account.AccountName,
	}
	pl := jwt.NewPayload(in.AccountId, int(config.JwtExpiresSeconds), claims)
	encoded, err := jwt.EncodeToken(pl)
	if err != nil {
		c.Error(err)
	}

	c.SetCookie(common.COOKIE_KEY_ACCESS_TOKEN, encoded, res.ExpiresIn, "/", config.AppHost, false, true)
	c.JSON(200, response.Login{
		AccessToken: encoded,
		ExpiresIn: int(config.JwtExpiresSeconds),
		Account: response.Account{
			AccountId: account.AccountId,
			AccountName: account.AccountName,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		}
	})
}

// GET /api/accounts/me
func (ctr *AccountController) GetOne(c *gin.Context) {
	accountId := common.GetAccountId(c)
	account, err := ctr.accountService.GetOne(input.Account{AccountId: accountId})
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
func (ctr *AccountController) PutPassword(c *gin.Context) {
	accountName := common.GetAccountName(c)

	var req request.PutAccountPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctr.accountService.Login(input.Login{
		AccountName: accountName, 
		AccountPassword: req.OldAccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	_, err := ctr.accountService.UpdateOne(input.Account{
		AccountId: account.AccountId,
		AccountPassword = req.NewAccountPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

// PUT /api/accounts/me
func (ctr *AccountController) Put(c *gin.Context) {
	accountId := common.GetAccountId(c)

	var req request.PutAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	account, err := ctr.accountService.UpdateOne(input.Account{
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
func (ctr *AccountController) Delete(c *gin.Context) {
	accountId := common.GetAccountId(c)
	if err := ctr.accountService.DeleteOne(input.Account{AccountId: accountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}