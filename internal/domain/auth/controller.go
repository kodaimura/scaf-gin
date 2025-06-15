package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/helper"
)

type Controller interface {
	SignupPage(c *gin.Context)
	LoginPage(c *gin.Context)
	Logout(c *gin.Context)

	ApiSignup(c *gin.Context)
	ApiLogin(c *gin.Context)
	ApiRefresh(c *gin.Context)
	ApiLogout(c *gin.Context)

	ApiPutMePassword(c *gin.Context)
}

type controller struct {
	db      *gorm.DB
	service Service
}

func NewController(db *gorm.DB, service Service) Controller {
	return &controller{
		db:      db,
		service: service,
	}
}

// GET /signup
func (ctrl *controller) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

// GET /login
func (ctrl *controller) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

// GET /logout
func (ctrl *controller) Logout(c *gin.Context) {
	core.Auth.RevokeRefreshToken(helper.GetRefreshToken(c))
	helper.SetAccessTokenCookie(c, "")
	helper.SetRefreshTokenCookie(c, "")
	c.Redirect(303, "/login")
}

// POST /api/accounts/signup
func (ctrl *controller) ApiSignup(c *gin.Context) {
	var req SignupRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	_, err := ctrl.service.Signup(SignupDto(req), ctrl.db)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{})
}

// POST /api/accounts/login
func (ctrl *controller) ApiLogin(c *gin.Context) {
	var req LoginRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.service.Login(LoginDto(req), ctrl.db)
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

	c.JSON(200, LoginResponse{
		AccountId:        account.Id,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  config.AccessTokenExpiresSeconds,
		RefreshExpiresIn: config.RefreshTokenExpiresSeconds,
	})
}

// POST /api/accounts/refresh
func (ctrl *controller) ApiRefresh(c *gin.Context) {
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

	c.JSON(200, RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   config.AccessTokenExpiresSeconds,
	})
}

// POST /api/accounts/logout
func (ctrl *controller) ApiLogout(c *gin.Context) {
	core.Auth.RevokeRefreshToken(helper.GetRefreshToken(c))
	helper.SetAccessTokenCookie(c, "")
	helper.SetRefreshTokenCookie(c, "")
	c.JSON(200, gin.H{})
}

// PUT /api/accounts/me/password
func (ctrl *controller) ApiPutMePassword(c *gin.Context) {
	accountName := helper.GetAccountName(c)

	var req PutMePasswordRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.service.Login(LoginDto{
		Name:     accountName,
		Password: req.OldPassword,
	}, ctrl.db)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = ctrl.service.UpdatePassword(UpdatePasswordDto{
		Id:       account.Id,
		Password: req.NewPassword,
	}, ctrl.db)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}
