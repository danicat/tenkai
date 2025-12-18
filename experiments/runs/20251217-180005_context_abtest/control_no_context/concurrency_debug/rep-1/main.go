package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Val   int
}

// ProcessJobs processes n jobs using m workers.
// Use this function signature for testing.
// It returns the number of processed jobs and an error if something goes wrong.
func ProcessJobs(numJobs, numWorkers int) (int, error) {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100) // Small buffer creates pressure


	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Send result to the results channel
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	var resultsWg sync.WaitGroup
	var finalProcessedCount int64 // Use a separate counter for collected results

	// Start a goroutine to collect results
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			atomic.AddInt64(&finalProcessedCount, 1)
		}
	}()

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Wait for all workers to finish
	wg.Wait()
	// All workers are done, no more results will be sent
	close(results)
	// Wait for the results collector to finish
	resultsWg.Wait()

	return int(finalProcessedCount), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
