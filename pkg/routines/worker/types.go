package worker

import "github.com/deepfabric/vectorsql/pkg/routines/task"

type Worker interface {
	Run()
	Stop()
	AddTask(task.Task)
}

type worker struct {
	ch chan struct{}
	ts chan task.Task
}
