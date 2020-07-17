package index

import (
	"github.com/deepfabric/thinkkv/pkg/engine"
	"github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/storage/cache"
	"github.com/deepfabric/vectorsql/pkg/storage/metadata"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/pilosa/pilosa/roaring"
)

const (
	ID  = "id"
	SEQ = "seq" // seq强制在第一列
)

const (
	Ffrac = 1000   // float fraction
	Dfrac = 100000 // double fraction
)

type Index interface {
	IdMap() (bsi.Bsi, error)

	AddTuples([]interface{}) error

	Eq(string, value.Value) (*roaring.Bitmap, error)
	Ne(string, value.Value) (*roaring.Bitmap, error)
	Lt(string, value.Value) (*roaring.Bitmap, error)
	Le(string, value.Value) (*roaring.Bitmap, error)
	Gt(string, value.Value) (*roaring.Bitmap, error)
	Ge(string, value.Value) (*roaring.Bitmap, error)
}

// id.attr's name.I       	-> bitmap -- bsi, bitmap
// id.attr's name.U      	-> bitmap -- ubsi bitmap
type index struct {
	isE   bool
	id    string // uid.database.table
	db    engine.DB
	lc    cache.Cache
	attrs []metadata.Attribute
}
