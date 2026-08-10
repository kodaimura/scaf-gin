package account

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"scaf-gin/internal/core"
	handlerutil "scaf-gin/internal/handler"
	usecase "scaf-gin/internal/usecase/account"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	ApiGetAccounts(c *gin.Context)
	ApiPostAccount(c *gin.Context)
	ApiGetCurrentAccount(c *gin.Context)
	ApiGetAccount(c *gin.Context)
	ApiPutAccount(c *gin.Context)
	ApiPutAccountDisable(c *gin.Context)
	ApiPutAccountEnable(c *gin.Context)
}

type handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) Handler {
	return &handler{
		usecase: usecase,
	}
}

// -----------------------------
// Handler Implementations
// -----------------------------

// GET /api/accounts
func (h *handler) ApiGetAccounts(c *gin.Context) {
	accounts, err := h.usecase.List(usecase.ListDto{})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.OK(c, GetAccountsResponse{Accounts: ToAccountResponseList(accounts)})
}

// POST /api/accounts
func (h *handler) ApiPostAccount(c *gin.Context) {
	var req PostAccountRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Create(usecase.CreateDto{
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

	handlerutil.Created(c, PostAccountResponse{Account: ToAccountResponse(account)})
}

// GET /api/accounts/me
func (h *handler) ApiGetCurrentAccount(c *gin.Context) {
	account, err := h.usecase.Get(usecase.GetDto{ID: handlerutil.GetAccountID(c)})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.OK(c, GetCurrentAccountResponse{Account: ToAccountResponse(account)})
}

// GET /api/accounts/:target_account_id
func (h *handler) ApiGetAccount(c *gin.Context) {
	targetAccountID, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Get(usecase.GetDto{ID: targetAccountID})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.OK(c, GetAccountResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id
func (h *handler) ApiPutAccount(c *gin.Context) {
	targetAccountID, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req PutAccountRequest
	if err := handlerutil.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Update(usecase.UpdateDto{
		ID:        targetAccountID,
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

	handlerutil.OK(c, PutAccountResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id/disable
func (h *handler) ApiPutAccountDisable(c *gin.Context) {
	targetAccountID, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Disable(usecase.DisableDto{ID: targetAccountID})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.OK(c, PutAccountDisableResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id/enable
func (h *handler) ApiPutAccountEnable(c *gin.Context) {
	targetAccountID, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.Enable(usecase.EnableDto{ID: targetAccountID})
	if err != nil {
		c.Error(err)
		return
	}

	handlerutil.OK(c, PutAccountEnableResponse{Account: ToAccountResponse(account)})
}

func parseAccountID(c *gin.Context) (int64, error) {
	rawAccountID := c.Param("target_account_id")
	if rawAccountID == "me" {
		return handlerutil.GetAccountID(c), nil
	}
	accountID, err := strconv.ParseInt(rawAccountID, 10, 64)
	if err != nil {
		return 0, core.ErrBadRequest
	}
	return accountID, nil
}
