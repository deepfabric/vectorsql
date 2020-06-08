package relation

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/bsi"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

type Relation interface {
	Destroy() error

	String() string

	IsEvent() bool

	IdBitmap() (bsi.Bsi, error)

	AddTuplesByJson([]map[string]interface{}) error

	Eq(string, value.Value) (*roaring.Bitmap, error)
	Ne(string, value.Value) (*roaring.Bitmap, error)
	Lt(string, value.Value) (*roaring.Bitmap, error)
	Le(string, value.Value) (*roaring.Bitmap, error)
	Gt(string, value.Value) (*roaring.Bitmap, error)
	Ge(string, value.Value) (*roaring.Bitmap, error)
}
