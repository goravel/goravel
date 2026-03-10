package utils

import (
	"time"

	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/foundation"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/support/convert"
)

type Task struct {
	Job
	UUID  string `json:"uuid"`
	Chain []Job  `json:"chain"`
}

type Job struct {
	Delay     *time.Time           `json:"delay"`
	Signature string               `json:"signature"`
	Args      []contractsqueue.Arg `json:"args"`
}

func TaskToJson(task contractsqueue.Task, json foundation.Json) (string, error) {
	var chain []Job
	for _, taskData := range task.Chain {
		for j, arg := range taskData.Args {
			// To avoid converting []uint8 to base64
			if arg.Type == "[]uint8" {
				taskData.Args[j].Value = cast.ToIntSlice(arg.Value)
			}
		}

		job := Job{
			Signature: taskData.Job.Signature(),
			Args:      taskData.Args,
		}

		if !taskData.Delay.IsZero() {
			job.Delay = &taskData.Delay
		}

		chain = append(chain, job)
	}

	var args []contractsqueue.Arg
	for _, arg := range task.Args {
		if arg.Type == "[]uint8" {
			arg.Value = cast.ToIntSlice(arg.Value)
		}
		args = append(args, arg)
	}

	job := Job{
		Signature: task.Job.Signature(),
		Args:      args,
	}

	if !task.Delay.IsZero() {
		job.Delay = &task.Delay
	}

	t := Task{
		UUID:  task.UUID,
		Job:   job,
		Chain: chain,
	}

	payload, err := json.MarshalString(t)
	if err != nil {
		return "", err
	}

	return payload, nil
}

func JsonToTask(payload string, jobStorer contractsqueue.JobStorer, json foundation.Json) (contractsqueue.Task, error) {
	var task Task
	if err := json.UnmarshalString(payload, &task); err != nil {
		return contractsqueue.Task{}, err
	}

	var chain []contractsqueue.ChainJob
	for _, item := range task.Chain {
		job, err := jobStorer.Get(item.Signature)
		if err != nil {
			return contractsqueue.Task{}, err
		}

		jobs := contractsqueue.ChainJob{
			Job:  job,
			Args: item.Args,
		}

		if item.Delay != nil && !item.Delay.IsZero() {
			jobs.Delay = *item.Delay
		}

		chain = append(chain, jobs)
	}

	job, err := jobStorer.Get(task.Signature)
	if err != nil {
		return contractsqueue.Task{}, err
	}

	jobs := contractsqueue.ChainJob{
		Job:  job,
		Args: task.Args,
	}

	if task.Delay != nil && !task.Delay.IsZero() {
		jobs.Delay = *task.Delay
	}

	return contractsqueue.Task{
		UUID:     task.UUID,
		ChainJob: jobs,
		Chain:    chain,
	}, nil
}

func ConvertArgs(args []contractsqueue.Arg) []any {
	realArgs := make([]any, 0, len(args))
	for _, arg := range args {
		switch arg.Type {
		case "bool":
			realArgs = append(realArgs, cast.ToBool(arg.Value))
		case "int":
			realArgs = append(realArgs, cast.ToInt(arg.Value))
		case "int8":
			realArgs = append(realArgs, cast.ToInt8(arg.Value))
		case "int16":
			realArgs = append(realArgs, cast.ToInt16(arg.Value))
		case "int32":
			realArgs = append(realArgs, cast.ToInt32(arg.Value))
		case "int64":
			realArgs = append(realArgs, cast.ToInt64(arg.Value))
		case "uint":
			realArgs = append(realArgs, cast.ToUint(arg.Value))
		case "uint8":
			realArgs = append(realArgs, cast.ToUint8(arg.Value))
		case "uint16":
			realArgs = append(realArgs, cast.ToUint16(arg.Value))
		case "uint32":
			realArgs = append(realArgs, cast.ToUint32(arg.Value))
		case "uint64":
			realArgs = append(realArgs, cast.ToUint64(arg.Value))
		case "float32":
			realArgs = append(realArgs, cast.ToFloat32(arg.Value))
		case "float64":
			realArgs = append(realArgs, cast.ToFloat64(arg.Value))
		case "string":
			realArgs = append(realArgs, cast.ToString(arg.Value))
		case "[]bool":
			realArgs = append(realArgs, cast.ToBoolSlice(arg.Value))
		case "[]int":
			realArgs = append(realArgs, cast.ToIntSlice(arg.Value))
		case "[]int8":
			realArgs = append(realArgs, convert.ToSlice[int8](arg.Value))
		case "[]int16":
			realArgs = append(realArgs, convert.ToSlice[int16](arg.Value))
		case "[]int32":
			realArgs = append(realArgs, convert.ToSlice[int32](arg.Value))
		case "[]int64":
			realArgs = append(realArgs, convert.ToSlice[int64](arg.Value))
		case "[]uint":
			realArgs = append(realArgs, convert.ToSlice[uint](arg.Value))
		case "[]uint8":
			realArgs = append(realArgs, convert.ToSlice[uint8](arg.Value))
		case "[]uint16":
			realArgs = append(realArgs, convert.ToSlice[uint16](arg.Value))
		case "[]uint32":
			realArgs = append(realArgs, convert.ToSlice[uint32](arg.Value))
		case "[]uint64":
			realArgs = append(realArgs, convert.ToSlice[uint64](arg.Value))
		case "[]float32":
			realArgs = append(realArgs, convert.ToSlice[float32](arg.Value))
		case "[]float64":
			realArgs = append(realArgs, convert.ToSlice[float64](arg.Value))
		case "[]string":
			realArgs = append(realArgs, cast.ToStringSlice(arg.Value))
		}
	}
	return realArgs
}
