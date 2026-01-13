package main

import (
	"concurrency/context-cancellation/dto"
	"concurrency/context-cancellation/processor"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func hierarchyCancel() {

	// Parent context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Child contexts
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	defer childCancel1()

	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel2()

	// Worker group 1
	tasks1 := make(chan dto.Task, 5)
	go processor.CtxWorker(childCtx1, 1, tasks1, make(chan string, 5))

	// Worker group 2
	tasks2 := make(chan dto.Task, 5)
	go processor.CtxWorker(childCtx2, 2, tasks2, make(chan string, 5))

	// Send tasks
	go func() {
		for i := 1; i <= 3; i++ {
			tasks1 <- dto.Task{ID: i, Name: "Group1-Task"}
			tasks2 <- dto.Task{ID: i, Name: "Group2-Task"}
		}
	}()

	time.Sleep(800 * time.Millisecond)
	fmt.Println("\n🛑 Cancelling child context 1...")
	childCancel1()

	time.Sleep(800 * time.Millisecond)
	fmt.Println("\n🛑 Cancelling PARENT context (affects all children)...")
	parentCancel()

	time.Sleep(500 * time.Millisecond)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("hierarchyCancel started")
	hierarchyCancel()
	fmt.Println("hierarchyCancel completed")
}
