// from http://blog.golang.org/pipelines

package main

import (
	"fmt"
)

func main() {
	out := sq(gen(1, 2, 3, 4, 5, 6, 7, 9))
	for v := range out {
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
