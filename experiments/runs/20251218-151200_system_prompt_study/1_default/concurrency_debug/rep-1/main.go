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
	var processedCount uint32         // Fixed: Use atomic for thread-safe counter

	var wg sync.WaitGroup

		var resultsWg sync.WaitGroup

	

		// Start workers

		for i := 0; i < numWorkers; i++ {

			wg.Add(1)

			go func(id int) {

				defer wg.Done()

				for job := range jobs {

					// Simulate some work

					time.Sleep(time.Millisecond)

	

					// Send result. This might block if results channel is full.

					// The resultCollector goroutine will relieve this pressure.

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

	

		// Start a goroutine to collect results and increment the processed count

		// This runs concurrently to prevent workers from blocking indefinitely

		// if the results channel buffer fills up.

		// The resultCollector will terminate when the results channel is closed.

		resultsWg.Add(1)

		go func() {

			defer resultsWg.Done()

			for range results {

				atomic.AddUint32(&processedCount, 1) // Atomically increment processedCount

			}

		}()

	

		// Wait for all workers to finish.

		// This implicitly means all results have been sent to the results channel.

		wg.Wait()

	

		// All workers are done, and thus no more results will be sent.

		// It's now safe to close the results channel, which will unblock the resultCollector.

		close(results)

	

		// Wait for the result collector to finish processing all results.

		resultsWg.Wait()

	

		return int(processedCount), nil}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
