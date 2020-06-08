package tree

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

type ExprStatement interface {
	Statement
	exprStatement()
}

func (*Value) exprStatement() {}

func (*OrExpr) exprStatement()  {}
func (*AndExpr) exprStatement() {}
func (*NotExpr) exprStatement() {}

func (*DivExpr) exprStatement()   {}
func (*ModExpr) exprStatement()   {}
func (*MultExpr) exprStatement()  {}
func (*PlusExpr) exprStatement()  {}
func (*MinusExpr) exprStatement() {}

func (*UnaryMinusExpr) exprStatement() {}

func (*EqExpr) exprStatement() {}
func (*NeExpr) exprStatement() {}
func (*LtExpr) exprStatement() {}
func (*LeExpr) exprStatement() {}
func (*GtExpr) exprStatement() {}
func (*GeExpr) exprStatement() {}

func (*Subquery) exprStatement() {}

func (*BetweenExpr) exprStatement()    {}
func (*NotBetweenExpr) exprStatement() {}

func (*IsNullExpr) exprStatement()    {}
func (*IsNotNullExpr) exprStatement() {}

func (*ParenExpr) exprStatement() {}

func (*FuncExpr) exprStatement() {}

func (ExprStatements) exprStatement() {}

func (ColunmName) exprStatement()     {}
func (ColunmNameList) exprStatement() {}

func (e *Value) String() string { return e.E.String() }

func (e *NotExpr) String() string { return fmt.Sprintf("NOT %s", e.E) }
func (e *OrExpr) String() string  { return fmt.Sprintf("%s OR %s", e.Left, e.Right) }
func (e *AndExpr) String() string { return fmt.Sprintf("%s AND %s", e.Left, e.Right) }

func (e *DivExpr) String() string   { return fmt.Sprintf("%s / %s", e.Left, e.Right) }
func (e *ModExpr) String() string   { return fmt.Sprintf("%s %% %s", e.Left, e.Right) }
func (e *MultExpr) String() string  { return fmt.Sprintf("%s * %s", e.Left, e.Right) }
func (e *PlusExpr) String() string  { return fmt.Sprintf("%s + %s", e.Left, e.Right) }
func (e *MinusExpr) String() string { return fmt.Sprintf("%s - %s", e.Left, e.Right) }

func (e *UnaryMinusExpr) String() string { return fmt.Sprintf("-%s", e.E) }

func (e *EqExpr) String() string { return fmt.Sprintf("%s = %s", e.Left, e.Right) }
func (e *NeExpr) String() string { return fmt.Sprintf("%s <> %s", e.Left, e.Right) }
func (e *LtExpr) String() string { return fmt.Sprintf("%s < %s", e.Left, e.Right) }
func (e *LeExpr) String() string { return fmt.Sprintf("%s <= %s", e.Left, e.Right) }
func (e *GtExpr) String() string { return fmt.Sprintf("%s > %s", e.Left, e.Right) }
func (e *GeExpr) String() string { return fmt.Sprintf("%s >= %s", e.Left, e.Right) }

func (e *BetweenExpr) String() string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", e.E, e.From, e.To)
}

func (e *NotBetweenExpr) String() string {
	return fmt.Sprintf("%s NOT BETWEEN %s AND %s", e.E, e.From, e.To)
}

func (e *IsNullExpr) String() string    { return fmt.Sprintf("%s IS NULL", e.E) }
func (e *IsNotNullExpr) String() string { return fmt.Sprintf("%s IS NOT NULL", e.E) }

func (e *ParenExpr) String() string { return fmt.Sprintf("(%s)", e.E) }

func (e *FuncExpr) String() string {
	return fmt.Sprintf("%s(%s)", e.Name, e.Es)
}

func (es ExprStatements) String() string {
	var s string

	for i := range es {
		if i > 0 {
			s += ", "
		}
		s += es[i].String()
	}
	return s
}

type Value struct {
	E value.Value
}

type NotExpr struct {
	E ExprStatement
}

type OrExpr struct {
	Left, Right ExprStatement
}

type AndExpr struct {
	Left, Right ExprStatement
}

type DivExpr struct {
	Left, Right ExprStatement
}

type ModExpr struct {
	Left, Right ExprStatement
}

type MultExpr struct {
	Left, Right ExprStatement
}

type PlusExpr struct {
	Left, Right ExprStatement
}

type MinusExpr struct {
	Left, Right ExprStatement
}

type UnaryMinusExpr struct {
	E ExprStatement
}

type EqExpr struct {
	Left, Right ExprStatement
}

type NeExpr struct {
	Left, Right ExprStatement
}

type LtExpr struct {
	Left, Right ExprStatement
}

type LeExpr struct {
	Left, Right ExprStatement
}

type GtExpr struct {
	Left, Right ExprStatement
}

type GeExpr struct {
	Left, Right ExprStatement
}

type BetweenExpr struct {
	E        ExprStatement
	From, To ExprStatement
}

type NotBetweenExpr struct {
	E        ExprStatement
	From, To ExprStatement
}

type IsNullExpr struct {
	E ExprStatement
}

type IsNotNullExpr struct {
	E ExprStatement
}

type ParenExpr struct {
	E ExprStatement
}

type FuncExpr struct {
	Name string
	Es   ExprStatements
}

type ExprStatements []ExprStatement
