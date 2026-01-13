package main

import (
	"concurrency/context-cancellation/dto"
	"concurrency/context-cancellation/processor"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func manualCancel() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan dto.Task, 10)
	results := make(chan string, 10)

	for i := 1; i <= 2; i++ {
		go processor.CtxWorker(ctx, i, tasks, results)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			tasks <- dto.Task{ID: i, Name: fmt.Sprintf("Task-%d", i)}
		}
		close(tasks)
	}()

	go func() {
		time.Sleep(1500 * time.Millisecond)
		fmt.Println("CANCELLING ALL WORKERS...")
		cancel()
	}()

	time.Sleep(4 * time.Second)
	close(results)

	fmt.Println("Results")
	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("manualCancel started")
	manualCancel()
	fmt.Println("manualCancel completed")
}
