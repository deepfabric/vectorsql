package build

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/deepfabric/vectorsql/pkg/sql/tree"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func (b *build) buildExpr(n tree.ExprStatement) (extend.Extend, error) {
	switch e := n.(type) {
	case *tree.Value:
		return e.E, nil
	case *tree.OrExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Or, left, right}, nil
	case *tree.AndExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.And, left, right}, nil
	case *tree.NotExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		return &extend.UnaryExtend{overload.Not, ext}, nil
	case *tree.DivExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Div, left, right}, nil
	case *tree.ModExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Mod, left, right}, nil
	case *tree.MultExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Mult, left, right}, nil
	case *tree.PlusExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Plus, left, right}, nil
	case *tree.MinusExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.Minus, left, right}, nil
	case *tree.UnaryMinusExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		return &extend.UnaryExtend{overload.UnaryMinus, ext}, nil
	case *tree.EqExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.EQ, left, right}, nil
	case *tree.NeExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.NE, left, right}, nil
	case *tree.LtExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.LT, left, right}, nil
	case *tree.LeExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.LE, left, right}, nil
	case *tree.GtExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.GT, left, right}, nil
	case *tree.GeExpr:
		left, err := b.buildExpr(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.GE, left, right}, nil
	case *tree.Subquery:
		return nil, errors.New("subquery not support now")
	case *tree.BetweenExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		to, err := b.buildExpr(e.To)
		if err != nil {
			return nil, err
		}
		from, err := b.buildExpr(e.From)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{
			Op: overload.And,
			Left: &extend.BinaryExtend{
				Op:    overload.GE,
				Left:  ext,
				Right: from,
			},
			Right: &extend.BinaryExtend{
				Op:    overload.LE,
				Left:  ext,
				Right: to,
			},
		}, nil
	case *tree.NotBetweenExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		to, err := b.buildExpr(e.To)
		if err != nil {
			return nil, err
		}
		from, err := b.buildExpr(e.From)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{
			Op: overload.Or,
			Left: &extend.BinaryExtend{
				Op:    overload.LT,
				Left:  ext,
				Right: from,
			},
			Right: &extend.BinaryExtend{
				Op:    overload.GT,
				Left:  ext,
				Right: to,
			},
		}, nil
	case *tree.IsNullExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.EQ, ext, value.ConstNull}, nil
	case *tree.IsNotNullExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		return &extend.BinaryExtend{overload.NE, ext, value.ConstNull}, nil
	case *tree.FuncExpr:
		return b.buildExprFunc(e)
	case *tree.ParenExpr:
		ext, err := b.buildExpr(e.E)
		if err != nil {
			return nil, err
		}
		return &extend.ParenExtend{ext}, nil
	case tree.ColunmNameList:
		name, err := b.buildExprColumn(e)
		if err != nil {
			return nil, err
		}
		return &extend.Attribute{name}, nil
	default:
		return nil, fmt.Errorf("unexpected expression '%s'", n)
	}
}

func (b *build) buildExprFunc(n *tree.FuncExpr) (extend.Extend, error) {
	n.Name = strings.ToLower(n.Name)
	if _, ok := AggFuncs[n.Name]; ok {
		return nil, fmt.Errorf("unexpected aggregate expression '%s' in where clause", n)
	}
	op, ok := ExtendFuncs[n.Name]
	if !ok {
		return nil, fmt.Errorf("unimplemented functions: %s", n.Name)
	}
	switch overload.OperatorType(op) {
	case overload.Unary:
		if len(n.Es) < 1 {
			return nil, fmt.Errorf("not enough arguments in call to '%s'", n.Name)
		}
		e, err := b.buildExpr(n.Es[0])
		if err != nil {
			return nil, err
		}
		return &extend.UnaryExtend{
			E:  e,
			Op: op,
		}, nil
	case overload.Binary:
		if len(n.Es) < 2 {
			return nil, fmt.Errorf("not enough arguments in call to '%s'", n.Name)
		}
		left, err := b.buildExpr(n.Es[0])
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(n.Es[1])
		if err != nil {
			return nil, err
		}
		e := &extend.BinaryExtend{
			Op:    op,
			Left:  left,
			Right: right,
		}
		return reduce(e), nil
	default: // multi
		var args []extend.Extend

		for i := range n.Es {
			arg, err := b.buildExpr(n.Es[i])
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		return &extend.MultiExtend{
			Op:   op,
			Args: args,
		}, nil
	}
}

func (b *build) buildExprColumn(ns tree.ColunmNameList) (string, error) {
	var name string

	for i := range ns {
		if i > 0 {
			name += "."
		}
		name += string(ns[i].Path)
		if ns[i].Index != nil {
			if idx, err := b.buildExprIntConstant(ns[i].Index); err != nil {
				return "", err
			} else {
				name += fmt.Sprintf("._%v", idx)
			}
		}
	}
	return name, nil
}

func (b *build) buildExprIntConstant(n tree.ExprStatement) (int64, error) {
	switch e := n.(type) {
	case *tree.Value:
		if i, err := value.GetInt(e.E); err != nil {
			return 0, err
		} else {
			return int64(i), nil
		}
	case *tree.ModExpr:
		x, err := b.buildExprIntConstant(e.Left)
		if err != nil {
			return 0, err
		}
		y, err := b.buildExprIntConstant(e.Right)
		if err != nil {
			return 0, err
		}
		return x % y, nil
	case *tree.MultExpr:
		x, err := b.buildExprIntConstant(e.Left)
		if err != nil {
			return 0, err
		}
		y, err := b.buildExprIntConstant(e.Right)
		if err != nil {
			return 0, err
		}
		return x * y, nil
	case *tree.PlusExpr:
		x, err := b.buildExprIntConstant(e.Left)
		if err != nil {
			return 0, err
		}
		y, err := b.buildExprIntConstant(e.Right)
		if err != nil {
			return 0, err
		}
		return x + y, nil
	case *tree.MinusExpr:
		x, err := b.buildExprIntConstant(e.Left)
		if err != nil {
			return 0, err
		}
		y, err := b.buildExprIntConstant(e.Right)
		if err != nil {
			return 0, err
		}
		return x - y, nil
	case *tree.UnaryMinusExpr:
		x, err := b.buildExprIntConstant(e.E)
		if err != nil {
			return 0, err
		}
		return x * -1, nil
	default:
		return 0, fmt.Errorf("'%s' is not integer", n)
	}
}

func reduce(e extend.Extend) extend.Extend {
	if be, ok := e.(*extend.BinaryExtend); ok {
		switch be.Op {
		case overload.Typecast:
			return buildTypecast(be)
		}
	}
	return e
}

func buildTypecast(e *extend.BinaryExtend) extend.Extend {
	if lv, ok := e.Left.(value.Value); ok {
		if rv, ok := e.Right.(*value.String); ok {
			switch typ := value.MustBeString(rv); strings.ToLower(typ) {
			case "int":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewInt(0)); err == nil {
					return v
				}
			case "uint8":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewUint8(0)); err == nil {
					return v
				}
			case "bool":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewBool(true)); err == nil {
					return v
				}
			case "time":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewTime(time.Now())); err == nil {
					return v
				}
			case "float":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewFloat(0.0)); err == nil {
					return v
				}
			case "string":
				if v, err := overload.BinaryEval(e.Op, lv, value.NewString("")); err == nil {
					return v
				}
			}
		}
	}
	return e
}

var ExtendFuncs map[string]int = map[string]int{
	"abs":     overload.Abs,
	"ceil":    overload.Ceil,
	"sign":    overload.Sign,
	"floor":   overload.Floor,
	"round":   overload.Round,
	"lower":   overload.Lower,
	"upper":   overload.Upper,
	"length":  overload.Length,
	"typeof":  overload.Typeof,
	"concat":  overload.Concat,
	"cast":    overload.Typecast,
	"like":    overload.Like,
	"notlike": overload.NotLike,
}
