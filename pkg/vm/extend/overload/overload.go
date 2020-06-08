package overload

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/util/arith"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func IsLogical(op int) bool {
	switch op {
	case Or:
		return true
	case Not:
		return true
	case And:
		return true
	case Like, NotLike:
		return true
	case EQ, LT, GT, LE, GE, NE:
		return true
	default:
		return false
	}
}

func OperatorType(op int) int {
	switch op {
	case UnaryMinus:
		return Unary
	case Not:
		return Unary
	case Abs:
		return Unary
	case Ceil:
		return Unary
	case Sign:
		return Unary
	case Floor:
		return Unary
	case Lower:
		return Unary
	case Round:
		return Unary
	case Upper:
		return Unary
	case Length:
		return Unary
	case Typeof:
		return Unary
	case Or:
		return Binary
	case And:
		return Binary
	case Plus:
		return Binary
	case Minus:
		return Binary
	case Mult:
		return Binary
	case Div:
		return Binary
	case Mod:
		return Binary
	case Typecast:
		return Binary
	case EQ:
		return Binary
	case LT:
		return Binary
	case GT:
		return Binary
	case LE:
		return Binary
	case GE:
		return Binary
	case NE:
		return Binary
	case Like, NotLike:
		return Binary
	case Concat:
		return Multi
	}
	return -1
}

func UnaryEval(op int, v value.Value) (value.Value, error) {
	if os, ok := UnaryOps[op]; ok {
		for _, o := range os {
			if o.Typ.Oid == types.T_any || v.ResolvedType().Oid == o.Typ.Oid {
				return o.Fn(v)
			}
		}
	}
	return nil, fmt.Errorf("%s not yet implemented for %s", OpName[op], v)
}

func BinaryEval(op int, left, right value.Value) (value.Value, error) {
	if os, ok := BinOps[op]; ok {
		for _, o := range os {
			if (o.LeftType.Oid == types.T_any || left.ResolvedType().Oid == o.LeftType.Oid) &&
				(o.RightType.Oid == types.T_any || right.ResolvedType().Oid == o.RightType.Oid) {
				return o.Fn(left, right)
			}
		}
	}
	return nil, fmt.Errorf("%s not yet implemented for %s, %s", OpName[op], left, right)
}

func MultiEval(op int, args []value.Value) (value.Value, error) {
	if os, ok := MultiOps[op]; ok {
		for _, o := range os {
			if n := len(args); n >= o.Min || (o.Max == -1 || n <= o.Max) {
				if n == 0 || n > 0 && (o.Typ.Oid == types.T_any || args[0].ResolvedType().Oid == o.Typ.Oid) {
					return o.Fn(args)
				}
			}
		}
	}
	return nil, fmt.Errorf("%s not yet implemented for %s", OpName[op], value.Array(args))
}

var (
	// ErrIntOutOfRange is reported when integer arithmetic overflows.
	ErrIntOutOfRange = errors.New("integer out of range")

	// ErrDivByZero is reported on a division by zero.
	ErrDivByZero = errors.New("division by zero")
	// ErrZeroModulus is reported when computing the rest of a division by zero.
	ErrZeroModulus   = errors.New("zero modulus")
	errAbsOfMinInt64 = errors.New("abs of min integer value (-9223372036854775808) not defined")
)

var UnaryOps = map[int][]*UnaryOp{
	UnaryMinus: {
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Int,
			Fn: func(v value.Value) (value.Value, error) {
				i := value.MustBeInt(v)
				if i == math.MinInt64 {
					return nil, ErrIntOutOfRange
				}
				return value.NewInt(-i), nil
			},
		},
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(-value.MustBeFloat(v)), nil
			},
		},
	},
	Abs: {
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Int,
			Fn: func(v value.Value) (value.Value, error) {
				x := value.MustBeInt(v)
				switch {
				case x == math.MinInt64:
					return nil, errAbsOfMinInt64
				case x < 0:
					return value.NewInt(-x), nil
				}
				return v, nil
			},
		},
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(math.Abs(value.MustBeFloat(v))), nil
			},
		},
	},
	Not: {
		&UnaryOp{
			Typ:        types.Bool,
			ReturnType: types.Bool,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewBool(!value.MustBeBool(v)), nil
			},
		},
	},
	Ceil: {
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(float64(value.MustBeInt(v))), nil
			},
		},
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(math.Ceil(value.MustBeFloat(v))), nil
			},
		},
	},
	Sign: {
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				x := value.MustBeFloat(v)
				switch {
				case x < 0.0:
					return value.NewFloat(-1.0), nil
				case x == 0.0:
					return value.NewFloat(0.0), nil
				}
				return value.NewFloat(1.0), nil
			},
		},
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Int,
			Fn: func(v value.Value) (value.Value, error) {
				x := value.MustBeInt(v)
				switch {
				case x < 0:
					return value.NewInt(-1), nil
				case x == 0:
					return value.NewInt(0), nil
				}
				return value.NewInt(1), nil
			},
		},
	},
	Floor: {
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(float64(value.MustBeInt(v))), nil
			},
		},
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(math.Floor(value.MustBeFloat(v))), nil
			},
		},
	},
	Lower: {
		&UnaryOp{
			Typ:        types.String,
			ReturnType: types.String,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewString(strings.ToLower(value.MustBeString(v))), nil
			},
		},
	},
	Round: {
		&UnaryOp{
			Typ:        types.Int,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(float64(value.MustBeInt(v))), nil
			},
		},
		&UnaryOp{
			Typ:        types.Float,
			ReturnType: types.Float,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewFloat(math.RoundToEven(value.MustBeFloat(v))), nil
			},
		},
	},
	Upper: {
		&UnaryOp{
			Typ:        types.String,
			ReturnType: types.String,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewString(strings.ToUpper(value.MustBeString(v))), nil
			},
		},
	},
	Length: {
		&UnaryOp{
			Typ:        types.String,
			ReturnType: types.Int,
			Fn: func(v value.Value) (value.Value, error) {
				n := utf8.RuneCountInString(value.MustBeString(v))
				return value.NewInt(int64(n)), nil
			},
		},
	},
	Typeof: {
		&UnaryOp{
			Typ:        types.Any,
			ReturnType: types.String,
			Fn: func(v value.Value) (value.Value, error) {
				return value.NewString(strings.ToLower(v.ResolvedType().String())), nil
			},
		},
	},
}

// BinOps contains the binary operations indexed by operation type.
var BinOps = map[int][]*BinOp{
	Or: {
		&BinOp{
			LeftType:   types.Bool,
			RightType:  types.Bool,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeBool(left), value.MustBeBool(right)
				return value.NewBool(a || b), nil
			},
		},
	},
	And: {
		&BinOp{
			LeftType:   types.Bool,
			RightType:  types.Bool,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeBool(left), value.MustBeBool(right)
				return value.NewBool(a && b), nil
			},
		},
	},
	Plus: {
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeInt(right)
				r, ok := arith.AddWithOverflow(int64(a), int64(b))
				if !ok {
					return nil, ErrIntOutOfRange
				}
				return value.NewInt(r), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeFloat(right)
				return value.NewFloat(a + b), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeFloat(right)
				return value.NewFloat(float64(a) + b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Int,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeInt(right)
				return value.NewFloat(a + float64(b)), nil
			},
		},
		&BinOp{
			LeftType:   types.Time,
			RightType:  types.Int,
			ReturnType: types.Time,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeTime(left), value.MustBeInt(right)
				return value.NewTime(a.Add(time.Duration(b))), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Time,
			ReturnType: types.Time,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeTime(right)
				return value.NewTime(b.Add(time.Duration(a))), nil
			},
		},
	},

	Minus: {
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeInt(right)
				r, ok := arith.SubWithOverflow(int64(a), int64(b))
				if !ok {
					return nil, ErrIntOutOfRange
				}
				return value.NewInt(r), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeFloat(right)
				return value.NewFloat(a - b), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeFloat(right)
				return value.NewFloat(float64(a) - b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Int,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeInt(right)
				return value.NewFloat(a - float64(b)), nil
			},
		},
		&BinOp{
			LeftType:   types.Time,
			RightType:  types.Time,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeTime(left), value.MustBeTime(right)
				return value.NewInt(int64(a.Sub(b))), nil
			},
		},
	},

	Mult: {
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				// See Rob Pike's implementation from
				// https://groups.google.com/d/msg/golang-nuts/h5oSN5t3Au4/KaNQREhZh0QJ

				a, b := value.MustBeInt(left), value.MustBeInt(right)
				c := a * b
				if a == 0 || b == 0 || a == 1 || b == 1 {
					// ignore
				} else if a == math.MinInt64 || b == math.MinInt64 {
					// This test is required to detect math.MinInt64 * -1.
					return nil, ErrIntOutOfRange
				} else if c/b != a {
					return nil, ErrIntOutOfRange
				}
				return value.NewInt(c), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeFloat(right)
				return value.NewFloat(a * b), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeFloat(right)
				return value.NewFloat(float64(a) * b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Int,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeInt(right)
				return value.NewFloat(a * float64(b)), nil
			},
		},
	},

	Div: {
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeInt(right)
				if b == 0 {
					return nil, ErrDivByZero
				}
				return value.NewInt(a / b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeFloat(right)
				if b == 0 {
					return nil, ErrDivByZero
				}
				return value.NewFloat(a / b), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeFloat(right)
				if b == 0 {
					return nil, ErrDivByZero
				}
				return value.NewFloat(float64(a) / b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Int,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeInt(right)
				if b == 0 {
					return nil, ErrDivByZero
				}
				return value.NewFloat(a / float64(b)), nil
			},
		},
	},

	Mod: {
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeInt(right)
				if b == 0 {
					return nil, ErrZeroModulus
				}
				return value.NewInt(a % b), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeFloat(right)
				return value.NewFloat(math.Mod(a, b)), nil
			},
		},
		&BinOp{
			LeftType:   types.Int,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeInt(left), value.MustBeFloat(right)
				return value.NewFloat(math.Mod(float64(a), b)), nil
			},
		},
		&BinOp{
			LeftType:   types.Float,
			RightType:  types.Int,
			ReturnType: types.Float,
			Fn: func(left, right value.Value) (value.Value, error) {
				a, b := value.MustBeFloat(left), value.MustBeInt(right)
				return value.NewFloat(math.Mod(a, float64(b))), nil
			},
		},
	},
	Typecast: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Uint8,
			ReturnType: types.Uint8,
			Fn: func(v, _ value.Value) (value.Value, error) {
				switch v.ResolvedType().Oid {
				case types.T_int:
					return value.NewUint8(uint8(value.MustBeInt(v) & 0xFF)), nil
				case types.T_uint8:
					return v, nil
				case types.T_bool:
					if value.MustBeBool(v) {
						return value.NewUint8(1), nil
					}
					return value.NewUint8(0), nil
				case types.T_time:
					return value.NewUint8(uint8(value.MustBeTime(v).Unix() & 0xFF)), nil
				case types.T_float:
					f := value.MustBeFloat(v)
					if math.IsNaN(f) || f <= float64(math.MinInt64) || f >= float64(math.MaxInt64) {
						return nil, ErrIntOutOfRange
					}
					return value.NewUint8(uint8(f)), nil
				case types.T_string:
					return value.ParseUint8(value.MustBeString(v))
				}
				return value.ConstNull, fmt.Errorf("cannot convert type %s to type %s", v.ResolvedType(), types.Uint8)
			},
		},
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Int,
			ReturnType: types.Int,
			Fn: func(v, _ value.Value) (value.Value, error) {
				switch v.ResolvedType().Oid {
				case types.T_int:
					return v, nil
				case types.T_uint8:
					return value.NewInt(int64(value.MustBeUint8(v))), nil
				case types.T_bool:
					if value.MustBeBool(v) {
						return value.NewInt(1), nil
					}
					return value.NewInt(0), nil
				case types.T_time:
					return value.NewInt(value.MustBeTime(v).Unix()), nil
				case types.T_float:
					f := value.MustBeFloat(v)
					if math.IsNaN(f) || f <= float64(math.MinInt64) || f >= float64(math.MaxInt64) {
						return nil, ErrIntOutOfRange
					}
					return value.NewInt(int64(f)), nil
				case types.T_string:
					return value.ParseInt(value.MustBeString(v))
				}
				return value.ConstNull, fmt.Errorf("cannot convert type %s to type %s", v.ResolvedType(), types.Int)
			},
		},
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Bool,
			ReturnType: types.Bool,
			Fn: func(v, _ value.Value) (value.Value, error) {
				switch v.ResolvedType().Oid {
				case types.T_int:
					return value.NewBool(value.MustBeInt(v) != 0), nil
				case types.T_uint8:
					return value.NewBool(value.MustBeUint8(v) != 0), nil
				case types.T_bool:
					return v, nil
				case types.T_float:
					return value.NewBool(value.MustBeFloat(v) != 0), nil
				case types.T_string:
					return value.ParseBool(value.MustBeString(v))
				}
				return value.ConstNull, fmt.Errorf("cannot convert type %s to type %s", v.ResolvedType(), types.Bool)
			},
		},
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Time,
			ReturnType: types.Time,
			Fn: func(v, _ value.Value) (value.Value, error) {
				switch v.ResolvedType().Oid {
				case types.T_int:
					return value.NewTime(time.Unix(value.MustBeInt(v), 0)), nil
				case types.T_time:
					return v, nil
				case types.T_float:
					return value.NewTime(time.Unix(int64(value.MustBeFloat(v)), 0)), nil
				case types.T_string:
					return value.ParseTime(value.MustBeString(v))
				}
				return value.ConstNull, fmt.Errorf("cannot convert type %s to type %s", v.ResolvedType(), types.Time)
			},
		},
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Float,
			ReturnType: types.Float,
			Fn: func(v, _ value.Value) (value.Value, error) {
				switch v.ResolvedType().Oid {
				case types.T_int:
					return value.NewFloat(float64(value.MustBeInt(v))), nil
				case types.T_uint8:
					return value.NewFloat(float64(value.MustBeUint8(v))), nil
				case types.T_bool:
					if value.MustBeBool(v) {
						return value.NewFloat(1), nil
					}
					return value.NewFloat(0), nil
				case types.T_time:
					return value.NewFloat(float64(value.MustBeTime(v).Unix())), nil
				case types.T_float:
					return v, nil
				case types.T_string:
					return value.ParseFloat(value.MustBeString(v))
				}
				return value.ConstNull, fmt.Errorf("cannot convert type %s to type %s", v.ResolvedType(), types.Float)
			},
		},
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.String,
			ReturnType: types.String,
			Fn: func(v, _ value.Value) (value.Value, error) {
				return value.NewString(v.String()), nil
			},
		},
	},
	EQ: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) == 0), nil
			},
		},
	},
	LT: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) < 0), nil
			},
		},
	},
	GT: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) > 0), nil
			},
		},
	},
	LE: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) <= 0), nil
			},
		},
	},
	GE: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) >= 0), nil
			},
		},
	},
	NE: {
		&BinOp{
			LeftType:   types.Any,
			RightType:  types.Any,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				return value.NewBool(value.Compare(left, right) != 0), nil
			},
		},
	},
	Like: {
		&BinOp{
			LeftType:   types.String,
			RightType:  types.String,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				rgp, err := regexp.Compile(value.MustBeString(right))
				if err != nil {
					return nil, err
				}
				return value.NewBool(rgp.MatchString(value.MustBeString(left))), nil
			},
		},
	},
	NotLike: {
		&BinOp{
			LeftType:   types.String,
			RightType:  types.String,
			ReturnType: types.Bool,
			Fn: func(left, right value.Value) (value.Value, error) {
				rgp, err := regexp.Compile("!" + value.MustBeString(right))
				if err != nil {
					return nil, err
				}
				return value.NewBool(rgp.MatchString(value.MustBeString(left))), nil
			},
		},
	},
}

var MultiOps = map[int][]*MultiOp{
	// concat concatenates the text representations of all the arguments.
	// NULL and Table arguments are ignored.
	Concat: {
		&MultiOp{
			Min:        1,
			Max:        -1,
			Typ:        types.String,
			ReturnType: types.String,
			Fn: func(args []value.Value) (value.Value, error) {
				var buffer bytes.Buffer

				for _, arg := range args {
					if oid := arg.ResolvedType().Oid; oid == types.T_null {
						continue
					}
					buffer.WriteString(value.MustBeString(arg))
				}
				return value.NewString(buffer.String()), nil
			},
		},
	},
}
