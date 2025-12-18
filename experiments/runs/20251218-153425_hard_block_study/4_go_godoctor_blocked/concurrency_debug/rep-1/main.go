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
	var processedCount atomic.Int64   // Fixed: Now thread-safe

	var wg sync.WaitGroup        // For workers and job sender
	var resultsWg sync.WaitGroup // For results collector

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Fixed: Race Condition using atomic
				processedCount.Add(1)

				// Send result. This might block if results channel is full,
				// but the collector goroutine will drain it.
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Start job sender
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs) // Close jobs channel after all jobs are sent
	}()

	// Start results collector
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		// Drain results channel. This goroutine ensures workers don't block indefinitely.
		for range results {
			// We are counting processed jobs directly in the worker goroutines
			// using atomic operations, so we just need to drain the channel here.
		}
	}()

	// Wait for all jobs to be sent and all workers to finish processing
	wg.Wait()

	// Close the results channel AFTER all workers are done
	// This signals the results collector to stop
	close(results)

	// Wait for the results collector to finish draining the channel
	resultsWg.Wait()

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
