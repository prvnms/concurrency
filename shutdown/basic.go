package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	/*
		This tells the Go runtime: "If the operating system sends a SIGINT (Ctrl+C) or a SIGTERM (kill command), don't crash the program. Instead, put that signal into my quit channel
		os.Interrupt	SIGINT	Ctrl + C
		syscall.SIGTERM	SIGTERM	kill command
		syscall.SIGKILL	SIGKILL	kill -9
	*/

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("shuting down")
				time.Sleep(500 * time.Millisecond)
				fmt.Println("cleaned up")
				done <- true
				return
			default:
				fmt.Println("Working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		quit <- syscall.SIGINT
	}()

	<-done
	fmt.Println("shutdown gracefully")
}
