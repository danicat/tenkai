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

	var processedCount atomic.Int32 // FIX 1: Use atomic for thread-safe counter

	var wg sync.WaitGroup   // For workers
	var rwg sync.WaitGroup // For results consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// FIX 1: Atomic increment
				processedCount.Add(1)

				// Send result. Now the results channel will be consumed concurrently.
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Start results consumer
	rwg.Add(1)
	go func() {
		defer rwg.Done()
		for range results {
			// Just drain the results channel
		}
	}()


	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Wait for all workers to finish processing jobs and sending results
	wg.Wait()
	// All workers are done sending results, so it's safe to close the results channel
	close(results)

	// Wait for the results consumer to finish draining the channel
	rwg.Wait()

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
