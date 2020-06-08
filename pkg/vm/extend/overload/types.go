package overload

import (
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

const (
	Multi = iota
	Unary
	Binary
)

const (
	// unary operator
	UnaryMinus = iota
	Abs
	Not // logical operator
	Ceil
	Sign
	Floor
	Lower
	Round
	Upper
	Length
	Typeof

	// binary operator
	Or  // logical operator
	And // logical operator
	Plus
	Minus
	Mult
	Div
	Mod
	Typecast
	Like
	NotLike

	// binary operator - comparison operator
	EQ
	LT
	GT
	LE
	GE
	NE

	// multiple operator
	Concat
)

var OpName = [...]string{
	UnaryMinus: "-",
	Abs:        "abs",
	Not:        "not",
	Ceil:       "ceil",
	Sign:       "sign",
	Floor:      "floor",
	Lower:      "lower",
	Round:      "round",
	Upper:      "upper",
	Length:     "length",
	Typeof:     "typeof",

	Or:       "or",
	And:      "and",
	Plus:     "+",
	Minus:    "-",
	Mult:     "*",
	Div:      "/",
	Mod:      "%",
	Typecast: "typecast",
	Like:     "like",
	NotLike:  "not like",

	EQ: "=",
	LT: "<",
	GT: ">",
	LE: "<=",
	GE: ">=",
	NE: "<>",

	Concat: "concat",
}

// UnaryOp is a unary operator.
type UnaryOp struct {
	Typ        *types.T
	ReturnType *types.T
	Fn         func(value.Value) (value.Value, error)
}

// BinOp is a binary operator.
type BinOp struct {
	LeftType   *types.T
	RightType  *types.T
	ReturnType *types.T

	Fn func(value.Value, value.Value) (value.Value, error)
}

// MultiOp is a multiple operator.
type MultiOp struct {
	Min        int // minimum number of parameters
	Max        int // maximum number of parameters, -1 means unlimited
	Typ        *types.T
	ReturnType *types.T

	Fn func([]value.Value) (value.Value, error)
}
