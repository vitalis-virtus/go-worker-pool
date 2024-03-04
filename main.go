package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vitalis-virtus/go-worker-pool/job"
	"github.com/vitalis-virtus/go-worker-pool/pool"
)

func getMockJobs(count int) []job.Job {
	var jobs = make([]job.Job, count)
	for i := 0; i < count; i++ {
		jobs[i] = job.Job{
			Descriptor: job.Descriptor{
				JobID:       uuid.NewString(),
				JobType:     "mock_job",
				JobMetadata: nil,
			},
			ExecFunc: mockExecFunc,
			Args:     i,
		}
	}

	return jobs
}

func mockExecFunc(ctx context.Context, args interface{}) (interface{}, error) {
	argVal, ok := args.(int)
	if !ok {
		return nil, errors.New("wrong argument type")
	}

	return argVal * 2, nil
}

func main() {
	jobs := getMockJobs(100)

	pool := pool.New(5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.GenerateFrom(jobs)

	go pool.Run(ctx)

	for {
		select {
		case r, ok := <-pool.Results():
			if !ok {
				continue
			}

			fmt.Println(r.Value)
		case <-pool.Done():
			return
		default:
		}
	}
}
