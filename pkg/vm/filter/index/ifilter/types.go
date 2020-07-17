package ifilter

import (
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/pilosa/pilosa/roaring"
)

const (
	EQ = iota
	NE
	LT
	LE
	GT
	GE
)

type Filter interface {
	String() string
	Bitmap() (*roaring.Bitmap, error)
}

type Condition struct {
	Op   int // eq, ne, lt, le, gt, ge
	Name string
	Val  value.Value
}

type filter struct {
	cs []*Condition
	r  storage.Relation
}
