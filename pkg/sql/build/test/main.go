package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/sql/build"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	bc, err := lru.New(100)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := lru.New(10000)
	if err != nil {
		log.Fatal(err)
	}
	db := pb.New("test.db", nil, 1024*1024*1024, false, false)
	stg := storage.New(10, db, bc, rc)
	o, err := build.New(os.Args[1], context.New(nil), stg).Build()
	if err != nil {
		log.Fatal(err)
	}
	{
		fmt.Printf("T: %v\n", o.T)
	}
	{
		fmt.Printf("N: %v\n", o.N)
	}
	{
		fmt.Printf("CF: %v\n", o.Cf)
	}
	{
		fmt.Printf("IF: %v\n", o.If)
	}
}
