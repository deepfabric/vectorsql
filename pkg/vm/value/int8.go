package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt8(v int8) *Int8 {
	r := Int8(v)
	return &r
}

func AsInt8(v interface{}) (Int8, bool) {
	switch t := v.(type) {
	case *Int8:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeInt8 attempts to retrieve a Int8 from a value, panicking if the
// assertion fails.
func MustBeInt8(v interface{}) int8 {
	i, ok := AsInt8(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return int8(i)
}

func GetInt8(v Value) (Int8, error) {
	if i, ok := v.(*Int8); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Int8)
}

func (a *Int8) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int8) ResolvedType() *types.T {
	return types.Int8
}

// ParseInt8 parses and returns the *Int8 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt8(s string) (*Int8, error) {
	i, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return nil, makeParseError(s, types.Int8, err)
	}
	return NewInt8(int8(i)), nil
}

func (a *Int8) Compare(v Value) int {
	var x, y int8

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = int8(*a)
	switch b := v.(type) {
	case *Int:
		y = int8(*b)
	case *Int8:
		y = int8(*b)
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

func (_ *Int8) Size() int                              { return 2 }
func (_ *Int8) IsLogical() bool                        { return false }
func (_ *Int8) IsAndOnly() bool                        { return true }
func (_ *Int8) Attributes() []string                   { return []string{} }
func (a *Int8) Eval(_ map[string]Value) (Value, error) { return a, nil }
