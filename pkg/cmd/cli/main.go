package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/deepfabric/vectorsql/pkg/vector"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	images := make(map[string][]byte)
	images["a"] = data
	req, err := vector.NewRequest("http://127.0.0.1:8888/query", images)
	if err != nil {
		log.Fatal(err)
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	{
		var body bytes.Buffer

		if _, err := body.ReadFrom(resp.Body); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("body: %s\n", string(body.Bytes()))
	}
}
