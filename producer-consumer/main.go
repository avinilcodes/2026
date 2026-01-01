package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello, Producer-Consumer!")
	var jobs = make(chan int, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go consumer(i, jobs, &wg)
	}

	go producer(jobs)

	wg.Wait()
}

func producer(ch chan<- int) {
	// Produce items and send them to a channel
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		ch <- i
	}
	close(ch)
}

func consumer(id int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println("Consumer: ", id, "job consumed: ", job)
		time.Sleep(200 * time.Millisecond)
	}
}
