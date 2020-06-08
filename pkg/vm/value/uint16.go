package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint16(v uint16) *Uint16 {
	r := Uint16(v)
	return &r
}

func AsUint16(v interface{}) (Uint16, bool) {
	switch t := v.(type) {
	case *Uint16:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeUint16 attempts to retrieve a Uint16 from a value, panicking if the
// assertion fails.
func MustBeUint16(v interface{}) uint16 {
	i, ok := AsUint16(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return uint16(i)
}

func GetUint16(v Value) (Uint16, error) {
	if i, ok := v.(*Uint16); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Uint16)
}

func (a *Uint16) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint16) ResolvedType() *types.T {
	return types.Uint16
}

// ParseUint16 parses and returns the *Uint16 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint16(s string) (*Uint16, error) {
	i, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return nil, makeParseError(s, types.Uint16, err)
	}
	return NewUint16(uint16(i)), nil
}

func (a *Uint16) Compare(v Value) int {
	var x, y uint16

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = uint16(*a)
	switch b := v.(type) {
	case *Int:
		y = uint16(*b)
	case *Uint16:
		y = uint16(*b)
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

func (_ *Uint16) Size() int                              { return 2 }
func (_ *Uint16) IsLogical() bool                        { return false }
func (_ *Uint16) IsAndOnly() bool                        { return true }
func (_ *Uint16) Attributes() []string                   { return []string{} }
func (a *Uint16) Eval(_ map[string]Value) (Value, error) { return a, nil }
