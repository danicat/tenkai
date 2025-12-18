package main

import (
	"fmt"
	"sync"
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
	var processedCount int // This will be incremented safely by the result collector
	
	var workerWg sync.WaitGroup // Use a separate WaitGroup for workers

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func(id int) {
			defer workerWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Goroutine to close results channel once all workers are done
	go func() {
		workerWg.Wait() // Wait for all workers to finish processing jobs
		close(results)  // Close the results channel after all results have been sent
	}()

	// Consume results concurrently and safely increment processedCount
	for range results {
		processedCount++
	}

	// All results have been consumed, and the results channel is closed.
	// The program can now safely exit.

	return processedCount, nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
