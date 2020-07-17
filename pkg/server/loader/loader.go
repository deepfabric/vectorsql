package loader

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/server"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/storage/index"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/valyala/fasthttp"
)

func New(port int, log logger.Log, cli client.Client, stg storage.Storage) server.Server {
	return &loader{
		cli:  cli,
		log:  log,
		stg:  stg,
		port: port,
	}
}

func (s *loader) Run() {
	s.srv = &fasthttp.Server{Handler: func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/create":
			s.dealCreate(ctx)
		default:
			ctx.Error("Unsupport Path", fasthttp.StatusNotFound)
		}
	}}
	s.srv.ListenAndServe(fmt.Sprintf(":%v", s.port))
}

func (s *loader) Stop() {
	s.srv.Shutdown()
}

func (s *loader) dealCreate(ctx *fasthttp.RequestCtx) {
	var req Create

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	if n := len(req.Item); n > 0 {
		var md metadata.Metadata

		md.IsE = false
		id := req.Name + "_item"
		sql := fmt.Sprintf("CREATE TABLE %s ", id)
		md.Attrs = make([]metadata.Attribute, n)
		for i := 0; i < n; i++ {
			md.Attrs[i].Name = req.Item[i].Name
			name, typ := stringToType(req.Item[i].Type)
			if len(name) == 0 {
				ctx.Response.SetStatusCode(400)
				ctx.Write([]byte(fmt.Sprintf("unsupport type '%s'", req.Item[i].Type)))
				return
			}
			md.Attrs[i].Type = typ
			if i == 0 {
				sql += fmt.Sprintf("(%s %s", req.Item[i].Name, name)
			} else {
				sql += fmt.Sprintf(", %s %s", req.Item[i].Name, name)
			}
		}
		if md.Attrs[0].Name != index.SEQ && md.Attrs[0].Type != types.T_uint64 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute seq(uint64)"))
			return
		}
		sql += fmt.Sprintf(") engine=ReplacingMergeTree() PARTITION BY intDiv(seq, 1000000) ORDER BY (seq);")
		if err := s.stg.NewRelation(id, md); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
		if err := s.cli.Exec(sql, nil); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
	}
	if n := len(req.Event); n > 0 {
		var md metadata.Metadata

		md.IsE = true
		id := req.Name + "_event"
		sql := fmt.Sprintf("CREATE TABLE %s ", id)
		md.Attrs = make([]metadata.Attribute, n)
		for i := 0; i < n; i++ {
			md.Attrs[i].Name = req.Event[i].Name
			name, typ := stringToType(req.Event[i].Type)
			if len(name) == 0 {
				ctx.Response.SetStatusCode(400)
				ctx.Write([]byte(fmt.Sprintf("unsupport type '%s'", req.Event[i].Type)))
				return
			}
			md.Attrs[i].Type = typ
			if i == 0 {
				sql += fmt.Sprintf("(%s %s", req.Event[i].Name, name)
			} else {
				sql += fmt.Sprintf(", %s %s", req.Event[i].Name, name)
			}
		}
		if md.Attrs[0].Name != index.SEQ && md.Attrs[0].Type != types.T_uint64 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute seq(uint64)"))
			return
		}
		if md.Attrs[1].Name != index.ID && md.Attrs[1].Type != types.T_uint64 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute id(uint64)"))
			return
		}
		sql += fmt.Sprintf(") engine=ReplacingMergeTree() PARTITION BY intDiv(seq, 1000000) ORDER BY (seq);")
		if err := s.stg.NewRelation(id, md); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
		if err := s.cli.Exec(sql, nil); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
	}
	ctx.Write([]byte("success"))
}

func stringToType(name string) (string, uint32) {
	switch strings.ToLower(name) {
	case "int8":
		return "Uint8", types.T_int8
	case "int16":
		return "Int16", types.T_int16
	case "int32":
		return "Int32", types.T_int32
	case "int64":
		return "Int64", types.T_int64
	case "uint8":
		return "UInt8", types.T_uint8
	case "uint16":
		return "UInt16", types.T_uint16
	case "uint32":
		return "UInt32", types.T_uint32
	case "uint64":
		return "UInt64", types.T_uint64
	case "float32":
		return "Float32", types.T_float32
	case "float64":
		return "Float64", types.T_float64
	case "string":
		return "String", types.T_string
	case "datetime":
		return "Datetime", types.T_timestamp
	}
	return "", 0
}
