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
	results := make(chan Result, 100)
	var processedCount atomic.Int32 // FIX 1: Use atomic.Int32 for thread-safe counter

	var workerWg sync.WaitGroup // WaitGroup for workers
	var resultWg sync.WaitGroup // WaitGroup for result consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func(id int) {
			defer workerWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				processedCount.Add(1) // FIX 1: Atomically increment

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
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for range results {
			// Just drain the results channel.
			// The processedCount is already handled by workers.
		}
	}()

	// Wait for all workers to finish processing jobs and sending results
	workerWg.Wait()
	close(results) // Close the results channel AFTER all workers have finished sending

	// Wait for the result consumer to finish draining the results channel
	resultWg.Wait()

	return int(processedCount.Load()), nil // FIX 1: Load atomic value
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
