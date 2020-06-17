package main

import (
	"log"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	lru "github.com/hashicorp/golang-lru"
)

var stg storage.Storage

func init() {
	bc, err := lru.New(100)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := lru.New(10000)
	if err != nil {
		log.Fatal(err)
	}
	db = pb.New("idx.db", nil, 1024*1024*1024, false, false)
	stg = storage.New(48, db, bc, rc)
	if err := stg.NewRelation("people", storage.MetaData{
		IsEvent: false,
		Attrs:   []string{"seq", "sex", "age", "area"},
		Types:   []int32{types.T_uint32, types.T_uint8, types.T_uint8, types.T_string},
	}); err != nil {
		log.Fatal(err)
	}
}

func inject(ts []map[string]interface{}) error {
	r, err := stg.Relation("people")
	if err != nil {
		return err
	}
	return r.AddTuplesByJson(ts)
}
