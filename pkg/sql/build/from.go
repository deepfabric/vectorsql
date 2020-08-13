package build

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/sql/tree"
)

func (b *build) buildFrom(n *tree.From) (string, error) {
	for i := range n.Tables {
		switch t := n.Tables[i].(type) {
		case *tree.AliasedTable:
			return b.buildAliasedTable(t)
		default:
			return "", fmt.Errorf("illegal table '%s'", n.Tables[i])
		}
	}
	return "", nil
}

func (b *build) buildAliasedTable(n *tree.AliasedTable) (string, error) {
	switch t := n.Tbl.(type) {
	case *tree.Subquery:
		return "", fmt.Errorf("'%s' unsupport now", n)
	case *tree.TableName:
		return b.buildTableName(t)
	default:
		return "", fmt.Errorf("illegal aliased table '%s'", n)
	}
}

func (b *build) buildTableName(n *tree.TableName) (string, error) {
	return b.buildExprColumn(n.N)
}
