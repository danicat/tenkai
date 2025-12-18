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
	var processedCount int32          // RACE CONDITION: Not thread-safe, using atomic for fix

	var wg sync.WaitGroup // For workers
	var sendWg sync.WaitGroup // For job sender

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Fix 1: Race Condition using atomic
				atomic.AddInt32(&processedCount, 1)

				// Fix 2: Send results to channel
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}


	// Send jobs
	sendWg.Add(1) // Add to the new WaitGroup
	go func() {
		defer sendWg.Done() // Mark as done when this goroutine exits
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Start a goroutine to consume results concurrently
	// This prevents workers from blocking when sending to a full results channel
	go func() {
		for range results {
			// Just drain the results. processedCount is already handled atomically by workers.
		}
	}()

	// Wait for all jobs to be sent and all workers to finish processing jobs
	sendWg.Wait() // Wait for job sender to finish and close 'jobs' channel
	wg.Wait()     // Wait for all workers to finish processing jobs

	// Close the results channel after all workers have completed and sent their results
	close(results)

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
