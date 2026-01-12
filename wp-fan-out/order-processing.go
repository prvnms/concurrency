package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fan-Out: Distribute workload to controlled number of workers

const (
	WORKERS = 3
	ORDERS  = 10
)

type Order struct {
	Id     int
	Amount float32
}
type OrderResult struct {
	OrderId int
	Success bool
	Message string
}

func worker(id int, jobs <-chan Order, results chan<- OrderResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for order := range jobs {
		fmt.Printf("Worker %d: Processing order #%d\n", id, order.Id)
		res := processOrder(order)
		results <- res

	}
}

func processOrder(order Order) OrderResult {
	pTime := time.Duration(300+rand.Intn(500)) * time.Millisecond
	time.Sleep(pTime)

	success := rand.Float32() > 0.1
	result := OrderResult{
		OrderId: order.Id,
		Success: success,
	}
	if success {
		result.Message = fmt.Sprintf("Order #%d processed successfully", order.Id)
	} else {
		result.Message = fmt.Sprintf("Order #%d failed: Payment declined", order.Id)
	}

	return result
}

func main() {

	jobs := make(chan Order, ORDERS)
	results := make(chan OrderResult, ORDERS)

	var wg sync.WaitGroup

	fmt.Println("Init workers")
	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		for i := 1; i <= ORDERS; i++ {
			order := Order{
				Id:     i,
				Amount: float32(100 + rand.Float64()*900),
			}
			jobs <- order
			fmt.Printf("Produced order %d\n", i)
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result.Message)
		if !result.Success {
			fmt.Println("failed")
		}
	}

}
