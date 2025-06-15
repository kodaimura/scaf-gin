package account

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Controller interface {
	ApiGetMe(c *gin.Context)
	ApiPutMe(c *gin.Context)
	ApiDeleteMe(c *gin.Context)
}

type controller struct {
	db *gorm.DB
	service Service
}

func NewController(db *gorm.DB, service Service) Controller {
	return &controller{
		db:      db,
		service: service,
	}
}

// GET /api/accounts/me
func (ctrl *controller) ApiGetMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	account, err := ctrl.service.GetOne(GetOneDto{Id: accountId}, ctrl.db)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountResponse(account))
}

// PUT /api/accounts/me
func (ctrl *controller) ApiPutMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)

	var req PutMeRequest
	if err := helper.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	account, err := ctrl.service.UpdateOne(UpdateOneDto{
		Id:   accountId,
		Name: req.Name,
	}, ctrl.db)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ToAccountResponse(account))
}

// DELETE /api/accounts/me
func (ctrl *controller) ApiDeleteMe(c *gin.Context) {
	accountId := helper.GetAccountId(c)
	if err := ctrl.service.DeleteOne(DeleteOneDto{Id: accountId}, ctrl.db); err != nil {
		c.Error(err)
		return
	}

	c.JSON(204, nil)
}
