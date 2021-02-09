package worker

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type (
	JobQueue    chan chan Job
	JobReceiver chan Job

	Job func() (id string, err error)

	// Worker executes enqueued Jobs that follows Job interface spec
	Worker struct {
		ID          int
		JobReceiver JobReceiver
		Queue       *JobQueue
		SharedGroup *sync.WaitGroup
		Quit        chan bool
	}
)

const (
	CtxKey   = "context"
	CtxValue = "Worker"
)

// NewWorker creates a new Worker with a shared JobQueue
// The JobQueue should be created by the Dispatcher and shared by all Workers
func NewWorker(id int, queue *JobQueue, group *sync.WaitGroup) *Worker {
	return &Worker{
		ID:          id,
		SharedGroup: group,
		Queue:       queue,
		JobReceiver: make(JobReceiver),
		Quit:        make(chan bool)}
}

// Start starts listening to incoming jobs and a Quit signal
func (w *Worker) Start() {
	info(fmt.Sprintf("Starting Worker %v", w.ID))

	w.SharedGroup.Add(1)

	go func() {
		*w.Queue <- w.JobReceiver
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Timestamp().
					Msg(fmt.Sprintf("Panic on Worker %v. %v", w.ID, r))
			}
		}()

		for {
			select {

			case job := <-w.JobReceiver:
				info(fmt.Sprintf("Worker %v. Will execute a Job ...", w.ID))

				id, err := job()

				if err != nil {
					log.Err(err).
						Str(CtxKey, CtxValue).
						Timestamp().
						Msg(fmt.Sprintf("Worker %v. [%v] failed.", w.ID, id))
				} else {
					info(fmt.Sprintf("Worker %v. [%v] executed with success", w.ID, id))
				}

			case <-w.Quit:
				info(fmt.Sprintf("Stopping Worker %v", w.ID))
				w.SharedGroup.Done()
				return
			}
		}
	}()
}

// Stop signals the worker to Quit listening for work requests
func (w *Worker) Stop() {
	w.Quit <- true
}
