package main

import (
	"fmt"
	"time"
)

// create a job struct
type Job struct {
	Name   string
	Delay  time.Duration
	Number int
}

// create a worker struct
type Worker struct {
	Id         int
	JobQueue   chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

// create a dispatcher struct
type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

// create constructor for worker
func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

// starting the worker with goroutine and select
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				fmt.Printf("Worker with id %d Started\n", w.Id)
				fib := Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d Finished Job %s with result %d\n", w.Id, job.Name, fib)
			case <-w.QuitChan:
				fmt.Printf("Worker with id %d Finished\n", w.Id)
				return
			}
		}
	}()
}

// stopping the worker with goroutine
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// function to asing the job to the worker
func Fibonacci(a int) int {
	if a <= 1 {
		return a
	}
	return Fibonacci(a-1) + Fibonacci(a-2)
}

// create constructor for dispatcher
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: make(chan chan Job, maxWorkers),
	}
}

// start dispatch with dispatcher
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func() {
				WorkerJobQueue := <-d.WorkerPool
				WorkerJobQueue <- job
			}()
		}
	}
}
