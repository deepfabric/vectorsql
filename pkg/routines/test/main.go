package main

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/routines"
	"github.com/deepfabric/vectorsql/pkg/routines/task"
)

type testTask struct {
	x, y int
	ch   chan task.TaskResult
}

type testTaskResult struct {
	z   int
	err error
}

func (t *testTask) Stop(r task.TaskResult) {
	t.ch <- r
}

func (t *testTask) Execute() task.TaskResult {
	return &testTaskResult{t.x + t.y, nil}
}

func (t *testTaskResult) Error() error {
	return t.err
}

func (t *testTaskResult) Result() interface{} {
	return t.z
}

func main() {
	var ts []*testTask

	r := routines.New(10)
	go r.Run()
	for i, j := 0, 10000; i < j; i++ {
		ts = append(ts, &testTask{i, j, make(chan task.TaskResult, 1)})
	}
	for i, j := 0, 10000; i < j; i++ {
		r.AddTask(ts[i])
	}
	for i, j := 0, 10000; i < j; i++ {
		select {
		case r := <-ts[i].ch:
			fmt.Printf("Result: %d\n", r.Result())
		}
	}
	r.Stop()
}
