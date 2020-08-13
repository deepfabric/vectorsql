package routines

import (
	"sync"

	"github.com/deepfabric/vectorsql/pkg/routines/task"
	"github.com/deepfabric/vectorsql/pkg/routines/worker"
)

func New(num int) Routines {
	r := &routines{
		num: uint(num),
		ch:  make(chan struct{}),
		ws:  make([]worker.Worker, num),
	}
	for i := 0; i < num; i++ {
		r.ws[i] = worker.New()
	}
	return r
}

func (r *routines) Run() {
	var wg sync.WaitGroup

	for i, j := 0, len(r.ws); i < j; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r.ws[idx].Run()
		}(i)
	}
	for {
		select {
		case <-r.ch:
			for _, w := range r.ws {
				w.Stop()
			}
			wg.Wait()
			r.ch <- struct{}{}
			return
		}
	}
}

func (r *routines) Stop() {
	r.ch <- struct{}{}
	<-r.ch
}

func (r *routines) AddTask(t task.Task) {
	r.Lock()
	r.ws[r.cnt%r.num].AddTask(t)
	r.cnt++
	r.Unlock()
}
