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
	var processedCount int32          // RACE CONDITION: Not thread-safe

	var wg sync.WaitGroup
	var resultsWG sync.WaitGroup // New WaitGroup for results consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&processedCount, 1)
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Consume results in a separate goroutine. Start this BEFORE waiting for workers.
	resultsWG.Add(1)
	go func() {
		defer resultsWG.Done()
		for range results {
			// Just drain the channel
		}
	}()

	// Send jobs
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Wait for all workers to finish processing jobs.
	wg.Wait()
	close(results) // Close results channel AFTER all workers are done sending.

	// Wait for results to be consumed.
	resultsWG.Wait()
	return int(atomic.LoadInt32(&processedCount)), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
