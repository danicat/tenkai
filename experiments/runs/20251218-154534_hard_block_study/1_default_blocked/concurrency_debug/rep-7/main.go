package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Val   int
}

type JobCounter struct {
	mu    sync.Mutex
	count int
}

func (jc *JobCounter) Inc() {
	jc.mu.Lock()
	defer jc.mu.Unlock()
	jc.count++
}

func (jc *JobCounter) Count() int {
	jc.mu.Lock()
	defer jc.mu.Unlock()
	return jc.count
}

// ProcessJobs processes n jobs using m workers.
// Use this function signature for testing.
// It returns the number of processed jobs and an error if something goes wrong.
func ProcessJobs(numJobs, numWorkers int) (int, error) {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100) // Small buffer creates pressure
	processedCounter := &JobCounter{}

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}
	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			processedCounter.Inc()
		}
	}()

	wg.Wait()
	close(results)

	resultsWg.Wait()

	return processedCounter.Count(), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
