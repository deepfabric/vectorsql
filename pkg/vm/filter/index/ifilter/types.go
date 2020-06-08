package ifilter

import (
	"github.com/deepfabric/vectorsql/pkg/vm/container/relation"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

const (
	EQ = iota
	NE
	LT
	LE
	GT
	GE
)

type Condition struct {
	Op   int // eq, ne, lt, le, gt, ge
	Name string
	Val  value.Value
}

type filter struct {
	cs []*Condition
	r  relation.Relation
}
