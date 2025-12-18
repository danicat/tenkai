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
// It returns the number of processed jobs and an error if something goes wrong.
func ProcessJobs(numJobs, numWorkers int) (int64, error) {
	jobs := make(chan Job, numJobs)       // Buffer for all jobs to avoid blocking sender
	results := make(chan Result, numJobs) // Buffer for all results to avoid blocking workers

	var processedCount atomic.Int64 // Thread-safe counter

	var wg sync.WaitGroup // For workers

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				processedCount.Add(1) // Increment atomically

				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Send jobs in a goroutine
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs) // Close the jobs channel after all jobs are sent
	}()

	// Goroutine to close the results channel after all workers are done
	// This is crucial to unblock the main goroutine's results consumption loop
	go func() {
		wg.Wait()      // Wait for all worker goroutines to finish
		close(results) // Close the results channel
	}()

	// Consume results in the main goroutine
	// This loop will run until the 'results' channel is closed
	for range results {
		// Draining the channel. The actual count is handled by processedCount.
	}

	return processedCount.Load(), nil // Return the atomically loaded count
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
