package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/deepfabric/vectorsql/pkg/request"
	"github.com/valyala/fasthttp"
)

func main() {
	if len(os.Args[2]) < 3 {
		fmt.Printf("Usage: cli query images[...]")
		return
	}
	fs := make(map[string]*request.Part)
	for i, j := 2, len(os.Args); i < j; i++ {
		data, err := ioutil.ReadFile(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}
		fs[os.Args[i]] = &request.Part{"image/jepg", data}
	}
	{
		fs["query"] = &request.Part{"application/json", []byte(fmt.Sprintf("{\"query\":\"%s\"}", os.Args[1]))}
	}
	req, err := request.NewRequest("http://127.0.0.1:8888/query", fs)
	if err != nil {
		log.Fatal(err)
	}
	var resp fasthttp.Response
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(&resp)
	if err := fasthttp.Do(req, &resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("body: %s\n", string(resp.Body()))
}
