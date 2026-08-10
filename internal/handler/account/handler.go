package account

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"scaf-gin/internal/core"
	"scaf-gin/internal/helper"
	usecase "scaf-gin/internal/usecase/account"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	ApiGetAccounts(c *gin.Context)
	ApiPostAccount(c *gin.Context)
	ApiGetMe(c *gin.Context)
	ApiPutMe(c *gin.Context)
	ApiDeleteMe(c *gin.Context)
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
	accounts, err := h.usecase.Get(usecase.GetDto{})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, GetAccountsResponse{Accounts: ToAccountResponseList(accounts)})
}

// POST /api/accounts
func (h *handler) ApiPostAccount(c *gin.Context) {
	var req PostAccountRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.CreateOne(usecase.CreateOneDto{
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

	c.JSON(201, PostAccountResponse{Account: ToAccountResponse(account)})
}

// GET /api/accounts/me
func (h *handler) ApiGetMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	account, err := h.usecase.GetOne(usecase.GetOneDto{Id: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, GetCurrentAccountResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/me
func (h *handler) ApiPutMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)

	var req PutMeRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.UpdateOne(usecase.UpdateOneDto{
		Id:        accountId,
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

	c.JSON(200, PutAccountResponse{Account: ToAccountResponse(account)})
}

// DELETE /api/accounts/me
func (h *handler) ApiDeleteMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	if err := h.usecase.DeleteOne(usecase.DeleteOneDto{Id: accountId}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(204, nil)
}

// GET /api/accounts/:target_account_id
func (h *handler) ApiGetAccount(c *gin.Context) {
	targetAccountId, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.GetOne(usecase.GetOneDto{Id: targetAccountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, GetAccountResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id
func (h *handler) ApiPutAccount(c *gin.Context) {
	targetAccountId, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req PutAccountRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.UpdateOne(usecase.UpdateOneDto{
		Id:        targetAccountId,
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

	c.JSON(200, PutAccountResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id/disable
func (h *handler) ApiPutAccountDisable(c *gin.Context) {
	targetAccountId, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.DisableOne(usecase.DisableOneDto{Id: targetAccountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, PutAccountDisableResponse{Account: ToAccountResponse(account)})
}

// PUT /api/accounts/:target_account_id/enable
func (h *handler) ApiPutAccountEnable(c *gin.Context) {
	targetAccountId, err := parseAccountID(c)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := h.usecase.EnableOne(usecase.EnableOneDto{Id: targetAccountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, PutAccountEnableResponse{Account: ToAccountResponse(account)})
}

func parseAccountID(c *gin.Context) (int, error) {
	rawAccountID := c.Param("target_account_id")
	if rawAccountID == "me" {
		return helper.GetAccountId(c), nil
	}
	accountID, err := strconv.Atoi(rawAccountID)
	if err != nil {
		return 0, core.ErrBadRequest
	}
	return accountID, nil
}
