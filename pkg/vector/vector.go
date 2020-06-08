package vector

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/valyala/fasthttp"
)

func New(url string) *vector {
	return &vector{url}
}

func (v *vector) GetVector(fs map[string]*request.Part) ([]float32, error) {
	var mp map[string][]string
	var resp fasthttp.Response

	req, err := request.NewRequest(v.url, fs)
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
	return mean(mp), nil
}

func mean(mp map[string][]string) []float32 {
	var y float64
	var xs []float32

	cnt := float32(0)
	for _, v := range mp {
		if len(v) == 0 {
			continue
		}
		if ys := strings2Floats(v); len(ys) > 0 {
			cnt++
			if len(xs) == 0 {
				xs = ys
			} else {
				sum(xs, ys)
			}
		}
	}
	for i, x := range xs {
		xs[i] = x / cnt
		y += math.Pow(float64(x), 2)
	}
	y = math.Sqrt(y)
	for i, x := range xs {
		xs[i] = float32(float64(x) / y)
	}
	return xs
}

func sum(xs, ys []float32) {
	if len(xs) != len(ys) {
		return
	}
	for i, _ := range xs {
		xs[i] += ys[i]
	}
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
