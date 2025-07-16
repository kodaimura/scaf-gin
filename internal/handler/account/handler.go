package account

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/helper"
	usecase "scaf-gin/internal/usecase/account"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	ApiGetMe(c *gin.Context)
	ApiPutMe(c *gin.Context)
	ApiDeleteMe(c *gin.Context)
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

// GET /api/accounts/me
func (h *handler) ApiGetMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	account, err := h.usecase.GetOne(usecase.GetOneDto{Id: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountResponse(account))
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
		Id:   accountId,
		Name: req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountResponse(account))
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
