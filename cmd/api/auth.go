package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"time"
)

type sessionManager struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
	ttl      time.Duration
}

func newSessionManager(ttl time.Duration) *sessionManager {
	return &sessionManager{
		sessions: make(map[string]time.Time),
		ttl:      ttl,
	}
}

func (sm *sessionManager) create() (string, time.Time, error) {
	sm.cleanupExpired()

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(tokenBytes)
	expiry := time.Now().Add(sm.ttl)

	sm.mu.Lock()
	sm.sessions[token] = expiry
	sm.mu.Unlock()

	return token, expiry, nil
}

func (sm *sessionManager) validate(token string) bool {
	if token == "" {
		return false
	}

	sm.mu.RLock()
	expiry, ok := sm.sessions[token]
	sm.mu.RUnlock()

	if !ok {
		return false
	}

	if time.Now().After(expiry) {
		sm.mu.Lock()
		delete(sm.sessions, token)
		sm.mu.Unlock()
		return false
	}

	return true
}

func (sm *sessionManager) revoke(token string) {
	if token == "" {
		return
	}

	sm.mu.Lock()
	delete(sm.sessions, token)
	sm.mu.Unlock()
}

func (sm *sessionManager) cleanupExpired() {
	now := time.Now()

	sm.mu.Lock()
	for token, expiry := range sm.sessions {
		if now.After(expiry) {
			delete(sm.sessions, token)
		}
	}
	sm.mu.Unlock()
}

func secureCompare(given, actual string) bool {
	given = strings.TrimSpace(given)
	actual = strings.TrimSpace(actual)
	if len(given) == 0 || len(actual) == 0 {
		return false
	}
	if len(given) != len(actual) {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1 {
		return true
	}
	return false
}

func parseBearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header is not a bearer token")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("authorization token missing")
	}
	return token, nil
}
