package value

import (
	"fmt"
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewUint64(v uint64) *Uint64 {
	r := Uint64(v)
	return &r
}

func AsUint64(v interface{}) (Uint64, bool) {
	switch t := v.(type) {
	case *Uint64:
		return *t, true
	default:
		return 0, false
	}
}

// MustBeUint64 attempts to retrieve a Uint64 from a value, panicking if the
// assertion fails.
func MustBeUint64(v interface{}) uint64 {
	i, ok := AsUint64(v)
	if !ok {
		panic(fmt.Errorf("expected *Int, found %T", v))
	}
	return uint64(i)
}

func GetUint64(v Value) (Uint64, error) {
	if i, ok := v.(*Uint64); ok {
		return *i, nil
	}
	return 0, fmt.Errorf("cannot convert %s to type %s", v.ResolvedType(), types.Uint64)
}

func (a *Uint64) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Uint64) ResolvedType() *types.T {
	return types.Uint64
}

// ParseUint64 parses and returns the *Uint64 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseUint64(s string) (*Uint64, error) {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, makeParseError(s, types.Uint64, err)
	}
	return NewUint64(uint64(i)), nil
}

func (a *Uint64) Compare(v Value) int {
	var x, y uint64

	if v == ConstNull {
		return 1 // NULL is less than any non-NULL value
	}
	x = uint64(*a)
	switch b := v.(type) {
	case *Int:
		y = uint64(*b)
	case *Uint64:
		y = uint64(*b)
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

func (_ *Uint64) Size() int                              { return 2 }
func (_ *Uint64) IsLogical() bool                        { return false }
func (_ *Uint64) IsAndOnly() bool                        { return true }
func (_ *Uint64) Attributes() []string                   { return []string{} }
func (a *Uint64) Eval(_ map[string]Value) (Value, error) { return a, nil }
