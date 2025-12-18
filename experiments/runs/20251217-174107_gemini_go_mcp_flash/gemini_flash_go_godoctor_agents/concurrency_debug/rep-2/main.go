package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import the atomic package
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
	var processedCount atomic.Int64   // RACE CONDITION FIX: Use atomic for thread-safe increments

	var workerWg sync.WaitGroup   // For workers to finish
	var resultReaderWg sync.WaitGroup // For result consumer to finish

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func(id int) {
			defer workerWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// BUG 1 FIX: Race Condition
				processedCount.Add(1)

				// BUG 2: Potential Deadlock if buffer fills
				// Workers block here if no one reads 'results'
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Goroutine to send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Goroutine to consume results
	resultReaderWg.Add(1)
	go func() {
		defer resultReaderWg.Done()
		for range results {
			// Just drain the results to prevent workers from blocking
		}
	}()

	// Wait for all workers to finish
	workerWg.Wait()
	// All workers are done sending results, so it's safe to close the results channel
	close(results)

	// Wait for the result reader to finish draining the channel
	resultReaderWg.Wait()

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