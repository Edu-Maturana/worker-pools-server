package main

import (
	"fmt"
	"time"
)

// Job represents a job to be executed, with a name and a number and a delay
type Job struct {
	Name   string
	Delay  time.Duration
	Number int
}

// Worker will be our concurrency-friendly worker
type Worker struct {
	Id         int
	JobQueue   chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

// Dispatcher is a dispatcher that will dispatch jobs to workers
type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

// NewWorker returns a new Worker with the provided id and workerpool
func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

// Start method starts all workers
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobQueue // add job to pool

			// Multiplexing
			select {
			case job := <-w.JobQueue: // get job from queue
				fmt.Printf("Worker with id %d started\n", w.Id)
				fib := Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d finished with result %d\n", w.Id, fib)
			case <-w.QuitChan: // quit if worker is told to do so
				fmt.Printf("Worker with id %d stopped\n", w.Id)
			}
		}
	}()
}

// Stop method stop the worker
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// Fibonacci calculates the fibonacci sequence
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// NewDispatcher returns a new Dispatcher with the provided maxWorkers
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	worker := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: worker,
	}
}

// Dispatch will dispatch jobs to workers
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue: // get job from queue
			// Asign the job to a worker
			go func() {
				workerJobQueue := <-d.WorkerPool // get worker from pool
				workerJobQueue <- job            // Workers will read from this channel
			}()
		}
	}
}
