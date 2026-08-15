package core

import (
	"testing"

	"scaf-gin/config"
)

func TestJwtAuthCreatesVerifiableTokenPair(t *testing.T) {
	auth := NewJwtAuth(config.Config{
		AccessTokenSecret:                    "access-secret",
		AccessTokenExpiresSeconds:            900,
		RefreshTokenSecret:                   "refresh-secret",
		RefreshTokenExpiresSeconds:           43_200,
		RefreshTokenRememberMeExpiresSeconds: 2_592_000,
	})
	payload := AuthPayload{AccountID: 42, TokenVersion: 3}

	accessToken, err := auth.CreateAccessToken(payload)
	if err != nil {
		t.Fatalf("create access token: %v", err)
	}
	refreshToken, err := auth.CreateRefreshToken(payload, false)
	if err != nil {
		t.Fatalf("create refresh token: %v", err)
	}

	accessPayload, err := auth.VerifyAccessToken(accessToken)
	if err != nil {
		t.Fatalf("verify access token: %v", err)
	}
	refreshPayload, err := auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("verify refresh token: %v", err)
	}
	if accessPayload != payload {
		t.Fatalf("unexpected access payload: %#v", accessPayload)
	}
	if refreshPayload != payload {
		t.Fatalf("unexpected refresh payload: %#v", refreshPayload)
	}
	if _, err := auth.VerifyAccessToken(refreshToken); err == nil {
		t.Fatal("refresh token must not verify as an access token")
	}
}
