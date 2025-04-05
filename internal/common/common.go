package common

import (
	"github.com/gin-gonic/gin"
	
	"goscaf/pkg/jwt"
)

func GetPayload(c *gin.Context) jwt.Payload {
	pl := c.Keys[CONTEXT_KEY_JWT_PAYLOAD]
	if pl == nil {
		return jwt.Payload{}
	}
	return pl.(jwt.Payload)
}

func GetCustomClaims(c *gin.Context) map[string]interface{} {
	pl := GetPayload(c)
	return pl.CustomClaims
}

func GetAccountId(c *gin.Context) int {
	value, _ := GetCustomClaims(c)["account_id"]
	return value.(int)
}

func GetAccountName(c *gin.Context) string {
	value, _ := GetCustomClaims(c)["account_name"]
	return value.(string)
}