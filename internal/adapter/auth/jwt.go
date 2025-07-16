package auth

import (
	"encoding/json"
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
}

// CreateAccessToken creates a signed JWT containing the given AuthPayload.
func (j *JwtAuth) CreateAccessToken(payload core.AuthPayload) (string, error) {
	return j.createToken(
		payload,
		config.AccessTokenSecret,
		time.Second*time.Duration(config.AccessTokenExpiresSeconds))
}

// CreateRefreshToken creates a signed JWT containing the given AuthPayload.
func (j *JwtAuth) CreateRefreshToken(payload core.AuthPayload) (string, error) {
	return j.createToken(
		payload,
		config.RefreshTokenSecret,
		time.Second*time.Duration(config.RefreshTokenExpiresSeconds),
	)
}

// Common function to generate a JWT token (access or refresh)
func (j *JwtAuth) createToken(payload core.AuthPayload, secretKey string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	jp := jwtPayload{
		AuthPayload: payload,
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
	return j.verifyToken(token, config.AccessTokenSecret)
}

// VerifyRefreshToken verifies the given refresh token and extracts the AuthPayload.
func (j *JwtAuth) VerifyRefreshToken(token string) (core.AuthPayload, error) {
	return j.verifyToken(token, config.RefreshTokenSecret)
}

// Common function to validate a JWT token (access or refresh)
func (j *JwtAuth) verifyToken(token string, secretKey string) (core.AuthPayload, error) {
	parsedToken, err := jwtpackage.Parse(token, func(t *jwtpackage.Token) (any, error) {
		if _, ok := t.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return core.AuthPayload{}, err
	}

	return tokenToAuthPayload(parsedToken)
}

func tokenToAuthPayload(token *jwtpackage.Token) (core.AuthPayload, error) {
	var jp jwtPayload

	claimsMap, ok := token.Claims.(jwtpackage.MapClaims)
	if !ok {
		return core.AuthPayload{}, errors.New("invalid claims format")
	}

	jsonBytes, err := json.Marshal(claimsMap)
	if err != nil {
		return core.AuthPayload{}, err
	}

	if err := json.Unmarshal(jsonBytes, &jp); err != nil {
		return core.AuthPayload{}, err
	}

	return jp.AuthPayload, nil
}

// RevokeRefreshToken is a no-op in JWT-based authentication.
func (j *JwtAuth) RevokeRefreshToken(token string) error {
	return nil
}
