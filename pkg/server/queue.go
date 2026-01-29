package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

// JobQueue manages job queuing and execution with a worker pool.
type JobQueue struct {
	store        *JobStore
	jobChan      chan *Job
	maxWorkers   int
	wg           sync.WaitGroup
	running      atomic.Bool
	ctx          context.Context
	cancel       context.CancelFunc
	workerCancel []context.CancelFunc
}

// NewJobQueue creates a new job queue with the specified number of workers.
func NewJobQueue(store *JobStore, maxWorkers int) *JobQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &JobQueue{
		store:      store,
		jobChan:    make(chan *Job, maxWorkers*2), // Buffer for pending jobs
		maxWorkers: maxWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts the worker pool.
func (q *JobQueue) Start(processor JobProcessor) error {
	if !q.running.CompareAndSwap(false, true) {
		return fmt.Errorf("queue already running")
	}

	log.Printf("Starting job queue with %d workers", q.maxWorkers)

	q.workerCancel = make([]context.CancelFunc, q.maxWorkers)

	for i := 0; i < q.maxWorkers; i++ {
		workerCtx, workerCancel := context.WithCancel(q.ctx)
		q.workerCancel[i] = workerCancel

		q.wg.Add(1)
		go q.worker(workerCtx, processor, i)
	}

	return nil
}

// Stop gracefully shuts down the queue.
func (q *JobQueue) Stop() {
	if !q.running.CompareAndSwap(true, false) {
		return
	}

	log.Println("Stopping job queue...")

	// Cancel all workers
	for _, cancel := range q.workerCancel {
		if cancel != nil {
			cancel()
		}
	}

	// Close the job channel to stop accepting new jobs
	close(q.jobChan)

	// Wait for all workers to finish
	q.wg.Wait()

	q.cancel()
	log.Println("Job queue stopped")
}

// Enqueue adds a job to the queue.
// Returns an error if the queue is not running.
func (q *JobQueue) Enqueue(job *Job) error {
	if !q.running.Load() {
		return fmt.Errorf("queue is not running")
	}

	select {
	case q.jobChan <- job:
		return nil
	default:
		return fmt.Errorf("queue is full, max workers: %d", q.maxWorkers)
	}
}

// RunningCount returns the approximate number of running jobs.
func (q *JobQueue) RunningCount() int {
	return len(q.jobChan)
}

// JobProcessor is the function that processes a single job.
type JobProcessor func(ctx context.Context, job *Job) error

// worker processes jobs from the queue.
func (q *JobQueue) worker(ctx context.Context, processor JobProcessor, workerNum int) {
	defer q.wg.Done()

	log.Printf("Worker %d started", workerNum)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", workerNum)
			return
		case job, ok := <-q.jobChan:
			if !ok {
				log.Printf("Worker %d: channel closed, exiting", workerNum)
				return
			}

			log.Printf("Worker %d processing job %s", workerNum, job.ID)

			// Update job status to running
			q.store.UpdateStatus(job.ID, JobStatusRunning)

			// Process the job
			if err := processor(ctx, job); err != nil {
				q.store.SetError(job.ID, err)
				log.Printf("Worker %d: job %s failed: %v", workerNum, job.ID, err)
			} else {
				// Job completed successfully, store the result
				q.store.SetResult(job.ID, job.Result)
				log.Printf("Worker %d: job %s completed successfully", workerNum, job.ID)
			}
		}
	}
}
