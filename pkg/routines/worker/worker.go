package worker

import "github.com/deepfabric/vectorsql/pkg/routines/task"

func New() Worker {
	return &worker{
		ch: make(chan struct{}),
		ts: make(chan task.Task),
	}
}

func (w *worker) Run() {
	for {
		select {
		case <-w.ch:
			w.ch <- struct{}{}
			return
		case t := <-w.ts:
			t.Stop(t.Execute())
		}
	}
}

func (w *worker) Stop() {
	w.ch <- struct{}{}
	<-w.ch
	close(w.ch)
	close(w.ts)
}

func (w *worker) AddTask(t task.Task) {
	w.ts <- t
}
