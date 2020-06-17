package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/deepfabric/vectorsql/pkg/server"
	"github.com/deepfabric/vectorsql/pkg/sql/build"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/valyala/fasthttp"
)

func New(port, mcpu int, b bv.BV, log logger.Log, cli client.Client, vec vector.Vector, stg storage.Storage) server.Server {
	return &query{
		b:    b,
		cli:  cli,
		log:  log,
		vec:  vec,
		stg:  stg,
		mcpu: mcpu,
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
	qr, vec, err := q.extractParameters(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		hr.Msg = err.Error()
		data, _ := json.Marshal(hr)
		ctx.Write(data)
		return
	}
	start := time.Now()
	o, err := build.New(qr, context.New(q.cli), q.stg).Build()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		hr.Msg = err.Error()
		data, _ := json.Marshal(hr)
		ctx.Write(data)
		return
	}
	{
		q.log.Debugf("T: %v\n", o.T)
	}
	{
		q.log.Debugf("N: %v\n", o.N)
	}
	{
		q.log.Debugf("CF: %v\n", o.Cf)
	}
	{
		q.log.Debugf("IF: %v\n", o.If)
	}
	rows, err := o.Result(q.log, q.mcpu, q.b, q.cli, vec)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		hr.Msg = err.Error()
		data, _ := json.Marshal(hr)
		ctx.Write(data)
		return
	}
	{
		data, _ := json.Marshal(rows)
		hr.Sql = qr
		hr.Msg = string(data)
		hr.ProcessTime = fmt.Sprintf("%v", time.Now().Sub(start))
	}
	data, _ := json.Marshal(hr)
	ctx.Write(data)
}

func (q *query) extractParameters(ctx *fasthttp.RequestCtx) (string, []float32, error) {
	var typ string
	var body []byte
	var mp map[string]interface{}

	form, err := ctx.MultipartForm()
	if err != nil {
		return "", nil, err
	}
	fs := make(map[string]*request.Part)
	for k, v := range form.File {
		for i, h := range v {
			if i == 0 {
				typ = h.Header.Get("Content-Type")
			}
			fp, err := h.Open()
			if err != nil {
				return "", nil, err
			}
			data, err := ioutil.ReadAll(fp)
			if err != nil {
				fp.Close()
				return "", nil, err
			}
			fp.Close()
			body = append(body, data...)
		}
		if len(body) > 0 {
			if typ == "application/json" {
				if err := json.Unmarshal(body, &mp); err != nil {
					return "", nil, err
				}
			} else {
				fs[k] = &request.Part{Typ: typ, Data: body}
			}
			body = []byte{}
		}
	}
	qr, err := q.getSqlQuery(mp)
	if err != nil {
		return "", nil, err
	}
	vec, err := q.vec.GetVector(fs)
	if err != nil {
		return "", nil, err
	}
	return qr, vec, nil
}

func (q *query) getSqlQuery(mp map[string]interface{}) (string, error) {
	return getString("query", mp)
}

func getString(k string, mp map[string]interface{}) (string, error) {
	v, ok := mp[k]
	if !ok {
		return "", errors.New("Not Exist")
	}
	if _, ok := v.(string); !ok {
		return "", errors.New("Not String")
	}
	return v.(string), nil
}
