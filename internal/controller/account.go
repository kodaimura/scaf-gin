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
	c.SetCookie(COOKIE_KEY_JWT, "", 0, "/", config.AppHost, false, true)
	c.Redirect(303, "/login")
}

// POST /api/signup
func (ctr *AccountController) Signup(c *gin.Context) {
	var req request.Signup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	var in input.Signup
	utils.MapFields(&in, req)

	account, err := ctr.accountService.Signup(in)
	if err != nil {
		c.Error(err)
		return
	}

	var res response.Account
	utils.MapFields(&res, account)
	c.JSON(200, res)
}

// POST /api/login
func (ctr *AccountController) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	var in input.Login
	utils.MapFields(&in, req)

	account, err := ctr.accountService.Login(in)
	if err != nil {
		c.Error(err)
		return
	}

	pl, err := ctr.accountService.GenerateJwtPayload(dto.AccountPK{Id: account.Id})
	if err != nil {
		c.Error(err)
		return
	}

	jwtStr, err := jwt.EncodeJwt(pl)
	if err != nil {
		c.Error(err)
	}
	var res response.Login
	res.AccessToken = jwtStr
	res.ExpiresIn = 3600 //TODO
	res.Account = account

	c.SetCookie(COOKIE_KEY_JWT, jwtStr, int(JWT_EXPIRES), "/", config.AppHost, false, true)
	c.JSON(200, res)
}

// GET /api/accounts/me
func (ctr *AccountController) GetOne(c *gin.Context) {
	pl := jwt.GetPayload(c)
	result, err := ctr.accountService.GetOne(dto.AccountPK{Id: pl.AccountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, result)
}

// PUT /api/accounts/me/password
func (ctr *AccountController) PutPassword(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var req request.PutPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	_, err := ctr.accountService.Login(dto.Login{Name: pl.AccountName, Password: req.OldPassword})
	if err != nil {
		c.Error(err)
		return
	}

	var input dto.UpdateAccount
	utils.MapFields(&input, req)
	input.Id = pl.AccountId

	if err := ctr.accountService.Update(input); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

// PUT /api/accounts/me
func (ctr *AccountController) Put(c *gin.Context) {
	pl := jwt.GetPayload(c)

	var req request.PutAccountName
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBadRequestError(err.Error()))
		return
	}

	var input dto.UpdateAccount
	utils.MapFields(&input, req)
	input.Id = pl.AccountId

	if err := ctr.accountService.Update(input); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}

// DELETE /api/accounts/me
func (ctr *AccountController) Delete(c *gin.Context) {
	pl := jwt.GetPayload(c)
	if err := ctr.accountService.Delete(dto.AccountPK{Id: pl.AccountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}