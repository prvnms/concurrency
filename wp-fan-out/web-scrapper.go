package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ScraperResult struct {
	URL       string
	Status    int
	Bytes     int64
	WordCount int
	Error     string
	Duration  time.Duration
	Retries   int
}

var urls = []string{
	"https://httpbin.org/html",
	"https://example.com",
	"https://books.toscrape.com",
	"https://quotes.toscrape.com",
	"https://news.ycombinator.com",
	"https://www.iana.org/domains/reserved",
	"https://go.dev/doc/",
	"https://www.scrapethissite.com/pages/simple/",
	"https://www.scrapethissite.com/pages/ajax-javascript/",
	"https://www.basketball-reference.com/",
	"https://coinmarketcap.com/",
	"https://www.worldometers.info/world-population/",
	"https://www.gov.uk/search/all",
	"https://quotes.toscrape.com/page/1/",
	"https://www.w3schools.com/html/",
}

const (
	WORKERNUM      = 4
	maxRetries     = 2
	userAgent      = "Mozilla/5.0 (compatible; RealWebScrapper/1.0)"
	requestTimeout = 10 * time.Second
)

func doWork(dataChan <-chan string, resultChan chan<- ScraperResult, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: requestTimeout,
	}

	for url := range dataChan {
		res := scrape(client, url)
		resultChan <- res
	}
}

func scrape(client *http.Client, url string) ScraperResult {
	start := time.Now()
	result := ScraperResult{
		URL: url,
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result.Retries = attempt

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			result.Error = err.Error()
			break
		}
		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(jitterBackoff(attempt))
				continue
			}
			result.Error = err.Error()
			break
		}

		func() {
			defer resp.Body.Close()

			result.Status = resp.StatusCode

			if resp.StatusCode != http.StatusOK {
				if attempt < maxRetries && shouldRetry(resp.StatusCode) {
					time.Sleep(jitterBackoff(attempt))
					return
				}
				result.Error = fmt.Sprintf("status %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				if attempt < maxRetries {
					time.Sleep(jitterBackoff(attempt))
					return
				}
				result.Error = err.Error()
				return
			}

			result.Bytes = int64(len(body))
			result.WordCount = countWords(string(body))
			result.Duration = time.Since(start)
		}()

		if result.Duration > 0 {
			return result
		}
	}

	if result.Duration == 0 {
		result.Duration = time.Since(start)
	}
	return result
}

func main() {
	fmt.Println("Starting web scraper...")

	workerChan := make(chan string, len(urls))
	resultChan := make(chan ScraperResult)
	var wg sync.WaitGroup

	for i := 1; i <= WORKERNUM; i++ {
		wg.Add(1)
		go doWork(workerChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	go func() {
		for res := range resultChan {
			if res.Error != "" {
				fmt.Printf("[FAIL] URL=%s Status=%d Retries=%d Duration=%v Error=%s\n",
					res.URL, res.Status, res.Retries, res.Duration, res.Error)
			} else {
				fmt.Printf("[OK]   URL=%s Status=%d Words=%d Bytes=%d Retries=%d Duration=%v\n",
					res.URL, res.Status, res.WordCount, res.Bytes, res.Retries, res.Duration)
			}
		}
	}()

	for _, url := range urls {
		workerChan <- url
	}
	close(workerChan)
	time.Sleep(2 * time.Second)
}

func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func shouldRetry(status int) bool {
	switch status {
	case 429, 503, 504, 502, 408:
		return true
	default:
		return status >= 500
	}
}

func jitterBackoff(attempt int) time.Duration {
	base := 300 * time.Millisecond
	maxJitter := 400 * time.Millisecond
	backoff := base * time.Duration(1<<attempt) // exponential
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return backoff + jitter
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
