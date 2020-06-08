package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/sql/build"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/storage/cache/table/mem"
	"github.com/deepfabric/vectorsql/pkg/storage/user"
	"github.com/deepfabric/vectorsql/pkg/vm/bv"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	lc, err := lru.New(10000)
	if err != nil {
		log.Fatal(err)
	}
	db := pb.New("test.db", nil, 1024*1024*1024, false, false)
	defer db.Close()
	stg := user.New(10, db, mem.New(), lc)
	o, err := build.New(os.Args[1], context.New(nil), stg).Build()
	if err != nil {
		log.Fatal(err)
	}
	{
		fmt.Printf("\tT: %v\n", o.T)
	}
	{
		fmt.Printf("\tN: %v\n", o.N)
	}
	{
		fmt.Printf("\tCF: %v\n", o.Cf)
	}
	{
		fmt.Printf("\tIF: %v\n", o.If)
	}
	r, err := o.Result(bv.New(), cli)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", r)
}
