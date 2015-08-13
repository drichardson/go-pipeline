// from http://blog.golang.org/pipelines

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("calling runTest")
	runTest()
	fmt.Println("runTest returned")

	fmt.Println("Sleeping 1 seconds to see if goroutines cleanup...")
	time.Sleep(1 * time.Second)
	fmt.Println("Done sleeping")
}

func runTest() {
	done := make(chan struct{}, 2)
	defer close(done)

	in := gen(done, 1, 2, 3, 4, 5, 6, 7, 9)

	// Distribute the sq work across two goroutines that both read from in
	c1 := sq(done, in)
	c2 := sq(done, in)

	// Consume the first value from input
	out := merge(done, c1, c2)
	fmt.Println(<-out)

	// done will be closed by the deferred call
}

func gen(done chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range nums {
			select {
			case out <- v:
			case <-done:
				fmt.Println("gen got done signal")
				return
			}
		}
	}()
	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				fmt.Println("sq got done signal")
				return
			}
		}
	}()
	return out
}

func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start and output goroutine for each input channel in cs. output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				fmt.Println("merge output got done signal")
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
