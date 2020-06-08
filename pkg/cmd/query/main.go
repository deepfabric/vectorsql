package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/server/query"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatal(err)
	}
	bc, err := lru.New(100)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := lru.New(10000)
	if err != nil {
		log.Fatal(err)
	}
	db := pb.New("test.db", nil, 1024*1024*1024, false, false)
	stg := storage.New(2, db, bc, rc)
	srv := query.New(8888, 2, bv.New(), logger.New(os.Stderr, "vectorsql"), cli, vector.New("http://172.19.0.17:6933/face_emb"), stg)
	{
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-ch
			srv.Stop()
			os.Exit(0)
		}()
	}
	srv.Run()
}
