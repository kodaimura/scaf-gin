package core

import (
	"errors"
	"strconv"
	"time"

	jwtpackage "github.com/golang-jwt/jwt/v5"

	"scaf-gin/config"
)

type AuthI interface {
	CreateAccessToken(payload AuthPayload) (string, error)
	CreateRefreshToken(payload AuthPayload, rememberMe bool) (string, error)
	VerifyAccessToken(token string) (AuthPayload, error)
	VerifyRefreshToken(token string) (AuthPayload, error)
	RevokeRefreshToken(token string) error
}

type AuthPayload struct {
	AccountId    int    `json:"account_id"`
	LoginID      string `json:"login_id"`
	TokenVersion int    `json:"token_version"`
}

// JwtAuth implements AuthI using JWT.
type JwtAuth struct{}

func NewJwtAuth() AuthI {
	return &JwtAuth{}
}

type jwtPayload struct {
	jwtpackage.RegisteredClaims
	AuthPayload
	TokenType string `json:"type"`
}

func (j *JwtAuth) CreateAccessToken(payload AuthPayload) (string, error) {
	return j.createToken(
		payload,
		config.AccessTokenSecret,
		"access",
		time.Second*time.Duration(config.AccessTokenExpiresSeconds),
	)
}

func (j *JwtAuth) CreateRefreshToken(payload AuthPayload, rememberMe bool) (string, error) {
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

func (j *JwtAuth) createToken(payload AuthPayload, secretKey string, tokenType string, expiresIn time.Duration) (string, error) {
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

func (j *JwtAuth) VerifyAccessToken(token string) (AuthPayload, error) {
	return j.verifyToken(token, config.AccessTokenSecret, "access")
}

func (j *JwtAuth) VerifyRefreshToken(token string) (AuthPayload, error) {
	return j.verifyToken(token, config.RefreshTokenSecret, "refresh")
}

func (j *JwtAuth) verifyToken(token string, secretKey string, expectedType string) (AuthPayload, error) {
	if token == "" {
		return AuthPayload{}, errors.New("missing token")
	}

	parsedToken, err := jwtpackage.ParseWithClaims(token, &jwtPayload{}, func(t *jwtpackage.Token) (any, error) {
		if _, ok := t.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return AuthPayload{}, err
	}

	claims, ok := parsedToken.Claims.(*jwtPayload)
	if !ok {
		return AuthPayload{}, errors.New("invalid claims format")
	}
	if claims.TokenType != expectedType {
		return AuthPayload{}, errors.New("invalid token type")
	}
	if claims.AccountId == 0 || claims.TokenVersion == 0 {
		return AuthPayload{}, errors.New("invalid token payload")
	}

	return claims.AuthPayload, nil
}

// RevokeRefreshToken is a no-op in JWT-based authentication.
func (j *JwtAuth) RevokeRefreshToken(token string) error {
	return nil
}
