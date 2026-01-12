package service

import (
	"context"
	"errors"
	"sync"
)

type jobResult struct {
	val interface{}
	err error
}

type Job struct {
	ctx  context.Context
	task func(context.Context) (interface{}, error)
	res  chan jobResult
}

type WorkerPool struct {
	jobs chan *Job
	wg   sync.WaitGroup
}

var defaultWorkerPool *WorkerPool

// NewWorkerPool creates and starts a worker pool
func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		jobs: make(chan *Job, queueSize),
	}
	wp.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobs {
				// execute task synchronously inside worker goroutine
				resVal, err := job.task(job.ctx)
				// best-effort send (if receiver gone, drop)
				select {
				case job.res <- jobResult{val: resVal, err: err}:
				default:
				}
			}
		}()
	}
	return wp
}

// Enqueue adds a job to the queue, respecting the provided ctx for enqueue/wait
func (wp *WorkerPool) Enqueue(ctx context.Context, task func(context.Context) (interface{}, error)) (interface{}, error) {
	job := &Job{
		ctx:  ctx,
		task: task,
		res:  make(chan jobResult, 1),
	}

	select {
	case wp.jobs <- job:
		// queued
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case r := <-job.res:
		return r.val, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
}

// StopWithContext stops the pool but respects a context timeout.
func (wp *WorkerPool) StopWithContext(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		wp.Stop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// InitDefaultWorkerPool initializes the package-level default worker pool.
// Returns an error if the default pool is already initialized.
func InitDefaultWorkerPool(numWorkers, queueSize int) error {
	if defaultWorkerPool != nil {
		return errors.New("default worker pool already initialized")
	}
	defaultWorkerPool = NewWorkerPool(numWorkers, queueSize)
	return nil
}

// SetDefaultWorkerPool sets the package-level default worker pool, taking ownership.
// Caller is responsible for closing a previously set pool if needed.
func SetDefaultWorkerPool(wp *WorkerPool) {
	defaultWorkerPool = wp
}

// GetDefaultPool returns the default pool or an error if it hasn't been initialized.
func GetDefaultPool() (*WorkerPool, error) {
	if defaultWorkerPool == nil {
		return nil, errors.New("default worker pool not initialized; call InitDefaultWorkerPool")
	}
	return defaultWorkerPool, nil
}

// ShutdownDefaultPool stops and clears the default pool.
func ShutdownDefaultPool(ctx context.Context) error {
	if defaultWorkerPool == nil {
		return nil
	}
	err := defaultWorkerPool.StopWithContext(ctx)
	defaultWorkerPool = nil
	return err
}
