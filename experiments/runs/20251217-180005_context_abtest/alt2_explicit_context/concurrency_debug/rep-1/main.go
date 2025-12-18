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
	results := make(chan Result, numJobs) // Buffer for all results to prevent deadlock.
	var processedCount atomic.Int64       // Thread-safe counter using sync/atomic

	var workersWg sync.WaitGroup
	var resultsWg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func(id int) {
			defer workersWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				processedCount.Add(1)

				// Workers send results. With numJobs buffer, this won't block.
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Start a goroutine to send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Start a goroutine to collect results
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			// Just drain the results
		}
	}()

	// Wait for all workers to finish processing jobs
	workersWg.Wait()
	// Close the results channel ONLY after all workers have finished
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
