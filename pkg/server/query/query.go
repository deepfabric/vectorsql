package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/server"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/valyala/fasthttp"
)

func New(port int, log logger.Log, vec vector.Vector) server.Server {
	return &query{
		log:  log,
		vec:  vec,
		port: port,
	}
}

func (q *query) Run() {
	q.srv = &fasthttp.Server{Handler: func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/query":
			q.dealQuery(ctx)
		default:
			ctx.Error("Unsupport Path", fasthttp.StatusNotFound)
		}
	}}
	q.srv.ListenAndServe(fmt.Sprintf(":%v", q.port))
}

func (q *query) Stop() {
	q.srv.Shutdown()
}

func (q *query) dealQuery(ctx *fasthttp.RequestCtx) {
	var hr HttpResult

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		hr.Msg = err.Error()
		data, _ := json.Marshal(hr)
		ctx.Write(data)
		return
	}
	body := []byte{}
	images := make(map[string][]byte)
	for k, v := range form.File {
		for _, h := range v {
			fp, err := h.Open()
			if err != nil {
				ctx.Response.SetStatusCode(400)
				hr.Msg = err.Error()
				data, _ := json.Marshal(hr)
				ctx.Write(data)
				return
			}
			data, err := ioutil.ReadAll(fp)
			if err != nil {
				fp.Close()
				ctx.Response.SetStatusCode(400)
				hr.Msg = err.Error()
				data, _ := json.Marshal(hr)
				ctx.Write(data)
				return
			}
			fp.Close()
			body = append(body, data...)
		}
		if len(body) > 0 {
			images[k] = body
			body = []byte{}
		}
	}
	v, err := q.vec.GetVector(images)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		hr.Msg = err.Error()
		data, _ := json.Marshal(hr)
		ctx.Write(data)
		return
	}
	{
		data, _ := json.Marshal(v)
		hr.Msg = string(data)
	}
	data, _ := json.Marshal(hr)
	ctx.Write(data)
}
