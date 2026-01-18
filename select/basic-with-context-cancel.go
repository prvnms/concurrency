package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- "from ch2"
	}()
	for {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		case <-ctx.Done():
			fmt.Println("ctx timeout ", ctx.Err())
			return
		}
	}

}
