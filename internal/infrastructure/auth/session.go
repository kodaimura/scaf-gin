package auth

import (
	"errors"
	"strconv"
	"time"

	"scaf-gin/internal/core"
)

// SessionAuth implements AuthI using in-memory sessions.
// Note: This is not suitable for production without a persistent session store (e.g., Redis).
type SessionAuth struct {
	sessions map[string]core.AuthPayload
}

func NewSessionAuth() core.AuthI {
	return &SessionAuth{
		sessions: make(map[string]core.AuthPayload),
	}
}

// GenerateToken creates a new session ID and stores the AuthPayload.
func (s *SessionAuth) GenerateToken(payload core.AuthPayload) (string, error) {
	sessionID := generateSessionID(payload.AccountId)
	s.sessions[sessionID] = payload
	return sessionID, nil
}

// ValidateToken checks if the session ID exists and returns the associated AuthPayload.
func (s *SessionAuth) ValidateToken(token string) (core.AuthPayload, error) {
	payload, exists := s.sessions[token]
	if !exists {
		return core.AuthPayload{}, errors.New("session not found or expired")
	}
	return payload, nil
}

// RevokeToken deletes the session entry associated with the token.
func (s *SessionAuth) RevokeToken(token string) error {
	delete(s.sessions, token)
	return nil
}

func generateSessionID(accountId int) string {
	return strconv.Itoa(accountId) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
}
