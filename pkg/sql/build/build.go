package build

import (
	"github.com/deepfabric/vectorsql/pkg/sql/parser"
	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/storage"
	"github.com/deepfabric/vectorsql/pkg/vm/context"
	"github.com/deepfabric/vectorsql/pkg/vm/op"
	"github.com/deepfabric/vectorsql/pkg/vm/opt"
)

func New(sql string, c context.Context, stg storage.Storage) *build {
	return &build{
		c:   c,
		sql: sql,
		stg: stg,
	}
}

func (b *build) Build() (*op.OP, error) {
	n, err := parser.Parse(b.sql)
	if err != nil {
		return nil, err
	}
	return b.buildStatement(n)
}

func (b *build) buildStatement(n *tree.Select) (*op.OP, error) {
	var o op.OP

	id, e, sc, err := b.buildRelation(n.Relation)
	if err != nil {
		return nil, err
	}
	o.N = n
	sc.Where = nil
	n.Relation = sc
	if e != nil {
		c, i, err := opt.New(b.c, b.stg).Optimize(e, id)
		if err != nil {
			return nil, err
		}
		o.Cf = c
		o.If = i
	}
	if n.Order != nil {
		t, err := b.buildOrder(n, n.Order)
		if err != nil {
			return nil, err
		}
		o.T = t
	}
	return &o, nil
}

func (b *build) buildOrder(n *tree.Select, ord tree.OrderStatement) (*op.Top, error) {
	switch t := ord.(type) {
	case *tree.Top:
		n.Order = nil
		return b.buildTop(t)
	case *tree.Ftop:
		n.Order = nil
		return b.buildFtop(t)
	}
	return nil, nil
}

func (b *build) buildTop(ord *tree.Top) (*op.Top, error) {
	if n, err := b.buildExprIntConstant(ord.N); err != nil {
		return nil, err
	} else {
		return &op.Top{Num: int(n), IsF: false}, nil
	}
}

func (b *build) buildFtop(ord *tree.Ftop) (*op.Top, error) {
	if n, err := b.buildExprIntConstant(ord.N); err != nil {
		return nil, err
	} else {
		return &op.Top{Num: int(n), IsF: true}, nil
	}
}
