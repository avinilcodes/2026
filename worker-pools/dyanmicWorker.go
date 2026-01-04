package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var jobs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i, job := range jobs {
		wg.Add(1)

		go worker(i+1, job, &wg)
	}
	wg.Wait()
	fmt.Println("All jobs completed")
}

func worker(id, job int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d started job %d\n", id, job)

	time.Sleep(time.Second * 2)
	fmt.Printf("Worker %d finished job %d\n", id, job)

}
