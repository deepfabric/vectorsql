package value

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func (_ Array) ResolvedType() *types.T                 { return types.Array }
func (_ Array) IsLogical() bool                        { return false }
func (_ Array) IsAndOnly() bool                        { return true }
func (_ Array) Attributes() []string                   { return []string{} }
func (a Array) Eval(_ map[string]Value) (Value, error) { return a, nil }

func (a Array) Size() int {
	size := 0
	for _, v := range a {
		size += v.Size()
	}
	return size
}

func (a Array) String() string {
	s := "["
	for i, v := range a {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%s", v)
	}
	s += "]"
	return s
}

func (a Array) Compare(v Value) int {
	if v == ConstNull {
		return 1
	}
	b, ok := v.(Array)
	if !ok {
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	if r := len(a) - len(b); r != 0 {
		if r < 0 {
			return -1
		}
		if r > 0 {
			return 1
		}

	}
	for i := range a {
		if r := int(a[i].ResolvedType().Oid - b[i].ResolvedType().Oid); r != 0 {
			return r
		}
		if r := a[i].Compare(b[i]); r != 0 {
			return r
		}
	}
	return 0
}
