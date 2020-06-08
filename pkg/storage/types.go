package storage

import (
	"errors"
	"sync"

	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/vm/container/relation"
	lru "github.com/hashicorp/golang-lru"
)

var (
	ID       = "id"
	SEQ      = "seq" // row number
	NotExist = errors.New("Not Exist")
)

type Storage interface {
	Close() error
	NewRelation(string, MetaData) error
	Relation(string) (relation.Relation, error)
}

type MetaData struct {
	IsEvent bool
	Types   []int32  // type of attributes
	Attrs   []string // attributes
}

type storage struct {
	mcpu int
	db   engine.DB
	bc   *lru.Cache // cache for bitmap
	rc   *lru.Cache // cache for relation
}

// M.id -> metadata
// id.attr's name        	-> bitmap -- null bitmap
// id.attr's name.I 		-> bitmap -- bsi, bitmap
// id.attr's name.U 		-> bitmap -- ubsi bitmap
// id.attr's name.value  	-> bitmap -- bool, string bitmap
// id.attr's name.I.value  	-> bitmap -- int8 bitmap
// id.attr's name.U.value  	-> bitmap -- uint8 bitmap
type index struct {
	sync.RWMutex
	isE  bool
	mcpu int
	id   string // relation id
	db   engine.DB
	lc   *lru.Cache // cache for bitmap
	mp   map[string]int32
}
