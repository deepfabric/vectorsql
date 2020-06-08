package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint32(v uint32) *Uint32 {
	r := Uint32(v)
	return &r
}

func AsUint32(v interface{}) (Uint32, bool) {
	switch t := v.(type) {
	case *Uint32:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeUint32 attempts to retrieve a Uint32 from a value, panicking if the
// assertion fails.
func MustBeUint32(v interface{}) uint32 {
	i, ok := AsUint32(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return uint32(i)
}

func GetUint32(v Value) (Uint32, error) {
	if i, ok := v.(*Uint32); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Uint32)
}

func (a *Uint32) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint32) ResolvedType() *types.T {
	return types.Uint32
}

// ParseUint32 parses and returns the *Uint32 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint32(s string) (*Uint32, error) {
	i, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return nil, makeParseError(s, types.Uint32, err)
	}
	return NewUint32(uint32(i)), nil
}

func (a *Uint32) Compare(v Value) int {
	var x, y uint32

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = uint32(*a)
	switch b := v.(type) {
	case *Int:
		y = uint32(*b)
	case *Uint32:
		y = uint32(*b)
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

func (_ *Uint32) Size() int                              { return 2 }
func (_ *Uint32) IsLogical() bool                        { return false }
func (_ *Uint32) IsAndOnly() bool                        { return true }
func (_ *Uint32) Attributes() []string                   { return []string{} }
func (a *Uint32) Eval(_ map[string]Value) (Value, error) { return a, nil }
