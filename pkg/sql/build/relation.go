package build

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
)

func (b *build) buildRelation(n tree.RelationStatement) (extend.Extend, *tree.SelectClause, error) {
	switch t := n.(type) {
	case *tree.AliasedTable:
		return nil, nil, fmt.Errorf("'%s' unsupport now", n)
	case *tree.JoinClause:
		return nil, nil, fmt.Errorf("'%s' unsupport now", n)
	case *tree.UnionClause:
		return nil, nil, fmt.Errorf("'%s' unsupport now", n)
	case *tree.SelectClause:
		return b.buildSelect(t)
	case *tree.AliasedSelect:
		return nil, nil, fmt.Errorf("'%s' unsupport now", n)
	default:
		return nil, nil, fmt.Errorf("unknown relation statement '%s'", n)
	}
}
