package main

import (
	"fmt"
	"time"
)

func greet(phrase string, doneChan chan bool) {
	fmt.Println("Hello!", phrase)
	doneChan <- true
}

func slowGreet(phrase string, doneChan chan bool) {
	time.Sleep(3 * time.Second) // simulate a slow, long-taking task
	fmt.Println("Hello!", phrase)
	doneChan <- true
	close(doneChan)
}

func main() {
	// dones := make([]chan bool, 4)

	// for i := 0; i < 4; i++ {
	// 	dones[i] = make(chan bool)
	// }

	done := make(chan bool)

	go greet("Nice to meet you!", done)
	go greet("How are you?", done)
	go slowGreet("How ... are ... you ...?", done)
	go greet("I hope you're liking the course!", done)

	// for _, done := range dones {
	// 	<-done
	// }

	for doneChan := range done {
		fmt.Println(doneChan)
	}

	fmt.Println("complete")
}
