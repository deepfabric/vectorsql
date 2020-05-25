package vector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

func New(url string) *vector {
	return &vector{url}
}

func (v *vector) GetVector(images map[string][]byte) ([]float32, error) {
	var body bytes.Buffer
	var mp map[string][]string

	req, err := NewRequest(v.url, images)
	if err != nil {
		return nil, nil
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if _, err := body.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body.Bytes(), &mp); err != nil {
		return nil, err
	}
	return mean(mp), nil
}

func NewRequest(url string, images map[string][]byte) (*http.Request, error) {
	var body bytes.Buffer

	w := multipart.NewWriter(&body)
	for k, v := range images {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, k, k))
		h.Set("Content-Type", "image/jpg")
		p, _ := w.CreatePart(h)
		io.Copy(p, bytes.NewReader(v))
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
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
