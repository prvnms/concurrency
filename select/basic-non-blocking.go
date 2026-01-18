package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)

	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No message available")
	}

	go func() {
		messages <- "Hello!"
	}()

	time.Sleep(100 * time.Millisecond)
	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No message")
	}
}
