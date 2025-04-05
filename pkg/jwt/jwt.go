package jwt

import (
	"strconv"
	"time"
	"encoding/json"
	"errors"

	jwtpackage "github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	jwtpackage.RegisteredClaims
	CustomClaims map[string]interface{}
}

func NewPayload(sub int, expires int, claims map[string]interface{}) Payload {
	var pl Payload

	pl.CustomClaims = claims
	pl.Subject = strconv.Itoa(sub)
	pl.ExpiresAt = jwtpackage.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires)))
	pl.NotBefore = jwtpackage.NewNumericDate(time.Now())
	pl.IssuedAt = jwtpackage.NewNumericDate(time.Now())

	return pl
}

func ExpireToken (pl Payload) Payload {
	pl.IssuedAt =  jwtpackage.NewNumericDate(time.Now())
	pl.ExpiresAt = jwtpackage.NewNumericDate(time.Now())
	return pl
}  


func EncodeToken (pl Payload, secret string) (string, error) {
	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, pl)
	return token.SignedString([]byte(secret))
}


func DecodeToken (encoded string, secret string) (Payload, error) {
	token, err := jwtpackage.Parse(encoded, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return Payload{}, err
	}

	return convertToPayload(token)
}


func convertToPayload (token *jwtpackage.Token) (Payload, error) {
	var pl Payload

	jsonString, err := json.Marshal(token.Claims.(jwtpackage.MapClaims))

	if err == nil {
		err = json.Unmarshal(jsonString, &pl)
	}

	return pl, err
}