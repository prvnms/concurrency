package main

import (
	"fmt"
	"time"
)

func main() {
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

	// for { // will cause dead lock
	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
	// }
}
