package op

import (
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/filter"
)

type Top struct {
	Num int
	IsF bool
}

type OP struct {
	T  *Top
	N  *tree.Select
	Cf filter.Filter
	If filter.Filter
}
