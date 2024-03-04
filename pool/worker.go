package pool

import (
	"context"
	"fmt"
	"github.com/vitalis-virtus/go-worker-pool/job"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan job.Job, results chan<- job.Result) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			results <- job.Execute(ctx)
		case <-ctx.Done():
			fmt.Printf("\ncanceled worker, error details: %v\n", ctx.Err())
			results <- job.Result{
				Error: ctx.Err(),
			}
			return
		}
	}
}
