package auth

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	handlerutil "scaf-gin/internal/handler"
	usecase "scaf-gin/internal/usecase/auth"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	ApiSignup(c *gin.Context)
	ApiLogin(c *gin.Context)
	ApiRefresh(c *gin.Context)
	ApiLogout(c *gin.Context)
	ApiForgotPassword(c *gin.Context)
	ApiVerifyResetPasswordToken(c *gin.Context)
	ApiResetPassword(c *gin.Context)

	ApiPutMePassword(c *gin.Context)
}

type handler struct {
	usecase usecase.Usecase
	logger  core.Logger
}

func NewHandler(usecase usecase.Usecase, logger core.Logger) Handler {
	return &handler{
		usecase: usecase,
		logger:  logger,
	}
}

// -----------------------------
// Handler Implementations
// -----------------------------

// POST /api/auth/signup
func (h *handler) ApiSignup(c *gin.Context) {
	if !config.Current.EnableSignup {
		c.Error(core.ErrForbidden)
		return
	}

	var req SignupRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Signup(usecase.SignupDto{
		LoginID:   req.LoginID,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.Created(c, SignupResponse{Account: ToAccountResponse(account)})
}

// POST /api/auth/login
func (h *handler) ApiLogin(c *gin.Context) {
	var req LoginRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	acct, accessToken, refreshToken, err := h.usecase.Login(usecase.LoginDto(req))
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.SetRefreshTokenCookie(c, refreshToken, req.RememberMe)
	h.logger.Info("account login: id=%d login_id=%s", acct.ID, acct.LoginID)

	handlerutil.OK(c, LoginResponse{
		Account:     ToAccountResponse(acct),
		AccessToken: accessToken,
	})
}

// POST /api/auth/refresh
func (h *handler) ApiRefresh(c *gin.Context) {
	refreshToken := handlerutil.GetRefreshToken(c)

	payload, accessToken, err := h.usecase.Refresh(refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("access token refreshed: id=%d", payload.AccountID)

	handlerutil.OK(c, RefreshResponse{
		AccessToken: accessToken,
	})
}

// POST /api/auth/logout
func (h *handler) ApiLogout(c *gin.Context) {
	handlerutil.SetRefreshTokenCookie(c, "", false)
	handlerutil.NoContent(c)
}

// POST /api/auth/forgot-password
func (h *handler) ApiForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	if err := h.usecase.ForgotPassword(usecase.ForgotPasswordDto{Email: req.Email}); err != nil {
		c.Error(err)
		return
	}

	handlerutil.NoContent(c)
}

// GET /api/auth/reset-password/verify
func (h *handler) ApiVerifyResetPasswordToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Error(core.ErrBadRequest)
		return
	}

	if err := h.usecase.VerifyResetPasswordToken(usecase.VerifyResetPasswordTokenDto{Token: token}); err != nil {
		c.Error(err)
		return
	}

	handlerutil.NoContent(c)
}

// POST /api/auth/reset-password
func (h *handler) ApiResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	if err := h.usecase.ResetPassword(usecase.ResetPasswordDto{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}); err != nil {
		c.Error(err)
		return
	}

	handlerutil.NoContent(c)
}

// PUT /api/accounts/me/password
func (h *handler) ApiPutMePassword(c *gin.Context) {
	accountID := handlerutil.GetAccountID(c)

	var req PutMePasswordRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	err := h.usecase.UpdatePassword(usecase.UpdatePasswordDto{
		ID:          accountID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.NoContent(c)
}
