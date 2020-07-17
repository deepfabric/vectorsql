package overload

import (
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
	Match
	NotMatch

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
	Match:    "match",
	NotMatch: "not match",

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
	Typ        uint32
	ReturnType uint32
	Fn         func(value.Values) (value.Values, error)
}

// BinOp is a binary operator.
type BinOp struct {
	LeftType   uint32
	RightType  uint32
	ReturnType uint32

	Fn func(value.Values, value.Values) (value.Values, error)
}

// MultiOp is a multiple operator.
type MultiOp struct {
	Min        int // minimum number of parameters
	Max        int // maximum number of parameters, -1 means unlimited
	Typ        uint32
	ReturnType uint32

	Fn func([]value.Values) (value.Values, error)
}
