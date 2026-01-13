package processor

import (
	"concurrency/context-cancellation/dto"
	"context"
	"fmt"
	"time"
)

func CtxWorker(ctx context.Context, id int, tasks <-chan dto.Task, results chan<- string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: received cancel signal (%v), shutting down...\n",
				id, ctx.Err())
			return

		case task, ok := <-tasks:
			if !ok {
				fmt.Printf("worker %d: channel closed, exiting\n", id)
				return
			}

			if err := processTask(ctx, id, task); err != nil {
				results <- fmt.Sprintf("worker %d: task %d CANCELLED", id, task.ID)
			} else {
				results <- fmt.Sprintf("woker %d: task %d completed", id, task.ID)
			}
		}
	}
}

func processTask(ctx context.Context, workerID int, task dto.Task) error {
	fmt.Printf("worker %d :: task %d (%s)\n", workerID, task.ID, task.Name)

	processingTime := 5
	for i := 0; i < processingTime; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: task %d interrupted at step %d/%d\n",
				workerID, task.ID, i+1, processingTime)
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			fmt.Printf("worker %d: task %d progress: %d/%d\n",
				workerID, task.ID, i+1, processingTime)
		}
	}

	return nil
}
