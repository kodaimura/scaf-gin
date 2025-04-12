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

// GenerateToken creates a signed JWT containing the given AuthPayload.
func (j *JwtAuth) GenerateToken(payload core.AuthPayload) (string, error) {
	now := time.Now()

	jp := jwtPayload{
		AuthPayload: payload,
		RegisteredClaims: jwtpackage.RegisteredClaims{
			Subject:   strconv.Itoa(payload.AccountId),
			IssuedAt:  jwtpackage.NewNumericDate(now),
			NotBefore: jwtpackage.NewNumericDate(now),
			ExpiresAt: jwtpackage.NewNumericDate(now.Add(time.Second * time.Duration(config.AuthExpiresSeconds))),
		},
	}

	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, jp)
	return token.SignedString([]byte(config.JwtSecretKey))
}

// ValidateToken verifies the given JWT and extracts the AuthPayload.
// Returns an error if the token is invalid or cannot be parsed.
func (j *JwtAuth) ValidateToken(token string) (core.AuthPayload, error) {
	parsedToken, err := jwtpackage.Parse(token, func(t *jwtpackage.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JwtSecretKey), nil
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

// RevokeToken is a no-op in JWT-based authentication.
// In a real-world app, this might involve blacklisting the token.
func (j *JwtAuth) RevokeToken(token string) error {
	return nil
}
