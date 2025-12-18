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
	var processedCount int32            // RACE CONDITION: Not thread-safe

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 1: Race Condition
				atomic.AddInt32(&processedCount, 1)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Goroutine to send jobs and close the jobs channel
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Goroutine to collect results concurrently
	// This goroutine must also wait for the results channel to be closed
	var resultsCollectorWg sync.WaitGroup // Use a separate WaitGroup for the results collector
	resultsCollectorWg.Add(1)
	go func() {
		defer resultsCollectorWg.Done()
		for range results {
			// Just drain the results
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	// All workers are done, no more results will be sent.
	// Now it's safe to close the results channel.
	close(results)

	// Wait for the results collector to finish draining
	resultsCollectorWg.Wait()

	return int(processedCount), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
