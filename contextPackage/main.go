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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	cancel()
	go work(ctx)

	time.Sleep(3 * time.Second)
}
