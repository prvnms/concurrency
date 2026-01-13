package main

import (
	"concurrency/context-cancellation/dto"
	"concurrency/context-cancellation/processor"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func timeoutCancel() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tasks := make(chan dto.Task, 10)
	results := make(chan string, 10)

	for i := 1; i <= 2; i++ {
		go processor.CtxWorker(ctx, i, tasks, results)
	}

	go func() {
		for i := 1; i <= 6; i++ {
			tasks <- dto.Task{ID: i, Name: fmt.Sprintf("Task-%d", i)}
			time.Sleep(100 * time.Millisecond)
		}
		close(tasks)
	}()

	<-ctx.Done()
	fmt.Printf("Timeout reached: %v\n", ctx.Err())

	time.Sleep(500 * time.Millisecond)
	close(results)

	fmt.Println("Results")
	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("timeoutCancel started")
	timeoutCancel()
	fmt.Println("timeoutCancel completed")
}
