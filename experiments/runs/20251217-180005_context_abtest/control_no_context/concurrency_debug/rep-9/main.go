package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import sync/atomic
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
	var processedCount int32          // Fix: Use int32 for atomic operations

	var wgWorkers sync.WaitGroup // WaitGroup for workers
	var wgResults sync.WaitGroup // WaitGroup for result consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wgWorkers.Add(1)
		go func(id int) {
			defer wgWorkers.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				atomic.AddInt32(&processedCount, 1) // Fix: Atomic increment
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

	// Start a goroutine to consume results concurrently
	wgResults.Add(1)
	go func() {
		defer wgResults.Done()
		for range results {
			// Just drain the results. processedCount is already incremented by workers.
		}
	}()

	wgWorkers.Wait() // Wait for all workers to finish
	close(results)   // Close the results channel after all workers are done producing
	wgResults.Wait() // Wait for the result consumer to finish draining

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
