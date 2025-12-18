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
	var jobsWg sync.WaitGroup // WaitGroup for workers to finish processing jobs
	var resultsWg sync.WaitGroup // WaitGroup for the results collector to finish

	var finalProcessedCount int32 // Thread-safe counter for processed jobs from results collector

	// Start workers
	for i := 0; i < numWorkers; i++ {
		jobsWg.Add(1)
		go func(id int) {
			defer jobsWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// No longer incrementing processedCount here.
				// The results collector will count the successfully processed jobs.

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

	// Start a goroutine to collect results concurrently
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			// Increment a thread-safe counter for truly processed jobs
			atomic.AddInt32(&finalProcessedCount, 1)
		}
	}()

	// Wait for all workers to finish processing jobs.
	// This ensures all results have been sent to the 'results' channel.
	jobsWg.Wait()
	close(results) // Close the results channel after all workers are done.

	// Wait for the results collector to finish draining all results.
	resultsWg.Wait()

	return int(finalProcessedCount), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
