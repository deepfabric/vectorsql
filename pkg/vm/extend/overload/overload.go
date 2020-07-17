package overload

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/deepfabric/vectorsql/pkg/lru"
	"github.com/deepfabric/vectorsql/pkg/match"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
	"github.com/deepfabric/vectorsql/pkg/vm/value/dynamic"
	"github.com/deepfabric/vectorsql/pkg/vm/value/static"
	"github.com/pilosa/pilosa/roaring"
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
	case Match, NotMatch:
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
	case Match, NotMatch:
		return Binary
	case Concat:
		return Binary
	}
	return -1
}

func UnaryEval(op int, typ uint32, vs value.Values) (value.Values, uint32, error) {
	if os, ok := UnaryOps[op]; ok {
		for _, o := range os {
			if unaryCheck(op, o.Typ, typ) {
				rs, err := o.Fn(vs)
				return rs, o.ReturnType, err
			}
		}
	}
	return nil, 0, fmt.Errorf("%s not yet implemented for %s", OpName[op], types.T(typ))
}

func BinaryEval(op int, ltyp, rtyp uint32, as, bs value.Values) (value.Values, uint32, error) {
	if os, ok := BinOps[op]; ok {
		for _, o := range os {
			if binaryCheck(op, o.LeftType, o.RightType, ltyp, rtyp) {
				rs, err := o.Fn(as, bs)
				return rs, o.ReturnType, err
			}
		}
	}
	return nil, 0, fmt.Errorf("%s not yet implemented for %s, %s", OpName[op], types.T(ltyp), types.T(rtyp))
}

func MultiEval(op int, typ uint32, vs []value.Values) (value.Values, uint32, error) {
	if os, ok := MultiOps[op]; ok {
		for _, o := range os {
			if n := vs[0].Count(); n >= o.Min && (o.Max == -1 || n <= o.Max) {
				if multiCheck(op, o.Typ, typ) {
					rs, err := o.Fn(vs)
					return rs, o.ReturnType, err
				}
			}
		}
	}
	return nil, 0, fmt.Errorf("%s not yet implemented for %s", OpName[op], types.T(typ))
}

func unaryCheck(op int, arg uint32, val uint32) bool {
	switch op {
	case Typeof:
		return true
	}
	return arg == val
}

func binaryCheck(op int, arg0, arg1 uint32, val0, val1 uint32) bool {
	return arg0 == val0 && arg1 == val1
}

func multiCheck(op int, arg uint32, val uint32) bool {
	return false
}

var (
	ErrDivByZero   = errors.New("division by zero")
	ErrZeroModulus = errors.New("zero modulus")
)

var UnaryOps = map[int][]*UnaryOp{
	UnaryMinus: {
		&UnaryOp{
			Typ:        types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = -a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = -v
					}
				}
				return r, nil
			},
		},
	},
	Abs: {
		&UnaryOp{
			Typ:        types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int64

				a := vs.(*static.Ints)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int8

				a := vs.(*static.Int8s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int16

				a := vs.(*static.Int16s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int32

				a := vs.(*static.Int32s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int64

				a := vs.(*static.Int64s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float64

				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float32

				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float64

				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -y
						} else {
							r.Vs[o] = y
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -v
						} else {
							r.Vs[i] = v
						}
					}
				}
				return r, nil
			},
		},
	},
	Not: {
		&UnaryOp{
			Typ:        types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = !a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = !v
					}
				}
				return r, nil
			},
		},
	},
	Ceil: {
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.Ceil(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.Ceil(v)
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(math.Ceil(float64(a.Vs[o])))
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(math.Ceil(float64(v)))
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.Ceil(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.Ceil(v)
					}
				}
				return r, nil
			},
		},
	},
	Sign: {
		&UnaryOp{
			Typ:        types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int64

				a := vs.(*static.Ints)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1
						} else {
							r.Vs[o] = 1
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1
						} else {
							r.Vs[i] = 1
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int8

				a := vs.(*static.Int8s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1
						} else {
							r.Vs[o] = 1
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1
						} else {
							r.Vs[i] = 1
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int16

				a := vs.(*static.Int16s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1
						} else {
							r.Vs[o] = 1
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1
						} else {
							r.Vs[i] = 1
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int32

				a := vs.(*static.Int32s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1
						} else {
							r.Vs[o] = 1
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1
						} else {
							r.Vs[i] = 1
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs value.Values) (value.Values, error) {
				var y int64

				a := vs.(*static.Int64s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1
						} else {
							r.Vs[o] = 1
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1
						} else {
							r.Vs[i] = 1
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float64

				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1.0
						} else {
							r.Vs[o] = 1.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1.0
						} else {
							r.Vs[i] = 1.0
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float32

				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1.0
						} else {
							r.Vs[o] = 1.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1.0
						} else {
							r.Vs[i] = 1.0
						}
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				var y float64

				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if y = a.Vs[o]; y < 0 {
							r.Vs[o] = -1.0
						} else {
							r.Vs[o] = 1.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v < 0 {
							r.Vs[i] = -1.0
						} else {
							r.Vs[i] = 1.0
						}
					}
				}
				return r, nil
			},
		},
	},
	Floor: {
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.Floor(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.Floor(v)
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(math.Floor(float64(a.Vs[o])))
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(math.Floor(float64(v)))
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.Floor(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.Floor(v)
					}
				}
				return r, nil
			},
		},
	},
	Lower: {
		&UnaryOp{
			Typ:        types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strings.ToLower(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strings.ToLower(v)
					}
				}
				return r, nil
			},
		},
	},
	Round: {
		&UnaryOp{
			Typ:        types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.RoundToEven(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.RoundToEven(v)
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(math.RoundToEven(float64(a.Vs[o])))
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(math.RoundToEven(float64(v)))
					}
				}
				return r, nil
			},
		},
		&UnaryOp{
			Typ:        types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = math.RoundToEven(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = math.RoundToEven(v)
					}
				}
				return r, nil
			},
		},
	},
	Upper: {
		&UnaryOp{
			Typ:        types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strings.ToUpper(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strings.ToUpper(v)
					}
				}
				return r, nil
			},
		},
	},
	Length: {
		&UnaryOp{
			Typ:        types.T_string,
			ReturnType: types.T_int,
			Fn: func(vs value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(len(a.Vs[o]))
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(len(v))
					}
				}
				return r, nil
			},
		},
	},
	Typeof: {
		&UnaryOp{
			ReturnType: types.T_string,
			Fn: func(vs value.Values) (value.Values, error) {
				switch t := vs.(type) {
				case *static.Ints:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_int).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Int8s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_int8).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Int16s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_int16).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Int32s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_int32).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Int64s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_int64).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Uint8s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_uint8).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Uint16s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_uint16).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Uint32s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_uint32).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Uint64s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_uint64).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Bools:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_bool).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Timestamps:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_timestamp).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Floats:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_float).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Float32s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_float32).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *static.Float64s:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_float64).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				case *dynamic.Strings:
					r := &dynamic.Strings{
						Np: t.Np,
						Dp: t.Dp,
						Is: t.Is,
						Vs: make([]string, len(t.Vs)),
					}
					ts := strings.ToLower(types.T(types.T_string).String())
					if len(r.Is) > 0 {
						for _, o := range r.Is {
							r.Vs[o] = ts
						}
					} else {
						for i := range r.Vs {
							r.Vs[i] = ts
						}
					}
					return r, nil
				}
				return nil, nil
			},
		},
	},
}

// BinOps contains the binary operations indexed by operation type.
var BinOps = map[int][]*BinOp{
	Or: {
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] || b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v || b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	And: {
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] && b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v && b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	Plus: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Ints{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Floats{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	Minus: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Ints{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Floats{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] - b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v - b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	Mult: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Ints{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Floats{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] * b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v * b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	Div: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Ints{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !!fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = int8(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = int8(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = int8(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = int8(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = int8(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = int8(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = int16(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = int16(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = int16(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = int16(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = int16(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = int16(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = int32(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = int32(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = int32(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = int32(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = int32(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = int32(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = int64(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = int64(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = int64(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = int64(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = int64(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = int64(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = uint8(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = uint8(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = uint8(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = uint8(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = uint8(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = uint8(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = uint16(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = uint16(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = uint16(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = uint16(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = uint16(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = uint16(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = uint32(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = uint32(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = uint32(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = uint32(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = uint32(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = uint32(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = uint64(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = uint64(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = uint64(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = uint64(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = uint64(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = uint64(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Floats{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = float32(b.Vs[o])
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = float32(b.Vs[i])
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = float32(a.Vs[o]) / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = float32(a.Vs[o]) / v
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = float32(v) / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = float32(v) / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[o] = a.Vs[o] / v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[o] = a.Vs[o] / v
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrDivByZero
							}
							r.Vs[i] = v / w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrDivByZero
								}
								r.Vs[i] = v / w
							}
						}
					}
				}
				return r, nil
			},
		},
	},
	Mod: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Ints{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = int8(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = int8(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Int8s{
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = int8(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = int8(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w int8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = int8(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = int8(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = int16(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = int16(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Int16s{
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = int16(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = int16(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w int16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = int16(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = int16(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = int32(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = int32(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Int32s{
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = int32(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = int32(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w int32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = int32(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = int32(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = int64(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = int64(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Int64s{
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v int64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = int64(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = int64(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w int64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = int64(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = int64(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = uint8(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = uint8(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Uint8s{
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint8

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = uint8(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = uint8(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w uint8

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = uint8(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = uint8(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = uint16(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = uint16(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Uint16s{
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint16

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = uint16(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = uint16(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w uint16

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = uint16(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = uint16(v) % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = uint32(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = uint32(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Uint32s{
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = uint32(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = uint32(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w uint32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = uint32(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = uint32(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = uint64(b.Vs[o])
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = a.Vs[o] % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = a.Vs[o] % v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = uint64(b.Vs[i])
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = v % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = v % w
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Uint64s{
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v uint64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = uint64(a.Vs[o]) % v
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = uint64(a.Vs[o]) % v
							}
						}
					}
				} else {
					var w uint64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = uint64(v) % w
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = uint64(v) % b.Vs[i]
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Floats{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = math.Mod(a.Vs[o], v)
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = math.Mod(a.Vs[o], v)
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = math.Mod(v, w)
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = math.Mod(v, b.Vs[i])
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = float32(math.Mod(float64(a.Vs[o]), float64(v)))
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = float32(math.Mod(float64(a.Vs[o]), float64(v)))
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = float32(math.Mod(float64(v), float64(w)))
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = float32(math.Mod(float64(v), float64(w)))
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = math.Mod(a.Vs[o], v)
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = math.Mod(a.Vs[o], v)
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = math.Mod(v, w)
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = math.Mod(v, b.Vs[i])
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = float32(b.Vs[o])
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = float32(math.Mod(float64(a.Vs[o]), float64(v)))
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = float32(math.Mod(float64(a.Vs[o]), float64(v)))
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = float32(b.Vs[i])
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = float32(math.Mod(float64(v), float64(w)))
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = float32(math.Mod(float64(v), float64(w)))
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Float32s{
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float32

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = float32(math.Mod(a.Vs[o], float64(v)))
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = float32(math.Mod(a.Vs[o], float64(v)))
							}
						}
					}
				} else {
					var w float32

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = float32(math.Mod(v, float64(w)))
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = float32(math.Mod(v, float64(w)))
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = math.Mod(a.Vs[o], v)
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = math.Mod(a.Vs[o], v)
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = math.Mod(v, w)
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = math.Mod(v, b.Vs[i])
							}
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Float64s{
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				var fp *roaring.Bitmap
				switch {
				case r.Dp == nil && r.Np != nil:
					fp = r.Np
				case r.Dp != nil && r.Np == nil:
					fp = r.Dp
				case r.Dp != nil && r.Np != nil:
					fp = r.Np.Union(r.Dp)
				}
				if len(r.Is) > 0 {
					var v float64

					for _, o := range r.Is {
						v = b.Vs[o]
						if fp == nil {
							if v == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[o] = math.Mod(a.Vs[o], v)
						} else {
							if !fp.Contains(o) {
								if v == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[o] = math.Mod(a.Vs[o], v)
							}
						}
					}
				} else {
					var w float64

					for i, v := range a.Vs {
						w = b.Vs[i]
						if fp == nil {
							if w == 0.0 {
								return nil, ErrZeroModulus
							}
							r.Vs[i] = math.Mod(v, w)
						} else {
							if !fp.Contains(uint64(i)) {
								if w == 0 {
									return nil, ErrZeroModulus
								}
								r.Vs[i] = math.Mod(v, w)
							}
						}
					}
				}
				return r, nil
			},
		},
	},
	EQ: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	LT: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = !a.Vs[o] && b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = !v && b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) < 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] < b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v < b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	GT: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] && !b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v && !b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] > b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v > b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	LE: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o] || (!a.Vs[o] && b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i] || (!v && b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) <= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] <= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v <= b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	GE: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] && !b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v && !b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},

		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] == b.Vs[o] && (a.Vs[o] && !b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v == b.Vs[i] && (v && !b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) >= 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] >= b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v >= b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	NE: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != int8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != int8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != int16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != int16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int32,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != int32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != int32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Int64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Int64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint8s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != uint8(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != uint8(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint8s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint16s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != uint16(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != uint16(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint16s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint32s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != uint32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != uint32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Uint64s), bs.(*static.Ints)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != uint64(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != uint64(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Ints), bs.(*static.Uint64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float32s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != float32(b.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != float32(b.Vs[i])
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float32s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o]) != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v) != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Float64s), bs.(*static.Floats)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Floats), bs.(*static.Float64s)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Bools), bs.(*static.Bools)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*static.Timestamps), bs.(*static.Timestamps)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) != 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] != b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v != b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
	Like: {
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				m := match.New(lru.New(1 << 10))
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rgp, err := m.Compile(b.Vs[o], true)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rgp.MatchString(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						rgp, err := m.Compile(b.Vs[i], true)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rgp.MatchString(v)
					}
				}
				return r, nil
			},
		},
	},
	NotLike: {
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				m := match.New(lru.New(1 << 10))
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rgp, err := m.Compile(b.Vs[o], true)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = !rgp.MatchString(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						rgp, err := m.Compile(b.Vs[i], true)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = !rgp.MatchString(v)
					}
				}
				return r, nil
			},
		},
	},
	Match: {
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				m := match.New(lru.New(1 << 10))
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rgp, err := m.Compile(b.Vs[o], false)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rgp.MatchString(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						rgp, err := m.Compile(b.Vs[i], false)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rgp.MatchString(v)
					}
				}
				return r, nil
			},
		},
	},
	NotMatch: {
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_bool,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &static.Bools{
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				m := match.New(lru.New(1 << 10))
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rgp, err := m.Compile(b.Vs[o], false)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = !rgp.MatchString(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						rgp, err := m.Compile(b.Vs[i], false)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = !rgp.MatchString(v)
					}
				}
				return r, nil
			},
		},
	},
	Typecast: {
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					copy(r.Vs, a.Vs)
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_int,
			ReturnType: types.T_int,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Ints{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rv
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rv
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Int8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 8)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = int8(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 8)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = int8(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_int16,
			ReturnType: types.T_int16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Int16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 16)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = int16(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 16)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = int16(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_int32,
			ReturnType: types.T_int32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Int32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 32)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = int32(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 32)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = int32(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					copy(r.Vs, a.Vs)
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_int64,
			ReturnType: types.T_int64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Int64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = int64(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = int64(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil

			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_int8,
			ReturnType: types.T_int8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint8(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint8(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_uint8,
			ReturnType: types.T_uint8,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Uint8s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint8, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 8)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = uint8(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 8)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = uint8(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint16(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint16(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_uint16,
			ReturnType: types.T_uint16,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Uint16s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint16, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 16)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = uint16(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 16)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = uint16(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint32(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_uint32,
			ReturnType: types.T_uint32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Uint32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 32)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = uint32(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 32)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = uint32(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1
						} else {
							r.Vs[o] = 0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1
						} else {
							r.Vs[i] = 0
						}
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = uint64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = uint64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_uint64,
			ReturnType: types.T_uint64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Uint64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]uint64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseInt(a.Vs[0], 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = uint64(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseInt(v, 0, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = uint64(rv)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0.0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0.0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0.0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0.0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] == 0.0 {
							r.Vs[o] = false
						} else {
							r.Vs[o] = true
						}
					}
				} else {
					for i, v := range a.Vs {
						if v == 0.0 {
							r.Vs[i] = false
						} else {
							r.Vs[i] = true
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_bool,
			ReturnType: types.T_bool,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Bools{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]bool, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := value.ParseBool(a.Vs[o])
						if err != nil {
							return nil, err
						}
						r.Vs[o] = bool(*rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := value.ParseBool(v)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = bool(*rv)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1.0
						} else {
							r.Vs[o] = 0.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1.0
						} else {
							r.Vs[i] = 0.0
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_float,
			ReturnType: types.T_float,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseFloat(a.Vs[o], 64)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rv
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rv
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1.0
						} else {
							r.Vs[o] = 0.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1.0
						} else {
							r.Vs[i] = 0.0
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float32(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float32(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_float32,
			ReturnType: types.T_float32,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Float32s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float32, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseFloat(a.Vs[o], 32)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = float32(rv)
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseFloat(v, 32)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = float32(rv)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						if a.Vs[o] {
							r.Vs[o] = 1.0
						} else {
							r.Vs[o] = 0.0
						}
					}
				} else {
					for i, v := range a.Vs {
						if v {
							r.Vs[i] = 1.0
						} else {
							r.Vs[i] = 0.0
						}
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &static.Float64s{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = float64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = float64(v)
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_float64,
			ReturnType: types.T_float64,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Floats{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]float64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := strconv.ParseFloat(a.Vs[o], 64)
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rv
					}
				} else {
					for i, v := range a.Vs {
						rv, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rv
					}
				}
				return r, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = int64(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = int64(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_timestamp,
			ReturnType: types.T_timestamp,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*dynamic.Strings)
				r := &static.Timestamps{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]int64, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						rv, err := time.Parse(value.TimestampOutputFormat, a.Vs[o])
						if err != nil {
							return nil, err
						}
						r.Vs[o] = rv.Unix()
					}
				} else {
					for i, v := range a.Vs {
						rv, err := time.Parse(value.TimestampOutputFormat, v)
						if err != nil {
							return nil, err
						}
						r.Vs[i] = rv.Unix()
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Ints)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatInt(a.Vs[o], 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatInt(v, 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int8,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int8s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatInt(int64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatInt(int64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int16,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int16s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatInt(int64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatInt(int64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int32,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int32s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatInt(int64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatInt(int64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_int64,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Int64s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatInt(a.Vs[o], 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatInt(v, 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint8,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint8s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatUint(uint64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatUint(uint64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint16,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint16s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatUint(uint64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatUint(uint64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint32,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint32s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatUint(uint64(a.Vs[o]), 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatUint(uint64(v), 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_uint64,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Uint64s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatUint(a.Vs[o], 10)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatUint(v, 10)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Floats)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatFloat(a.Vs[o], 'f', -1, 64)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatFloat(v, 'f', -1, 64)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float32,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float32s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatFloat(float64(a.Vs[o]), 'f', -1, 32)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_float64,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Float64s)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatFloat(a.Vs[o], 'f', -1, 64)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatFloat(v, 'f', -1, 64)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_bool,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Bools)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = strconv.FormatBool(a.Vs[o])
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = strconv.FormatBool(v)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_timestamp,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				a := vs.(*static.Timestamps)
				r := &dynamic.Strings{
					Np: a.Np,
					Dp: a.Dp,
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = time.Unix(a.Vs[o], 0).UTC().Format(value.TimestampOutputFormat)
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = time.Unix(v, 0).UTC().Format(value.TimestampOutputFormat)
					}
				}
				return a, nil
			},
		},
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(vs, _ value.Values) (value.Values, error) {
				return vs, nil
			},
		},
	},
	Concat: {
		&BinOp{
			LeftType:   types.T_string,
			RightType:  types.T_string,
			ReturnType: types.T_string,
			Fn: func(as, bs value.Values) (value.Values, error) {
				a, b := as.(*dynamic.Strings), bs.(*dynamic.Strings)
				r := &dynamic.Strings{
					Is: a.Is,
					Vs: make([]string, len(a.Vs)),
				}
				{
					switch {
					case a.Np == nil && b.Np != nil:
						r.Np = b.Np
					case a.Np != nil && b.Np == nil:
						r.Np = a.Np
					case a.Np != nil && b.Np != nil:
						r.Np = a.Np.Union(b.Np)
					}
				}
				{
					switch {
					case a.Dp == nil && b.Dp != nil:
						r.Dp = b.Dp
					case a.Dp != nil && b.Dp == nil:
						r.Dp = a.Np
					case a.Dp != nil && b.Dp != nil:
						r.Dp = a.Dp.Union(b.Dp)
					}
				}
				if len(r.Is) > 0 {
					for _, o := range r.Is {
						r.Vs[o] = a.Vs[o] + b.Vs[o]
					}
				} else {
					for i, v := range a.Vs {
						r.Vs[i] = v + b.Vs[i]
					}
				}
				return r, nil
			},
		},
	},
}

var MultiOps = map[int][]*MultiOp{}
