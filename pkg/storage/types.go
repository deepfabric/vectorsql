package storage

import (
	"sync"

	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/storage/index"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/pilosa/pilosa/roaring"
)

type Storage interface {
	Close() error
	Relation(string) (Relation, error)
	NewRelation(string, metadata.Metadata) error
}

type Relation interface {
	Destroy() error

	IsEvent() bool

	Metadata() metadata.Metadata

	AddTuples([]interface{}) error

	Eq(string, value.Value) (*roaring.Bitmap, error)
	Ne(string, value.Value) (*roaring.Bitmap, error)
	Lt(string, value.Value) (*roaring.Bitmap, error)
	Le(string, value.Value) (*roaring.Bitmap, error)
	Gt(string, value.Value) (*roaring.Bitmap, error)
	Ge(string, value.Value) (*roaring.Bitmap, error)
}

type storage struct {
	sync.RWMutex
	rc lru.LRU // cache for relation
	db engine.DB
	lc cache.Cache
}

type relation struct {
	sync.RWMutex
	id  string
	db  engine.DB
	idx index.Index
	md  metadata.Metadata
}
