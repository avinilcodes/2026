package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// producerConsumerWithContext demonstrates a producer-consumer pattern using channels and context cancellation
func main() {
	fmt.Println("Hello, Producer-Consumer!")
	var jobs = make(chan int, 10)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go consumer(ctx, i, jobs, &wg)
	}

	go producer(ctx, jobs)

	wg.Wait()
}

func producer(ctx context.Context, ch chan<- int) {
	// Produce items and send them to a channel
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Producer stopping due to context cancellation")
			return
		default:
			time.Sleep(100 * time.Millisecond)
			ch <- i
		}
	}
	close(ch)
}

func consumer(ctx context.Context, id int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer stopping due to context cancellation")
			return
		default:
			fmt.Println("Consumer: ", id, "job consumed: ", job)
			time.Sleep(200 * time.Millisecond)

		}
	}
}
