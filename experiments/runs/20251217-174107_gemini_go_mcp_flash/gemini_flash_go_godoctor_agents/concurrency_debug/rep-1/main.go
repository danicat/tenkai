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

// JobCounter is a thread-safe counter for processed jobs.
type JobCounter struct {
	mu    sync.Mutex
	count int
}

// Increment increases the job counter by one.
func (c *JobCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Get returns the current value of the job counter.
func (c *JobCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ProcessJobs processes n jobs using m workers.
// Use this function signature for testing.
// It returns the number of processed jobs and an error if something goes wrong.
func ProcessJobs(numJobs, numWorkers int) (int, error) {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var jobCounter JobCounter // Use thread-safe counter

	var wg sync.WaitGroup // For workers

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate some work
				time.Sleep(time.Millisecond)

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

	// Start a goroutine to consume results concurrently
	var resultWG sync.WaitGroup // For result consumer
	resultWG.Add(1)
	go func() {
		defer resultWG.Done()
		for range results {
			jobCounter.Increment()
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	// Close the results channel ONLY after all workers have finished sending
	close(results)

	// Wait for the result consumer to finish processing all results
	resultWG.Wait()

	return jobCounter.Get(), nil
}

func main() {
	count, err := ProcessJobs(1000, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Successfully processed %d jobs\n", count)
}
