package common

import (
	"strings"
	"github.com/gin-gonic/gin"
	
	"goscaf/internal/core"
)

func GetAccessToken (c *gin.Context) string {
	token, err := c.Cookie(COOKIE_KEY_ACCESS_TOKEN)
	if err == nil {
		return token
	}

	bearer := c.Request.Header.Get("Authorization")
	if bearer != "" && !strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimSpace(bearer[7:])
	}

	return ""
}

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