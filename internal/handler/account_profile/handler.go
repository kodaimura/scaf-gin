package account_profile

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/helper"
	usecase "scaf-gin/internal/usecase/account_profile"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	ApiGetMe(c *gin.Context)
	ApiPutMe(c *gin.Context)
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

// GET /api/accounts/me/profile
func (h *handler) ApiGetMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)

	profile, err := h.usecase.GetOne(usecase.GetOneDto{AccountId: accountId})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountProfileResponse(profile))
}

// PUT /api/accounts/me/profile
func (h *handler) ApiPutMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)

	var req PutMeRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	profile, err := h.usecase.UpdateOne(usecase.UpdateOneDto{
		AccountId:   accountId,
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountProfileResponse(profile))
}
