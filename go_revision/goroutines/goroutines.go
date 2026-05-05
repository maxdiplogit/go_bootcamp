package goroutines

import (
	"fmt"
	"sync"
	"time"
)

func sayHello(name string) {
	fmt.Println("Hello ", name)
}

func ExecuteExample() {
	go sayHello("Harshmeet")
	go sayHello("Gautham")
	go sayHello("Dil")

	time.Sleep(1 * time.Second)
}

func ChannelExample() (val int) {
	ch := make(chan int)

	go func() {
		ch <- 10
		// If we dont close the channel and try to receive a value from it again without sending in this unbuffered chan, the whole thing will panic
		// The convention: only the sender closes the channel, never the receiver.
		// If multiple goroutines send on the same channel, none of them should close it (because another might still be sending).
		// Closing is for signaling "no more values are coming."
		close(ch)
	}()

	val = <-ch
	fmt.Println("Val after already extracting from channel: ", <-ch)
	return
}

func BufferedChannelExample() {
	ch := make(chan int, 3)

	ch <- 1
	fmt.Println("Inserted 1 into ch...")
	ch <- 2
	fmt.Println("Inserted 2 into ch...")
	ch <- 3
	fmt.Println("Inserted 3 into ch...")
	// ch <- 4 // BLOCKS -- buffer is full, waits for someone to receive
	// fmt.Println("Inserted 4 into ch...")
	// ch <- 5 // BLOCKS -- buffer is full, waits for someone to receive
	// fmt.Println("Inserted 5 into ch...")
}

func RangeOverChannelExample() {
	ch := make(chan int, 3)

	go func() {
		for i := 1; i < 5; i++ {
			ch <- i
		}

		close(ch)
	}()

	x := <-ch

	fmt.Println("x: ", x)

	for v := range ch {
		fmt.Println("value from channel: ", v)
	}
}

func SelectChannelExample() string {
	ch_1 := make(chan string)
	ch_2 := make(chan string)

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch_1 <- "from ch_1"
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch_2 <- "from ch_2"
	}()

	select {
	case msg := <-ch_1:
		return msg
	case msg := <-ch_2:
		return msg
	}
}

func TimeoutSelectChannelExample() string {
	ch := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- "got value"
	}()

	select {
	case v := <-ch:
		return v
	case <-time.After(50 * time.Millisecond):
		return "timed out"
	}
}

func DefaultSelectChannelExample() string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		ch <- "got value"
	}()

	// time.Sleep(time.Second * 3)

	select {
	case v := <-ch:
		return v
	default:
		return "nothing ready, did not block"
	}
}

func WaitgroupExample() {
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Println("worker ", id, " running")
		}(i)
	}

	wg.Wait()
	fmt.Println("WaitgroupExample completed!")
}

func MutexExample() (counter int) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := range 5 {
		fmt.Print("Mutex example: ", i, "\n")
		wg.Add(1)
		go func() {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	return
}

func RWMutexExample() {
	// Multiple readers but only one writer
}

// This worker only takes in the tasks channel as a receiving end of values, not sending end (unidirectional and not bi-directional)
func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf("Worker %d processing job %d\n", id, task)
	}
}

func WorkerPoolExample() {
	// We have a fixed number of workers that have to pick up tasks from a queue
	// So, fixed number of workers -> fixed number of goroutines
	// Tasks put in a queue -> Tasks put in a Buffered Channel

	const MAX_CONCURRENT_WORKERS int = 5
	const MAX_TASKS_IN_CHANNEL int = 10

	tasks := make(chan int, MAX_TASKS_IN_CHANNEL)
	var wg sync.WaitGroup

	// Step 1: Start Workers
	for i := range MAX_CONCURRENT_WORKERS {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	// Step 2: Send tasks
	for i := range 5 {
		tasks <- i
	}

	close(tasks)

	wg.Wait()

	fmt.Println("All tasks completed!")
}
