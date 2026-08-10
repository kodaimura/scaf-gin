package auth

import (
	"errors"
	"strconv"
	"time"

	jwtpackage "github.com/golang-jwt/jwt/v5"

	"scaf-gin/config"
	"scaf-gin/internal/core"
)

// JwtAuth implements the AuthI interface using JWT for authentication.
type JwtAuth struct{}

func NewJwtAuth() core.AuthI {
	return &JwtAuth{}
}

type jwtPayload struct {
	jwtpackage.RegisteredClaims
	core.AuthPayload
	TokenType string `json:"type"`
}

// CreateAccessToken creates a signed JWT containing the given AuthPayload.
func (j *JwtAuth) CreateAccessToken(payload core.AuthPayload) (string, error) {
	return j.createToken(
		payload,
		config.AccessTokenSecret,
		"access",
		time.Second*time.Duration(config.AccessTokenExpiresSeconds))
}

// CreateRefreshToken creates a signed JWT containing the given AuthPayload.
func (j *JwtAuth) CreateRefreshToken(payload core.AuthPayload, rememberMe bool) (string, error) {
	expiresSeconds := config.RefreshTokenExpiresSeconds
	if rememberMe {
		expiresSeconds = config.RefreshTokenRememberMeExpiresSeconds
	}
	return j.createToken(
		payload,
		config.RefreshTokenSecret,
		"refresh",
		time.Second*time.Duration(expiresSeconds),
	)
}

// Common function to generate a JWT token (access or refresh)
func (j *JwtAuth) createToken(payload core.AuthPayload, secretKey string, tokenType string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	jp := jwtPayload{
		AuthPayload: payload,
		TokenType:   tokenType,
		RegisteredClaims: jwtpackage.RegisteredClaims{
			Subject:   strconv.Itoa(payload.AccountId),
			IssuedAt:  jwtpackage.NewNumericDate(now),
			NotBefore: jwtpackage.NewNumericDate(now),
			ExpiresAt: jwtpackage.NewNumericDate(now.Add(expiresIn)),
		},
	}

	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, jp)
	return token.SignedString([]byte(secretKey))
}

// VerifyAccessToken verifies the given JWT and extracts the AuthPayload.
func (j *JwtAuth) VerifyAccessToken(token string) (core.AuthPayload, error) {
	return j.verifyToken(token, config.AccessTokenSecret, "access")
}

// VerifyRefreshToken verifies the given refresh token and extracts the AuthPayload.
func (j *JwtAuth) VerifyRefreshToken(token string) (core.AuthPayload, error) {
	return j.verifyToken(token, config.RefreshTokenSecret, "refresh")
}

// Common function to validate a JWT token (access or refresh)
func (j *JwtAuth) verifyToken(token string, secretKey string, expectedType string) (core.AuthPayload, error) {
	if token == "" {
		return core.AuthPayload{}, errors.New("missing token")
	}

	parsedToken, err := jwtpackage.ParseWithClaims(token, &jwtPayload{}, func(t *jwtpackage.Token) (any, error) {
		if _, ok := t.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return core.AuthPayload{}, err
	}

	claims, ok := parsedToken.Claims.(*jwtPayload)
	if !ok {
		return core.AuthPayload{}, errors.New("invalid claims format")
	}
	if claims.TokenType != expectedType {
		return core.AuthPayload{}, errors.New("invalid token type")
	}
	if claims.AccountId == 0 || claims.TokenVersion == 0 {
		return core.AuthPayload{}, errors.New("invalid token payload")
	}

	return claims.AuthPayload, nil
}

// RevokeRefreshToken is a no-op in JWT-based authentication.
func (j *JwtAuth) RevokeRefreshToken(token string) error {
	return nil
}
