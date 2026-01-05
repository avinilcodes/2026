package main

import (
	"context"
	"fmt"
	"time"
)

func work(ctx context.Context) {
	for i := 0; i < 5; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("working step", i)
		case <-ctx.Done():
			fmt.Println("work cancelled")
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go work(ctx)
	// This example uses context.WithTimeout to automatically cancel the context after 2 seconds.
	// The work function checks ctx.Done() (like a done channel) to know when to stop.
	// Cox-Buday’s chapter emphasizes that context can be used “to add timeouts and cancellations” to patterns.
	// In short, use context to manage lifecycles instead of custom channels when possible: it’s idiomatic, integrates with many libraries, and makes cancellation explicit.
	time.Sleep(2 * time.Second)
}
