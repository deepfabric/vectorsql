package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/valyala/fasthttp"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: loader ip:port table-name images-directory csv...\n")
		return
	}
	dir := os.Args[3]
	name := os.Args[2]
	log := logger.New(os.Stderr, "loader:")
	for i, j := 4, len(os.Args); i < j; i++ {
		{
			var req fasthttp.Request
			var resp fasthttp.Response

			req.Header.SetMethod("POST")
			t := time.Now()
			req.SetBody(loadCsv(os.Args[i], dir, log))
			log.Debugf("loadCsv process %v\n", time.Now().Sub(t))
			req.Header.SetMethod("POST")
			req.SetRequestURI(fmt.Sprintf("http://%s/insert?name=%s", os.Args[1], name))
			t = time.Now()
			if err := fasthttp.Do(&req, &resp); err != nil {
				log.Fatal(err)
			}
			log.Debugf("%s's process %v: %s\n", os.Args[i], time.Now().Sub(t), string(resp.Body()))
		}
	}
}

func loadCsv(name string, dir string, log logger.Log) []byte {
	var buf bytes.Buffer

	fp, err := os.Open(name)
	if err != nil {
		log.Fatalf("Failed to open csv file '%s': %v\n", name, err)
	}
	defer fp.Close()
	vs, err := csv.NewReader(fp).ReadAll()
	if err != nil {
		log.Fatalf("Failed to load csv file '%s': %v\n", name, err)
	}
	for i, j := 0, len(vs); i < j; i++ {
		data, err := ioutil.ReadFile(path.Join(dir, vs[i][2]))
		if err != nil {
			log.Fatalf("Failed to load image file '%s': %v\n", vs[i][2], err)
		}
		vs[i][2] = string(data)
	}
	if err := csv.NewWriter(bufio.NewWriter(&buf)).WriteAll(vs); err != nil {
		log.Fatalf("Failed to generate request body: %v\n", err)
	}
	return buf.Bytes()
}
