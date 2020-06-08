package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/deepfabric/thinkkv/pkg/engine/pb"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/container/relation"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
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
	if err := stg.NewRelation("people", storage.MetaData{
		IsEvent: false,
		Attrs:   []string{"gender", "age"},
		Types:   []int32{types.T_uint8, types.T_uint8},
	}); err != nil {
		log.Fatal(err)
	}
	r, err := stg.Relation("people")
	if err != nil {
		log.Fatal(err)
	}
	if err := inject(r, 10); err != nil {
		log.Fatal(err)
	}
	stg.Close()
}

//var areas = []string{"上海", "北京", "成都", "苏州", "日本", "非洲"}

func inject(r relation.Relation, n int) error {
	ts := make([]map[string]interface{}, n)
	for i := 0; i < n; i++ {
		ts[i] = make(map[string]interface{})
	}
	for i := 0; i < n; i++ {
		if i%1000000 == 0 {
			fmt.Printf("process %v\n", i/1000000)
		}
		ts[i]["seq"] = uint32(i)
		if i%2 == 0 {
			ts[i]["gender"] = uint8(0)
		} else {
			ts[i]["gender"] = uint8(1)
		}
		ts[i]["age"] = uint8(rand.Intn(100))
	}
	if err := r.AddTuplesByJson(ts); err != nil {
		return err
	}
	return nil
}
