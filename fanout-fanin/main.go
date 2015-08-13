// from http://blog.golang.org/pipelines

package main

import (
	"fmt"
	"sync"
)

func main() {
	in := gen(1, 2, 3, 4, 5, 6, 7, 9)

	// Distribute the sq work across two goroutines that both read from in
	c1 := sq(in)
	c2 := sq(in)

	// Consume the merged output from c1 and c2
	for v := range merge(c1, c2) {
		fmt.Println(v)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range nums {
			out <- v
		}
		close(out)
	}()
	return out
}

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

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start and output goroutine for each input channel in cs. output
	// copies values from c to out until c is closed, then calls wg.Done.
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

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
