package concurrency_patterns

import (
	"fmt"
)

// Pipeline Pattern
// produce -> transform -> consume
// A pipeline breaks the work into stages. Each stage runs in its own goroutine, connected to the next by a channel.
// ----------------------------------------------------------------------------------------------------------------------

// Stage 1 -> produce
func produce() <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := range 5 {
			out <- i
		}
	}()

	return out
}

// Stage 2 -> transform
func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for v := range in {
			out <- v * v
		}
	}()

	return out
}

// Stage 3 -> consume
func PipelinePatternExample() {
	nums := produce()
	squares := square(nums)

	for v := range squares {
		fmt.Println("PipelinePatternExample: ", v)
	}
}

// ----------------------------------------------------------------------------------------------------------------------
