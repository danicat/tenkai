package main

import (
	"fmt"
	"sync" // Added for atomic operations
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
func ProcessJobs(numJobs, numWorkers int) (int32, error) {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100) // Small buffer creates pressure
	var processedCount int32          // RACE CONDITION: Now atomic

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 1: Race Condition fixed
				atomic.AddInt32(&processedCount, 1)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Consume results concurrently
	var resultsCollectorWg sync.WaitGroup
	resultsCollectorWg.Add(1)
	go func() {
		defer resultsCollectorWg.Done()
		for range results {
			// Just drain the results to prevent deadlock
		}
	}()

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	wg.Wait()      // Wait for all workers to finish
	close(results) // Close results channel after all workers are done and no more results will be sent

	resultsCollectorWg.Wait() // Wait for the results collector to finish draining

	return atomic.LoadInt32(&processedCount), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
