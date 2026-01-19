package main

import (
	"context"
	"fmt"
	"time"
)

// Retry attempts to execute a function with exponential backoff retries.
// It takes a context, number of attempts, and a function to retry.
// The function returns an error if all attempts fail, or nil if it succeeds.
func Retry(ctx context.Context, attempts int, fn func(context.Context) error) error {
	if attempts <= 0 {
		return fmt.Errorf("attempts must be greater than 0")
	}

	var lastErr error
	backoff := time.Millisecond * 100

	for attempt := 1; attempt <= attempts; attempt++ {
		// Check if context is already cancelled
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		// Try executing the function
		err := fn(ctx)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// If this is the last attempt, return the error
		if attempt == attempts {
			return fmt.Errorf("max retries exceeded: %w", lastErr)
		}

		// Wait before retrying with exponential backoff
		select {
		case <-time.After(backoff):
			// Continue to next attempt
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		}

		// Exponential backoff: double the backoff time for next attempt (capped at 30 seconds)
		backoff *= 2
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
	}

	return fmt.Errorf("retry failed: %w", lastErr)
}

func main() {
	ctx := context.Background()
	err := Retry(ctx, 5, func(ctx context.Context) error {
		// Simulate a function that may fail
		fmt.Println("Attempting operation...")
		return fmt.Errorf("operation failed")
	})
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	}
}
