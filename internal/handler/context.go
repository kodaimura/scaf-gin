package handler

import (
	"scaf-gin/internal/core"

	"github.com/gin-gonic/gin"
)

const ContextKeyPayload = "payload"

func SetPayload(c *gin.Context, payload core.AuthPayload) {
	c.Set(ContextKeyPayload, payload)
}

func GetPayload(c *gin.Context) core.AuthPayload {
	pl, ok := c.Get(ContextKeyPayload)
	if !ok {
		return core.AuthPayload{}
	}

	if payload, ok := pl.(core.AuthPayload); ok {
		return payload
	}
	return core.AuthPayload{}
}

func GetAccountID(c *gin.Context) int {
	return GetPayload(c).AccountId
}
