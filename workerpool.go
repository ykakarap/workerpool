package workerpool

import "sync"

// WorkerPool is a struct holding the information about the worker pool
type WorkerPool struct {
	maxWorkers int
	job        chan Job
	errors     chan error
	wg         *sync.WaitGroup
}

// Request submits a new job to the job channel which will eventually
// get picked by a free worker
func (wp *WorkerPool) Request(j Job) {
	wp.job <- j
}

// Errors will read from the error channel and return the errors
func (wp *WorkerPool) Errors() error {
	return <-wp.errors
}

// ShutDown will send the channe close signal to all the worker
// wait till all the workers finished working and then return
func (wp *WorkerPool) ShutDown() {
	close(wp.job)

	// TO-DO:
	//	How to handle the error channel?
	//	Use waitGroup to wait before the function returns

	wp.wg.Wait()
	return
}

// Job - Jobs will be picked up by workers and executed using the Run() function
// error returned by the Run() function will be sent to the workerpool's error channel
type Job interface {
	Run() error
}

// CreatePool creates a worker pool of 'workers' workers
func CreatePool(workers int) *WorkerPool {
	wp := &WorkerPool{
		maxWorkers: workers,
		job:        make(chan Job),
		errors:     make(chan error),
		wg:         &sync.WaitGroup{},
	}

	wp.wg.Add(wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		go handle(wp.job, wp.errors, wp.wg)
	}

	return wp
}

func handle(jc chan Job, ec chan error, wg *sync.WaitGroup) {
	var job Job
	var e error
	var done bool

	// start listening to jobs in the 'job' channel
	for {
		job, done = <-jc
		if done {
			wg.Done()
			return
		}
		e = job.Run()
		if e != nil {
			ec <- e
		}
	}
}
