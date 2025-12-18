package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import atomic
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
	var processedCount int32
    var collectedResultsCount int32 // To verify all results are collected

	var workerWg sync.WaitGroup // New WaitGroup specifically for workers
	var overallWg sync.WaitGroup // Overall WaitGroup for all coordinating goroutines

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		overallWg.Add(1) // Workers are part of the overall orchestration
		go func(id int) {
			defer workerWg.Done()
			defer overallWg.Done()
			for job := range jobs {
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&processedCount, 1)
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Send jobs
	overallWg.Add(1) // Job sender is part of the overall orchestration
	go func() {
		defer overallWg.Done()
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Goroutine to close results channel AFTER all workers are done
	overallWg.Add(1)
	go func() {
		defer overallWg.Done()
		workerWg.Wait() // Wait for all worker goroutines to finish
		close(results)
	}()

	// Consume results concurrently
	overallWg.Add(1)
	go func() {
		defer overallWg.Done()
		for range results {
			atomic.AddInt32(&collectedResultsCount, 1) // Count collected results
		}
	}()

	overallWg.Wait() // Wait for all coordinating goroutines to finish

	// We should also verify that all results were collected if numJobs > 0
    if numJobs > 0 && collectedResultsCount != int32(numJobs) {
        return int(processedCount), fmt.Errorf("collected %d results, expected %d", collectedResultsCount, numJobs)
    }

	return int(processedCount), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
