package jwt

import (
	"time"
	"encoding/json"
	"errors"
	"strings"

	jwtpackage "github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	jwt.StandardClaims
	CustomClaims map[string]interface{}
}

func NewPayload(sub string, expires int64, claims map[string]interface{}) Payload {
	var pl Payload

	pl.CustomClaims = claims
	pl.Subject = sub
	pl.ExpiresAt = time.Now().Add(expires).Unix()
	pl.NotBefore = time.Now().Unix()
	pl.IssuedAt = time.Now().Unix()

	return pl
}

func ExpireToken (pl Payload) Payload {
	pl.IssuedAt =  time.Now().Unix()
	pl.ExpiresAt = time.Now().Unix()
	return pl
}  


func EncodeToken (pl Payload) (string, error) {
	cf := config.GetConfig()
	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, pl)
	return token.SignedString([]byte(cf.JwtSecretKey))
}


func DecodeToken (encoded string) (Payload, error) {
	cf := config.GetConfig()
	token, err := jwtpackage.Parse(encoded, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(cf.JwtSecretKey), nil
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