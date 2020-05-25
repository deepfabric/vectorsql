package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/deepfabric/vectorsql/pkg/logger"
	"github.com/deepfabric/vectorsql/pkg/server/query"
	"github.com/deepfabric/vectorsql/pkg/vector"
)

func main() {
	srv := query.New(8888, logger.New(os.Stderr, "vectorsql"), vector.New("http://172.19.0.17:6933/face_emb"))
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
