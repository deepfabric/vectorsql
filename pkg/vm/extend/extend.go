package extend

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/types"
	"github.com/deepfabric/vectorsql/pkg/vm/util"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func (e *UnaryExtend) IsLogical() bool {
	return overload.IsLogical(e.Op)
}

func (e *UnaryExtend) IsAndOnly() bool {
	return !overload.IsLogical(e.Op)
}

func (e *UnaryExtend) Attributes() []string {
	return e.E.Attributes()
}

func (e *UnaryExtend) ReturnType() uint32 {
	switch e.Op {
	case overload.Not:
		return e.E.ReturnType()
	case overload.Abs:
		return e.E.ReturnType()
	case overload.Ceil:
		return e.E.ReturnType()
	case overload.Sign:
		return e.E.ReturnType()
	case overload.Floor:
		return e.E.ReturnType()
	case overload.Lower:
		return e.E.ReturnType()
	case overload.Round:
		return e.E.ReturnType()
	case overload.Upper:
		return e.E.ReturnType()
	case overload.Length:
		return types.T_int
	case overload.Typeof:
		return types.T_string
	case overload.UnaryMinus:
		return e.E.ReturnType()
	}
	return 0
}

func (e *UnaryExtend) Eval(mp map[string]value.Values) (value.Values, uint32, error) {
	vs, typ, err := e.E.Eval(mp)
	if err != nil {
		return nil, 0, err
	}
	return overload.UnaryEval(e.Op, typ, vs)
}

func (e *UnaryExtend) String() string {
	switch e.Op {
	case overload.Not:
		return fmt.Sprintf("not %s", e.E.String())
	case overload.Abs:
		return fmt.Sprintf("abs(%s)", e.E.String())
	case overload.Ceil:
		return fmt.Sprintf("ceil(%s)", e.E.String())
	case overload.Sign:
		return fmt.Sprintf("sign(%s)", e.E.String())
	case overload.Floor:
		return fmt.Sprintf("floor(%s)", e.E.String())
	case overload.Lower:
		return fmt.Sprintf("lower(%s)", e.E.String())
	case overload.Round:
		return fmt.Sprintf("round(%s)", e.E.String())
	case overload.Upper:
		return fmt.Sprintf("upper(%s)", e.E.String())
	case overload.Length:
		return fmt.Sprintf("length(%s)", e.E.String())
	case overload.Typeof:
		return fmt.Sprintf("typeof(%s)", e.E.String())
	case overload.UnaryMinus:
		return fmt.Sprintf("-%s", e.E.String())
	}
	return ""
}

func (e *BinaryExtend) IsLogical() bool {
	return overload.IsLogical(e.Op)
}

func (e *BinaryExtend) IsAndOnly() bool {
	return e.Op != overload.Or && e.Left.IsAndOnly() && e.Right.IsAndOnly()
}

func (e *BinaryExtend) Attributes() []string {
	return util.MergeAttributes(e.Left.Attributes(), e.Right.Attributes())
}

func (e *BinaryExtend) ReturnType() uint32 {
	switch e.Op {
	case overload.EQ:
		return types.T_bool
	case overload.LT:
		return types.T_bool
	case overload.GT:
		return types.T_bool
	case overload.LE:
		return types.T_bool
	case overload.GE:
		return types.T_bool
	case overload.NE:
		return types.T_bool
	case overload.Or:
		return types.T_bool
	case overload.And:
		return types.T_bool
	case overload.Div:
		lt, rt := e.Left.ReturnType(), e.Right.ReturnType()
		return returnType(lt, rt)
	case overload.Mod:
		lt, rt := e.Left.ReturnType(), e.Right.ReturnType()
		return returnType(lt, rt)
	case overload.Plus:
		lt, rt := e.Left.ReturnType(), e.Right.ReturnType()
		return returnType(lt, rt)
	case overload.Mult:
		lt, rt := e.Left.ReturnType(), e.Right.ReturnType()
		return returnType(lt, rt)
	case overload.Minus:
		lt, rt := e.Left.ReturnType(), e.Right.ReturnType()
		return returnType(lt, rt)
	case overload.Typecast:
		return e.Right.ReturnType()
	case overload.Like:
		return types.T_bool
	case overload.NotLike:
		return types.T_bool
	case overload.Match:
		return types.T_bool
	case overload.NotMatch:
		return types.T_bool
	case overload.Concat:
		return types.T_string
	}
	return 0
}

func (e *BinaryExtend) Eval(mp map[string]value.Values) (value.Values, uint32, error) {
	l, lt, err := e.Left.Eval(mp)
	if err != nil {
		return nil, 0, err
	}
	r, rt, err := e.Right.Eval(mp)
	if err != nil {
		return nil, 0, err
	}
	return overload.BinaryEval(e.Op, lt, rt, l, r)
}

func (e *BinaryExtend) String() string {
	switch e.Op {
	case overload.EQ:
		return fmt.Sprintf("%s = %s", e.Left.String(), e.Right.String())
	case overload.LT:
		return fmt.Sprintf("%s < %s", e.Left.String(), e.Right.String())
	case overload.GT:
		return fmt.Sprintf("%s > %s", e.Left.String(), e.Right.String())
	case overload.LE:
		return fmt.Sprintf("%s <= %s", e.Left.String(), e.Right.String())
	case overload.GE:
		return fmt.Sprintf("%s >= %s", e.Left.String(), e.Right.String())
	case overload.NE:
		return fmt.Sprintf("%s <> %s", e.Left.String(), e.Right.String())
	case overload.Or:
		return fmt.Sprintf("%s or %s", e.Left.String(), e.Right.String())
	case overload.And:
		return fmt.Sprintf("%s and %s", e.Left.String(), e.Right.String())
	case overload.Div:
		return fmt.Sprintf("%s / %s", e.Left.String(), e.Right.String())
	case overload.Mod:
		return fmt.Sprintf("%s %% %s", e.Left.String(), e.Right.String())
	case overload.Plus:
		return fmt.Sprintf("%s + %s", e.Left.String(), e.Right.String())
	case overload.Mult:
		return fmt.Sprintf("%s * %s", e.Left.String(), e.Right.String())
	case overload.Minus:
		return fmt.Sprintf("%s - %s", e.Left.String(), e.Right.String())
	case overload.Typecast:
		return fmt.Sprintf("cast(%s as %s)", e.Left.String(), e.Right.String())
	case overload.Like:
		return fmt.Sprintf("like(%s, %s)", e.Left.String(), e.Right.String())
	case overload.NotLike:
		return fmt.Sprintf("notLike(%s, %s)", e.Left.String(), e.Right.String())
	case overload.Match:
		return fmt.Sprintf("match(%s, %s)", e.Left.String(), e.Right.String())
	case overload.NotMatch:
		return fmt.Sprintf("notMatch(%s, %s)", e.Left.String(), e.Right.String())
	case overload.Concat:
		return fmt.Sprintf("%s ++ %s", e.Left.String(), e.Right.String())
	}
	return ""
}

func (e *MultiExtend) IsLogical() bool {
	return overload.IsLogical(e.Op)
}

func (e *MultiExtend) IsAndOnly() bool {
	for _, arg := range e.Args {
		if !arg.IsAndOnly() {
			return false
		}
	}
	return true
}

func (e *MultiExtend) Attributes() []string {
	var rs []string

	mp := make(map[string]struct{})
	for _, arg := range e.Args {
		attrs := arg.Attributes()
		for i, j := 0, len(attrs); i < j; i++ {
			if _, ok := mp[attrs[i]]; !ok {
				mp[attrs[i]] = struct{}{}
				rs = append(rs, attrs[i])
			}
		}
	}
	return rs
}

func (e *MultiExtend) ReturnType() uint32 {
	return 0
}

func (e *MultiExtend) Eval(mp map[string]value.Values) (value.Values, uint32, error) {
	var err error
	var typ uint32
	var arg value.Values
	var args []value.Values

	for _, v := range e.Args {
		if arg, typ, err = v.Eval(mp); err != nil {
			return nil, 0, err
		}
		args = append(args, arg)
	}
	return overload.MultiEval(e.Op, typ, args)
}

func (e *MultiExtend) String() string {
	return ""
}

func (e *ParenExtend) IsLogical() bool {
	return e.E.IsLogical()
}

func (e *ParenExtend) IsAndOnly() bool {
	return e.E.IsAndOnly()
}

func (e *ParenExtend) Attributes() []string {
	return e.E.Attributes()
}

func (e *ParenExtend) Eval(mp map[string]value.Values) (value.Values, uint32, error) {
	return e.E.Eval(mp)
}

func (e *ParenExtend) String() string {
	return "(" + e.E.String() + ")"
}

func (a *Attribute) IsLogical() bool {
	return false
}

func (a *Attribute) IsAndOnly() bool {
	return true
}

func (a *Attribute) Attributes() []string {
	return []string{a.Name}
}

func (a *Attribute) ReturnType() uint32 {
	return a.Type
}

func (a *Attribute) Eval(mp map[string]value.Values) (value.Values, uint32, error) {
	if vs, ok := mp[a.Name]; ok {
		return vs, a.Type, nil
	}
	return nil, 0, fmt.Errorf("attribute '%s' not exist", a.Name)
}

func (a *Attribute) String() string {
	return a.Name
}

func returnType(x, y uint32) uint32 {
	if x == y {
		return x
	}
	if x == types.T_int || x == types.T_float {
		return y
	}
	return x
}
