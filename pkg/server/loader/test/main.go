package main

import (
	"log"
	"os"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/server/loader"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
)

func main() {
	db := pb.New("test.db", nil, 0, false, false)
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatal(err)
	}
	srv := loader.New(8888, logger.New(os.Stderr, "loader:"), cli, storage.New(db, lru.New(10), cache.New(1<<20)))
	srv.Run()
}
