package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Id   int
	Name string
}

const numWorkers int = 10

func worker(jobs chan Job, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("done called")
			return
		case job, ok := <-jobs:
			if !ok {
				fmt.Println("No jobs")
				return
			}
			fmt.Println("job is ", job)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	jobs := make(chan Job, 10)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for range numWorkers {
		wg.Add(1)
		go worker(jobs, &wg, ctx)
	}

	go func() {
		for i := 0; i < 10; i++ {
			jobs <- Job{Id: i, Name: fmt.Sprintf("Data %d", i)}
			time.Sleep(100 * time.Millisecond)
		}
		close(jobs)
	}()

	time.Sleep(1500 * time.Millisecond)
	fmt.Println("shutdown started")
	cancel()
	wg.Wait()
	fmt.Println("shutdown done")

}
