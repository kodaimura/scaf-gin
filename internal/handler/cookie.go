package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"scaf-gin/config"
)

const (
	CookieKeyAccessToken  = "access_token"
	CookieKeyRefreshToken = "refresh_token"
)

func GetAccessToken(c *gin.Context) string {
	token, err := c.Cookie(CookieKeyAccessToken)
	if err == nil {
		return token
	}

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

func SetAccessTokenCookie(c *gin.Context, accessToken string) {
	maxAge := config.Current.AccessTokenExpiresSeconds
	if accessToken == "" {
		maxAge = -1
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		CookieKeyAccessToken,
		accessToken,
		maxAge,
		"/", "",
		config.Current.AppEnv == "production",
		true,
	)
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
