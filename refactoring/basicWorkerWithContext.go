package main

import (
	"context"
	"fmt"
	"time"
)

// basicWorker to be refactored for using context package
func main() {
	var jobs = make(chan int, 10)
	var results = make(chan int, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for w := 1; w <= 3; w++ {
		go worker(ctx, w, jobs, results)
	}

	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 10; a++ {
		fmt.Println(<-results)
	}
	close(results)
}

func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
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
