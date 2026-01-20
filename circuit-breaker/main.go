package main

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState string

const (
	Closed   CircuitState = "CLOSED"    // Normal operation
	Open     CircuitState = "OPEN"      // Failing, reject requests
	HalfOpen CircuitState = "HALF_OPEN" // Testing if service recovered
)

// CircuitBreaker manages the circuit breaker pattern
type CircuitBreaker struct {
	maxFailures     int           // Number of failures before opening
	openTimeout     time.Duration // Time to wait before trying again
	lastFailureTime time.Time
	failureCount    int
	state           CircuitState
	mu              sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		openTimeout:  openTimeout,
		state:        Closed,
		failureCount: 0,
	}
}

// Call wraps an external/DB call with circuit breaker logic
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if circuit should transition from Open to HalfOpen
	if cb.state == Open {
		if time.Since(cb.lastFailureTime) > cb.openTimeout {
			cb.state = HalfOpen
		} else {
			return fmt.Errorf("circuit breaker is open, service unavailable (retry in %v)",
				cb.openTimeout-time.Since(cb.lastFailureTime))
		}
	}

	// Execute the function
	err := fn()

	if err != nil {
		// Record failure
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Open circuit if threshold exceeded
		if cb.failureCount >= cb.maxFailures {
			cb.state = Open
			return fmt.Errorf("circuit breaker opened after %d failures: %w",
				cb.maxFailures, err)
		}
		return fmt.Errorf("operation failed (failures: %d/%d): %w",
			cb.failureCount, cb.maxFailures, err)
	}

	// Success - reset on HalfOpen or maintain Closed state
	if cb.state == HalfOpen {
		cb.state = Closed
		cb.failureCount = 0
	}

	return nil
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failureCount
}

// Reset manually resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = Closed
	cb.failureCount = 0
	cb.lastFailureTime = time.Time{}
}

// ============ Example Usage ============

// Simulated external service
type ExternalService struct {
	requestCount int
	failUntil    int // Fail first N requests
}

func (es *ExternalService) FetchData() error {
	es.requestCount++
	if es.requestCount <= es.failUntil {
		return fmt.Errorf("service temporarily unavailable")
	}
	return nil
}

func main() {
	// Create circuit breaker: max 3 failures, 2 second timeout before retry
	cb := NewCircuitBreaker(3, 2*time.Second)
	service := &ExternalService{failUntil: 4}

	fmt.Println("=== Circuit Breaker Demo ===\n")

	// Make 6 requests
	for i := 1; i <= 6; i++ {
		fmt.Printf("Request %d: ", i)
		err := cb.Call(func() error {
			return service.FetchData()
		})

		if err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Println("✓ Success")
		}

		fmt.Printf("   State: %v, Failures: %d\n\n",
			cb.GetState(), cb.GetFailureCount())

		time.Sleep(300 * time.Millisecond)
	}

	// Wait and observe circuit recovery
	fmt.Println("\n=== Waiting for circuit to timeout ===")
	fmt.Printf("Circuit state: %v\n", cb.GetState())
	time.Sleep(2500 * time.Millisecond)

	// Try again - circuit should transition to HalfOpen
	fmt.Printf("\nRequest 7 (after timeout): ")
	err := cb.Call(func() error {
		return service.FetchData()
	})

	if err != nil {
		fmt.Printf("❌ %v\n", err)
	} else {
		fmt.Println("✓ Success - Circuit recovered!")
	}
	fmt.Printf("Final State: %v, Failures: %d\n", cb.GetState(), cb.GetFailureCount())
}
