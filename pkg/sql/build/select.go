package build

import (
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/rewrite/not"
)

func (b *build) buildSelect(n *tree.SelectClause) (extend.Extend, *tree.SelectClause, error) {
	if err := b.buildFrom(n.From); err != nil {
		return nil, nil, err
	}
	e, err := b.buildWhere(n.Where)
	if err != nil {
		return nil, nil, err
	}
	return not.New().Rewrite(e), n, nil
}
