package controller

import (
	"github.com/gin-gonic/gin"
	
	"goscaf/internal/common"
)


type IndexController struct {}


func NewIndexController() *IndexController {
	return &IndexController{}
}


//GET /
func (ctrl *IndexController) IndexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"account_name": common.GetAccountName(c),
	})
}