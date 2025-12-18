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
	var processedCount int            // RACE CONDITION: Not thread-safe
	var mu sync.Mutex

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

				// BUG 1: Race Condition - FIXED
				mu.Lock()
				processedCount++
				mu.Unlock()

				// BUG 2: Potential Deadlock if buffer fills - Will be fixed by consuming results concurrently
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

	// Start a goroutine to consume results concurrently
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			// Just drain the results
		}
	}()

	// Wait for all workers to finish processing jobs and sending results
	workersWg.Wait()
	close(results) // Close results channel AFTER all workers have finished sending

	// Wait for the results consumer to finish draining the results channel
	resultsWg.Wait()

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
