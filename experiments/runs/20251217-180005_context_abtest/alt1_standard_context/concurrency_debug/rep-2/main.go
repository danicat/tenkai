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
	var processedCount atomic.Int64   // FIX: Use atomic for thread-safe counter

	var workerWG sync.WaitGroup
	var resultCollectorWG sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWG.Add(1)
		go func(id int) {
			defer workerWG.Done()
			for job := range jobs {
				time.Sleep(time.Millisecond)

				// FIX 1: Race Condition
				processedCount.Add(1)

				// Send result to results channel
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Start a goroutine to send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs) // Close jobs channel after all jobs are sent
	}()

	// Start a goroutine to collect results
	resultCollectorWG.Add(1)
	go func() {
		defer resultCollectorWG.Done()
		for range results {
			// Just drain the results channel
		}
	}()

	// Wait for all workers to finish processing jobs
	workerWG.Wait()
	close(results) // Close results channel after all workers are done sending
	resultCollectorWG.Wait() // Wait for results collector to finish draining

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
