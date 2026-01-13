package main

import (
	"concurrency/context-cancellation/dto"
	"concurrency/context-cancellation/processor"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func deadlineCancel() {

	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Deadline set for: %s (in 1.5 seconds)\n", deadline.Format("15:04:05"))

	tasks := make(chan dto.Task, 10)
	results := make(chan string, 10)

	for i := 1; i <= 2; i++ {
		go processor.CtxWorker(ctx, i, tasks, results)
	}

	go func() {
		for i := 1; i <= 4; i++ {
			tasks <- dto.Task{ID: i, Name: fmt.Sprintf("Task-%d", i)}
		}
		close(tasks)
	}()

	<-ctx.Done()
	fmt.Printf("Deadline exceeded: %v\n", ctx.Err())

	time.Sleep(500 * time.Millisecond)
	close(results)

	fmt.Println("Results")
	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("deadlineCancel started")
	deadlineCancel()
	fmt.Println("deadlineCancel completed")
}
