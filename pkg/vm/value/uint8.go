package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint8(v uint8) *Uint8 {
	r := Uint8(v)
	return &r
}

func AsUint8(v interface{}) (Uint8, bool) {
	switch t := v.(type) {
	case *Uint8:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeUint8 attempts to retrieve a Uint8 from a value, panicking if the
// assertion fails.
func MustBeUint8(v interface{}) uint8 {
	i, ok := AsUint8(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return uint8(i)
}

func GetUint8(v Value) (Uint8, error) {
	if i, ok := v.(*Uint8); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Uint8)
}

func (a *Uint8) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint8) ResolvedType() *types.T {
	return types.Uint8
}

// ParseUint8 parses and returns the *Uint8 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint8(s string) (*Uint8, error) {
	i, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return nil, makeParseError(s, types.Uint8, err)
	}
	return NewUint8(uint8(i)), nil
}

func (a *Uint8) Compare(v Value) int {
	var x, y uint8

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = uint8(*a)
	switch b := v.(type) {
	case *Int:
		y = uint8(*b)
	case *Uint8:
		y = uint8(*b)
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

func (_ *Uint8) Size() int                              { return 2 }
func (_ *Uint8) IsLogical() bool                        { return false }
func (_ *Uint8) IsAndOnly() bool                        { return true }
func (_ *Uint8) Attributes() []string                   { return []string{} }
func (a *Uint8) Eval(_ map[string]Value) (Value, error) { return a, nil }
