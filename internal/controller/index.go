package controller

import (
	"github.com/gin-gonic/gin"
	
	"scaf-gin/internal/helper"
)


type IndexController struct {}


func NewIndexController() *IndexController {
	return &IndexController{}
}

//GET /
func (ctrl *IndexController) IndexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"account_name": helper.GetAccountName(c),
	})
}