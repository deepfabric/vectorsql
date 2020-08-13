package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/config"
	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/routines"
	"github.com/deepfabric/vectorsql/pkg/server"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/vector"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
)

var (
	path = flag.String("c", "./config/server.toml", "config file path")
)

func main() {
	var cfg config.Config

	flag.Parse()
	if _, err := toml.DecodeFile(*path, &cfg); err != nil {
		fmt.Printf("Failed to parse config file '%s': %v\n", *path, err)
		os.Exit(0)
	}
	log := logger.New(os.Stderr, cfg.LogConfig.Prefix)
	switch strings.ToLower(cfg.LogConfig.Level) {
	case "debug":
		log.SetLevel(logger.DEBUG)
	case "info":
		log.SetLevel(logger.INFO)
	case "warn":
		log.SetLevel(logger.WARN)
	case "error":
		log.SetLevel(logger.ERROR)
	case "fatal":
		log.SetLevel(logger.FATAL)
	case "panic":
		log.SetLevel(logger.PANIC)
	}
	db := pb.New(cfg.Db, nil, 0, false, false)
	defer db.Close()
	cli, err := client.New(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	vec := vector.New(cfg.Url)
	b := bv.New(cfg.Addrs, log)
	stg := storage.New(db, lru.New(100), cache.New(cfg.CacheSize))
	defer stg.Close()
	scfg := &server.Config{
		B:   b,
		Log: log,
		Cli: cli,
		Vec: vec,
		Stg: stg,
		Ctx: context.New(cli, stg),
		Rts: routines.New(cfg.Routines),
	}
	srv := server.New(cfg.Port, scfg)
	srv.Run()
}
