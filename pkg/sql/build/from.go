package build

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/sql/tree"
)

func (b *build) buildFrom(n *tree.From) error {
	for i := range n.Tables {
		switch t := n.Tables[i].(type) {
		case *tree.AliasedTable:
			if err := b.buildAliasedTable(t); err != nil {
				return err
			}
		default:
			return fmt.Errorf("illegal table '%s'", n.Tables[i])
		}
	}
	return nil
}

func (b *build) buildAliasedTable(n *tree.AliasedTable) error {
	switch t := n.Tbl.(type) {
	case *tree.Subquery:
		return fmt.Errorf("'%s' unsupport now", n)
	case *tree.TableName:
		return b.buildTableName(t)
	default:
		return fmt.Errorf("illegal aliased table '%s'", n)
	}
}

func (b *build) buildTableName(n *tree.TableName) error {
	name, err := b.buildExprColumn(n.N)
	if err != nil {
		return err
	}
	if name != DefaultTable {
		return fmt.Errorf("table '%s' not exist", name)
	}
	return nil
}
