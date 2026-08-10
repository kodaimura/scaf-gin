package core

import (
	"errors"
	"strconv"
	"time"

	jwtpackage "github.com/golang-jwt/jwt/v5"

	"scaf-gin/config"
)

type Auth interface {
	CreateAccessToken(payload AuthPayload) (string, error)
	CreateRefreshToken(payload AuthPayload, rememberMe bool) (string, error)
	VerifyAccessToken(token string) (AuthPayload, error)
	VerifyRefreshToken(token string) (AuthPayload, error)
}

type AuthPayload struct {
	AccountId    int
	TokenVersion int
}

// JwtAuth implements Auth using JWT.
type JwtAuth struct {
	cfg config.Config
}

func NewJwtAuth(cfg config.Config) Auth {
	return &JwtAuth{cfg: cfg}
}

type jwtPayload struct {
	jwtpackage.RegisteredClaims
	TokenVersion int    `json:"token_version"`
	TokenType    string `json:"type"`
}

func (j *JwtAuth) CreateAccessToken(payload AuthPayload) (string, error) {
	return j.createToken(
		payload,
		j.cfg.AccessTokenSecret,
		"access",
		time.Second*time.Duration(j.cfg.AccessTokenExpiresSeconds),
	)
}

func (j *JwtAuth) CreateRefreshToken(payload AuthPayload, rememberMe bool) (string, error) {
	expiresSeconds := j.cfg.RefreshTokenExpiresSeconds
	if rememberMe {
		expiresSeconds = j.cfg.RefreshTokenRememberMeExpiresSeconds
	}
	return j.createToken(
		payload,
		j.cfg.RefreshTokenSecret,
		"refresh",
		time.Second*time.Duration(expiresSeconds),
	)
}

func (j *JwtAuth) createToken(payload AuthPayload, secretKey string, tokenType string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	jp := jwtPayload{
		TokenVersion: payload.TokenVersion,
		TokenType:    tokenType,
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
	return j.verifyToken(token, j.cfg.AccessTokenSecret, "access")
}

func (j *JwtAuth) VerifyRefreshToken(token string) (AuthPayload, error) {
	return j.verifyToken(token, j.cfg.RefreshTokenSecret, "refresh")
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
	accountID, err := strconv.Atoi(claims.Subject)
	if err != nil || accountID == 0 || claims.TokenVersion == 0 {
		return AuthPayload{}, errors.New("invalid token payload")
	}

	return AuthPayload{
		AccountId:    accountID,
		TokenVersion: claims.TokenVersion,
	}, nil
}
