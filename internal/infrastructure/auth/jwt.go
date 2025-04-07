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

func (j *JwtAuth) GenerateCredential(payload core.AuthPayload) (string, error) {
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

func (j *JwtAuth) ValidateCredential(credential string) (core.AuthPayload, error) {
	token, err := jwtpackage.Parse(credential, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(config.JwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return core.AuthPayload{}, err
	}

	return tokenToAuthPayload(token)
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

func (j *JwtAuth) RevokeCredential(credential string) error {
	return nil
}