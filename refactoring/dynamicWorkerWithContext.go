package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// dynamicWorkerWithContext demonstrates a dynamic worker pool with context cancellation
func main() {
	var wg sync.WaitGroup
	var jobs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, job := range jobs {
		wg.Add(1)

		go worker(ctx, i+1, job, &wg)
	}
	wg.Wait()
	fmt.Println("All jobs completed")
}

func worker(ctx context.Context, id, job int, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		fmt.Printf("Worker %d stopping due to context cancellation\n", id)
		return
	default:
		fmt.Printf("Worker %d started job %d\n", id, job)

		time.Sleep(time.Second * 2)
		fmt.Printf("Worker %d finished job %d\n", id, job)
	}

}
