package pool

import (
	"context"
	"github.com/vitalis-virtus/go-worker-pool/job"
	"sync"
)

type Pool interface {
	GenerateFrom(jobsBulk []job.Job)
	Run(ctx context.Context)
	Done() chan struct{}
	Results() chan job.Result
}

type pool struct {
	count   int
	jobs    chan job.Job
	results chan job.Result
	done    chan struct{}
}

func New(count int) Pool {
	return &pool{
		count:   count,
		jobs:    make(chan job.Job),
		results: make(chan job.Result),
		done:    make(chan struct{}),
	}
}

func (p *pool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < p.count; i++ {
		wg.Add(1)

		go worker(ctx, &wg, p.jobs, p.results)
	}

	wg.Wait()
	close(p.done)
	close(p.results)
}

func (p *pool) GenerateFrom(jobsBulk []job.Job) {
	for i, _ := range jobsBulk {
		p.jobs <- jobsBulk[i]
	}
	close(p.jobs)
}

func (p *pool) Results() chan job.Result {
	return p.results
}

func (p *pool) Done() chan struct{} {
	return p.done
}
