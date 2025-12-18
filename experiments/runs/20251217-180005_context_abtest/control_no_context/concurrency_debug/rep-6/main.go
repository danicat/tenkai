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
	var atomicProcessedCount atomic.Int64 // FIX 1: Use atomic for thread-safe counter

	var wg sync.WaitGroup
	var resultWg sync.WaitGroup // New WaitGroup for results collector

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 1: Race Condition
				atomicProcessedCount.Add(1)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Start results collector goroutine
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for range results {
			// Just drain the channel. atomicProcessedCount is updated by workers.
		}
	}()

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs) // Close jobs channel after all jobs are sent
	}()

	// Wait for all workers to finish processing jobs and sending results
	wg.Wait()
	close(results) // FIX 2: Close results channel ONLY after all workers are done
	// This signals the results collector goroutine to stop.

	// Wait for the results collector to finish draining all results
	resultWg.Wait() // FIX 2: Ensure all results are consumed

	return int(atomicProcessedCount.Load()), nil // Return the atomically loaded value
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
