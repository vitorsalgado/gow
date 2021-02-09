package worker

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// Dispatcher represents a Job listener and dispatcher
type Dispatcher struct {
	Queue        JobQueue
	JobReceiver  JobReceiver
	Workers      []*Worker
	WorkersGroup *sync.WaitGroup
	InnerGroup   sync.WaitGroup
	quit         chan bool
}

// NewDispatcher creates a dispatcher with all required configurations
// Call Run() to Wake the Dispatcher
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		Queue:        make(JobQueue),
		JobReceiver:  make(JobReceiver),
		Workers:      make([]*Worker, maxWorkers),
		WorkersGroup: &sync.WaitGroup{},
		InnerGroup:   sync.WaitGroup{},
		quit:         make(chan bool)}
}

// Run starts the dispatcher initializing all Worker instances
func (d *Dispatcher) Run() {
	info("Starting Dispatcher")

	for i := 0; i < len(d.Workers); i++ {
		d.Workers[i] = NewWorker(i, &d.Queue, d.WorkersGroup)
		d.Workers[i].Start()
	}

	d.InnerGroup.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Timestamp().
					Msg(fmt.Sprintf("Panic on Dispatcher. %v", r))

				d.stopJobs()
				d.WorkersGroup.Wait()
				d.InnerGroup.Done()
			}
		}()

		for {
			select {
			case job := <-d.JobReceiver:
				d.InnerGroup.Add(1)
				go func(j Job) {
					defer d.InnerGroup.Done()
					receiver := <-d.Queue
					receiver <- j
				}(job)
			case <-d.quit:
				d.InnerGroup.Done()
				return
			}
		}
	}()
}

// Dispatch dispatches a Job instance to the JobReceiver channel.
// The will be pickup by a available Worker and executed.
func (d *Dispatcher) Dispatch(job Job) {
	d.JobReceiver <- job
}

// Quit signals all workers to Quit processing Job instances
func (d *Dispatcher) Quit() {
	info("Dispatcher will Quit. Stopping Workers ...")

	d.stopJobs()
	d.WorkersGroup.Wait()

	info("Workers stopped")

	d.quit <- true

	d.InnerGroup.Wait()

	info("Dispatcher finished")
}

func (d *Dispatcher) stopJobs() {
	for i := 0; i < len(d.Workers); i++ {
		d.Workers[i].Stop()
	}
}
