package leader

import (
	"fmt"
	"sync"
	"time"
)

// MockDB provides an in-memory database implementation for testing
type MockDB struct {
	locks map[string]time.Time // owner_id -> expires_at
	mu    sync.RWMutex
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		locks: make(map[string]time.Time),
	}
}

// GetMockLeaderElection creates a LeaderElection with mock database
func GetMockLeaderElection(ownerID string) *LeaderElection {
	return &LeaderElection{
		db:       nil, // Will use mock methods
		ownerID:  ownerID,
		lockTTL:  5 * time.Second,
		stopChan: make(chan struct{}),
		isLeader: false,
		mockDB:   NewMockDB(),
	}
}

// TryAcquireLock simulates lock acquisition in mock database
func (m *MockDB) TryAcquireLock(ownerID string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	expiresAt := now.Add(ttl)

	// Check if any lock exists and is still valid
	for owner, expiry := range m.locks {
		if expiry.After(now) {
			// Lock still valid, cannot acquire
			if owner != ownerID {
				return false, nil
			}
		} else {
			// Lock expired, remove it
			delete(m.locks, owner)
		}
	}

	// Acquire lock
	m.locks[ownerID] = expiresAt
	return true, nil
}

// RenewLock simulates lock renewal in mock database
func (m *MockDB) RenewLock(ownerID string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// Check if this owner has a valid lock
	if expiry, exists := m.locks[ownerID]; exists && expiry.After(now) {
		m.locks[ownerID] = now.Add(ttl)
		return true, nil
	}

	return false, nil
}

// ReleaseLock simulates lock release in mock database
func (m *MockDB) ReleaseLock(ownerID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.locks, ownerID)
	return nil
}

// GetCurrentLeader returns the current leader
func (m *MockDB) GetCurrentLeader() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()

	for owner, expiry := range m.locks {
		if expiry.After(now) {
			return owner, nil
		}
	}

	return "", fmt.Errorf("no active leader")
}

// GetLockStatus returns lock details
func (m *MockDB) GetLockStatus() (*LeaderLock, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()

	for owner, expiry := range m.locks {
		if expiry.After(now) {
			return &LeaderLock{
				OwnerID:   owner,
				ExpiresAt: expiry,
			}, nil
		}
	}

	return nil, fmt.Errorf("no active lock")
}

// CleanupExpiredLocks removes expired locks
func (m *MockDB) CleanupExpiredLocks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for owner, expiry := range m.locks {
		if !expiry.After(now) {
			delete(m.locks, owner)
		}
	}
}
