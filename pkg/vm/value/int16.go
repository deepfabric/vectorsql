package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt16(v int16) *Int16 {
	r := Int16(v)
	return &r
}

func AsInt16(v interface{}) (Int16, bool) {
	switch t := v.(type) {
	case *Int16:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeInt16 attempts to retrieve a Int16 from a value, panicking if the
// assertion fails.
func MustBeInt16(v interface{}) int16 {
	i, ok := AsInt16(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return int16(i)
}

func GetInt16(v Value) (Int16, error) {
	if i, ok := v.(*Int16); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Int16)
}

func (a *Int16) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int16) ResolvedType() *types.T {
	return types.Int16
}

// ParseInt16 parses and returns the *Int16 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt16(s string) (*Int16, error) {
	i, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return nil, makeParseError(s, types.Int16, err)
	}
	return NewInt16(int16(i)), nil
}

func (a *Int16) Compare(v Value) int {
	var x, y int16

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = int16(*a)
	switch b := v.(type) {
	case *Int:
		y = int16(*b)
	case *Int16:
		y = int16(*b)
	default:
		panic(makeUnsupportedComparisonMessage(a, v))
	}
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	default:
		return 0
	}
}

func (_ *Int16) Size() int                              { return 2 }
func (_ *Int16) IsLogical() bool                        { return false }
func (_ *Int16) IsAndOnly() bool                        { return true }
func (_ *Int16) Attributes() []string                   { return []string{} }
func (a *Int16) Eval(_ map[string]Value) (Value, error) { return a, nil }
