package main

import "fmt"

func main() {
	jobs := make(chan int, 5)
	done := make(chan struct{})

	go func() {
		for i := 0; i <= 5; i++ {
			jobs <- i
		}
		done <- struct{}{}
	}()
loop:
	for {
		select {
		case job := <-jobs:
			fmt.Println("job ", job)
		case <-done:
			fmt.Println("done invoked")
			break loop
		}
	}

}
