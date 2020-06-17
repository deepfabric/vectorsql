package main

import (
	"log"

	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/vm/util/encoding"
)

var db engine.DB

func loadId() uint32 {
	v, err := db.Get([]byte("_id.people"))
	switch {
	case err == nil:
		return encoding.DecodeUint32(v)
	case err == engine.NotExist:
		return uint32(0)
	default:
		log.Fatalf("Failet to load id: %v\n", err)
	}
	return 0
}

func storeId(id uint32) error {
	return db.Set([]byte("_id.people"), encoding.EncodeUint32(id))
}
