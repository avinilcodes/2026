package leader

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// LeaderLock represents the leader lock entry in database
type LeaderLock struct {
	OwnerID   string
	ExpiresAt time.Time
}

// LeaderElection handles distributed leader election with DB-based lock
type LeaderElection struct {
	db              *sql.DB
	ownerID         string
	lockTTL         time.Duration
	heartbeatTicker *time.Ticker
	mu              sync.RWMutex
	isLeader        bool
	stopChan        chan struct{}
	mockDB          *MockDB // For testing without real database
}

// NewLeaderElection creates a new leader election instance
func NewLeaderElection(db *sql.DB, ownerID string) *LeaderElection {
	return &LeaderElection{
		db:       db,
		ownerID:  ownerID,
		lockTTL:  5 * time.Second,
		stopChan: make(chan struct{}),
		isLeader: false,
	}
}

// InitializeDatabase creates the leader_lock table
func (le *LeaderElection) InitializeDatabase(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS leader_lock (
		owner_id VARCHAR(255) PRIMARY KEY,
		expires_at TIMESTAMP NOT NULL,
		INDEX idx_expires_at (expires_at)
	);
	`

	_, err := le.db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create leader_lock table: %w", err)
	}

	// Clean up expired locks
	_, err = le.db.ExecContext(ctx, `
		DELETE FROM leader_lock 
		WHERE expires_at < NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired locks: %w", err)
	}

	return nil
}

// TryAcquireLock attempts to acquire the leadership lock
func (le *LeaderElection) TryAcquireLock(ctx context.Context) (bool, error) {
	// Use mock database if no real DB connection
	if le.mockDB != nil {
		acquired, err := le.mockDB.TryAcquireLock(le.ownerID, le.lockTTL)
		le.mu.Lock()
		le.isLeader = acquired
		le.mu.Unlock()
		if acquired {
			log.Printf("[%s] Successfully acquired leadership lock (mock)", le.ownerID)
		}
		return acquired, err
	}

	expiresAt := time.Now().Add(le.lockTTL)

	// Try to acquire lock only if it's expired or doesn't exist
	_, err := le.db.ExecContext(ctx, `
		INSERT INTO leader_lock (owner_id, expires_at)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			owner_id = IF(expires_at < NOW(), ?, owner_id),
			expires_at = IF(expires_at < NOW(), ?, expires_at)
	`, le.ownerID, expiresAt, le.ownerID, expiresAt)

	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	// Verify we acquired the lock
	acquired, err := le.verifyLockOwnership(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to verify lock ownership: %w", err)
	}

	le.mu.Lock()
	le.isLeader = acquired
	le.mu.Unlock()

	if acquired {
		log.Printf("[%s] Successfully acquired leadership lock", le.ownerID)
	}

	return acquired, nil
}

// RenewLock renews the leadership lock (heartbeat)
func (le *LeaderElection) RenewLock(ctx context.Context) (bool, error) {
	// Use mock database if no real DB connection
	if le.mockDB != nil {
		renewed, err := le.mockDB.RenewLock(le.ownerID, le.lockTTL)
		le.mu.Lock()
		le.isLeader = renewed
		le.mu.Unlock()
		if renewed {
			log.Printf("[%s] Lock renewed successfully (mock)", le.ownerID)
		} else {
			log.Printf("[%s] Failed to renew lock - leadership lost (mock)", le.ownerID)
		}
		return renewed, err
	}

	expiresAt := time.Now().Add(le.lockTTL)

	// Only renew if we own the lock
	result, err := le.db.ExecContext(ctx, `
		UPDATE leader_lock 
		SET expires_at = ?
		WHERE owner_id = ? AND expires_at > NOW()
	`, expiresAt, le.ownerID)

	if err != nil {
		return false, fmt.Errorf("failed to renew lock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	renewed := rowsAffected > 0

	le.mu.Lock()
	le.isLeader = renewed
	le.mu.Unlock()

	if renewed {
		log.Printf("[%s] Lock renewed successfully", le.ownerID)
	} else {
		log.Printf("[%s] Failed to renew lock - leadership lost", le.ownerID)
	}

	return renewed, nil
}

// verifyLockOwnership checks if this instance owns the lock
func (le *LeaderElection) verifyLockOwnership(ctx context.Context) (bool, error) {
	// Use mock database if no real DB connection
	if le.mockDB != nil {
		lock, err := le.mockDB.GetLockStatus()
		if err != nil {
			return false, nil
		}
		return lock.OwnerID == le.ownerID, nil
	}

	var ownerID string

	err := le.db.QueryRowContext(ctx, `
		SELECT owner_id FROM leader_lock 
		WHERE owner_id = ? AND expires_at > NOW()
	`, le.ownerID).Scan(&ownerID)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return ownerID == le.ownerID, nil
}

// StartHeartbeat starts the renewal heartbeat loop
func (le *LeaderElection) StartHeartbeat(ctx context.Context, heartbeatInterval time.Duration) {
	le.heartbeatTicker = time.NewTicker(heartbeatInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				le.heartbeatTicker.Stop()
				return
			case <-le.stopChan:
				le.heartbeatTicker.Stop()
				return
			case <-le.heartbeatTicker.C:
				// Try to renew the lock
				_, err := le.RenewLock(ctx)
				if err != nil {
					log.Printf("[%s] Heartbeat error: %v", le.ownerID, err)
				}
			}
		}
	}()

	log.Printf("[%s] Heartbeat started with interval: %v", le.ownerID, heartbeatInterval)
}

// Stop releases the lock and stops the heartbeat
func (le *LeaderElection) Stop(ctx context.Context) error {
	close(le.stopChan)

	if le.heartbeatTicker != nil {
		le.heartbeatTicker.Stop()
	}

	// Release the lock
	if le.mockDB != nil {
		le.mockDB.ReleaseLock(le.ownerID)
	} else {
		_, err := le.db.ExecContext(ctx, `
			DELETE FROM leader_lock 
			WHERE owner_id = ?
		`, le.ownerID)

		if err != nil {
			return fmt.Errorf("failed to release lock: %w", err)
		}
	}

	le.mu.Lock()
	le.isLeader = false
	le.mu.Unlock()

	log.Printf("[%s] Leadership lock released", le.ownerID)
	return nil
}

// IsLeader returns whether this instance is the current leader
func (le *LeaderElection) IsLeader() bool {
	le.mu.RLock()
	defer le.mu.RUnlock()
	return le.isLeader
}

// GetCurrentLeader returns the current leader's owner_id
func (le *LeaderElection) GetCurrentLeader(ctx context.Context) (string, error) {
	// Use mock database if no real DB connection
	if le.mockDB != nil {
		return le.mockDB.GetCurrentLeader()
	}

	var ownerID string

	err := le.db.QueryRowContext(ctx, `
		SELECT owner_id FROM leader_lock 
		WHERE expires_at > NOW()
		LIMIT 1
	`).Scan(&ownerID)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no active leader")
	}
	if err != nil {
		return "", err
	}

	return ownerID, nil
}

// WaitForLeadership blocks until this instance becomes the leader
func (le *LeaderElection) WaitForLeadership(ctx context.Context, maxWait time.Duration) error {
	deadline := time.Now().Add(maxWait)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-le.stopChan:
			return fmt.Errorf("leader election stopped")
		case <-ticker.C:
			if le.IsLeader() {
				return nil
			}

			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for leadership after %v", maxWait)
			}
		}
	}
}

// GetLockStatus returns the current lock status from database
func (le *LeaderElection) GetLockStatus(ctx context.Context) (*LeaderLock, error) {
	// Use mock database if no real DB connection
	if le.mockDB != nil {
		return le.mockDB.GetLockStatus()
	}

	var lock LeaderLock

	err := le.db.QueryRowContext(ctx, `
		SELECT owner_id, expires_at FROM leader_lock 
		WHERE expires_at > NOW()
		LIMIT 1
	`).Scan(&lock.OwnerID, &lock.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active lock")
	}
	if err != nil {
		return nil, err
	}

	return &lock, nil
}
