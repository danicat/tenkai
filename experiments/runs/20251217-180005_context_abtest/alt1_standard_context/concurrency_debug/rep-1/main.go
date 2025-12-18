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
	var processedCount int64          // RACE CONDITION: Not thread-safe, using int64 for atomic operations

	var wg sync.WaitGroup
	var resultsWg sync.WaitGroup // New WaitGroup for results consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Race Condition Fix: Use atomic operations
				atomic.AddInt64(&processedCount, 1)

				// Send result
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
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			// Just drain the results channel to prevent workers from blocking
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(results) // Close the results channel ONLY after all workers are done sending
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
