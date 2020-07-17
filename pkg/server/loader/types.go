package loader

import (
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/valyala/fasthttp"
)

type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Create struct {
	Name  string      `json:"name"`
	Item  []Attribute `json:"item"`
	Event []Attribute `json:"event"`
}

type loader struct {
	port int
	log  logger.Log
	cli  client.Client
	stg  storage.Storage
	srv  *fasthttp.Server
}
