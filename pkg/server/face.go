package server

import "github.com/deepfabric/vectorsql/pkg/routines/task"

func (t *faceTask) Stop(r task.TaskResult) {
	t.ch <- r
}

func (t *faceTask) Execute() task.TaskResult {
	xb, err := t.vec.GetVector(t.req)
	return &faceResult{err, xb}
}

func (t *faceResult) Error() error {
	return t.err
}

func (t *faceResult) Result() interface{} {
	return t.xb
}
