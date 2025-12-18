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

func ProcessJobs(numJobs, numWorkers int) (int, error) {
	jobs := make(chan Job, numJobs)    // Buffer for all jobs to avoid blocking producer
	results := make(chan Result, numJobs) // Buffer for all results to avoid blocking workers
	var processedCount atomic.Int32 // THREAD-SAFE counter

	var workerWg sync.WaitGroup     // WaitGroup for workers
	var resultConsumerWg sync.WaitGroup // WaitGroup for the result consumer

	// Start workers
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func(id int) {
			defer workerWg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

				// Send result
				results <- Result{JobID: job.ID, Val: job.ID * 2}
			}
		}(i)
	}

	// Job producer
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	// Start result consumer
	resultConsumerWg.Add(1)
	go func() {
		defer resultConsumerWg.Done()
		for range results {
			processedCount.Add(1)
		}
	}()

	// Wait for all workers to finish processing jobs
	workerWg.Wait()
	// All workers are done, so no more results will be sent.
	// It's safe to close the results channel.
	close(results)

	// Wait for the result consumer to finish draining the results channel
	resultConsumerWg.Wait()

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
