package main

// Fan-In: Aggregate data from multiple sources
import (
	"fmt"
	"sync"
	"time"
)

type Log struct {
	Source  string
	Message string
	Time    time.Time
}

func logGenerator(source string, interval time.Duration) <-chan Log {
	logs := make(chan Log)

	go func() {
		defer close(logs)
		for i := 1; i <= 5; i++ {
			log := Log{
				Source:  source,
				Message: fmt.Sprintf("Event %d from %s", i, source),
				Time:    time.Now(),
			}
			logs <- log
			time.Sleep(interval)
		}
	}()

	return logs
}

func fanIn(channels ...<-chan Log) <-chan Log {
	merged := make(chan Log)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Log) {
			defer wg.Done()
			for log := range c {
				merged <- log
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {

	// producers
	webSvcLogs := logGenerator("WebServer", 300*time.Millisecond)
	dbLogs := logGenerator("Database", 500*time.Millisecond)
	authSvcLogs := logGenerator("AuthService", 400*time.Millisecond)

	logs := fanIn(webSvcLogs, dbLogs, authSvcLogs)

	fmt.Println("logs from sources")
	for log := range logs {
		fmt.Printf("[%s] %s - %s\n",
			log.Time.Format("15:04:05.000"),
			log.Source,
			log.Message)
	}

	fmt.Println("logs collected")
}
