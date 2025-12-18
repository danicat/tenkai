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
	var processedCount atomic.Int64

	var workersWg sync.WaitGroup // Use a dedicated WaitGroup for workers
	var resultsWg sync.WaitGroup // New WaitGroup for the results consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func(id int) {
			defer workersWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Send results. This might block if results channel is full.
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs) // Close jobs channel after all jobs are sent
	}()

	// Start a goroutine to consume results and update the counter
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			processedCount.Add(1)
		}
	}()

	// Wait for all workers to finish processing jobs
	workersWg.Wait()
	close(results) // Close the results channel after all workers have finished sending results

	// Wait for all results to be consumed
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
