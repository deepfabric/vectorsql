package build

import (
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/rewrite/not"
)

func (b *build) buildSelect(n *tree.SelectClause) (string, extend.Extend, *tree.SelectClause, error) {
	id, err := b.buildFrom(n.From)
	if err != nil {
		return "", nil, nil, err
	}
	e, err := b.buildWhere(n.Where, id)
	if err != nil {
		return "", nil, nil, err
	}
	return id, not.New().Rewrite(e), n, nil
}
