package index

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/helper"
)

type Controller interface {
	IndexPage(c *gin.Context)
}

type controller struct{}

func NewController() Controller {
	return &controller{}
}

// GET /
func (ctrl *controller) IndexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"name": helper.GetAccountName(c),
	})
}
