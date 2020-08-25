package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/valyala/fasthttp"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: queryer ip:port sql vector\n")
		return
	}
	log := logger.New(os.Stderr, "queryer:")
	fs := make(map[string]*request.Part)
	{
		fs["a"] = &request.Part{Data: []byte(os.Args[3])}
	}
	{
		fs["query"] = &request.Part{"application/json", []byte(fmt.Sprintf("{\"query\":\"%s\"}", os.Args[2]))}
	}
	req, err := request.NewRequest(fmt.Sprintf("http://%s/queryWithVector", os.Args[1]), fs)
	if err != nil {
		log.Fatal(err)
	}
	var resp fasthttp.Response
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(&resp)
	tm := time.Now()
	if err := fasthttp.Do(req, &resp); err != nil {
		log.Fatal(err)
	}
	ts, err := csv.NewReader(bytes.NewReader(resp.Body())).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	{
		fmt.Printf("process: %v\n", time.Now().Sub(tm))
	}
	for i, t := range ts {
		n := len(t) - 1
		if err := ioutil.WriteFile(fmt.Sprintf("%v.jpg", i), []byte(t[n]), os.FileMode(0664)); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", t[:n])
	}
}
