package job

import "context"

type Descriptor struct {
	JobID       string
	JobType     string
	JobMetadata map[string]interface{}
}

type ExecutionFunc func(ctx context.Context, args interface{}) (interface{}, error)

type Job struct {
	Descriptor Descriptor
	ExecFunc   ExecutionFunc
	Args       interface{}
}

func (j *Job) Execute(ctx context.Context) Result {
	value, err := j.ExecFunc(ctx, j.Args)
	if err != nil {
		return Result{
			Error:    err,
			Metadata: j.Descriptor.JobMetadata,
		}
	}

	return Result{
		Value:    value,
		Metadata: j.Descriptor.JobMetadata,
	}
}
