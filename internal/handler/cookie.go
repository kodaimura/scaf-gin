package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
)

const (
	CookieKeyRefreshToken = "refresh_token"
)

func GetAccessToken(c *gin.Context) string {
	bearer := c.GetHeader("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimSpace(bearer[7:])
	}

	return ""
}

func GetRefreshToken(c *gin.Context) string {
	token, err := c.Cookie(CookieKeyRefreshToken)
	if err == nil {
		return token
	}
	return ""
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string, rememberMe bool) {
	maxAge := config.Current.RefreshTokenExpiresSeconds
	if rememberMe {
		maxAge = config.Current.RefreshTokenRememberMeExpiresSeconds
	}
	if refreshToken == "" {
		maxAge = -1
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		CookieKeyRefreshToken,
		refreshToken,
		maxAge,
		"/", "",
		config.Current.AppEnv == "production",
		true,
	)
}
