package general

import (
	"github.com/gin-gonic/gin"

	"scaf-gin/internal/helper"
)

// -----------------------------
// Handler Interface
// -----------------------------

type Handler interface {
	IndexPage(c *gin.Context)
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}

// -----------------------------
// Handler Implementations
// -----------------------------

// GET /
func (h *handler) IndexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"name": helper.GetAccountName(c),
	})
}
