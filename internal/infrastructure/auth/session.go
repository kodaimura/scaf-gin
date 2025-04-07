package auth

import (
	"strconv"
	"errors"
	"time"

	"goscaf/internal/core"
)

type SessionAuth struct {
	// sessions stores session data in memory.
	// You should implement a persistent session store, 
	// such as Redis or a database, for production environments.
	sessions map[string]core.AuthPayload
}

func NewSessionAuth() core.AuthI {
	return &SessionAuth{
		sessions: make(map[string]core.AuthPayload),
	}
}

func (s *SessionAuth) GenerateToken(payload core.AuthPayload) (string, error) {
	sessionID := generateSessionID(payload.AccountId)
	s.sessions[sessionID] = payload
	return sessionID, nil
}

func (s *SessionAuth) ValidateToken(token string) (core.AuthPayload, error) {
	payload, exists := s.sessions[token]
	if !exists {
		return core.AuthPayload{}, errors.New("invalid session")
	}
	return payload, nil
}

func (s *SessionAuth) RevokeToken(token string) error {
	delete(s.sessions, token)
	return nil
}

func generateSessionID(accountId int) string {
	return strconv.Itoa(accountId) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
}
