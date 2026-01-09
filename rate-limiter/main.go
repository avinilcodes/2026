package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func worker(ctx context.Context, limiter *rate.Limiter, jobs <-chan int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d context done, exiting\n", id)
			return
		case job, ok := <-jobs:
			time.Sleep(time.Second * 2)
			if !ok {
				return
			}
			if err := limiter.Wait(ctx); err != nil {
				return
			}

			process(id, job)
		}
	}
}

func process(id, job int) {
	// Simulate job processing
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Worker %d processed job: %d\n", id, job)
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var jobs = make(chan int, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		jobs <- i
	}

	limiter := rate.NewLimiter(5, 10)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		fmt.Printf("Starting worker: %d\n", i)
		go worker(ctx, limiter, jobs, &wg, i)
	}

	defer wg.Wait()

	close(jobs)
	fmt.Printf("All jobs submitted, waiting for workers to finish\n")
}
