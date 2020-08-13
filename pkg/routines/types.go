package routines

import (
	"sync"

	"github.com/deepfabric/vectorsql/pkg/routines/task"
	"github.com/deepfabric/vectorsql/pkg/routines/worker"
)

type Routines interface {
	Run()
	Stop()
	AddTask(task.Task)
}

type routines struct {
	sync.Mutex
	cnt uint
	num uint
	ch  chan struct{}
	ws  []worker.Worker
}
