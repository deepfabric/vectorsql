package not

import (
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func New() *not {
	return &not{}
}

func (n *not) Rewrite(e extend.Extend) extend.Extend {
	return n.rewriteNot(e)
}

func (n *not) rewriteNot(e extend.Extend) extend.Extend {
	switch v := e.(type) {
	case *extend.ParenExtend:
		return &extend.ParenExtend{n.rewriteNot(v.E)}
	case *extend.UnaryExtend:
		return n.rewriteNotUnary(v)
	case *extend.BinaryExtend:
		return n.rewriteNotBinary(v)
	}
	return e
}

func (n *not) rewriteNotUnary(e *extend.UnaryExtend) extend.Extend {
	if e.Op != overload.Not {
		return e
	}
	return n.negation(e.E, false)
}

func (n *not) rewriteNotBinary(e *extend.BinaryExtend) extend.Extend {
	if e.Op == overload.Or || e.Op == overload.And {
		return &extend.BinaryExtend{
			Op:    e.Op,
			Left:  n.rewriteNot(e.Left),
			Right: n.rewriteNot(e.Right),
		}
	}
	return e
}

func (n *not) negation(e extend.Extend, isParen bool) extend.Extend {
	switch v := e.(type) {
	case *value.Bool:
		if value.MustBeBool(v) {
			return value.NewBool(false)
		}
		return value.NewBool(true)
	case *extend.UnaryExtend:
		return n.negationUnary(v)
	case *extend.ParenExtend:
		return &extend.ParenExtend{n.negation(v.E, true)}
	case *extend.BinaryExtend:
		if !isParen && (v.Op == overload.And || v.Op == overload.Or) {
			v.Left = n.negation(v.Left, isParen)
			return v
		}
		return n.negationBinary(v, isParen)
	}
	return e
}

func (n *not) negationUnary(e *extend.UnaryExtend) extend.Extend {
	if e.Op == overload.Not {
		return e.E
	}
	return e
}

func (n *not) negationBinary(e *extend.BinaryExtend, isParen bool) extend.Extend {
	switch e.Op {
	case overload.EQ:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.NE,
		}
	case overload.NE:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.EQ,
		}
	case overload.LT:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.GE,
		}
	case overload.GT:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.LE,
		}
	case overload.LE:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.GT,
		}
	case overload.GE:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.LE,
		}
	case overload.Or:
		return &extend.BinaryExtend{
			Op:    overload.And,
			Left:  n.negation(e.Left, isParen),
			Right: n.negation(e.Right, isParen),
		}
	case overload.And:
		return &extend.BinaryExtend{
			Op:    overload.Or,
			Left:  n.negation(e.Left, isParen),
			Right: n.negation(e.Right, isParen),
		}
	case overload.Like:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.NotLike,
		}
	case overload.NotLike:
		return &extend.BinaryExtend{
			Left:  e.Left,
			Right: e.Right,
			Op:    overload.Like,
		}
	}
	return e
}
