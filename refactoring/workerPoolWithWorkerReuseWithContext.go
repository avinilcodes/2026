package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d stopping due to context cancellation\n", id)
			return
		default:
			fmt.Printf("worker %d started job %d\n", id, j)
			time.Sleep(time.Second * 2)
			results <- j * 2
			fmt.Printf("worker %d finished job %d\n", id, j)
		}
	}
}

// workerPoolWithWorkerReuseWithContext demonstrates a worker pool with worker reuse and context cancellation
func main() {
	const numWorkers = 3
	const numJobs = 5

	var jobs = make(chan int, numJobs)
	var results = make(chan int, numJobs)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println("Result:", result)
	}
}
