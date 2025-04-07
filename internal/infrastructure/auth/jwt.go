package auth

import (
	"strconv"
	"time"
	"encoding/json"
	"errors"
	jwtpackage "github.com/golang-jwt/jwt/v5"

	"goscaf/config"
	"goscaf/internal/core"
)

type JwtAuth struct {}

type JwtPayload struct {
	jwtpackage.RegisteredClaims
	core.AuthPayload
}

func NewJwtAuth() *JwtAuth {
    return &JwtAuth{}
}

func (j *JwtAuth) GenerateToken(payload core.AuthPayload) (string, error) {
	jwtPayload := JwtPayload{
		AuthPayload: payload,
		RegisteredClaims: jwtpackage.RegisteredClaims{
			Subject:   strconv.Itoa(payload.AccountId),
			ExpiresAt: jwtpackage.NewNumericDate(time.Now().Add(time.Second * time.Duration(config.AuthExpiresSeconds))),
			NotBefore: jwtpackage.NewNumericDate(time.Now()),
			IssuedAt: jwtpackage.NewNumericDate(time.Now()),
		},
	}

	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, jwtPayload)
	return token.SignedString([]byte(config.JwtSecretKey))
}

func (j *JwtAuth) ValidateToken(token string) (core.AuthPayload, error) {
	parsedToken, err := jwtpackage.Parse(token, func(parsedToken *jwtpackage.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(config.JwtSecretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return core.AuthPayload{}, err
	}

	return tokenToAuthPayload(parsedToken)
}


func tokenToAuthPayload (token *jwtpackage.Token) (core.AuthPayload, error) {
	var jwtPayload JwtPayload

	jsonString, err := json.Marshal(token.Claims.(jwtpackage.MapClaims))
	if err != nil {
		return core.AuthPayload{}, err
	}
	if err = json.Unmarshal(jsonString, &jwtPayload); err != nil {
		return core.AuthPayload{}, err
	}
	return jwtPayload.AuthPayload, nil
}

func (j *JwtAuth) RevokeToken(token string) error {
	return nil
}