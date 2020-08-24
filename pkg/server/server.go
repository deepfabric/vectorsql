package server

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/deepfabric/vectorsql/pkg/routines/task"
	"github.com/deepfabric/vectorsql/pkg/sql/build"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/valyala/fasthttp"
)

func New(port int, dsn string, cfg *Config) Server {
	return &server{
		dsn:  dsn,
		port: port,
		b:    cfg.B,
		cli:  cfg.Cli,
		log:  cfg.Log,
		stg:  cfg.Stg,
		vec:  cfg.Vec,
		ctx:  cfg.Ctx,
		rts:  cfg.Rts,
	}
}

func (s *server) Run() {
	go s.rts.Run()
	s.srv = &fasthttp.Server{
		MaxRequestBodySize: 4 << 30,
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/query":
				s.dealQuery(ctx)
			case "/queryWithVector":
				s.dealQueryWithVector(ctx)
			case "/create":
				s.dealCreate(ctx)
			case "/insert":
				s.dealInsert(ctx)
			case "/insertWithVector":
				s.dealInsertWithVector(ctx)
			default:
				ctx.Error("Unsupport Path", fasthttp.StatusNotFound)
			}
		},
	}
	if err := s.srv.ListenAndServe(fmt.Sprintf(":%v", s.port)); err != nil {
		s.log.Debugf("Failed to listen '%v': %v\n", s.port, err)
	}
}

func (s *server) Stop() {
	s.rts.Stop()
	s.srv.Shutdown()
}

func (s *server) dealQuery(ctx *fasthttp.RequestCtx) {
	var buf bytes.Buffer

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	qr, vec, err := s.extractParameters(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	o, err := build.New(qr, s.ctx, s.stg).Build()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	{
		s.log.Debugf("T: %v\n", o.T)
	}
	{
		s.log.Debugf("N: %v\n", o.N)
	}
	{
		s.log.Debugf("CF: %v\n", o.Cf)
	}
	{
		s.log.Debugf("IF: %v\n", o.If)
	}
	rows, err := o.Result(s.log, s.b, s.cli, vec)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	if err := csv.NewWriter(bufio.NewWriter(&buf)).WriteAll(rows); err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Write([]byte(err.Error()))
		return
	}
	ctx.Write(buf.Bytes())
}

func (s *server) dealQueryWithVector(ctx *fasthttp.RequestCtx) {
	var buf bytes.Buffer

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	qr, vec, err := s.extractParametersWithVector(ctx)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	o, err := build.New(qr, s.ctx, s.stg).Build()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	{
		s.log.Debugf("T: %v\n", o.T)
	}
	{
		s.log.Debugf("N: %v\n", o.N)
	}
	{
		s.log.Debugf("CF: %v\n", o.Cf)
	}
	{
		s.log.Debugf("IF: %v\n", o.If)
	}
	rows, err := o.Result(s.log, s.b, s.cli, vec)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	if err := csv.NewWriter(bufio.NewWriter(&buf)).WriteAll(rows); err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Write([]byte(err.Error()))
		return
	}
	ctx.Write(buf.Bytes())
}

func (s *server) dealCreate(ctx *fasthttp.RequestCtx) {
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

		if n < 3 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte(fmt.Sprintf("need at least uid, xid, pic attributes")))
			return
		}
		md.IsE = false
		id := metadata.Ikey(req.Name)
		sql := fmt.Sprintf("CREATE TABLE %s ", id)
		md.Attrs = make([]metadata.Attribute, n)
		for i := 0; i < n; i++ {
			md.Attrs[i].Name = req.Item[i].Name
			md.Attrs[i].Index = req.Item[i].Index
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
		if md.Attrs[0].Name != "uid" && md.Attrs[0].Type != types.T_uint64 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute uid(uint64)"))
			return
		}
		if md.Attrs[1].Name != "xid" && md.Attrs[1].Type != types.T_uint64 {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute xid(uint64)"))
			return
		}
		if md.Attrs[2].Name != "pic" && md.Attrs[2].Type != types.T_string {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte("need attribute pic(string)"))
			return
		}
		sql += fmt.Sprintf(") engine=ReplacingMergeTree() PARTITION BY intDiv(uid, 1000000) ORDER BY (uid);")
		if err := s.stg.NewRelation(id, md); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
		{
			s.log.Debugf("create table use '%s'\n", sql)
		}
		if err := s.cli.Exec(sql, nil); err != nil {
			ctx.Response.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return
		}
	}
	ctx.Write([]byte("success"))
}

func (s *server) dealInsert(ctx *fasthttp.RequestCtx) {
	var rids []string // removed remove id list

	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	name := ctx.QueryArgs().Peek("name")
	if len(name) == 0 {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte("need table's name"))
		return
	}
	ts, err := csv.NewReader(bytes.NewReader(ctx.PostBody())).ReadAll()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	id := metadata.Ikey(string(name))
	r, err := s.stg.Relation(id)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	attrs := r.Metadata().Attrs
	for len(ts) > 0 {
		{
			s.log.Debugf("tuples %v\n", len(ts))
		}
		n := len(ts)
		if n > 5000 {
			n = 5000
		}
		rs, xbs, xids, iargs, cargs, err := s.convert(ts[:n], attrs)
		if err != nil {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte(err.Error()))
			return
		}
		rids = append(rids, rs...)
		{
			query := fmt.Sprintf("insert into %s (", id)
			for i, attr := range attrs {
				if i == 0 {
					query += fmt.Sprintf("%s", attr.Name)
				} else {
					query += fmt.Sprintf(", %s", attr.Name)
				}
			}
			query += ") VALUES"
			{
				s.log.Debugf("convert ok %v: %v -> %v\n", query, n, len(cargs))
			}
			for _, _ = range cargs {
				query += " ("
				for i := range attrs {
					if i == 0 {
						query += "?"
					} else {
						query += ", ?"
					}
				}
				query += ")"
			}
			cli, err := client.New(s.dsn)
			if err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
			if err := cli.Exec(query, cargs); err != nil {
				cli.Close()
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
			cli.Close()
		}
		{
			if err := r.AddTuples(iargs); err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
		}
		{
			{
				s.log.Debugf("xbs: %v, xids: %v\n", len(xbs), len(xids))
			}
			if err := s.b.Add(xbs, xids); err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
		}
		ts = ts[n:]
	}
	ctx.Write([]byte(fmt.Sprintf("success: skip uid list: %v", rids)))
}

func (s *server) dealInsertWithVector(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "application/json")
	name := ctx.QueryArgs().Peek("name")
	if len(name) == 0 {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte("need table's name"))
		return
	}
	ts, err := csv.NewReader(bytes.NewReader(ctx.PostBody())).ReadAll()
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	id := metadata.Ikey(string(name))
	r, err := s.stg.Relation(id)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}
	attrs := r.Metadata().Attrs
	for len(ts) > 0 {
		{
			s.log.Debugf("tuples %v\n", len(ts))
		}
		n := len(ts)
		if n > 5000 {
			n = 5000
		}
		xbs, xids, iargs, cargs, err := s.convertWithVector(ts[:n], attrs)
		if err != nil {
			ctx.Response.SetStatusCode(400)
			ctx.Write([]byte(err.Error()))
			return
		}
		{
			query := fmt.Sprintf("insert into %s (", id)
			for i, attr := range attrs {
				if i == 0 {
					query += fmt.Sprintf("%s", attr.Name)
				} else {
					query += fmt.Sprintf(", %s", attr.Name)
				}
			}
			query += ") VALUES"
			{
				s.log.Debugf("convert ok %v: %v -> %v\n", query, n, len(cargs))
			}
			for _, _ = range cargs {
				query += " ("
				for i := range attrs {
					if i == 0 {
						query += "?"
					} else {
						query += ", ?"
					}
				}
				query += ")"
			}
			cli, err := client.New(s.dsn)
			if err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
			if err := cli.Exec(query, cargs); err != nil {
				cli.Close()
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
			cli.Close()
		}
		{
			if err := r.AddTuples(iargs); err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
		}
		{
			{
				s.log.Debugf("xbs: %v, xids: %v\n", len(xbs), len(xids))
			}
			if err := s.b.Add(xbs, xids); err != nil {
				ctx.Response.SetStatusCode(500)
				ctx.Write([]byte(err.Error()))
				return
			}
		}
		ts = ts[n:]
	}
	ctx.Write([]byte(fmt.Sprintf("success")))
}

func (s *server) convert(ts [][]string, attrs []metadata.Attribute) ([]string, []float32, []int64, []interface{}, [][]interface{}, error) {
	var rids []string // removed id list

	xbs := make([]float32, 0, len(ts))
	xids := make([]int64, 0, len(ts))
	iargs := make([]interface{}, len(attrs))
	cargs := make([][]interface{}, 0, len(ts))
	{
		for i := range attrs {
			iargs[i] = newSlice(attrs[i].Type, len(ts))
		}
	}
	fts := make([]*faceTask, len(ts))
	for i, t := range ts {
		req := make(map[string]*request.Part)
		req["a"] = &request.Part{Data: []byte(t[2])}
		fts[i] = &faceTask{s.vec, make(chan task.TaskResult, 1), req}
	}
	for _, ft := range fts {
		s.rts.AddTask(ft)
	}
	mp := make(map[int]interface{})
	t := time.Now()
	for i, ft := range fts {
		select {
		case r := <-ft.ch:
			if xb := r.Result().([]float32); len(xb) != 512 {
				mp[i] = nil
				rids = append(rids, ts[i][0])
				s.log.Debugf("uid = %s: vector not 512: %v\n", ts[i][0], len(xb))
			} else {
				xbs = append(xbs, xb...)
			}
		}
	}
	{
		s.log.Debugf("face process: %v\n", time.Now().Sub(t))
	}
	for j, t := range ts {
		if _, ok := mp[j]; ok {
			continue
		}
		arg := make([]interface{}, len(attrs))
		for i, attr := range attrs {
			v, rs, err := appendSlice(iargs[i], attr.Type, t[i])
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			arg[i] = v
			iargs[i] = rs
			if i == 1 {
				xids = append(xids, int64(v.(uint64)))
			}
		}
		cargs = append(cargs, arg)
	}
	return rids, xbs, xids, iargs, cargs, nil
}

func (s *server) convertWithVector(ts [][]string, attrs []metadata.Attribute) ([]float32, []int64, []interface{}, [][]interface{}, error) {
	xbs := make([]float32, 0, len(ts)*512)
	xids := make([]int64, 0, len(ts))
	iargs := make([]interface{}, len(attrs))
	cargs := make([][]interface{}, 0, len(ts))
	{
		for i := range attrs {
			iargs[i] = newSlice(attrs[i].Type, len(ts))
		}
	}
	for _, t := range ts {
		var vec []float32

		if err := json.Unmarshal([]byte(t[len(t)-1]), &vec); err != nil {
			return nil, nil, nil, nil, err
		}
		xbs = append(xbs, vec...)
	}
	for _, t := range ts {
		arg := make([]interface{}, len(attrs))
		for i, attr := range attrs {
			v, rs, err := appendSlice(iargs[i], attr.Type, t[i])
			if err != nil {
				return nil, nil, nil, nil, err
			}
			arg[i] = v
			iargs[i] = rs
			if i == 1 {
				xids = append(xids, int64(v.(uint64)))
			}
		}
		cargs = append(cargs, arg)
	}
	return xbs, xids, iargs, cargs, nil
}

func (s *server) extractParameters(ctx *fasthttp.RequestCtx) (string, []float32, error) {
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
	qr, err := s.getSqlQuery(mp)
	if err != nil {
		return "", nil, err
	}
	vec, err := s.vec.GetVector(fs)
	if err != nil {
		return "", nil, err
	}
	return qr, vec, nil
}

func (s *server) extractParametersWithVector(ctx *fasthttp.RequestCtx) (string, []float32, error) {
	var typ string
	var body []byte
	var vec []float32
	var mp map[string]interface{}

	form, err := ctx.MultipartForm()
	if err != nil {
		return "", nil, err
	}
	for _, v := range form.File {
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
			}
		}
	}
	qr, err := s.getSqlQuery(mp)
	if err != nil {
		return "", nil, err
	}
	if err := json.Unmarshal(body, &vec); err != nil {
		return "", nil, err
	}
	return qr, vec, nil
}

func (s *server) getSqlQuery(mp map[string]interface{}) (string, error) {
	return getString("query", mp)
}

func appendSlice(vs interface{}, typ uint32, s string) (interface{}, interface{}, error) {
	switch typ {
	case types.T_int8:
		rs := vs.([]int8)
		v, err := strconv.ParseInt(s, 0, 8)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, int8(v))
		return v, rs, nil
	case types.T_int16:
		rs := vs.([]int16)
		v, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, int16(v))
		return v, rs, nil
	case types.T_int32:
		rs := vs.([]int32)
		v, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, int32(v))
		return v, rs, nil
	case types.T_int64:
		rs := vs.([]int64)
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, v)
		return v, rs, nil
	case types.T_uint8:
		rs := vs.([]uint8)
		v, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, uint8(v))
		return v, rs, nil
	case types.T_uint16:
		rs := vs.([]uint16)
		v, err := strconv.ParseUint(s, 0, 16)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, uint16(v))
		return v, rs, nil
	case types.T_uint32:
		rs := vs.([]uint32)
		v, err := strconv.ParseUint(s, 0, 32)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, uint32(v))
		return v, rs, nil
	case types.T_uint64:
		rs := vs.([]uint64)
		v, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, v)
		return v, rs, nil
	case types.T_float32:
		rs := vs.([]float32)
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, float32(v))
		return v, rs, nil
	case types.T_float64:
		rs := vs.([]float64)
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, nil, err
		}
		rs = append(rs, v)
		return v, rs, nil
	case types.T_string:
		rs := vs.([]string)
		rs = append(rs, s)
		return s, rs, nil
	case types.T_timestamp:
		rs := vs.([]int64)
		v, err := value.ParseTimestamp(s)
		if err != nil {
			return nil, nil, err
		}
		t := value.MustBeTimestamp(v).Unix()
		rs = append(rs, t)
		return t - 8*3600, rs, nil
	}
	return nil, nil, nil
}

func newSlice(typ uint32, size int) interface{} {
	switch typ {
	case types.T_int8:
		return make([]int8, 0, size)
	case types.T_int16:
		return make([]int16, 0, size)
	case types.T_int32:
		return make([]int32, 0, size)
	case types.T_int64:
		return make([]int64, 0, size)
	case types.T_uint8:
		return make([]uint8, 0, size)
	case types.T_uint16:
		return make([]uint16, 0, size)
	case types.T_uint32:
		return make([]uint32, 0, size)
	case types.T_uint64:
		return make([]uint64, 0, size)
	case types.T_float32:
		return make([]float32, 0, size)
	case types.T_float64:
		return make([]float64, 0, size)
	case types.T_string:
		return make([]string, 0, size)
	case types.T_timestamp:
		return make([]int64, 0, size)
	}
	return nil
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
