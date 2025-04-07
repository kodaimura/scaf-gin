package common

import (
	"github.com/gin-gonic/gin"
	
	"goscaf/internal/core"
)

func GetPayload(c *gin.Context) core.AuthPayload {
	pl := c.Keys[CONTEXT_KEY_JWT_PAYLOAD]
	if pl == nil {
		return core.AuthPayload{}
	}
	return pl.(core.AuthPayload)
}

func GetAccountId(c *gin.Context) int {
	payload := GetPayload(c)
	return payload.AccountId
}

func GetAccountName(c *gin.Context) string {
	payload := GetPayload(c)
	return payload.AccountName
}