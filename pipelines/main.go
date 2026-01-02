package main

import (
	"fmt"
	"sync"
)

func main() {
	// Example usage of the pipeline functions
	c := gen(2, 3)
	s := sq(c)
	fmt.Println(<-s)
	fmt.Println(<-s)

	// Or more concisely:
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n)
	}

	// Merging multiple channels
	c1 := gen(2, 4)
	c2 := gen(3, 5)
	m := merge(sq(c1), sq(c2))
	for n := range m {
		fmt.Println(n)
	}
	done := make(chan struct{}, 2)
	// Consume the first value from the output.
	m = mergeDone(done, sq(c1), sq(c2))
	fmt.Println("Merged output:", <-m)
	// Since we didn't receive the second value from out,
	// one of the output goroutines is hung attempting to send it.
	done <- struct{}{}
	done <- struct{}{}
}

// sq squares numbers received from the input channel and sends them to the output channel.
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// gen generates a channel that emits the provided integers.
func gen(in ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range in {
			out <- n
		}
		close(out)
	}()
	return out
}

// merge combines multiple input channels into a single output channel.
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// merge combines multiple input channels into a single output channel.
func mergeDone(done <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			select {
			case out <- n:
			case <-done:
			}
			wg.Done()
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
