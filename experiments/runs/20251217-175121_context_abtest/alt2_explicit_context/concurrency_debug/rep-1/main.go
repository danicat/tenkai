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
	var processedCount atomic.Int64 // Use atomic.Int64 for thread-safe counter

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

				// Fix 1: Race Condition with atomic counter
				processedCount.Add(1)

				// Fix 2: Send results
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

	// Start a goroutine to consume results
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for range results {
			// Just drain the results channel.
			// The actual counting is done by processedCount in workers.
		}
	}()

	// Wait for all workers to finish processing jobs
	workerWg.Wait()
	// All workers are done, so no more results will be sent.
	// It's safe to close the results channel.
	close(results)

	// Wait for the result consumer to finish draining the channel
	resultWg.Wait()

	return int(processedCount.Load()), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
