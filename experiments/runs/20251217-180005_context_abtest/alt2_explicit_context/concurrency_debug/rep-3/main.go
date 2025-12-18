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
	var processedCount int
	var mu sync.Mutex

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 1: Race Condition
				mu.Lock()
				processedCount++
				mu.Unlock()

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
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
	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for range results {
			// Just drain the results channel.
			// The processedCount is already being updated by workers.
		}
	}()

	// Wait for all worker goroutines to finish processing jobs
	wg.Wait()

	// Close the results channel ONLY after all workers have finished sending results
	close(results)

	// Wait for the results collector goroutine to finish
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
