package value

import (
	"strconv"

	"github.com/deepfabric/vectorsql/pkg/vm/types"
)

func NewInt8(v int8) *Int8 {
	r := Int8(v)
	return &r
}

func (a *Int8) String() string {
	return strconv.FormatInt(int64(*a), 10)
}

func (_ *Int8) ResolvedType() types.T {
	return types.T_int8
}

// ParseInt8 parses and returns the *Int8 value represented by the provided
// string, or an error if parsing is unsuccessful.
func ParseInt8(s string) (*Int8, error) {
	i, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return nil, makeParseError(s, types.T_int8, err)
	}
	return NewInt8(int8(i)), nil
}

func (a *Int8) Compare(v Value) int {
	var x, y int8

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

func (_ *Int8) Size() int            { return 2 }
func (_ *Int8) IsLogical() bool      { return false }
func (_ *Int8) IsAndOnly() bool      { return true }
func (_ *Int8) Attributes() []string { return []string{} }
