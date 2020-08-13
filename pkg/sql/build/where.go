package build

import (
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
)

func (b *build) buildWhere(n *tree.Where, id string) (extend.Extend, error) {
	if n == nil {
		return nil, nil
	}
	return b.buildExpr(n.E, id)
}
