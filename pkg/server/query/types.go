package query

import (
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	"github.com/valyala/fasthttp"
)

type HttpResult struct {
	Msg         string `json:"msg"`
	Sql         string `json:"sql"`
	ProcessTime string `json:"process time"`
}

type query struct {
	port int
	mcpu int
	b    bv.BV
	log  logger.Log
	vec  vector.Vector
	cli  client.Client
	stg  storage.Storage
	srv  *fasthttp.Server
}
