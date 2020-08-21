package server

import (
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/deepfabric/vectorsql/pkg/routines"
	"github.com/deepfabric/vectorsql/pkg/routines/task"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/valyala/fasthttp"
)

type Server interface {
	Run()
	Stop()
}

type Attribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Index bool   `json:"index"`
}

type Create struct {
	Name  string      `json:"name"`
	Item  []Attribute `json:"item"`
	Event []Attribute `json:"event"`
}

type Config struct {
	B   bv.BV
	Log logger.Log
	Cli client.Client
	Vec vector.Vector
	Ctx context.Context
	Stg storage.Storage
	Rts routines.Routines
}

type faceTask struct {
	vec vector.Vector
	ch  chan task.TaskResult
	req map[string]*request.Part
}

type faceResult struct {
	err error
	xb  []float32
}

type server struct {
	port int
	b    bv.BV
	dsn  string
	log  logger.Log
	cli  client.Client
	vec  vector.Vector
	ctx  context.Context
	stg  storage.Storage
	srv  *fasthttp.Server
	rts  routines.Routines
}
