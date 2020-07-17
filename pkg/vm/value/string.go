package value

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewString(v string) *String {
	r := String(v)
	return &r
}

func (a *String) String() string {
	return fmt.Sprintf("'%v'", *a)
}

func (_ *String) ResolvedType() types.T {
	return types.T_string
}

func (a *String) Compare(v Value) int {
	b, ok := v.(*String)
	if !ok {
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	if *a < *b {
		return -1
	}
	if *a > *b {
		return 1
	}
	return 0
}

func (a *String) Size() int            { return 1 + len(*a) }
func (_ *String) IsLogical() bool      { return false }
func (_ *String) IsAndOnly() bool      { return true }
func (_ *String) Attributes() []string { return []string{} }
