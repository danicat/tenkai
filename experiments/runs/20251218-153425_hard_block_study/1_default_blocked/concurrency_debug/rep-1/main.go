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
	var processedCount int32 // RACE CONDITION: Now thread-safe using atomic operations

	var wg sync.WaitGroup
	var resultWg sync.WaitGroup // New WaitGroup for result collector

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// FIX 1: Race Condition resolved using atomic operations
				atomic.AddInt32(&processedCount, 1)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
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

	// Start result collector
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for range results {
			// Just drain
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	// Once all workers are done, close the results channel
	// This signals the resultCollector to stop after draining remaining results
	close(results)

	// Wait for the result collector to finish draining and closing the channel
	resultWg.Wait()

	return int(atomic.LoadInt32(&processedCount)), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
