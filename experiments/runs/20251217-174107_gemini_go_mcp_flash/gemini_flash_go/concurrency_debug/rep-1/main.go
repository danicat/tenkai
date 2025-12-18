package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import atomic
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
	jobs := make(chan Job, numJobs)           // Buffered for all jobs to avoid blocking the sender
	results := make(chan Result, numWorkers) // Buffered for workers to send results without immediate blocking

	var processedCount int64 // Use int64 for atomic operations

	var wgWorkers sync.WaitGroup // WaitGroup for workers
	var wgResults sync.WaitGroup // WaitGroup for results consumer

	// Start result collector goroutine
	wgResults.Add(1)
	go func() {
		defer wgResults.Done()
		for range results {
			// This goroutine drains the results channel, preventing workers from blocking
			// The actual processedCount is updated by workers via atomic operations.
		}
	}()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wgWorkers.Add(1)
		go func(id int) {
			defer wgWorkers.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Fix 1: Race Condition using atomic operation
				atomic.AddInt64(&processedCount, 1)

				// Send results to channel. This won't block indefinitely
				// because the result collector is draining it concurrently.
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Send jobs to the jobs channel
	for i := 0; i < numJobs; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs) // All jobs sent, close the jobs channel

	// Wait for all workers to finish processing jobs
	wgWorkers.Wait()

	// All workers are done, so no more results will be sent.
	// It's now safe to close the results channel.
	close(results)

	// Wait for the result collector to finish draining the results channel
	wgResults.Wait()

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
