package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
	"goscaf/internal/core"
)

// GetAccessToken retrieves the access token from cookie or Authorization header.
// If the cookie is not found, it attempts to extract a Bearer token from the header.
func GetAccessToken(c *gin.Context) string {
	token, err := c.Cookie(COOKIE_KEY_ACCESS_TOKEN)
	if err == nil {
		return token
	}

	bearer := c.GetHeader("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimSpace(bearer[7:])
	}

	return ""
}

// GetPayload retrieves the AuthPayload from the context.
// Returns an empty AuthPayload if the context value is not present or invalid.
func GetPayload(c *gin.Context) core.AuthPayload {
	pl, ok := c.Get(CONTEXT_KEY_PAYLOAD)
	if !ok {
		return core.AuthPayload{}
	}

	if payload, ok := pl.(core.AuthPayload); ok {
		return payload
	}
	return core.AuthPayload{}
}

// GetAccountId returns the account ID from the AuthPayload in the context.
func GetAccountId(c *gin.Context) int {
	return GetPayload(c).AccountId
}

// GetAccountName returns the account name from the AuthPayload in the context.
func GetAccountName(c *gin.Context) string {
	return GetPayload(c).AccountName
}
