package value

import (
	"fmt"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewString(v string) *String {
	r := String(v)
	return &r
}

func AsString(v interface{}) (String, bool) {
	switch t := v.(type) {
	case *String:
		return *t, true
	default:
		return "", false
	}
}

func MustBeString(v interface{}) string {
	s, ok := AsString(v)
	if !ok {
		panic(fmt.Errorf("expected *String, found %T", v))
	}
	return string(s)
}

func GetString(v Value) (String, error) {
	if s, ok := v.(*String); ok {
		return *s, nil
	}
	return "", fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.String)
}

func (a *String) String() string {
	return fmt.Sprintf("'%v'", *a)
}

func (_ *String) ResolvedType() *types.T {
	return types.String
}

func (a *String) Compare(v Value) int {
	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
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

func (a *String) Size() int                              { return 1 + len(*a) }
func (_ *String) IsLogical() bool                        { return false }
func (_ *String) IsAndOnly() bool                        { return true }
func (_ *String) Attributes() []string                   { return []string{} }
func (a *String) Eval(_ map[string]Value) (Value, error) { return a, nil }
