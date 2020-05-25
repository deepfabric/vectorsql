package query

import (
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/valyala/fasthttp"
)

type HttpResult struct {
	Msg         string `json:"msg"`
	ProcessTime string `json:"process time"`
}

type query struct {
	port int
	log  logger.Log
	vec  vector.Vector
	srv  *fasthttp.Server
}
