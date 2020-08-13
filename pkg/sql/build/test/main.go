package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/sql/build"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func main() {
	db := pb.New("test.db", nil, 0, false, false)
	defer db.Close()
	stg := storage.New(db, lru.New(10), cache.New(1<<20))
	{
		var attrs []metadata.Attribute

		attrs = append(attrs, metadata.Attribute{types.T_uint64, "uid"})
		attrs = append(attrs, metadata.Attribute{types.T_uint8, "age"})
		attrs = append(attrs, metadata.Attribute{types.T_string, "pic"})
		if err := stg.NewRelation("user_item", metadata.Metadata{
			IsE:   false,
			Attrs: attrs,
		}); err != nil {
			log.Fatal(err)
		}
	}
	o, err := build.New(os.Args[1], context.New(nil, stg), stg).Build()
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
