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
	var processedCount int32          // RACE CONDITION: Not thread-safe

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Fix: Race Condition resolved with atomic operations
				atomic.AddInt32(&processedCount, 1)

				// Send result. This might block if the results channel is full
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
			// Drain results. processedCount is already incremented by workers.
			// We just need to ensure the channel is drained to prevent deadlocks.
		}
	}()

	wg.Wait()        // Wait for all workers to finish processing jobs
	close(results)   // Close the results channel after all workers are done sending
	resultsWg.Wait() // Wait for the results consumer to finish draining

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
