package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leader-election/leader"
)

func main() {
	fmt.Println("Welcome to the Spam Detection Service - Leader Election Demo")
	fmt.Println("Using in-memory mock database (no MySQL required)\n")

	// Create leader election using mock database
	leaderElection := leader.GetMockLeaderElection("demo-instance-1")

	// Try to acquire leadership lock
	acquired, err := leaderElection.TryAcquireLock(context.Background())
	if err != nil {
		log.Printf("Error acquiring lock: %v", err)
	}

	if acquired {
		log.Println("✓ Successfully acquired leadership lock")
	} else {
		log.Println("✗ Failed to acquire leadership lock")
	}

	// Start heartbeat renewal (every 1.5 seconds to renew before 5 second TTL expires)
	ctx, cancel := context.WithCancel(context.Background())
	leaderElection.StartHeartbeat(ctx, 1500*time.Millisecond)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main application loop
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("\nLeadership monitoring started. Press Ctrl+C to stop.\n")

	for {
		select {
		case <-sigChan:
			log.Println("\n\nShutting down...")
			cancel()
			if err := leaderElection.Stop(context.Background()); err != nil {
				log.Printf("Error stopping leader election: %v", err)
			}
			log.Println("Shutdown complete")
			return

		case <-ticker.C:
			status, err := leaderElection.GetLockStatus(context.Background())
			if err != nil {
				log.Printf("Lock status: No active leader - %v", err)
			} else {
				timeUntilExpiry := time.Until(status.ExpiresAt)
				log.Printf("Lock status: Owner=%s | Expires in: %.1fs", status.OwnerID, timeUntilExpiry.Seconds())
			}

			if leaderElection.IsLeader() {
				log.Println("          → [LEADER] Executing leader-only tasks...")
			}
		}
	}
}
