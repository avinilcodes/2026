package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthManager handles JWT issuance and refresh token rotation.
type AuthManager struct {
	secret       []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
	refreshStore *inMemoryRefreshStore
}

// NewAuthManager returns a new AuthManager using env vars or defaults.
func NewAuthManager() *AuthManager {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "dev-secret-change-me"
	}
	am := &AuthManager{
		secret:       []byte(sec),
		accessTTL:    15 * time.Minute,
		refreshTTL:   7 * 24 * time.Hour,
		refreshStore: newInMemoryRefreshStore(),
	}
	return am
}

// IssueTokens issues an access token (JWT) and a refresh token (opaque string).
func (a *AuthManager) IssueTokens(ctx context.Context, subject string) (access string, refresh string, err error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	access, err = at.SignedString(a.secret)
	if err != nil {
		return "", "", err
	}
	refresh, err = generateRandomToken(64)
	if err != nil {
		return "", "", err
	}
	a.refreshStore.put(refresh, subject, time.Now().Add(a.refreshTTL))
	return access, refresh, nil
}

// RotateRefresh consumes the provided refresh token and returns a new pair.
// It enforces single-use: the old refresh token is deleted on successful rotation.
func (a *AuthManager) RotateRefresh(ctx context.Context, refresh string) (newAccess, newRefresh string, err error) {
	subj, ok := a.refreshStore.getAndDelete(refresh)
	if !ok {
		return "", "", errors.New("invalid or expired refresh token")
	}
	return a.IssueTokens(ctx, subj)
}

// ValidateAccess validates a JWT and returns the subject if valid.
func (a *AuthManager) ValidateAccess(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", errors.New("invalid token")
}

// generateRandomToken returns a URL-safe base64 string of the requested size in bytes.
func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// in-memory refresh store (single-process). Keys expire after TTL but explicit
// deletion on rotation is used to implement single-use refresh tokens.
type inMemoryRefreshStore struct {
	mu    sync.Mutex
	store map[string]refreshEntry
}

type refreshEntry struct {
	subject string
	expiry  time.Time
}

func newInMemoryRefreshStore() *inMemoryRefreshStore {
	s := &inMemoryRefreshStore{
		store: make(map[string]refreshEntry),
	}
	go s.evictLoop()
	return s
}

func (s *inMemoryRefreshStore) put(token, subject string, expiry time.Time) {
	s.mu.Lock()
	s.store[token] = refreshEntry{subject: subject, expiry: expiry}
	s.mu.Unlock()
}

func (s *inMemoryRefreshStore) getAndDelete(token string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.store[token]
	if !ok {
		return "", false
	}
	if time.Now().After(e.expiry) {
		delete(s.store, token)
		return "", false
	}
	delete(s.store, token)
	return e.subject, true
}

func (s *inMemoryRefreshStore) evictLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.store {
			if now.After(v.expiry) {
				delete(s.store, k)
			}
		}
		s.mu.Unlock()
	}
}
