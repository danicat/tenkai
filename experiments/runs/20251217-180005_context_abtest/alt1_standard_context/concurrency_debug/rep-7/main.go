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
	jobs := make(chan Job, numJobs) // Buffer for all jobs
	results := make(chan Result, numWorkers) // Buffer for results, enough for each worker to send one result
	var processedCount atomic.Int32

	var workerWg sync.WaitGroup // WaitGroup for workers
	var collectorWg sync.WaitGroup // WaitGroup for results collector

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func(id int) {
			defer workerWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				processedCount.Add(1)
				
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
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for range results {
			// Just drain the results channel
		}
	}()

	// Wait for all workers to finish processing jobs
	workerWg.Wait()
	close(results) // Close results channel after all workers are done sending

	// Wait for the results collector to drain the channel
	collectorWg.Wait()

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
