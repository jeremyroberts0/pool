package pool

// Job a job to execute in the thread pool
type Job func() error

// Pool a single thread pool.
type Pool struct {
	jobs []Job
	max  int
}

//New Create a new thread pool.  When run, jobs are executed in <max> goroutines
func New(max int) Pool {
	return Pool{
		jobs: []Job{},
		max:  max,
	}
}

// Add add a job to the pool
func (p *Pool) Add(j Job) error {
	p.jobs = append(p.jobs, j)
	return nil
}

// Run runs the thread pool, once a thread pool has started running jobs cannot be added to them.
func (p *Pool) Run() (errors []error) {
	queues := [][]Job{}
	for queueIndex := 0; len(queues) < p.max; queueIndex++ {
		queues = append(queues, []Job{})
	}

	currentQueue := 0

	for _, job := range p.jobs {
		queues[currentQueue] = append(queues[currentQueue], job)
		currentQueue++
		if currentQueue == p.max {
			currentQueue = 0
		}
	}

	doneChan := make(chan bool)
	errChan := make(chan error)
	jobsRemaining := len(p.jobs)

	for _, queue := range queues {
		go runQueue(queue, doneChan, errChan)
	}

	for {
		select {
		case _ = <-doneChan:
			jobsRemaining--
		case err := <-errChan:
			errors = append(errors, err)
		}
		if jobsRemaining == 0 {
			return
		}
	}
}

func runQueue(jobs []Job, doneChan chan bool, errChan chan error) {
	for _, job := range jobs {
		err := job()
		if err != nil {
			errChan <- err
		}
		doneChan <- true
	}
}
