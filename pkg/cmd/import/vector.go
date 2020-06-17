package main

import (
	"encoding/json"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/valyala/fasthttp"
)

const URL = "http://172.19.0.17:6930/face_emb"

func getVector(fs map[string]*request.Part) ([][]float32, error) {
	var rs [][]float32
	var mp map[string][]string
	var resp fasthttp.Response

	req, err := request.NewRequest(URL, fs)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(&resp)
	if err := fasthttp.Do(req, &resp); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &mp); err != nil {
		return nil, err
	}
	for _, v := range mp {
		rs = append(rs, strings2Floats(v))
	}
	return rs, nil
}

func strings2Floats(xs []string) []float32 {
	var rs []float32

	for _, x := range xs {
		y, err := strconv.ParseFloat(x, 32)
		if err != nil {
			return nil
		}
		rs = append(rs, float32(y))
	}
	return rs
}
