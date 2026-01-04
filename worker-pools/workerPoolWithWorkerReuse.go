package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d started job %d\n", id, j)
		time.Sleep(time.Second * 2)
		results <- j * 2
		fmt.Printf("worker %d finished job %d\n", id, j)
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 5

	var jobs = make(chan int, numJobs)
	var results = make(chan int, numJobs)
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
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
